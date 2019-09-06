package store

import (
	"context"

	"a.com/go-server/service/board/model"
)

type Store interface {
	// Comment
	ListComments(ctx context.Context, oid, cid string, stamp int64, limit int) (model.Comments, error)
	GetComments(ctx context.Context, ids []string) (model.Comments, error)
	GetComment(ctx context.Context, id string) (*model.Comment, error)
	NewComment(ctx context.Context, pitem *model.Comment) error
	DelComment(ctx context.Context, id string) error
	IncrCommReply(ctx context.Context, cid string, pitem *model.Comment) error
	DecrCommReply(ctx context.Context, cid, rid string) error
	IncrCommLike(ctx context.Context, cid string) error
	DecrCommLike(ctx context.Context, cid string) error

	// Comment Like
	ListUserCommLikes(ctx context.Context, uid, oid string) (model.CommentLikes, error)
	NewCommLike(ctx context.Context, pitem *model.CommentLike) error
	DelCommLike(ctx context.Context, id string) error

	// Like
	ListUserLikes(ctx context.Context, uid string) (model.Likes, error)
	ListLikes(ctx context.Context, oid string, stamp int64, limit int) (model.Likes, error)
	GetLike(ctx context.Context, id string) (*model.Like, error)
	NewLike(ctx context.Context, pitem *model.Like) error
	DelLike(ctx context.Context, id string) error

	// Summary
	GetSummaries(ctx context.Context, ids []string) (model.Summaries, error)
	GetSummary(ctx context.Context, id string) (*model.Summary, error)
	IncrSummaryComm(ctx context.Context, cid, id string) error
	DecrSummaryComm(ctx context.Context, cid, id string, count int) error
	IncrSummaryLike(ctx context.Context, id string) error
	DecrSummaryLike(ctx context.Context, id string) error
	IncrSummaryRepost(ctx context.Context, id string) error
	DecrSummaryRepost(ctx context.Context, id string) error
}
