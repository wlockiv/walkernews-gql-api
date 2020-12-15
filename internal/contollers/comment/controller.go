package comment

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErrors "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func Create(newComment model.NewComment, userKey string) (*model.Comment, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(
		f.Let().Bind(
			"linkRef",
			f.Select("ref",
				f.Get(f.MatchTerm(f.Index("link_ref_by_id"), newComment.LinkID))),
		).Bind(
			"commentRef",
			f.If(f.Exists(f.MatchTerm(f.Index("comment_ref_by_id"), newComment.ParentID)),
				f.Select("ref", f.Get(f.MatchTerm(f.Index("comment_ref_by_id"), newComment.ParentID))),
				f.Null()),
		).In(f.Create(
			f.Collection("comments"), f.Obj{
				"data": f.Obj{
					"id":      f.NewId(),
					"content": newComment.Content,
					"link":    f.Var("linkRef"),
					"parent":  f.Var("commentRef"),
					"user":    f.Identity(),
				},
			}),
		),
	)
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) Create", err)
	}

	var comment *model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
		return nil, internalErrors.NewUnmarshallError("comment", err)
	}

	return comment, nil
}

func GetByParentLinkId(linkId string) ([]*model.Comment, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(f.Map(
		f.Paginate(f.Join(
			f.MatchTerm(f.Index("link_ref_by_id"), linkId),
			f.Index("comment_ref_by_parent"))),
		f.Lambda("ref", f.Select("data", f.Get(f.Var("ref")))),
	))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByParentLinkId", err)
	}

	var comments []*model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comments); err != nil {
		return nil, internalErrors.NewUnmarshallError("comments", err)
	}

	return comments, nil
}

func GetByParentCommentId(commentId string) ([]*model.Comment, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(f.Map(
		f.Paginate(f.Join(
			f.MatchTerm(f.Index("comment_ref_by_id"), commentId),
			f.Index("comment_ref_by_parent"))),
		f.Lambda("ref", f.Select("data", f.Get(f.Var("ref")))),
	))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByParentCommentId", err)
	}

	var comments []*model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comments); err != nil {
		return nil, internalErrors.NewUnmarshallError("comments", err)
	}

	return comments, nil
}

func GetByParent(parentRef f.RefV) (*model.Comment, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(f.Map(
		f.Paginate(f.MatchTerm(f.Index("comment_ref_by_parent"), parentRef)),
		f.Lambda(
			"parentRef",
			f.Select("data", f.Select("data", f.Get(f.Var("parentRef"))))),
	))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByParent", err)
	}

	var comment *model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
		return nil, internalErrors.NewUnmarshallError("comment", err)
	}

	return comment, nil
}

func GetByRefV(refV f.RefV) (*model.Comment, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByRefV", err)
	}

	var comment *model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
		return nil, internalErrors.NewUnmarshallError("comment", err)
	}

	return comment, nil
}
