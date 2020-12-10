package tables

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/services"
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

// TODO: Make it so that this table takes all required fields as args?
func (ut *LinksTable) Create(input *model.NewLink) (*model.Link, error) {
	link := &model.Link{
		ID:      uuid.NewV4().String(),
		Title:   input.Title,
		Address: input.Address,
		UserID:  input.UserID,
	}

	av, err := dynamodbattribute.Marshal(link)
	if err != nil {
		fmt.Println("There was a problem marshalling a link: ", err)
		return nil, err
	}

	dynamoInput := &dynamodb.PutItemInput{
		Item:      av.M,
		TableName: &ut.tableName,
	}

	if _, err := ut.dynamodb.PutItem(dynamoInput); err != nil {
		fmt.Println("There was a problem putting a link to the table: ", err)
		return nil, err
	}

	return link, nil
}

func GetLinksTable() *LinksTable {
	table := LinksTable{
		tableName: "walkernews-links",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
