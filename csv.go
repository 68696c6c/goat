package goat

import (
	"bufio"
	"encoding/csv"
	"io"
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

func HandleCSVRows(path string, skipHeaderRow bool, breakOnEOF bool, callback func(line []string) error) error {
	reader, err := OpenCSV(path)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		line, err := reader.Read()
		if skipHeaderRow && i == 0 {
			continue
		}
		if breakOnEOF && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = callback(line)
		if err != nil {
			return err
		}
	}
	return nil
}
