package main

import (
	"a.com/go-server/common/utime"
	"a.com/go-server/common/xid"
	"a.com/go-server/proto/pb"
)

type CommentModel struct {
	Id       string `json:"id" bson:"_id,omitempty"`
	Oid      string `json:"oid" bson:"oid"`
	IsRepost bool   `json:"is_repost" bson:"is_repost"`
	Level    int    `json:"level" bson:"level"`
	Cid      string `json:"cid" bson:"cid"`

	Uname    string `json:"uname" bson:"uname"`
	Uid      string `json:"uid" bson:"uid"`
	AvatarId string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	Content string `json:"content" bson:"content"`
	ImgId   string `json:"img_id" bson:"img_id"`
	ImgEx   string `json:"img_ex" bson:"img_ex"`

	IsLiking   bool  `json:"is_liking" bson:"is_liking,omitempty"`
	LikesCount int64 `json:"likes_count" bson:"likes_count"`
	ReplyCount int64 `json:"reply_count" bson:"reply_count"`
	CreatedAt  int64 `json:"created_at" bson:"created_at"`

	Replys []CommentModel `json:"replys" bson:"replys"` // 展示两条二级评论
}

func (c *CommentModel) ConstructPb() *pb.CommentInfo {
	return &pb.CommentInfo{
		Id:         c.Id,
		Oid:        c.Oid,
		IsRepost:   c.IsRepost,
		Level:      int32(c.Level),
		Cid:        c.Cid,
		Content:    c.Content,
		ImgId:      c.ImgEx,
		ImgEx:      c.ImgId,
		LikesCount: int32(c.LikesCount),
		ReplyCount: int32(c.ReplyCount),
		Author: &pb.UserBaseInfo{
			Uname:    c.Uname,
			Uid:      c.Uid,
			AvatarId: c.AvatarId,
			AvatarEx: c.AvatarEx,
		},
		IsLiking:  c.IsLiking,
		CreatedAt: c.CreatedAt,
	}
}

func (c *CommentModel) DestructPb(in *pb.CommNewArgs) *CommentModel {
	c.Id = xid.New().String()
	c.Oid = in.Oid
	c.IsRepost = in.IsRepost
	c.Level = int(in.Level)
	c.Cid = in.Cid
	c.Uname = in.Author.Uname
	c.Uid = in.Author.Uid
	c.AvatarId = in.Author.AvatarId
	c.AvatarEx = in.Author.AvatarEx
	c.Content = in.Content
	c.ImgId = in.ImgId
	c.ImgEx = in.ImgEx
	c.CreatedAt = utime.Millisec()
	return c
}

type CommentModels []CommentModel

func (c CommentModels) Len() int {
	return len(c)
}

func (c CommentModels) Less(i, j int) bool {
	return c[i].CreatedAt < c[j].CreatedAt
}

func (c CommentModels) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CommentLikeModel struct {
	Id        string `json:"id" bson:"_id,omitempty"` // Id = uid + cid
	Oid       string `json:"oid" bson:"oid"`
	Uid       string `json:"uid" bson:"uid"`
	Cid       string `json:"cid" bson:"cid"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

func (c *CommentLikeModel) DestructPb(in *pb.CommLikeArgs) *CommentLikeModel {
	c.Id = in.Uid + in.Cid
	c.Oid = in.Oid
	c.Uid = in.Uid
	c.Cid = in.Cid
	c.CreatedAt = utime.Millisec()
	return c
}

type CommentLikeModels []CommentLikeModel

type LikeModel struct {
	Id  string `json:"id" bson:"_id,omitempty"` // Id = uid + oid
	Oid string `json:"oid" bson:"oid"`

	Uname    string `json:"uname" bson:"uname"`
	Uid      string `json:"uid" bson:"uid"`
	AvatarId string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
}

func (l *LikeModel) ConstructPb() *pb.LikeInfo {
	return &pb.LikeInfo{
		Id: l.Id,
		Author: &pb.UserBaseInfo{
			Uname:    l.Uname,
			Uid:      l.Uid,
			AvatarId: l.AvatarId,
			AvatarEx: l.AvatarEx,
		},
		Oid: l.Oid,
	}
}

func (l *LikeModel) DestructPb(in *pb.LikeNewArgs) *LikeModel {
	l.Id = in.Author.Uid + in.Oid
	l.Oid = in.Oid
	l.Uname = in.Author.Uname
	l.Uid = in.Author.Uid
	l.AvatarId = in.Author.AvatarId
	l.AvatarEx = in.Author.AvatarEx
	l.CreatedAt = utime.Millisec()
	return l
}

type LikeModels []LikeModel

type SummaryModel struct {
	Id              string `json:"id" bson:"_id,omitempty"` // Id = oid
	CommsCount      int    `json:"comms_count" bson:"comms_count"`
	CommsFirstCount int    `json:"comms_first_count" bson:"comms_first_count"`
	LikesCount      int    `json:"likes_count" bson:"likes_count"`
	RepostCount     int    `json:"repost_count" bson:"repost_count"`
	IsLiking        bool   `json:"is_liking" bson:"is_liking,omitempty"`
}

func (s *SummaryModel) ConstructPb() *pb.BoardSummaryInfo {
	return &pb.BoardSummaryInfo{
		Id:              s.Id,
		CommsCount:      int32(s.CommsCount),
		CommsFirstCount: int32(s.CommsFirstCount),
		LikesCount:      int32(s.LikesCount),
		RepostCount:     int32(s.RepostCount),
		IsLiking:        s.IsLiking,
	}
}

type SummaryModels []SummaryModel
