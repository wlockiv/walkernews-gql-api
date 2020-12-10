package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/wlockiv/walkernews/graph/generated"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/tables"
)

func (r *linkResolver) User(ctx context.Context, obj *model.Link) (*model.User, error) {
	table := tables.GetUserTable()

	user, err := table.GetById(obj.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	table := tables.GetLinksTable()
	newLink, err := table.Create(input)
	if err != nil {
		return nil, err
	}

	return newLink, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	table := tables.GetUserTable()

	if newUser, err := table.Create(input); err != nil {
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
	table := tables.GetLinksTable()
	links, err := table.GetAll()
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *queryResolver) Link(ctx context.Context, id string) (*model.Link, error) {
	table := tables.GetLinksTable()
	link, err := table.GetById(id)
	if err != nil {
		return nil, err
	}

	return link, err
}

// Link returns generated.LinkResolver implementation.
func (r *Resolver) Link() generated.LinkResolver { return &linkResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type linkResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
