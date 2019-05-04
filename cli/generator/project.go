package generator

type AuthorConfig struct {
	Name  string
	Email string
}

type ProjectConfig struct {
	Name    string
	License string
	Author  AuthorConfig

	Module     string
	DirName    string
	SRCPath    string
	AppPath    string
	CMDPath    string
	ModelsPath string
	ReposPath  string

	Models []*Model
	Repos  []*Repo
}
