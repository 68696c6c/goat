package goat

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

type CSVCallback func(line []string, eof bool) error

func OpenCSV(path string) (reader *csv.Reader, err error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader = csv.NewReader(bufio.NewReader(csvFile))
	return
}

func HandleCSVRows(path string, skipHeaderRow bool, breakOnEOF bool, callback CSVCallback) error {
	reader, err := OpenCSV(path)
	if err != nil {
		return err
	}
	eof := false
	for i := 0; true; i++ {
		line, err := reader.Read()
		if skipHeaderRow && i == 0 {
			continue
		}
		if err == io.EOF {
			if breakOnEOF {
				break
			} else {
				eof = true
			}
		} else if err != nil {
			return err
		}
		err = callback(line, eof)
		if err != nil {
			return err
		}
		if eof {
			break
		}
	}
	return nil
}
