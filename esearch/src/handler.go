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

func (s *SvrHandler) SearchByName(ctx context.Context,
	in *pb.SearchNameArgs) (*pb.SearchUserInfos, error) {
	res := pb.SearchUserInfos{Items: make([](*pb.SearchUserInfo), 0)}

	ptok := page.PageToken{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := ES.SearchByName(in.Name, int(ptok.Offset), ptok.Limit)
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

func (s *SvrHandler) SearchByNear(ctx context.Context,
	in *pb.SearchNearArgs) (*pb.SearchUserInfos, error) {
	res := pb.SearchUserInfos{Items: make([](*pb.SearchUserInfo), 0)}

	ptok := page.PageToken{}
	if err := ptok.Decode(in.PageToken); err != nil {
		return &res, err
	}

	rows, err := ES.SearchByNear(in.Lat, in.Lon, int(in.Gender),
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

func (s *SvrHandler) transUserInfo(res *pb.SearchUserInfos, rows []SearchItem) {
	for _, row := range rows {
		user := pb.SearchUserInfo{}
		if err := json.Unmarshal(row.Source, &user); err != nil {
			continue
		}
		res.Items = append(res.Items, &user)
	}
}
