package core

import (
	"bufio"
	"fmt"
	"os"
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

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
