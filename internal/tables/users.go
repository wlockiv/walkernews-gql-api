package tables

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/internal/services"
	util2 "github.com/wlockiv/walkernews/pkg/util"
	"strings"
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

	var transactWriteItems []*dynamodb.TransactWriteItem

	// For the actual user record
	if newUserAV, err := dynamodbattribute.Marshal(newUser); err != nil {
		return nil, err
	} else {
		transactWriteItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           &ut.tableName,
				ConditionExpression: aws.String("attribute_not_exists(id)"),
				Item:                newUserAV.M,
			},
		}

		transactWriteItems = append(transactWriteItems, transactWriteItem)
	}

	// For the username hash (enforcing uniqueness)
	if usernameHashAV, err := dynamodbattribute.Marshal(map[string]string{"id": "username#" + username}); err != nil {
		return nil, err
	} else {
		transactWriteItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           &ut.tableName,
				ConditionExpression: aws.String("attribute_not_exists(id)"),
				Item:                usernameHashAV.M,
			},
		}

		transactWriteItems = append(transactWriteItems, transactWriteItem)
	}

	transactWriteItemsInput := &dynamodb.TransactWriteItemsInput{TransactItems: transactWriteItems}

	if _, err := ut.dynamodb.TransactWriteItems(transactWriteItemsInput); err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFail") {
			return nil, errors.New(`the requested username, '` + username + `', has already been claimed`)
		}

		return nil, err
	} else {
		return &User{ID: userId, Username: username}, nil
	}

}

func GetUserTable() *UserTable {
	table := UserTable{
		tableName: "walkernews-users",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
