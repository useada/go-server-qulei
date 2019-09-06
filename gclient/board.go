package gclient

import (
	"context"

	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"
)

func (cli *Client) ListComments(c context.Context, args *pb.CommListArgs) (*pb.CommentInfos, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListComments(ctx, args)
}

func (cli *Client) GetComment(c context.Context, args *pb.CommGetArgs) (*pb.CommentInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).GetComment(ctx, args)
}

func (cli *Client) NewComment(c context.Context, args *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewComment(ctx, args)
}

func (cli *Client) DelComment(c context.Context, args *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelComment(ctx, args)
}

func (cli *Client) LikeComment(c context.Context, args *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).LikeComment(ctx, args)
}

func (cli *Client) UnLikeComment(c context.Context, args *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).UnLikeComment(ctx, args)
}

func (cli *Client) ListLikes(c context.Context, args *pb.LikeListArgs) (*pb.LikeInfos, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).ListLikes(ctx, args)
}

func (cli *Client) NewLike(c context.Context, args *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).NewLike(ctx, args)
}

func (cli *Client) DelLike(c context.Context, args *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).DelLike(ctx, args)
}

func (cli *Client) GetSummaries(c context.Context, args *pb.SummaryArgs) (*pb.SummaryInfos, error) {
	conn, err := cli.Get(constant.BoardService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewBoardClient(conn).GetSummaries(ctx, args)
}
