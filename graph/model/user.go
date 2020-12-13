package model

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"os"
)

type User struct {
	ID       string `fauna:"id"`
	Email    string `fauna:"email"`
	Username string `fauna:"username"`
}

func (u *User) Save(password string) error {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Create(
		f.Collection("users"), f.Obj{
			"credentials": f.Obj{"password": password},
			"data": f.Obj{
				"id":       f.NewId(),
				"email":    f.LowerCase(u.Email),
				"username": u.Username,
			},
		},
	))
	if err != nil {
		return err
	}

	var user User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return err
	} else {
		// Update the model
		u.ID = user.ID
		u.Username = user.Username
		u.Email = user.Email
	}

	return nil
}

func (u *User) GetByRefV(refV f.RefV) (*User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, err
	}

	var user User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetById(id string) (*User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("users_by_id"), id),
		),
	)
	if err != nil {
		return nil, err
	}

	var user User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("users_by_email"), email),
		),
	)
	if err != nil {
		return nil, err
	}

	var user User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetUserKey(email, password string) (string, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))

	res, err := client.Query(
		f.Login(
			f.MatchTerm(f.Index("users_by_email"), email),
			f.Obj{"password": password},
		),
	)
	if err != nil {
		return "", err
	}

	var token string
	if err := res.At(f.ObjKey("secret")).Get(&token); err != nil {
		return "", err
	}

	return token, nil
}
