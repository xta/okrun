package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	rerunTargetFile := true

	var errorBytes bytes.Buffer
	captureErr := bufio.NewWriter(&errorBytes)

	if len(os.Args) < 2 {
		log.Fatal("Error: please provide a target file path for your gofile.")
	}
	targetPath := os.Args[1]

	for rerunTargetFile {
		if !rerunTargetFile {
			return
		}

		errorBytes.Reset()

		cmd := exec.Command("go", "run", targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = captureErr
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		cmd.Wait()
		captureErr.Flush()

		addressableErrors := processError(errorBytes.String())

		if len(addressableErrors) > 0 {
			fixErrors(addressableErrors)
		} else {
			rerunTargetFile = false
		}
	}
}

type unusedImportError struct {
	filePath, pkg string
	fileLine      int
}

func processError(body string) (unusedImportErrors []unusedImportError) {
	fixableErrors := errorLines(body)
	unusedImportErrors = processLines(fixableErrors)
	return
}

func errorLines(body string) (lines []string) {
	allLines := strings.Split(body, "\n")
	for _, lineValue := range allLines {
		if strings.Contains(lineValue, "imported and not used") {
			lines = append(lines, lineValue)
		}
	}
	return
}

func processLines(lines []string) (errors []unusedImportError) {
	for _, line := range lines {
		sections := strings.Split(line, ":")

		filePath := sections[0]
		fileLine, err := strconv.Atoi(sections[1])
		if err != nil {
			panic(err)
		}
		pkg := strings.TrimSpace(sections[3])

		anUnusedImportError := &unusedImportError{
			filePath: filePath,
			fileLine: fileLine,
			pkg:      pkg,
		}
		errors = append(errors, *anUnusedImportError)
	}
	return
}

func fixErrors(errors []unusedImportError) (success bool) {
	for _, anUnusedImportError := range errors {

		file, err := os.Open(anUnusedImportError.filePath)
		if err != nil {
			log.Fatal(err)
			return false
		}
		defer file.Close()

		lines, err := readLines(anUnusedImportError.filePath)
		if err != nil {
			log.Fatalf("readLines: %s", err)
			return false
		}

		linesPosition := anUnusedImportError.fileLine - 1
		lines[linesPosition] = "//" + lines[linesPosition]

		if err := writeLines(lines, anUnusedImportError.filePath); err != nil {
			log.Fatalf("writeLines: %s", err)
			return false
		}
		formatFile(anUnusedImportError.filePath)
	}

	return true
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
