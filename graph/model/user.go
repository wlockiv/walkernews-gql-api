package model

type User struct {
	ID     string `fauna:"id"`
	Email  string `fauna:"email"`
	Handle string `fauna:"displayName"`
}
