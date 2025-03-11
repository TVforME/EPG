package model

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// LoadCSVRecords loads CSV records from a file
func loadCSVRecords(filename string) ([][]string, error) {
	return readCSVFile(filename)
}

// ReadCSVFile reads a CSV file and returns the records
func readCSVFile(filename string) ([][]string, error) {

	csvPath := fmt.Sprintf("./csv/%s", filename)

	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row and discard it
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// parseInt parses a string to an integer
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// parseBool parses a string to a boolean
func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
