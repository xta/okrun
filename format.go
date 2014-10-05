package main

import (
	"os"
	"os/exec"
)

func formatFile(filePath string) error {
	cmd := exec.Command("gofmt", "-w", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()

	return nil
}
