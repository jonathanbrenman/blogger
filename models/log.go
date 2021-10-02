package models

type StandardLog struct {
	CreatedAt string `json:"created_at"`
	Level     string `json:"level"`
	Text      string `json:"text"`
}