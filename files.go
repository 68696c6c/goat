package goat

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

func FileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func WriteFile(filePath, contents string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed create file '%s'", filePath)
	}

	if _, err = f.WriteString(contents); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	err = f.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to close file '%s'", filePath)
	}

	return nil
}

func ReadFileJSON[T any](filePath string, target *T) (*os.File, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read json file")
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(target)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json file")
	}
	return file, nil
}

func WriteFileJSON(filePath string, data any) error {
	contents, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	err = WriteFile(filePath, string(contents))
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}

func ReadFileCSV(path string) (reader *csv.Reader, err error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader = csv.NewReader(bufio.NewReader(csvFile))
	return
}

func WriteFileCSV(filePath string, rows [][]string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed create file '%s'", filePath)
	}

	w := csv.NewWriter(f)
	if err = w.WriteAll(rows); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	err = f.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to close file '%s'", filePath)
	}

	return nil
}

type CSVRowCallback func(line []string, eof bool) error

func HandleCSVRows(path string, skipHeaderRow bool, breakOnEOF bool, callback CSVRowCallback) error {
	reader, err := ReadFileCSV(path)
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
