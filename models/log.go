package models

type StandardLog struct {
	CreatedAt string `json:"created_at"`
	Host 	  string `json:"host"`
	Process   string `json:"process"`
	Level     string `json:"level"`
	Text      string `json:"text"`
}