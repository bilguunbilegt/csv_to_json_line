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
	tempFile := createTempCSVFile(t, [][]string{{"name", "age", "city"}, {"Bilguun B", "36", "Chicago"}, {"Kevin B", "5", "Chicago"}})
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
	tempFile := createTempCSVFile(t, [][]string{{"name", "age", "city"}, {"Bilguun B", "36", "Chicago"}, {"Kevin B", "5", "Chicago"}})
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
		{"name", "age", "city"},
		{"Bilguun B", "36", "Chicago"},
		{"Kevin B", "5", "Chicago"},
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
		{"name", "age", "city"},
		{"Bilguun B", "36", "Chicago"},
		{"Kevin B", "5", "Chicago"},
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

	want := []map[string]string{
		{"name": "Bilguun B", "age": "36", "city": "Chicago"},
		{"name": "Kevin B", "age": "5", "city": "Chicago"},
	}

	if len(got) != len(want) {
		t.Fatalf("unexpected number of records: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		if !equal(got[i], want[i]) {
			t.Fatalf("unexpected output at record %d: got %v, want %v", i, got[i], want[i])
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
