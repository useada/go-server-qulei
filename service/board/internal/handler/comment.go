package handler

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"

	"a.com/go-server/service/board/internal/model"
)

func (s *SvrHandler) ListComments(ctx context.Context, in *pb.CommListArgs) (*pb.CommentInfos, error) {
	if len(in.PageToken) == 0 {
		in.PageToken = page.Default(constant.TIME_INF_MAX, PAGE_COUNT)
	}

	ptk := page.Token{}
	if err := ptk.Decode(in.PageToken); err != nil {
		s.Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, errors.Wrap(err, "decode page token")
	}

	items, err := s.listCacheComms(ctx, in.Oid, in.Cid, ptk)
	if err == nil {
		sort.Sort(items)
		return s.packCommentInfos(ctx, items, ptk, in.Uid, in.Oid)
	}

	items, err = s.listStoreComms(ctx, in.Oid, in.Cid, ptk)
	if err != nil {
		s.Log.Error("list db comments oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}
	return s.packCommentInfos(ctx, items, ptk, in.Uid, in.Oid)
}

func (s *SvrHandler) GetComment(ctx context.Context, in *pb.CommGetArgs) (*pb.CommentInfo, error) {
	pitem, err := s.Cache.GetComment(ctx, in.Oid, in.Id)
	if err == nil {
		return s.packCommentInfo(ctx, pitem, in.Uid, in.Oid)
	}

	if pitem, err = s.Store.GetComment(ctx, in.Id); err == nil {
		s.Cache.SetComment(ctx, pitem)
	}
	return s.packCommentInfo(ctx, pitem, in.Uid, in.Oid)
}

func (s *SvrHandler) NewComment(ctx context.Context, in *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &model.Comment{}
	if err := s.Store.NewComment(ctx, pitem.DestructPb(in)); err == nil {
		s.Cache.PushComment(ctx, pitem)
	} else {
		s.Log.Error("new comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "new comment")
	}

	if len(in.Cid) != 0 { // 二级评论，更新一级评论数据
		if err := s.Store.IncrCommReply(ctx, in.Cid, pitem); err == nil {
			s.Cache.DelComment(ctx, in.Oid, in.Cid)
		} else {
			s.Log.Error("incr reply oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		}
	}

	if err := s.Store.IncrSummaryComm(ctx, in.Cid, in.Oid); err == nil {
		s.Cache.DelSummary(ctx, in.Oid)
	} else {
		s.Log.Error("incr summary oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) DelComment(ctx context.Context, in *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := s.Store.DelComment(ctx, in.Id); err == nil {
		s.Cache.PopComment(ctx, in.Oid, in.Cid, in.Id)
	} else {
		s.Log.Error("del comm oid:%s id:%s err:%v", in.Oid, in.Id, err)
		return nil, errors.Wrap(err, "del comment")
	}

	if len(in.Cid) != 0 { // 二级评论，更新一级评论数据
		if err := s.Store.DecrCommReply(ctx, in.Cid, in.Id); err == nil {
			s.Cache.DelComment(ctx, in.Oid, in.Cid)
		} else {
			s.Log.Error("decr reply oid:%s id:%s err:%v", in.Oid, in.Id, err)
		}
	}

	if err := s.Store.DecrSummaryComm(ctx, in.Cid, in.Oid, 1); err == nil {
		s.Cache.DelSummary(ctx, in.Oid)
	} else {
		s.Log.Error("decr summary oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Id}, nil
}

func (s *SvrHandler) listCacheComms(ctx context.Context, oid, cid string, ptk page.Token) (model.Comments, error) {
	ids, err := s.Cache.RangeComments(ctx, oid, cid, ptk.Offset, ptk.Limit+1)
	if err != nil {
		return nil, errors.Wrap(err, "list cache comments")
	}
	return s.getComments(ctx, oid, ids)
}

func (s *SvrHandler) listStoreComms(ctx context.Context, oid, cid string, ptk page.Token) (model.Comments, error) {
	count := ptk.Limit
	if ptk.Offset == constant.TIME_INF_MAX {
		count = COUNT_COMM_CACHE
	}
	items, err := s.Store.ListComments(ctx, oid, cid, ptk.Offset, count+1)
	if err != nil {
		return nil, errors.Wrap(err, "list db comments")
	}

	if count == COUNT_COMM_CACHE && len(items) > 0 {
		s.Cache.InitComments(ctx, oid, cid, items, len(items) < count)
	}
	return items, nil
}

func (s *SvrHandler) packCommentInfos(ctx context.Context, items model.Comments, ptk page.Token, uid, oid string) (*pb.CommentInfos, error) {
	res := &pb.CommentInfos{
		Items: make([]*pb.CommentInfo, 0),
	}

	if ptk.Limit+1 <= len(items) {
		ptk.Offset = items[ptk.Limit-1].CreatedAt
		res.PageToken = ptk.Encode()
	}

	xmap := s.listUserCommLikes(ctx, uid, oid)
	for _, item := range items {
		if len(res.Items) == ptk.Limit {
			break
		}
		if _, ok := xmap[item.Cid]; ok {
			item.IsLiking = true
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) packCommentInfo(ctx context.Context, pitem *model.Comment, uid, oid string) (*pb.CommentInfo, error) {
	xmap := s.listUserCommLikes(ctx, uid, oid)
	if _, ok := xmap[pitem.Cid]; ok {
		pitem.IsLiking = true
	}
	return pitem.ConstructPb(), nil
}

func (s *SvrHandler) getComments(ctx context.Context, oid string, ids []string) (model.Comments, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	citems, err := s.Cache.GetComments(ctx, oid, ids)
	if len(citems) == len(ids) {
		return citems, nil
	}
	cids := s.diffCommIds(ctx, citems, ids)

	ditems, err := s.Store.GetComments(ctx, cids)
	if err != nil {
		return nil, errors.Wrap(err, "multi get comments")
	}
	for _, item := range ditems {
		s.Cache.SetComment(ctx, &item)
	}
	return append(citems, ditems...), nil
}

func (s *SvrHandler) diffCommIds(ctx context.Context, items model.Comments, ids []string) []string {
	xmap := make(map[string]bool)
	for _, item := range items {
		xmap[item.ID] = true
	}

	diffids := make([]string, 0)
	for _, id := range ids {
		if _, ok := xmap[id]; !ok {
			diffids = append(diffids, id)
		}
	}
	return diffids
}
