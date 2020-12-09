package tables

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/services"
)

type LinksTable struct {
	tableName string
	dynamodb  *dynamodb.DynamoDB
}

// TODO: Make it so that this table takes all required fields as args?
func (ut *LinksTable) Put(link *model.Link) (err error) {
	av, err := dynamodbattribute.Marshal(link)
	if err != nil {
		fmt.Println("There was a problem marshalling a link: ", err)
		return err
	}

	dynamoInput := &dynamodb.PutItemInput{
		Item:      av.M,
		TableName: &ut.tableName,
	}

	if _, err := ut.dynamodb.PutItem(dynamoInput); err != nil {
		fmt.Println("There was a problem putting a link to the table: ", err)
		return err
	}

	return nil
}

func GetLinksTable() *LinksTable {
	table := LinksTable{
		tableName: "walkernews-links",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
