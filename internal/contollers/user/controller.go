package user

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErrors "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func Create(newUser model.NewUser) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(f.Call("create_user",
		f.Arr{newUser.Email, newUser.Username, newUser.Password}))
	if err != nil {
		return nil, internalErrors.NewDBError("(User) Create", err)
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, internalErrors.NewUnmarshallError("user", err)
	}

	return &user, nil
}

func GetByRefV(refV f.RefV) (*model.User, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, internalErrors.NewDBError("(User) GetByRef", err)
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, internalErrors.NewUnmarshallError("user", err)
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
		return nil, internalErrors.NewDBError("(User) GetById", err)
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, internalErrors.NewUnmarshallError("user", err)
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
		return nil, internalErrors.NewDBError("(User) GetByEmail", err)
	}

	var user model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		return nil, internalErrors.NewUnmarshallError("user", err)
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
		return "", internalErrors.NewDBError("Login", err)
	}

	var token string
	if err := res.At(f.ObjKey("secret")).Get(&token); err != nil {
		return "", internalErrors.NewUnmarshallError("fdb token", err)
	}

	return token, nil
}

func GetCurrent(userKey string) (*model.User, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(f.Get(f.Identity()))
	if err != nil {
		err = internalErrors.NewDBError("Identity", err)
		return nil, err
	}

	var user *model.User
	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		err := internalErrors.NewUnmarshallError("User", err)
		return nil, err
	}

	return user, nil
}
