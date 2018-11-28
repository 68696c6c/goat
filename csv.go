package goat

import (
	"bufio"
	"encoding/csv"
	"os"
)

func OpenCSV(path string) (reader *csv.Reader, err error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader = csv.NewReader(bufio.NewReader(csvFile))
	return
}
