package gclient

import (
	"context"
	"time"

	pb "a.com/go-server/proto/pb"
)

type BoardClient struct {
}

var Board *BoardClient

func (b *BoardClient) ListComments(args *pb.CommListArgs) (*pb.CommentInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListComments(ctx, args)
}

func (b *BoardClient) GetComment(args *pb.CommGetArgs) (*pb.CommentInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).GetComment(ctx, args)
}

func (b *BoardClient) NewComment(args *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewComment(ctx, args)
}

func (b *BoardClient) DelComment(args *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelComment(ctx, args)
}

func (b *BoardClient) LikeComment(args *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).LikeComment(ctx, args)
}

func (b *BoardClient) UnLikeComment(args *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).UnLikeComment(ctx, args)
}

func (b *BoardClient) ListLikes(args *pb.LikeListArgs) (*pb.LikeInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListLikes(ctx, args)
}

func (b *BoardClient) NewLike(args *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewLike(ctx, args)
}

func (b *BoardClient) DelLike(args *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelLike(ctx, args)
}

func (b *BoardClient) MutiGetSummary(args *pb.BoardSummaryArgs) (*pb.BoardSummaryInfos, error) {
	conn, err := GetConn(b.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).MutiGetSummary(ctx, args)
}

func (b *BoardClient) timeout() time.Duration {
	return 100 * time.Millisecond
}

func (b *BoardClient) service() string {
	return "board"
}
