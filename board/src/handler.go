package main

import (
	"context"
	"errors"
	"sort"

	"google.golang.org/grpc"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterBoardServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) ListComments(ctx context.Context,
	in *pb.CommListArgs) (*pb.CommentInfos, error) {
	ptk := page.PageToken{}
	if err := ptk.Decode(in.PageToken); err != nil {
		Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, err
	}

	items, err := s.listCacheComms(in.Oid, in.Cid, in.Direct, ptk)
	if err == nil {
		sort.Sort(items)
		return s.packCommentInfos(items, ptk, in.Direct, in.Uid)
	}

	items, err = s.listDBComms(in.Oid, in.Cid, in.Direct, ptk)
	if err != nil {
		Log.Error("list db comments oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}

	if len(items) > ptk.Limit {
		items = items[0:ptk.Limit]
	}
	return s.packCommentInfos(items, ptk, in.Direct, in.Uid)
}

func (s *SvrHandler) GetComment(ctx context.Context,
	in *pb.CommGetArgs) (*pb.CommentInfo, error) {
	pitem, err := Cache.GetHashComm(in.Oid, in.Id)
	if err == nil {
		return s.packCommentInfo(pitem, in.Uid)
	}

	if pitem, err = DB.GetComment(in.Id); err == nil {
		Cache.SetHashComm(pitem)
	}
	return s.packCommentInfo(pitem, in.Uid)
}

func (s *SvrHandler) NewComment(ctx context.Context,
	in *pb.CommNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &CommentModel{}
	if err := DB.NewComment(pitem.DestructPb(in)); err == nil {
		Cache.PushComment(pitem)
	} else {
		Log.Error("new comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}

	if err := DB.IncrCommReply(int(in.Level), in.Cid, pitem); err == nil {
		Cache.DelHashComm(in.Oid, in.Cid)
	} else {
		Log.Error("incr reply oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
	}

	if err := DB.IncrSummaryComm(int(in.Level), in.Oid); err == nil {
		Cache.DelSummary(in.Oid)
	} else {
		Log.Error("incr summary oid:%s level:%d err:%v", in.Oid, in.Level, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.Id}, nil
}

func (s *SvrHandler) DelComment(ctx context.Context,
	in *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelComment(in.Id); err == nil {
		Cache.PopComment(in.Oid, in.Cid, in.Id)
	} else {
		Log.Error("del comm oid:%s id:%s err:%v", in.Oid, in.Id, err)
		return nil, err
	}

	if err := DB.DecrCommReply(int(in.Level), in.Cid, in.Id); err == nil {
		Cache.DelHashComm(in.Oid, in.Cid)
	} else {
		Log.Error("decr reply oid:%s id:%s err:%v", in.Oid, in.Id, err)
	}

	if err := DB.DecrSummaryComm(int(in.Level), in.Oid, 1); err == nil {
		Cache.DelSummary(in.Oid)
	} else {
		Log.Error("decr summary oid:%s level:%v err:%v", in.Oid, in.Level, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Id}, nil
}

func (s *SvrHandler) LikeComment(ctx context.Context,
	in *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &CommentLikeModel{}
	if err := DB.NewCommLike(pitem.DestructPb(in)); err == nil {
		Cache.DelUserCommLikes(in.Uid)
	} else {
		Log.Error("like comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}

	if err := DB.IncrCommLike(in.Cid); err == nil {
		Cache.DelHashComm(in.Oid, in.Cid)
	} else {
		Log.Error("incr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.Id}, nil
}

func (s *SvrHandler) UnLikeComment(ctx context.Context,
	in *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelCommLike(in.Uid + in.Cid); err == nil {
		Cache.DelUserCommLikes(in.Uid)
	} else {
		Log.Error("unlike comm oid:%s cid:%s err:%v", in.Oid, in.Cid, err)
		return nil, err
	}

	if err := DB.DecrCommLike(in.Cid); err == nil {
		Cache.DelHashComm(in.Oid, in.Cid)
	} else {
		Log.Error("decr like comm oid:%s cid:%v err:%v", in.Oid, in.Cid, err)
	}
	return &pb.ReplyBaseInfo{Id: in.Uid + in.Cid}, nil
}

func (s *SvrHandler) ListLikes(ctx context.Context,
	in *pb.LikeListArgs) (*pb.LikeInfos, error) {
	ptk := page.PageToken{}
	if err := ptk.Decode(in.PageToken); err != nil {
		Log.Error("decode pagetoken token:%s err:%v", in.PageToken, err)
		return nil, err
	}

	items, err := DB.ListLikes(in.Oid, ptk.Offset, ptk.Limit)
	if err != nil {
		Log.Error("list likes oid:%s offset:%d err:%v", in.Oid, ptk.Offset, err)
		return nil, err
	}
	return s.packLikeInfos(items, ptk)
}

func (s *SvrHandler) NewLike(ctx context.Context,
	in *pb.LikeNewArgs) (*pb.ReplyBaseInfo, error) {
	pitem := &LikeModel{}
	if err := DB.NewLike(pitem.DestructPb(in)); err == nil {
		Cache.DelUserLikes(in.Author.Uid)
	} else {
		Log.Error("new like oid:%s uid:%s err:%v", in.Oid, in.Author.Uid, err)
		return nil, err
	}

	if err := DB.IncrSummaryLike(in.Oid); err == nil {
		Cache.DelSummary(in.Oid)
	} else {
		Log.Error("incr likes oid:%s err:%v", in.Oid, err)
	}
	return &pb.ReplyBaseInfo{Id: pitem.Id}, nil
}

func (s *SvrHandler) DelLike(ctx context.Context,
	in *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	if err := DB.DelLike(in.Uid + in.Oid); err == nil {
		Cache.DelUserLikes(in.Uid)
	} else {
		Log.Error("del like oid:%s uid:%s err:%v", in.Oid, in.Uid, err)
		return nil, err
	}

	if err := DB.DecrSummaryLike(in.Oid); err == nil {
		Cache.DelSummary(in.Oid)
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

	citems, err := Cache.MutiGetSummary(in.Oids)
	if len(citems) == len(in.Oids) {
		return s.packSummaryInfos(citems, in.Uid)
	}
	oids := s.diffSumIds(citems, in.Oids)

	ditems, err := DB.MutiGetSummary(oids)
	if err != nil {
		return s.packSummaryInfos(citems, in.Uid)
	}
	s.cacheSummary(ditems, oids)

	return s.packSummaryInfos(append(citems, ditems...), in.Uid)
}

func (s *SvrHandler) listCacheComms(oid, cid, direct string,
	ptk page.PageToken) (CommentModels, error) {
	ids, err := Cache.ListZsetComms(oid, cid, direct, ptk.Offset, ptk.Limit)
	if err != nil {
		return nil, err
	}
	return s.mutiGetComms(oid, ids)
}

func (s *SvrHandler) listDBComms(oid, cid, direct string,
	ptk page.PageToken) (CommentModels, error) {
	limit := ptk.Limit
	if ptk.Offset == 0 {
		limit = COUNT_COMM_CACHE
	}
	items, err := DB.ListComments(oid, cid, direct, ptk.Offset, limit)
	if err != nil {
		return nil, err
	}

	if limit == COUNT_COMM_CACHE && len(items) > 0 {
		Cache.InitComms(oid, cid, items, len(items) < limit)
	}
	return items, nil
}

func (s *SvrHandler) packCommentInfos(items CommentModels,
	ptk page.PageToken, direct, uid string) (*pb.CommentInfos, error) {
	res := &pb.CommentInfos{
		Items: make([]*pb.CommentInfo, 0),
	}

	if ptk.Limit == len(items) {
		if direct == "lt" || direct == "lte" {
			ptk.Offset = items[0].CreatedAt
		} else {
			ptk.Offset = items[ptk.Limit-1].CreatedAt
		}
		pagetoken, err := ptk.Encode()
		if err != nil {
			return res, err
		}
		res.PageToken = pagetoken
	}

	xmap := s.userCommLikes(uid)
	for _, item := range items {
		if _, ok := xmap[item.Cid]; ok {
			item.IsLiking = true
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) packCommentInfo(pitem *CommentModel,
	uid string) (*pb.CommentInfo, error) {
	xmap := s.userCommLikes(uid)
	if _, ok := xmap[pitem.Cid]; ok {
		pitem.IsLiking = true
	}
	return pitem.ConstructPb(), nil
}

func (s *SvrHandler) packLikeInfos(items LikeModels,
	ptk page.PageToken) (*pb.LikeInfos, error) {
	res := &pb.LikeInfos{
		Items: make([]*pb.LikeInfo, 0),
	}

	if ptk.Limit == len(items) {
		ptk.Offset = items[ptk.Limit-1].CreatedAt
		pagetoken, err := ptk.Encode()
		if err != nil {
			return res, err
		}
		res.PageToken = pagetoken
	}

	for _, item := range items {
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) packSummaryInfos(items SummaryModels,
	uid string) (*pb.BoardSummaryInfos, error) {
	res := &pb.BoardSummaryInfos{
		Items: make([]*pb.BoardSummaryInfo, 0),
	}

	xmap := s.userLikes(uid)
	for _, item := range items {
		if _, ok := xmap[item.Id]; ok {
			item.IsLiking = true
		}
		res.Items = append(res.Items, item.ConstructPb())
	}
	return res, nil
}

func (s *SvrHandler) mutiGetComms(oid string,
	ids []string) (CommentModels, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	citems, err := Cache.MutiGetHashComms(oid, ids)
	if len(citems) == len(ids) {
		return citems, nil
	}
	cids := s.diffCommIds(citems, ids)

	ditems, err := DB.MutiGetComments(cids)
	if err != nil {
		return nil, err
	}
	for _, item := range ditems {
		Cache.SetHashComm(&item)
	}
	return append(citems, ditems...), nil
}

func (s *SvrHandler) diffCommIds(items CommentModels,
	ids []string) []string {
	xmap := make(map[string]bool)
	for _, item := range items {
		xmap[item.Id] = true
	}

	diffids := make([]string, 0)
	for _, id := range ids {
		if _, ok := xmap[id]; !ok {
			diffids = append(diffids, id)
		}
	}
	return diffids
}

func (s *SvrHandler) diffSumIds(items SummaryModels,
	ids []string) []string {
	xmap := make(map[string]bool)
	for _, item := range items {
		xmap[item.Id] = true
	}

	diffids := make([]string, 0)
	for _, id := range ids {
		if _, ok := xmap[id]; !ok {
			diffids = append(diffids, id)
		}
	}
	return diffids
}

func (s *SvrHandler) cacheSummary(items SummaryModels, ids []string) {
	xmap := make(map[string]SummaryModel)
	for _, item := range items {
		xmap[item.Id] = item
	}

	for _, id := range ids {
		val, ok := xmap[id]
		if !ok {
			val = SummaryModel{Id: id}
		}
		Cache.NewSummary(&val)
	}
	return
}

func (s *SvrHandler) userCommLikes(uid string) map[string]CommentLikeModel {
	xmap := make(map[string]CommentLikeModel)
	if len(uid) == 0 {
		return xmap
	}

	citems, _ := Cache.ListUserCommLikes(uid)
	for _, item := range citems {
		xmap[item.Cid] = item
	}
	if len(xmap) > 0 {
		return xmap
	}

	ditems, err := DB.ListUserCommLikes(uid)
	if err == nil && len(ditems) == 0 {
		ditems = append(ditems, CommentLikeModel{Id: "GUARD-ID", Cid: "GUARD-Cid"})
	}
	for _, item := range ditems {
		xmap[item.Cid] = item
	}
	if len(ditems) > 0 {
		Cache.NewUserCommLikes(uid, ditems)
	}
	return xmap
}

func (s *SvrHandler) userLikes(uid string) map[string]LikeModel {
	xmap := make(map[string]LikeModel)
	if len(uid) == 0 {
		return xmap
	}

	citems, _ := Cache.ListUserLikes(uid)
	for _, item := range citems {
		xmap[item.Oid] = item
	}
	if len(xmap) > 0 {
		return xmap
	}

	ditems, err := DB.ListUserLikes(uid)
	if err == nil && len(ditems) == 0 {
		ditems = append(ditems, LikeModel{Id: "GUARD-ID", Oid: "GUARD-Oid"})
	}
	for _, item := range ditems {
		xmap[item.Oid] = item
	}
	if len(ditems) > 0 {
		Cache.NewUserLikes(uid, ditems)
	}
	return xmap
}
