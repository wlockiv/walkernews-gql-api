package tables

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/services"
	"time"
)

type LinksTable struct {
	tableName string
	dynamodb  *dynamodb.DynamoDB
}

type Link struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	UserId  string `json:"userId"`
}

func (lt *LinksTable) Create(input model.NewLink) (*model.Link, error) {
	link := &model.Link{
		ID:        uuid.NewV4().String(),
		Title:     input.Title,
		Address:   input.Address,
		UserID:    input.UserID,
		CreatedAt: time.Now().UTC(),
	}

	av, err := dynamodbattribute.Marshal(link)
	if err != nil {
		fmt.Println("There was a problem marshalling a link: ", err)
		return nil, err
	}

	dynamoInput := &dynamodb.PutItemInput{
		Item:      av.M,
		TableName: &lt.tableName,
	}

	if _, err := lt.dynamodb.PutItem(dynamoInput); err != nil {
		fmt.Println("There was a problem putting a link to the table: ", err)
		return nil, err
	}

	return link, nil
}

func (lt *LinksTable) GetById(linkId string) (*model.Link, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: &lt.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &linkId},
		},
	}

	result, err := lt.dynamodb.GetItem(getItemInput)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, NotFoundError
	}

	link := model.Link{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (lt *LinksTable) GetAll() ([]*model.Link, error) {

	scanInput := &dynamodb.ScanInput{
		TableName: &lt.tableName,
	}

	result, err := lt.dynamodb.Scan(scanInput)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, errors.New("the user could not be found")
	}

	var links []*model.Link
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func GetLinksTable() *LinksTable {
	table := LinksTable{
		tableName: "walkernews-links",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
