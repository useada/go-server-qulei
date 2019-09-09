package handler

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/pb"
	"a.com/go-server/service/esearch/internal/model"
	"a.com/go-server/service/esearch/internal/store"
)

func RegisterHandler(svr *grpc.Server, store store.Store, log *zap.SugaredLogger) {
	pb.RegisterEsearchServer(svr, &SvrHandler{
		Store: store,
		Log:   log,
	})
}

type SvrHandler struct {
	Store store.Store
	Log   *zap.SugaredLogger
}

func (s *SvrHandler) UsersByName(ctx context.Context, in *pb.UsersByNameArgs) (*pb.UserInfos, error) {
	res := pb.UserInfos{Items: make([](*pb.UserInfo), 0)}

	ptok := page.Token{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := s.Store.UsersByName(in.Name, int(ptok.Offset), ptok.Limit)
	if err != nil || len(rows) == 0 {
		return &res, err
	}

	// 生成新的page_token
	ptok.Offset += int64(len(rows))
	res.PageToken = ptok.Encode()

	s.transUserInfo(&res, rows)
	return &res, nil
}

func (s *SvrHandler) UsersByNear(ctx context.Context, in *pb.UsersByNearArgs) (*pb.UserInfos, error) {
	res := pb.UserInfos{Items: make([](*pb.UserInfo), 0)}

	ptok := page.Token{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := s.Store.UsersByNear(in.Lat, in.Lon, int(in.Gender),
		int(ptok.Offset), ptok.Limit)
	if err != nil || len(rows) == 0 {
		return &res, err
	}

	// 生成新的page_token
	ptok.Offset += int64(len(rows))
	res.PageToken = ptok.Encode()

	s.transUserInfo(&res, rows)
	return &res, nil
}

func (s *SvrHandler) transUserInfo(res *pb.UserInfos, rows []model.SearchInfo) {
	for _, row := range rows {
		user := pb.UserInfo{}
		if err := json.Unmarshal(row.Source, &user); err != nil {
			continue
		}
		res.Items = append(res.Items, &user)
	}
}
