package main

import (
	"fmt"
	"os"

	core "github.com/Bloodsick/go-reloaded/core"
)

func main() {
	//Checking if there are enough arguments.
	if len(os.Args) != 3 {
		printUsage()
		return
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	fileInfo, err := os.Stat(inputFile)
	//Checking if file exists.
	if os.IsNotExist(err) {
		fmt.Printf("Error: The input file '%s' does not exist.\n", inputFile)
		return
	}
	//Checking if it's a file and not a directory.
	if fileInfo.IsDir() {
		fmt.Printf("Error: '%s' is a directory, not a text file.\n", inputFile)
		return
	}
	core.ProcessAndSave(inputFile, outputFile)
}

// Helper function to help the user correct their mistake.
func printUsage() {
	fmt.Println("Usage Error: Incorrect number of arguments.")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("Correct Format: go run . <input_file.txt> <output_file.txt>")
	fmt.Println("Example:        go run . sample.txt result.txt")
}
