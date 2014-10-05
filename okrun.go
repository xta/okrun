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
	commandArgs := []string{"run"}
	commandArgs = append(commandArgs, os.Args[1:]...)

	for rerunTargetFile {
		if !rerunTargetFile {
			return
		}

		errorBytes.Reset()

		cmd := exec.Command("go", commandArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = captureErr
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		cmd.Wait()
		captureErr.Flush()

		handleableErrors, unhandleableLines := handleErrors(errorBytes.String())

		if len(unhandleableLines) > 0 {
			rerunTargetFile = false

			errorMsg := "Error: the following errors cannot be automatically corrected by okrun to run your file:\n"
			for i := 0; i < len(unhandleableLines); i++ {
				errorMsg = errorMsg + "* " + unhandleableLines[i] + "\n"
			}
			errorMsg = errorMsg + "Please fix the error(s) listed above and then use okrun."
			log.Fatal(errorMsg)
		} else if len(handleableErrors) > 0 {
			fixErrors(handleableErrors)
		} else {
			rerunTargetFile = false
		}
	}
}
