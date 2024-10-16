package models

type MoveFotoReq struct {
	IdNo       string `json:"idNo"`
	FilePath   string `json:"filePath"`
	RenameFile string `json:"renameFile"`
}

type MoveFotoRes struct {
	FileId   string `json:"fileId"`
	FilePath string `json:"filePath"`
}
