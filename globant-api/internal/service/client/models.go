package service

type FileModel struct {
	UserCode  string
	FileBytes []byte
	FileName  string
	FileType  string
}

type AuthResponse struct {
	UserName string
	Token    string
	UserCode string
}
