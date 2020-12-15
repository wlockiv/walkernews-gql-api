package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/wlockiv/walkernews/graph/generated"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/internal/auth"
	commentCtrl "github.com/wlockiv/walkernews/internal/contollers/comment"
	linkCtrl "github.com/wlockiv/walkernews/internal/contollers/link"
	userCtrl "github.com/wlockiv/walkernews/internal/contollers/user"
	internalErr "github.com/wlockiv/walkernews/internal/errors"
	"github.com/wlockiv/walkernews/pkg/jwt"
)

func (r *commentResolver) Parent(ctx context.Context, obj *model.Comment) (model.Post, error) {
	comment, err := commentCtrl.GetByRefV(obj.Parent)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	usr, err := userCtrl.GetByRefV(obj.User)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *commentResolver) Comments(ctx context.Context, obj *model.Comment) ([]*model.Comment, error) {
	comments, err := commentCtrl.GetByParentCommentId(obj.ID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *linkResolver) User(ctx context.Context, obj *model.Link) (*model.User, error) {
	usr, err := userCtrl.GetByRefV(obj.User)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *linkResolver) Comments(ctx context.Context, obj *model.Link) ([]*model.Comment, error) {
	comments, err := commentCtrl.GetByParentLinkId(obj.ID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	authCtx, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	} else if authCtx.User == nil {
		return nil, internalErr.NewAuthError(errors.New("must be logged in to create a link"))
	}

	if authCtx.UserKey == "" {
		return nil, errors.New("not authorized")
	}

	link, err := linkCtrl.Create(input, authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, id string) (*model.Link, error) {
	authCtx, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	} else if authCtx.User == nil {
		notLoggedInErr := errors.New("must be logged in to delete a link")
		return nil, internalErr.NewAuthError(notLoggedInErr)
	}

	link, err := linkCtrl.DeleteById(id, authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return link, nil
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

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	// TODO: Add conditional logic for PostType
	authCtx, err := auth.ForContext(ctx)
	if err != nil {
		return nil, internalErr.NewAuthError(err)
	} else if authCtx.User == nil {
		notLoggedInErr := errors.New("no user currently logged in")
		return nil, internalErr.NewAuthError(notLoggedInErr)
	}

	comment, err := commentCtrl.Create(input, authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return comment, nil
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
	authCtx, err := auth.ForContext(ctx)
	if err != nil {
		return nil, internalErr.NewAuthError(err)
	} else if authCtx.User == nil {
		notLoggedInErr := errors.New("no user currently logged in")
		return nil, internalErr.NewAuthError(notLoggedInErr)
	}

	usr, err := userCtrl.GetCurrent(authCtx.UserKey)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Link returns generated.LinkResolver implementation.
func (r *Resolver) Link() generated.LinkResolver { return &linkResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type commentResolver struct{ *Resolver }
type linkResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *commentResolver) Link(ctx context.Context, obj *model.Comment) (*model.Link, error) {
	link, err := linkCtrl.GetByRefV(obj.Link)
	if err != nil {
		return nil, err
	}

	return link, nil
}
