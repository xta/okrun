package main

import (
	"os"
	"strings"
	"testing"
)

func init() {
	testFilePath := "support/test/file.go"
	removeCommentLines(testFilePath)
}

func TestMain(t *testing.T) {
	testFilePath := "support/test/file.go"

	os.Args = []string{"./okrun", testFilePath}
	main()

	lines, err := readLines(testFilePath)
	if err != nil {
		t.Error("Error: was not able to read the test file.go.")
	}

	testLine := lines[3]
	expected := "\t//\t\"flag\""

	if testLine != expected {
		t.Error("Error: The file was not modified as expected by okrun.")
	}
}

func removeCommentLines(path string) error {
	lines, err := readLines(path)
	if err != nil {
		return err
	}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			lines[i] = strings.TrimLeft(line, "//")
		}
	}

	if err := writeLines(lines, path); err != nil {
		return err
	}

	formatFile(path)

	return nil
}
