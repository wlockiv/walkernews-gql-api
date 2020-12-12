package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/wlockiv/walkernews/graph/generated"
	"github.com/wlockiv/walkernews/graph/model"
)

func (r *linkResolver) User(ctx context.Context, obj *model.Link) (*model.User, error) {
	var user model.User
	res, err := user.GetById(obj.UserID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	link := model.NewLinkModel(input.Title, input.Address, input.UserID)

	if err := link.Save(); err != nil {
		return nil, err
	}

	return link, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := model.User{
		Email:    input.Email,
		Username: input.Username,
	}

	if err := user.Save(input.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user model.User
	token, err := user.GetToken(input.Email, input.Password)
	if err != nil {
		return "", nil
	}

	return token, nil

	//table, err := controllers.GetUserTable()
	//if err != nil {
	//	return "", err
	//}
	//
	//userId, err := table.Authenticate(input.Email, input.Password)
	//if err != nil {
	//	return "", err
	//}
	//
	//token, err := jwt.GenerateToken(userId)
	//if err != nil {
	//	return "", err
	//}
	//
	//return token, nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var link model.Link
	links, err := link.GetAll()
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *queryResolver) Link(ctx context.Context, id string) (*model.Link, error) {
	var link model.Link
	if err := link.GetById(id); err != nil {
		return nil, err
	}

	return &link, nil
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
