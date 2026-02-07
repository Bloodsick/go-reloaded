package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Returns a formatted way of errors that may occur.
func ProcessAndSave(inputFile, outputFile string) {
	//Read the input file.
	content, err := ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	//Process the content.
	processed := ExtraSpaces(content)

	//Create the output file.
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	//Write to the file.
	writer := bufio.NewWriter(file)
	for _, line := range processed {
		//Write the line and a newline character.
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
	//Flush ensures all data is written to the disk.
	writer.Flush()
	fmt.Println("File saved successfully to", outputFile)
}

// Read the .txt file line by line and save it as a list.
func ReadFile(txt string) ([]string, error) {
	file, err := os.Open(txt)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var lines []string
	for {
		// ReadString reads until the delimiter or end of file (EOF).
		line, err := reader.ReadString('\n')
		if err != nil && line == "" { // EOF with empty line.
			break
		}
		// Remove the newline for processing, it will be added back later.
		line = strings.TrimSuffix(line, "\n")
		lines = append(lines, line)

		if err != nil { // EOF reached.
			break
		}
	}
	return lines, nil
}
