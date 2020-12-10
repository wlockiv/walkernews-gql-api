package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Table int

func New(tableName string) (*dynamo.Table, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	db := dynamo.New(awsSession, &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table(tableName)

	return &table, nil
}
