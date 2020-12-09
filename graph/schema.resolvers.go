package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/satori/go.uuid"
	"github.com/wlockiv/walkernews/graph/generated"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/tables"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	table := tables.GetLinksTable()
	var err error

	newLink := &model.Link{
		ID:      uuid.NewV4().String(),
		Title:   input.Title,
		Address: input.Address,
		User:    &model.User{Username: "walker"},
	}

	fmt.Println(table, newLink)

	return newLink, err
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	table := tables.GetUserTable()

	if newUser, err := table.Put(input.Username, input.Password); err != nil {
		fmt.Println("There was a problem creating the user: ")
		fmt.Println(err.Error())
		return nil, err
	} else {
		return &model.User{ID: newUser.ID, Username: newUser.Username}, nil
	}
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
