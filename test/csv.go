package testhelper

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

// LoadCSV reads a CSV file and returns rows as []map[string]string
// using the first row as header keys.
func LoadCSV(t *testing.T, path string) []map[string]string {
	t.Helper()

	f, err := os.Open("../../../testdata/" + path)
	if err != nil {
		t.Fatalf("failed to open CSV %s: %v", path, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV %s: %v", path, err)
	}

	if len(records) < 2 {
		t.Fatalf("CSV %s must have header + at least 1 data row", path)
	}

	headers := records[0]
	rows := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(map[string]string, len(headers))
		for i, header := range headers {
			row[header] = record[i]
		}
		rows = append(rows, row)
	}
	return rows
}

// DecimalFrom parses a string into decimal.Decimal, failing the test on error.
func DecimalFrom(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("invalid decimal %q: %v", s, err)
	}
	return d
}
