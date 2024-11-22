package extractor

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

// Define the necessary structs to parse the JSON
type File struct {
	ListRecords []ListRecord `json:"list_records"`
}

type ListRecord struct {
	ID     string  `json:"id"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Text  string   `json:"text"`
	User  []string `json:"user"`
	Date  []string `json:"date"`
}

// ConvertJSONToCSV takes JSON data as input, extracts list records, and returns them as a CSV string
func ConvertJSONToCSV(jsonData []byte) (string, error) {
	// Unmarshal the JSON data into the File struct
	var fileData File
	if err := json.Unmarshal(jsonData, &fileData); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Create a buffer to store CSV data
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)

	// Write each record to the CSV
	for _, record := range fileData.ListRecords {
		// Initialize a map to store field values by key
		fieldMap := make(map[string]string)
		for _, field := range record.Fields {
			switch field.Key {
			case "name":
				fieldMap["name"] = field.Text
			case "date":
				if len(field.Date) > 0 {
					fieldMap["date"] = field.Date[0]
				}
			case "people":
				if len(field.User) > 0 {
					fieldMap["people"] = field.User[0]
				}
			}
		}

		// Write the record to the CSV
		row := []string{
			fieldMap["name"],
			fieldMap["date"],
			fieldMap["people"],
		}
		if err := csvWriter.Write(row); err != nil {
			return "", fmt.Errorf("error writing record to CSV: %w", err)
		}
	}

	// Flush the writer to ensure all data is written
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return buffer.String(), nil
}
