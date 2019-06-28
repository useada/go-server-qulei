package main

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterBoardServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) ListComments(ctx context.Context,
	in *pb.CommListArgs) (*pb.CommentInfos, error) {
	ptk := page.Token{}
	if err := ptk.Decode(in.PageToken); err != nil {
		Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, errors.Wrap(err, "decode page token")
	}

	items, err := s.listCacheComms(ctx, in.Oid, in.Cid, ptk)
	if err == nil {
		sort.Sort(items)
		return s.packCommentInfos(ctx, items, ptk, in.Uid)
	}

	items, err = s.listDBComms(ctx, in.Oid, in.Cid, ptk)
	if err != nil {
		Log.Error("list db comments oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}
	return s.packCommentInfos(ctx, items, ptk, in.Uid)
}

func (s *SvrHandler) GetComment(ctx context.Context,
	in *pb.CommGetArgs) (*pb.CommentInfo, error) {
	pitem, err := Cache.GetHashComm(ctx, in.Oid, in.Id)
	if err == nil {
		return s.packCommentInfo(ctx, pitem, in.Uid)
	}

	if pitem, err = DB.GetComment(ctx, in.Id); err == nil {
		Cache.SetHashComm(ctx, pitem)
	}
	return s.packCommentInfo(ctx, pitem, in.Uid)
}

func (s *SvrHandler) NewComment(ctx context.Context,
	in *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &CommentModel{}
	if err := DB.NewComment(ctx, pitem.DestructPb(in)); err == nil {
		Cache.PushComment(ctx, pitem)
	} else {
		Log.Error("new comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "new comment")
	}

	if len(in.Cid) != 0 { // 二级评论，更新一级评论数据
		if err := DB.IncrCommReply(ctx, in.Cid, pitem); err == nil {
			Cache.DelHashComm(ctx, in.Oid, in.Cid)
		} else {
			Log.Error("incr reply oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		}
	}

	if err := DB.IncrSummaryComm(ctx, in.Cid, in.Oid); err == nil {
		Cache.DelSummary(ctx, in.Oid)
	} else {
		Log.Error("incr summary oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) DelComment(ctx context.Context,
	in *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelComment(ctx, in.Id); err == nil {
		Cache.PopComment(ctx, in.Oid, in.Cid, in.Id)
	} else {
		Log.Error("del comm oid:%s id:%s err:%v", in.Oid, in.Id, err)
		return nil, errors.Wrap(err, "del comment")
	}

	if len(in.Cid) != 0 { // 二级评论，更新一级评论数据
		if err := DB.DecrCommReply(ctx, in.Cid, in.Id); err == nil {
			Cache.DelHashComm(ctx, in.Oid, in.Cid)
		} else {
			Log.Error("decr reply oid:%s id:%s err:%v", in.Oid, in.Id, err)
		}
	}

	if err := DB.DecrSummaryComm(ctx, in.Cid, in.Oid, 1); err == nil {
		Cache.DelSummary(ctx, in.Oid)
	} else {
		Log.Error("decr summary oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Id}, nil
}

func (s *SvrHandler) LikeComment(ctx context.Context,
	in *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &CommentLikeModel{}
	if err := DB.NewCommLike(ctx, pitem.DestructPb(in)); err == nil {
		Cache.DelUserCommLikes(ctx, in.Uid)
	} else {
		Log.Error("like comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "like comment")
	}

	if err := DB.IncrCommLike(ctx, in.Cid); err == nil {
		Cache.DelHashComm(ctx, in.Oid, in.Cid)
	} else {
		Log.Error("incr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) UnLikeComment(ctx context.Context,
	in *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelCommLike(ctx, in.Uid+in.Cid); err == nil {
		Cache.DelUserCommLikes(ctx, in.Uid)
	} else {
		Log.Error("unlike comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, errors.Wrap(err, "unlike comment")
	}

	if err := DB.DecrCommLike(ctx, in.Cid); err == nil {
		Cache.DelHashComm(ctx, in.Oid, in.Cid)
	} else {
		Log.Error("decr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Uid + in.Cid}, nil
}

func (s *SvrHandler) ListLikes(ctx context.Context,
	in *pb.LikeListArgs) (*pb.LikeInfos, error) {
	ptk := page.Token{}
	if err := ptk.Decode(in.PageToken); err != nil {
		Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, errors.Wrap(err, "decode page token")
	}

	items, err := DB.ListLikes(ctx, in.Oid, ptk.Offset, ptk.Limit+1)
	if err != nil {
		Log.Error("list likes oid:%s offset:%d err:%v", in.Oid, ptk.Offset, err)
		return nil, errors.Wrap(err, "list db likes")
	}
	return s.packLikeInfos(ctx, items, ptk)
}

func (s *SvrHandler) NewLike(ctx context.Context,
	in *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &LikeModel{}
	if err := DB.NewLike(ctx, pitem.DestructPb(in)); err == nil {
		Cache.DelUserLikes(ctx, in.Author.Uid)
	} else {
		Log.Error("new like oid:%s uid:%s err:%v", in.Oid, in.Author.Uid, err)
		return nil, errors.Wrap(err, "new like")
	}

	if err := DB.IncrSummaryLike(ctx, in.Oid); err == nil {
		Cache.DelSummary(ctx, in.Oid)
	} else {
		Log.Error("incr likes oid:%s err:%v", in.Oid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.ID}, nil
}

func (s *SvrHandler) DelLike(ctx context.Context,
	in *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelLike(ctx, in.Uid+in.Oid); err == nil {
		Cache.DelUserLikes(ctx, in.Uid)
	} else {
		Log.Error("del like oid:%s uid:%s err:%v", in.Oid, in.Uid, err)
		return nil, errors.Wrap(err, "del like")
	}

	if err := DB.DecrSummaryLike(ctx, in.Oid); err == nil {
		Cache.DelSummary(ctx, in.Oid)
	} else {
		Log.Error("decr likes oid:%s err:%v", in.Oid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Uid + in.Oid}, nil
}

func (s *SvrHandler) MutiGetSummary(ctx context.Context,
	in *pb.BoardSummaryArgs) (*pb.BoardSummaryInfos, error) {
	if len(in.Oids) == 0 {
		return nil, errors.New("oids empty")
	}

	citems, err := Cache.MutiGetSummary(ctx, in.Oids)
	if len(citems) == len(in.Oids) {
		return s.packSummaryInfos(ctx, citems, in.Uid)
	}
	oids := s.diffSumIds(ctx, citems, in.Oids)

	ditems, err := DB.MutiGetSummary(ctx, oids)
	if err != nil {
		return s.packSummaryInfos(ctx, citems, in.Uid)
	}
	s.cacheSummary(ctx, ditems, oids)

	return s.packSummaryInfos(ctx, append(citems, ditems...), in.Uid)
}

func (s *SvrHandler) listCacheComms(ctx context.Context, oid, cid string,
	ptk page.Token) (CommentModels, error) {
	ids, err := Cache.ListZsetComms(ctx, oid, cid, ptk.Offset, ptk.Limit+1)
	if err != nil {
		return nil, errors.Wrap(err, "list cache comments")
	}
	return s.mutiGetComms(ctx, oid, ids)
}

func (s *SvrHandler) listDBComms(ctx context.Context, oid, cid string,
	ptk page.Token) (CommentModels, error) {
	count := ptk.Limit
	if ptk.Offset == constant.TIME_INF_MAX {
		count = COUNT_COMM_CACHE
	}
	items, err := DB.ListComments(ctx, oid, cid, ptk.Offset, count+1)
	if err != nil {
		return nil, errors.Wrap(err, "list db comments")
	}

	if count == COUNT_COMM_CACHE && len(items) > 0 {
		Cache.InitComms(ctx, oid, cid, items, len(items) < count)
	}
	return items, nil
}

func (s *SvrHandler) packCommentInfos(ctx context.Context, items CommentModels,
	ptk page.Token, uid string) (*pb.CommentInfos, error) {
	res := &pb.CommentInfos{
		Items:     make([]*pb.CommentInfo, 0),
		PageToken: "",
	}

	if ptk.Limit+1 <= len(items) {
		ptk.Offset = items[ptk.Limit-1].CreatedAt
		pagetoken, err := ptk.Encode()
		if err != nil {
			return res, errors.Wrap(err, "encode page token")
		}
		res.PageToken = pagetoken
	}

	xmap := s.userCommLikes(ctx, uid)
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

func (s *SvrHandler) packCommentInfo(ctx context.Context, pitem *CommentModel,
	uid string) (*pb.CommentInfo, error) {
	xmap := s.userCommLikes(ctx, uid)
	if _, ok := xmap[pitem.Cid]; ok {
		pitem.IsLiking = true
	}
	return pitem.ConstructPb(), nil
}

func (s *SvrHandler) packLikeInfos(ctx context.Context, items LikeModels,
	ptk page.Token) (*pb.LikeInfos, error) {
	res := &pb.LikeInfos{
		Items:     make([]*pb.LikeInfo, 0),
		PageToken: "",
	}

	if ptk.Limit+1 <= len(items) {
		ptk.Offset = items[ptk.Limit-1].CreatedAt
		pagetoken, err := ptk.Encode()
		if err != nil {
			return res, errors.Wrap(err, "encode page token")
		}
		res.PageToken = pagetoken
	}

	for _, item := range items {
		if len(res.Items) == ptk.Limit {
			break
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) packSummaryInfos(ctx context.Context, items SummaryModels,
	uid string) (*pb.BoardSummaryInfos, error) {
	res := &pb.BoardSummaryInfos{
		Items: make([]*pb.BoardSummaryInfo, 0),
	}

	xmap := s.userLikes(ctx, uid)
	for _, item := range items {
		if _, ok := xmap[item.ID]; ok {
			item.IsLiking = true
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) mutiGetComms(ctx context.Context, oid string,
	ids []string) (CommentModels, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	citems, err := Cache.MutiGetHashComms(ctx, oid, ids)
	if len(citems) == len(ids) {
		return citems, nil
	}
	cids := s.diffCommIds(ctx, citems, ids)

	ditems, err := DB.MutiGetComments(ctx, cids)
	if err != nil {
		return nil, errors.Wrap(err, "multi get comments")
	}
	for _, item := range ditems {
		Cache.SetHashComm(ctx, &item)
	}
	return append(citems, ditems...), nil
}

func (s *SvrHandler) diffCommIds(ctx context.Context, items CommentModels,
	ids []string) []string {
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

func (s *SvrHandler) diffSumIds(ctx context.Context, items SummaryModels,
	ids []string) []string {
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

func (s *SvrHandler) cacheSummary(ctx context.Context, items SummaryModels, ids []string) {
	xmap := make(map[string]SummaryModel)
	for _, item := range items {
		xmap[item.ID] = item
	}

	for _, id := range ids {
		val, ok := xmap[id]
		if !ok {
			val = SummaryModel{ID: id}
		}
		Cache.NewSummary(ctx, &val)
	}
	return
}

func (s *SvrHandler) userCommLikes(ctx context.Context, uid string) map[string]CommentLikeModel {
	xmap := make(map[string]CommentLikeModel)
	if len(uid) == 0 {
		return xmap
	}

	citems, _ := Cache.ListUserCommLikes(ctx, uid)
	for _, item := range citems {
		xmap[item.Cid] = item
	}
	if len(xmap) > 0 {
		return xmap
	}

	ditems, err := DB.ListUserCommLikes(ctx, uid)
	if err == nil && len(ditems) == 0 {
		ditems = append(ditems, CommentLikeModel{ID: "GUARD-ID", Cid: "GUARD-Cid"})
	}
	Cache.NewUserCommLikes(ctx, uid, ditems)

	for _, item := range ditems {
		xmap[item.Cid] = item
	}
	return xmap
}

func (s *SvrHandler) userLikes(ctx context.Context, uid string) map[string]LikeModel {
	xmap := make(map[string]LikeModel)
	if len(uid) == 0 {
		return xmap
	}

	citems, _ := Cache.ListUserLikes(ctx, uid)
	for _, item := range citems {
		xmap[item.Oid] = item
	}
	if len(xmap) > 0 {
		return xmap
	}

	ditems, err := DB.ListUserLikes(ctx, uid)
	if err == nil && len(ditems) == 0 {
		ditems = append(ditems, LikeModel{ID: "GUARD-ID", Oid: "GUARD-Oid"})
	}
	Cache.NewUserLikes(ctx, uid, ditems)

	for _, item := range ditems {
		xmap[item.Oid] = item
	}
	return xmap
}
