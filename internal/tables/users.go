package tables

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
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

func (ut *UserTable) Create(input model.NewUser) (*model.User, error) {
	userId := uuid.NewV4().String()

	hashedPassword, err := util2.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:       userId,
		Username: input.Username,
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
	enforcementValue := map[string]string{"id": "username#" + strings.ToLower(input.Username)}
	if usernameHashAV, err := dynamodbattribute.Marshal(enforcementValue); err != nil {
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
			return nil, errors.New(`the requested username, '` + input.Username + `', has already been claimed`)
		}

		return nil, err
	} else {
		return &model.User{ID: userId, Username: input.Username}, nil
	}
}

func (ut *UserTable) GetById(userId string) (*model.User, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: &ut.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &userId},
		},
	}

	result, err := ut.dynamodb.GetItem(getItemInput)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, NotFoundError
	}

	user := model.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserTable() *UserTable {
	table := UserTable{
		tableName: "walkernews-users",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
