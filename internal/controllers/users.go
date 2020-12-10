package controllers

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/services"
	"github.com/wlockiv/walkernews/pkg/util"
	"log"
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

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:       strings.ToLower(input.Username),
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

func (ut *UserTable) Authenticate(username, password string) (userId string, err error) {
	id := strings.ToLower(username)
	input := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &id},
		},
		TableName: &ut.tableName,
	}

	result, err := ut.dynamodb.GetItem(&input)
	if err != nil {
		return "", err
	}
	if result.Item == nil {
		log.Println("the username was not found")
		return "", WrongUsernameOrPasswordError
	}

	user := &User{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return "", err
	}

	correct := util.CheckPasswordHash(password, user.Password)
	if !correct {
		log.Println("the password was incorrect")
		return "", WrongUsernameOrPasswordError
	}

	return user.ID, nil
}

func GetUserTable() *UserTable {
	table := UserTable{
		tableName: "walkernews-users",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
