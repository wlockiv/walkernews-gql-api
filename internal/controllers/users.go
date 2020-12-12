package controllers

import (
	"github.com/guregu/dynamo"
	"github.com/wlockiv/walkernews/graph/model"
)

type UsersController struct {
	table *dynamo.Table
}

type User struct {
	ID       string `json:"id" dynamo:"id"`
	Username string `json:"username" dynamo:"username"`
	Password string `json:"password" dynamo:"password"`
}

func (c *UsersController) GetById(userId string) (*model.User, error) {
	var result *model.User
	if err := c.table.Get("id", userId).One(&result); err != nil {
		return nil, err
	}

	return result, nil

}

func GetUserTable() (*UsersController, error) {
	dynamodbTable, err := New("walkernews-users")
	if err != nil {
		return nil, err
	}

	table := UsersController{
		table: dynamodbTable,
	}

	return &table, nil
}
