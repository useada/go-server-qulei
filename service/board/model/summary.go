package model

import (
	"a.com/go-server/proto/pb"
)

type Summary struct {
	ID              string `json:"id" bson:"_id,omitempty"` // Id = oid
	CommsCount      int    `json:"comms_count" bson:"comms_count"`
	CommsFirstCount int    `json:"comms_first_count" bson:"comms_first_count"`
	LikesCount      int    `json:"likes_count" bson:"likes_count"`
	RepostCount     int    `json:"repost_count" bson:"repost_count"`
	IsLiking        bool   `json:"is_liking" bson:"is_liking,omitempty"`
}

func (s *Summary) ConstructPb() *pb.SummaryInfo {
	return &pb.SummaryInfo{
		Id:              s.ID,
		CommsCount:      int32(s.CommsCount),
		CommsFirstCount: int32(s.CommsFirstCount),
		LikesCount:      int32(s.LikesCount),
		RepostCount:     int32(s.RepostCount),
		IsLiking:        s.IsLiking,
	}
}

type Summaries []Summary
