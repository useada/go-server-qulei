package model

import (
	"a.com/go-server/common/utime"
	"a.com/go-server/proto/pb"
)

type Like struct {
	ID  string `json:"id" bson:"_id,omitempty"` // Id = uid + oid
	Oid string `json:"oid" bson:"oid"`

	Uname    string `json:"uname" bson:"uname"`
	UID      string `json:"uid" bson:"uid"`
	AvatarID string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
}

func (l *Like) ConstructPb() *pb.LikeInfo {
	return &pb.LikeInfo{
		Id: l.ID,
		Author: &pb.UserBaseInfo{
			Uname:    l.Uname,
			Uid:      l.UID,
			AvatarId: l.AvatarID,
			AvatarEx: l.AvatarEx,
		},
		Oid: l.Oid,
	}
}

func (l *Like) DestructPb(in *pb.LikeNewArgs) *Like {
	l.ID = in.Author.Uid + in.Oid
	l.Oid = in.Oid
	l.Uname = in.Author.Uname
	l.UID = in.Author.Uid
	l.AvatarID = in.Author.AvatarId
	l.AvatarEx = in.Author.AvatarEx
	l.CreatedAt = utime.Millisec()
	return l
}

type Likes []Like
