package service

import (
	"context"

	"github.com/pkg/errors"

	"a.com/go-server/proto/pb"

	"a.com/go-server/service/board/model"
)

func (s *SvrHandler) LikeComment(ctx context.Context, in *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &model.CommentLike{}
	if err := s.Store.NewCommLike(ctx, pitem.DestructPb(in)); err == nil {
		s.Cache.DelUserCommLikes(ctx, in.Uid)
	} else {
		s.Log.Error("like comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "like comment")
	}

	if err := s.Store.IncrCommLike(ctx, in.Cid); err == nil {
		s.Cache.DelComment(ctx, in.Oid, in.Cid)
	} else {
		s.Log.Error("incr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) UnLikeComment(ctx context.Context, in *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	if err := s.Store.DelCommLike(ctx, in.Uid+in.Cid); err == nil {
		s.Cache.DelUserCommLikes(ctx, in.Uid)
	} else {
		s.Log.Error("unlike comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "unlike comment")
	}

	if err := s.Store.DecrCommLike(ctx, in.Cid); err == nil {
		s.Cache.DelComment(ctx, in.Oid, in.Cid)
	} else {
		s.Log.Error("decr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Uid + in.Cid}, nil
}
