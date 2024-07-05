package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "csvtojl <input CSV file> <output JSONL file>",
		Short: "Convert CSV to JSON Lines",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			inputFile := args[0]
			outputFile := args[1]
			if err := convertCSVToJSONLines(inputFile, outputFile); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func convertCSVToJSONLines(inputFile, outputFile string) error {
	csvFile, err := openFile(inputFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	jsonFile, err := createFile(outputFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	records, err := readCSV(csvFile)
	if err != nil {
		return err
	}

	headers := records[0]
	return writeJSONLines(jsonFile, headers, records[1:])
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	return file, nil
}

// Read csv in Go,
// https://gosamples.dev/read-csv/
// https://golangdocs.com/reading-and-writing-csv-files-in-golang

func readCSV(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %v", err)
	}
	return records, nil
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	return file, nil
}

func writeJSONLines(file *os.File, headers []string, records [][]string) error {
	for _, record := range records {
		// we could also use json.marshall but it was creating json line with different order in the column.
		// If we do not care about the order of the column, we can simply use the following commented out code, but I
		// want to create the json file which has same column order as csv file. Thus, I had to manually write the json.

		// data := make(map[string]string)
		// for i, value := range record {
		// 	  data[headers[i]] = value
		// }
		// jsonLine, err := json.Marshal(data)

		var jsonLine strings.Builder

		jsonLine.WriteString("{")

		for i, value := range record {
			jsonLine.WriteString(fmt.Sprintf("\"%s\":\"%s\"", headers[i], value))
			if i < len(headers)-1 {
				jsonLine.WriteString(",")
			}
		}

		jsonLine.WriteString("}\n")

		_, err := file.WriteString(jsonLine.String())

		if err != nil {
			return fmt.Errorf("failed to write JSON line: %v", err)
		}
	}
	return nil
}
