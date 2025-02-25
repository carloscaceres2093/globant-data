package config

type Config struct {
	AppName    string
	Env        string
	UploadFile string
	Auth       string
}

type Configs struct {
	Scope map[string]Config
}
