package tables

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/internal/services"
	util2 "github.com/wlockiv/walkernews/pkg/util"
)

type UserTable struct {
	tableName string
	dynamodb  *dynamodb.DynamoDB
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO: Should this be taking UN & Password directly?
func (ut *UserTable) Put(username, password string) (*User, error) {
	userId := uuid.NewV4().String()
	hashedPassword, err := util2.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:       userId,
		Username: username,
		Password: hashedPassword,
	}

	av, err := dynamodbattribute.Marshal(newUser)
	if err != nil {
		return nil, err
	}

	dynamoInput := &dynamodb.PutItemInput{
		Item:      av.M,
		TableName: &ut.tableName,
	}

	if _, err := ut.dynamodb.PutItem(dynamoInput); err != nil {
		return nil, err
	}

	return &User{ID: userId, Username: username}, nil
}

func GetUserTable() *UserTable {
	table := UserTable{
		tableName: "walkernews-users",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
