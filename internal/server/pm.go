package server

type Word struct {
	Word     string `json:"word"`
	ID       string `json:"id"`
	FileName string `json:"file_name"`
}

type WordCache struct {
	ID         string `json:"id"`
	Expiration int64  `json:"expiration"`
	Filepath   string `json:"filepath"`
}
