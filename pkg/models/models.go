package models

type TtsModel struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type TtsResp struct {
	Url string `json:"url"`
}
