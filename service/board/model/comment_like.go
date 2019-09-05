package model

import (
	"a.com/go-server/common/utime"
	"a.com/go-server/proto/pb"
)

type CommentLike struct {
	ID        string `json:"id" bson:"_id,omitempty"` // Id = uid + cid
	Oid       string `json:"oid" bson:"oid"`
	UID       string `json:"uid" bson:"uid"`
	Cid       string `json:"cid" bson:"cid"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

func (c *CommentLike) DestructPb(in *pb.CommLikeArgs) *CommentLike {
	c.ID = in.Uid + in.Cid
	c.Oid = in.Oid
	c.UID = in.Uid
	c.Cid = in.Cid
	c.CreatedAt = utime.Millisec()
	return c
}

type CommentLikes []CommentLike
