package server

type Word struct {
	Word     string `json:"word"`
	ID       string `json:"id"`
	FileName string `json:"file_name"`
}
