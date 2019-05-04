package generator

type AuthorConfig struct {
	Name  string
	Email string
}

type ProjectConfig struct {
	Path    string
	SRCPath string
	Name    string
	License string
	Author  AuthorConfig
	Models  []*Model
}
