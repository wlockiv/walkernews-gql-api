package user

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	"os"
)

func Create(newUser model.NewUser) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Create(
		f.Collection("users"), f.Obj{
			"credentials": f.Obj{"password": newUser.Password},
			"data": f.Obj{
				"id":       f.NewId(),
				"email":    f.LowerCase(newUser.Email),
				"username": newUser.Username,
			},
		},
	))
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetByRefV(refV f.RefV) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetById(id string) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("users_by_id"), id),
		),
	)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetByEmail(email string) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("users_by_email"), f.LowerCase(email)),
		),
	)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserKey(email, password string) (string, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Login(
			f.MatchTerm(f.Index("users_by_email"), f.LowerCase(email)),
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

func GetCurrent(userKey string) (*model.User, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(f.Get(f.Identity()))
	if err != nil {
		return nil, err
	}

	var user *model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, err
	}

	return user, nil
}
