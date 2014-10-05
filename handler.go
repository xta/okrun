package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type unusedImportError struct {
	filePath, pkg string
	fileLine      int
}

// handleErrors handles/fixes error(s) if it can in the gofile
func handleErrors(body string) (handleable []unusedImportError, unhandleableLines []string) {
	handleableLines, unhandleableLines := distinguishLines(body)
	if len(unhandleableLines) == 0 {
		handleable = buildHandleables(handleableLines)
	}
	return
}

func distinguishLines(body string) (handleableLines, unhandleableLines []string) {
	allLines := strings.Split(body, "\n")
	for _, lineValue := range allLines {
		lineValue = strings.TrimSpace(lineValue)
		if len(lineValue) == 0 {
			// empty line, do nothing
		} else if strings.HasPrefix(lineValue, "#") {
			// comment line, do nothing
		} else if strings.Contains(lineValue, "imported and not used") {
			handleableLines = append(handleableLines, lineValue)
		} else {
			unhandleableLines = append(unhandleableLines, lineValue)
		}
	}
	return
}

func buildHandleables(lines []string) (errors []unusedImportError) {
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

func fixErrors(errors []unusedImportError) error {
	for _, anUnusedImportError := range errors {
		lines, err := readLines(anUnusedImportError.filePath)
		if err != nil {
			return err
		}

		linesPosition := anUnusedImportError.fileLine - 1
		lines[linesPosition] = "//" + lines[linesPosition]

		if err := writeLines(lines, anUnusedImportError.filePath); err != nil {
			return err
		}
		formatFile(anUnusedImportError.filePath)
	}

	return nil
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
