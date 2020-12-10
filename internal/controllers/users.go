package controllers

import (
	"github.com/guregu/dynamo"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/pkg/util"
	"strings"
)

type UserTable struct {
	table *dynamo.Table
}

type User struct {
	ID       string `json:"id" dynamo:"id"`
	Username string `json:"username" dynamo:"username"`
	Password string `json:"password" dynamo:"password"`
}

func (ut *UserTable) Create(input model.NewUser) (*model.User, error) {
	hashedPassword, err := util.HashPassword(input.Password)
	userId := strings.ToLower(input.Username)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:       userId,
		Username: input.Username,
		Password: hashedPassword,
	}

	if err := ut.table.Put(newUser).If("attribute_not_exists(id)").Run(); err != nil {
		return nil, err
	}

	return &model.User{ID: userId, Username: input.Username}, nil

}

func (ut *UserTable) GetById(userId string) (*model.User, error) {
	var result *model.User
	if err := ut.table.Get("id", userId).One(&result); err != nil {
		return nil, err
	}

	return result, nil

}


func (ut *UserTable) Authenticate(username, password string) (userId string, err error) {
	id := strings.ToLower(username)

	var result *User
	if err := ut.table.Get("ID", id).One(&result); err != nil {
		return "", nil
	}

	passwordCorrect := util.CheckPasswordHash(password, result.Password)
	if !passwordCorrect {
		return "", WrongUsernameOrPasswordError
	}

	return result.ID, nil
}

func GetUserTable() (*UserTable, error) {
	dynamodbTable, err := New("walkernews-users")
	if err != nil {
		return nil, err
	}

	table := UserTable{
		table: dynamodbTable,
	}

	return &table, nil
}
