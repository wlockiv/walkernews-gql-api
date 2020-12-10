package model

type User struct {
	ID       string `json:"id" dynamo:"id"`
	Username string `json:"username" dynamo:"username"`
}
