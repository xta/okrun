package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
)

func main() {
	rerunTargetFile := true

	var errorBytes bytes.Buffer
	captureErr := bufio.NewWriter(&errorBytes)

	if len(os.Args) < 2 {
		log.Fatal("Error: please provide a target gofile path.\nUsage: okrun <path/to/file.go> [arg]*")
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

		addressableErrors := handleErrors(errorBytes.String())

		if len(addressableErrors) > 0 {
			fixErrors(addressableErrors)
		} else {
			rerunTargetFile = false
		}
	}
}
