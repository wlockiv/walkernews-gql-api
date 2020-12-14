package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/wlockiv/walkernews/graph/generated"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/auth"
	linkCtrl "github.com/wlockiv/walkernews/internal/contollers/link"
	userCtrl "github.com/wlockiv/walkernews/internal/contollers/user"
	"github.com/wlockiv/walkernews/pkg/jwt"
)

func (r *linkResolver) User(ctx context.Context, obj *model.Link) (*model.User, error) {
	usr, err := userCtrl.GetByRefV(obj.User)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	authCtx := auth.ForContext(ctx)
	if authCtx.UserKey == "" {
		return nil, errors.New("not authorized")
	}

	link, err := linkCtrl.Create(input, authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, id string) (string, error) {
	authCtx := auth.ForContext(ctx)
	err := linkCtrl.DeleteById(id, authCtx.UserKey)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	newUser, err := userCtrl.Create(input)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	userKey, err := userCtrl.GetUserKey(input.Email, input.Password)
	if err != nil {
		return "", nil
	}

	usr, err := userCtrl.GetByEmail(input.Email)
	if err != nil {
		return "", nil
	}

	token, err := jwt.GenerateToken(usr, userKey)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	links, err := linkCtrl.GetAll()
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *queryResolver) Link(ctx context.Context, id string) (*model.Link, error) {
	link, err := linkCtrl.GetById(id)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	authCtx := auth.ForContext(ctx)
	usr, err := userCtrl.GetCurrent(authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return usr, nil
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
