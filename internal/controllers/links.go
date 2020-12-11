package controllers

import (
	"github.com/guregu/dynamo"
	uuid "github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/model"
	"time"
)

type LinksController struct {
	table *dynamo.Table
}

func (c *LinksController) Create(input model.NewLink) (*model.Link, error) {
	newLink := &model.Link{
		ID:        uuid.NewV4().String(),
		Title:     input.Title,
		Address:   input.Address,
		UserID:    input.UserID,
		CreatedAt: time.Now().UTC(),
	}

	if err := c.table.Put(newLink).Run(); err != nil {
		return nil, err
	}

	return newLink, nil
}

func (c *LinksController) GetById(linkId string) (*model.Link, error) {
	var result *model.Link
	if err := c.table.Get("id", linkId).One(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *LinksController) GetAll() ([]*model.Link, error) {
	var results []*model.Link
	if err := c.table.Scan().All(&results); err != nil {
		return nil, err
	}

	return results, nil

}

func GetLinksTable() (*LinksController, error) {
	dynamodbTable, err := New("walkernews-links")
	if err != nil {
		return nil, err
	}

	table := LinksController{
		table: dynamodbTable,
	}

	return &table, nil
}
