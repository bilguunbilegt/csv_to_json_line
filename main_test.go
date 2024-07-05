package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"testing"
)

// Helper function to create a temporary CSV file for testing
func createTempCSVFile(t *testing.T, content [][]string) *os.File {
	t.Helper()

	tempFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	writer := csv.NewWriter(tempFile)
	err = writer.WriteAll(content)
	if err != nil {
		t.Fatalf("Failed to write to temporary CSV file: %v", err)
	}

	writer.Flush()
	tempFile.Close()

	return tempFile
}

func TestOpenFile(t *testing.T) {
	tempFile := createTempCSVFile(t, [][]string{{"name", "age", "city", "state"}, {"Bilguun B", "36", "Chicago", "IL"}, {"Kevin B", "5", "Chicago", "IL"}})
	defer os.Remove(tempFile.Name())

	_, err := openFile(tempFile.Name())
	if err != nil {
		t.Errorf("openFile() error = %v", err)
	}
}

func TestCreateFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_*.jsonl")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = createFile(tempFile.Name())
	if err != nil {
		t.Errorf("createFile() error = %v", err)
	}
}

func TestReadCSV(t *testing.T) {
	tempFile := createTempCSVFile(t, [][]string{{"name", "age", "city", "state"}, {"Bilguun B", "36", "Chicago", "IL"}, {"Kevin B", "5", "Chicago", "IL"}})
	defer os.Remove(tempFile.Name())

	file, err := openFile(tempFile.Name())
	if err != nil {
		t.Fatalf("openFile() error = %v", err)
	}
	defer file.Close()

	records, err := readCSV(file)
	if err != nil {
		t.Errorf("readCSV() error = %v", err)
	}
	if len(records) != 3 {
		t.Errorf("readCSV() expected 3 records, got %d", len(records))
	}
}

func TestWriteJSONLines(t *testing.T) {
	tempCSVContent := [][]string{
		{"name", "age", "city", "state"},
		{"Bilguun B", "36", "Chicago", "IL"},
		{"Kevin B", "5", "Chicago", "IL"},
	}
	headers := tempCSVContent[0]
	records := tempCSVContent[1:]

	tempFile, err := os.CreateTemp("", "test_*.jsonl")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	err = writeJSONLines(tempFile, headers, records)
	if err != nil {
		t.Errorf("writeJSONLines() error = %v", err)
	}
}

func TestConvertCSVToJSONLines(t *testing.T) {
	tempCSV := createTempCSVFile(t, [][]string{
		{"name", "age", "city", "state"},
		{"Bilguun B", "36", "Chicago", "IL"},
		{"Kevin B", "5", "Chicago", "IL"},
	})
	defer os.Remove(tempCSV.Name())

	tempJSONL, err := os.CreateTemp("", "test_*.jsonl")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempJSONL.Name())

	err = convertCSVToJSONLines(tempCSV.Name(), tempJSONL.Name())
	if err != nil {
		t.Errorf("convertCSVToJSONLines() error = %v", err)
	}

	// Read the JSONL file and compare the contents
	jsonFile, err := os.Open(tempJSONL.Name())
	if err != nil {
		t.Fatalf("Failed to open JSONL file: %v", err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var got []map[string]string
	for decoder.More() {
		var obj map[string]string
		err := decoder.Decode(&obj)
		if err != nil {
			t.Fatalf("Failed to decode JSONL object: %v", err)
		}
		got = append(got, obj)
	}

	expect := []map[string]string{
		{"name": "Bilguun B", "age": "36", "city": "Chicago", "state": "IL"},
		{"name": "Kevin B", "age": "5", "city": "Chicago", "state": "IL"},
	}

	if len(got) != len(expect) {
		t.Fatalf("unexpected number of records: got %d, want %d", len(got), len(expect))
	}

	for i := range got {
		if !equal(got[i], expect[i]) {
			t.Fatalf("unexpected output at record %d: got %v, want %v", i, got[i], expect[i])
		}
	}
}

// Helper function to compare two maps
func equal(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
