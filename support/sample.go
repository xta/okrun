package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	wordPointer := flag.String("word", "default", "a word test string")
	flag.Parse()
	fmt.Println("word:", *wordPointer)

	fmt.Println("hi")
}
