package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatal("Error: unable to find your Go executable binary.", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Error: please provide a target file path for your gofile.")
	}
	targetPath := os.Args[1]
	args := []string{"go", "run", targetPath}

	currentEnv := os.Environ()

	err = syscall.Exec(binary, args, currentEnv)
	if err != nil {
		log.Fatal("Error: an issue occurred while executing your gofile.", err)
	}

}
