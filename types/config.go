package types

type ConfigPathType int

const (
	ConfigPathTypeDefault ConfigPathType = iota
	ConfigPathTypeRel
	ConfigPathTypeAbs
)

type Config struct {
	utils    GoatUtilsInterface
	fileName string
	filePath string
}

type ConfigInterface interface {
	FileName() string
	FilePath() string
}

func NewConfig(u GoatUtilsInterface, file string, filePath string) *Config {
	c := &Config{
		utils:    u,
		fileName: file,
		filePath: filePath,
	}
	return c
}

func (c *Config) FileName() string {
	return c.fileName
}

func (c *Config) FilePath() string {
	return c.filePath
}
