package comment

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErrors "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func Create(newComment model.NewComment, userKey string) (*model.Comment, error) {
	client := f.NewFaunaClient(userKey)
	var res f.Value
	var err error
	switch newComment.ParentType {
	case model.PostTypeLink:
		res, err = client.Query(
			f.Call("create_comment_on_link",
				f.Arr{newComment.Content, newComment.ParentID}))
	case model.PostTypeComment:
		res, err = client.Query(
			f.Call("create_comment_on_comment",
				f.Arr{newComment.Content, newComment.ParentID}))
	}
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) Create on "+newComment.ParentType.String(), err)
	}

	var comment *model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
		return nil, internalErrors.NewUnmarshallError("comment", err)
	}

	return comment, nil
}

func GetByParentIdAndType(parentId string, parentType model.PostType) ([]*model.Comment, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(
		f.Call("get_comment_by_parent_id_and_type",
			f.Arr{parentId, parentType.String()}))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByParentIdAndType", err)
	}

	var comments []*model.Comment
	if err := res.At(f.ObjKey("data")).Get(&comments); err != nil {
		return nil, internalErrors.NewUnmarshallError("comments", err)
	}

	return comments, nil
}

// ! Deprecate
//func GetByParentLinkId(linkId string) ([]*model.Comment, error) {
//	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
//	res, err := client.Query(f.Map(
//		f.Paginate(f.Join(
//			f.MatchTerm(f.Index("link_ref_by_id"), linkId),
//			f.Index("comment_ref_by_parent"))),
//		f.Lambda("ref", f.Select("data", f.Get(f.Var("ref")))),
//	))
//	if err != nil {
//		return nil, internalErrors.NewDBError("(Comment) GetByParentLinkId", err)
//	}
//
//	var comments []*model.Comment
//	if err := res.At(f.ObjKey("data")).Get(&comments); err != nil {
//		return nil, internalErrors.NewUnmarshallError("comments", err)
//	}
//
//	return comments, nil
//}

// ! Deprecate
//func GetByParentCommentId(commentId string) ([]*model.Comment, error) {
//	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
//	res, err := client.Query(f.Map(
//		f.Paginate(f.Join(
//			f.MatchTerm(f.Index("comment_ref_by_id"), commentId),
//			f.Index("comment_ref_by_parent"))),
//		f.Lambda("ref", f.Select("data", f.Get(f.Var("ref")))),
//	))
//	if err != nil {
//		return nil, internalErrors.NewDBError("(Comment) GetByParentCommentId", err)
//	}
//
//	var comments []*model.Comment
//	if err := res.At(f.ObjKey("data")).Get(&comments); err != nil {
//		return nil, internalErrors.NewUnmarshallError("comments", err)
//	}
//
//	return comments, nil
//}

// ! Deprecate
//func GetByParent(parentRef f.RefV) (*model.Comment, error) {
//	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
//	res, err := client.Query(f.Map(
//		f.Paginate(f.MatchTerm(f.Index("comment_ref_by_parent"), parentRef)),
//		f.Lambda(
//			"parentRef",
//			f.Select("data", f.Select("data", f.Get(f.Var("parentRef"))))),
//	))
//	if err != nil {
//		return nil, internalErrors.NewDBError("(Comment) GetByParent", err)
//	}
//
//	var comment *model.Comment
//	if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
//		return nil, internalErrors.NewUnmarshallError("comment", err)
//	}
//
//	return comment, nil
//}

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
