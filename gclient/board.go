package gclient

import (
	"context"
	"time"

	pb "a.com/go-server/proto/pb"
)

type BoardClient struct {
}

var Board *BoardClient

func (b *BoardClient) ListComments(c context.Context, args *pb.CommListArgs) (*pb.CommentInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListComments(ctx, args)
}

func (b *BoardClient) GetComment(c context.Context, args *pb.CommGetArgs) (*pb.CommentInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).GetComment(ctx, args)
}

func (b *BoardClient) NewComment(c context.Context, args *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewComment(ctx, args)
}

func (b *BoardClient) DelComment(c context.Context, args *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelComment(ctx, args)
}

func (b *BoardClient) LikeComment(c context.Context, args *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).LikeComment(ctx, args)
}

func (b *BoardClient) UnLikeComment(c context.Context, args *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).UnLikeComment(ctx, args)
}

func (b *BoardClient) ListLikes(c context.Context, args *pb.LikeListArgs) (*pb.LikeInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListLikes(ctx, args)
}

func (b *BoardClient) NewLike(c context.Context, args *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewLike(ctx, args)
}

func (b *BoardClient) DelLike(c context.Context, args *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelLike(ctx, args)
}

func (b *BoardClient) GetSummaries(c context.Context, args *pb.SummaryArgs) (*pb.SummaryInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).GetSummaries(ctx, args)
}

func (b *BoardClient) timeout() time.Duration {
	return 100 * time.Millisecond
}

func (b *BoardClient) service() string {
	return "board"
}
