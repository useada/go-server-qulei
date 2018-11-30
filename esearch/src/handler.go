package main

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"a.com/go-server/common/page"
	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterEsearchServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) UsersByName(ctx context.Context,
	in *pb.UsersByNameArgs) (*pb.UserInfos, error) {
	res := pb.UserInfos{Items: make([](*pb.UserInfo), 0)}

	ptok := page.PageToken{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := ES.UsersByName(in.Name, int(ptok.Offset), ptok.Limit)
	if err != nil || len(rows) == 0 {
		return &res, err
	}

	// 生成新的page_token
	ptok.Offset += int64(len(rows))
	if res.PageToken, err = ptok.Encode(); err != nil {
		return &res, err
	}

	s.transUserInfo(&res, rows)
	return &res, nil
}

func (s *SvrHandler) UsersByNear(ctx context.Context,
	in *pb.UsersByNearArgs) (*pb.UserInfos, error) {
	res := pb.UserInfos{Items: make([](*pb.UserInfo), 0)}

	ptok := page.PageToken{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := ES.UsersByNear(in.Lat, in.Lon, int(in.Gender),
		int(ptok.Offset), ptok.Limit)
	if err != nil || len(rows) == 0 {
		return &res, err
	}

	// 生成新的page_token
	ptok.Offset += int64(len(rows))
	if res.PageToken, err = ptok.Encode(); err != nil {
		return &res, err
	}

	s.transUserInfo(&res, rows)
	return &res, nil
}

func (s *SvrHandler) transUserInfo(res *pb.UserInfos, rows []SearchModel) {
	for _, row := range rows {
		user := pb.UserInfo{}
		if err := json.Unmarshal(row.Source, &user); err != nil {
			continue
		}
		res.Items = append(res.Items, &user)
	}
}
