package handler

import (
	"context"

	"github.com/pkg/errors"

	"a.com/go-server/proto/pb"
	"a.com/go-server/service/board/internal/model"
)

func (s *SvrHandler) GetSummaries(ctx context.Context, in *pb.SummaryArgs) (*pb.SummaryInfos, error) {
	if len(in.Oids) == 0 {
		return nil, errors.New("oids empty")
	}

	citems, err := s.Cache.GetSummaries(ctx, in.Oids)
	if len(citems) == len(in.Oids) {
		return s.packSummaryInfos(ctx, citems, in.Uid)
	}
	oids := s.diffSumIds(ctx, citems, in.Oids)

	ditems, err := s.Store.GetSummaries(ctx, oids)
	if err != nil {
		return s.packSummaryInfos(ctx, citems, in.Uid)
	}
	s.cacheSummary(ctx, ditems, oids)

	return s.packSummaryInfos(ctx, append(citems, ditems...), in.Uid)
}

func (s *SvrHandler) packSummaryInfos(ctx context.Context, items model.Summaries, uid string) (*pb.SummaryInfos, error) {
	res := &pb.SummaryInfos{
		Items: make([]*pb.SummaryInfo, 0),
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

func (s *SvrHandler) diffSumIds(ctx context.Context, items model.Summaries, ids []string) []string {
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

func (s *SvrHandler) cacheSummary(ctx context.Context, items model.Summaries, ids []string) {
	xmap := make(map[string]model.Summary)
	for _, item := range items {
		xmap[item.ID] = item
	}

	for _, id := range ids {
		val, ok := xmap[id]
		if !ok {
			val = model.Summary{ID: id}
		}
		s.Cache.NewSummary(ctx, &val)
	}
	return
}

func (s *SvrHandler) userLikes(ctx context.Context, uid string) map[string]model.Like {
	xmap := make(map[string]model.Like)
	if len(uid) == 0 {
		return xmap
	}

	citems, _ := s.Cache.ListUserLikes(ctx, uid)
	for _, item := range citems {
		xmap[item.Oid] = item
	}
	if len(xmap) > 0 {
		return xmap
	}

	ditems, err := s.Store.ListUserLikes(ctx, uid)
	if err == nil && len(ditems) == 0 {
		ditems = append(ditems, model.Like{ID: "GUARD-ID", Oid: "GUARD-Oid"})
	}
	s.Cache.NewUserLikes(ctx, uid, ditems)

	for _, item := range ditems {
		xmap[item.Oid] = item
	}
	return xmap
}
