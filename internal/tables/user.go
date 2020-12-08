package tables

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/services"
)

type UserTable struct {
	tableName string
	dynamodb  *dynamodb.DynamoDB
}

func (ut *UserTable) Put(user *model.User) (err error) {
	av, err := dynamodbattribute.Marshal(user)
	if err != nil {
		fmt.Println("There was a problem marshalling a user: ", err)
		return err
	}

	dynamoInput := &dynamodb.PutItemInput{
		Item:      av.M,
		TableName: &ut.tableName,
	}

	if _, err := ut.dynamodb.PutItem(dynamoInput); err != nil {
		fmt.Println("There was a problem putting a user to the table: ", err)
		return err
	}

	return nil
}

func GetUserTable() *UserTable {
	table := UserTable{
		tableName: "walkernews-users",
		dynamodb:  services.NewDynamoService(),
	}

	return &table
}
