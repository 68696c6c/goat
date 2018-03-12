package goat

type configPathType int

const (
	configPathTypeDefault configPathType = iota
	configPathTypeRel
	configPathTypeAbs
)

type config struct {
	fileName string
	filePath string
}

type configInterface interface {
	FileName() string
	FilePath() string
}

func newConfig(file string, filePath string) *config {
	c := &config{
		fileName: file,
		filePath: filePath,
	}
	return c
}

func (c *config) FileName() string {
	return c.fileName
}

func (c *config) FilePath() string {
	return c.filePath
}
