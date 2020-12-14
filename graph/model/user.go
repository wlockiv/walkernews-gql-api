package model

type User struct {
	ID       string `fauna:"id"`
	Email    string `fauna:"email"`
	Username string `fauna:"username"`
}
