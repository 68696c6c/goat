package types

type ConfigPathType int

const (
	ConfigPathTypeDefault ConfigPathType = iota
	ConfigPathTypeRel
	ConfigPathTypeAbs
)

type Config struct {
	FilePath string
	FileName string
}

func NewConfig(file string, filePath string) *Config {
	c := &Config{
		FileName: file,
		FilePath: filePath,
	}
	return c
}
