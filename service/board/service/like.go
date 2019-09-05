package service

import (
	"context"

	"github.com/pkg/errors"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"

	"a.com/go-server/service/board/model"
)

func (s *SvrHandler) ListLikes(ctx context.Context, in *pb.LikeListArgs) (*pb.LikeInfos, error) {
	if len(in.PageToken) == 0 {
		in.PageToken = page.Default(constant.TIME_INF_MAX, PAGE_COUNT)
	}

	ptk := page.Token{}
	if err := ptk.Decode(in.PageToken); err != nil {
		s.Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, errors.Wrap(err, "decode page token")
	}

	items, err := s.Store.ListLikes(ctx, in.Oid, ptk.Offset, ptk.Limit+1)
	if err != nil {
		s.Log.Error("list likes oid:%s offset:%d err:%v", in.Oid, ptk.Offset, err)
		return nil, errors.Wrap(err, "list db likes")
	}
	return s.packLikeInfos(ctx, items, ptk)
}

func (s *SvrHandler) NewLike(ctx context.Context, in *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &model.Like{}
	if err := s.Store.NewLike(ctx, pitem.DestructPb(in)); err == nil {
		s.Cache.DelUserLikes(ctx, in.Author.Uid)
	} else {
		s.Log.Error("new like oid:%s uid:%s err:%v", in.Oid, in.Author.Uid, err)
		return nil, errors.Wrap(err, "new like")
	}

	if err := s.Store.IncrSummaryLike(ctx, in.Oid); err == nil {
		s.Cache.DelSummary(ctx, in.Oid)
	} else {
		s.Log.Error("incr likes oid:%s err:%v", in.Oid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) DelLike(ctx context.Context, in *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := s.Store.DelLike(ctx, in.Uid+in.Oid); err == nil {
		s.Cache.DelUserLikes(ctx, in.Uid)
	} else {
		s.Log.Error("del like oid:%s uid:%s err:%v", in.Oid, in.Uid, err)
		return nil, errors.Wrap(err, "del like")
	}

	if err := s.Store.DecrSummaryLike(ctx, in.Oid); err == nil {
		s.Cache.DelSummary(ctx, in.Oid)
	} else {
		s.Log.Error("decr likes oid:%s err:%v", in.Oid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Uid + in.Oid}, nil
}

func (s *SvrHandler) packLikeInfos(ctx context.Context, items model.Likes, ptk page.Token) (*pb.LikeInfos, error) {
	res := &pb.LikeInfos{
		Items:     make([]*pb.LikeInfo, 0),
		PageToken: "",
	}

	if ptk.Limit+1 <= len(items) {
		ptk.Offset = items[ptk.Limit-1].CreatedAt
		res.PageToken = ptk.Encode()
	}

	for _, item := range items {
		if len(res.Items) == ptk.Limit {
			break
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}
