package controllers

import (
	"github.com/guregu/dynamo"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
	"time"
)

type LinksTable struct {
	table *dynamo.Table
}

func (lt *LinksTable) Create(input model.NewLink) (*model.Link, error) {
	newLink := &model.Link{
		ID:        uuid.NewV4().String(),
		Title:     input.Title,
		Address:   input.Address,
		UserID:    input.UserID,
		CreatedAt: time.Now().UTC(),
	}

	if err := lt.table.Put(newLink).Run(); err != nil {
		return nil, err
	}

	return newLink, nil
}

func (lt *LinksTable) GetById(linkId string) (*model.Link, error) {
	var result *model.Link
	if err := lt.table.Get("id", linkId).One(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (lt *LinksTable) GetAll() ([]*model.Link, error) {
	var results []*model.Link
	if err := lt.table.Scan().All(&results); err != nil {
		return nil, err
	}

	return results, nil

}

func GetLinksTable() (*LinksTable, error) {
	dynamodbTable, err := New("walkernews-links")
	if err != nil {
		return nil, err
	}

	table := LinksTable{
		table: dynamodbTable,
	}

	return &table, nil
}
