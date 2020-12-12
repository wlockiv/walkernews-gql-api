// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewLink struct {
	Title   string `json:"title"`
	Address string `json:"address"`
}

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	Token string `json:"token"`
}
