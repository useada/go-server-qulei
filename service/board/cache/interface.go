package cache

import (
	"context"

	"a.com/go-server/service/board/model"
)

type Cache interface {
	// Comments
	InitComments(ctx context.Context, oid, cid string, items model.Comments, total bool) error
	GetComments(ctx context.Context, oid string, ids []string) (model.Comments, error)
	GetComment(ctx context.Context, oid, id string) (*model.Comment, error)
	SetComment(ctx context.Context, pitem *model.Comment) error
	DelComment(ctx context.Context, oid, id string) error

	PushComment(ctx context.Context, pitem *model.Comment) error
	PopComment(ctx context.Context, oid, cid, id string) error
	RangeComments(ctx context.Context, oid, cid string, stamp int64, limit int) ([]string, error)

	// Likes
	ListUserLikes(ctx context.Context, uid string) (model.Likes, error)
	NewUserLikes(ctx context.Context, uid string, items model.Likes) error
	DelUserLikes(ctx context.Context, uid string) error

	// Comment Likes
	ListUserCommLikes(ctx context.Context, uid string) (model.CommentLikes, error)
	NewUserCommLikes(ctx context.Context, uid string, items model.CommentLikes) error
	DelUserCommLikes(ctx context.Context, uid string) error

	// Summarys
	GetSummaries(ctx context.Context, oids []string) (model.Summaries, error)
	NewSummary(ctx context.Context, pitem *model.Summary) error
	DelSummary(ctx context.Context, oid string) error
}
