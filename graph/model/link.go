package model

type Link struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	UserID    string `json:"user"`
}
