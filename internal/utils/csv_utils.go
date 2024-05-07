package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func ReadCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %v", filename, err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading records from file %s: %v", filename, err)
		return nil, err
	}

	var books []string
	for _, record := range records {
		books = append(books, record[0])
	}
	return books, nil
}

func AppendToCSV(filename string, record []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// RemoveFromCSV removes a book from the specified CSV file by book name.
func RemoveFromCSV(filename, bookName string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read existing records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Open the same file in write mode to overwrite
	file, err = os.Create(filename) // Using os.Create to truncate the file
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write back only those records that do not match the book name
	for _, record := range records {
		if strings.ToLower(record[0]) != strings.ToLower(bookName) {
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}
