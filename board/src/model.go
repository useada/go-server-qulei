package main

import (
	"github.com/rs/xid"

	"a.com/go-server/common/utime"
	"a.com/go-server/proto/pb"
)

type CommentModel struct {
	ID  string `json:"id" bson:"_id,omitempty"`
	Oid string `json:"oid" bson:"oid"`
	Cid string `json:"cid" bson:"cid"`

	Uname    string `json:"uname" bson:"uname"`
	UID      string `json:"uid" bson:"uid"`
	AvatarID string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	IsRepost bool   `json:"is_repost" bson:"is_repost"`
	Content  string `json:"content" bson:"content"`
	ImgID    string `json:"img_id" bson:"img_id"`
	ImgEx    string `json:"img_ex" bson:"img_ex"`

	IsLiking   bool  `json:"is_liking" bson:"is_liking,omitempty"`
	LikesCount int64 `json:"likes_count" bson:"likes_count"`
	ReplyCount int64 `json:"reply_count" bson:"reply_count"`
	CreatedAt  int64 `json:"created_at" bson:"created_at"`

	Replys []CommentModel `json:"replys" bson:"replys"` // 展示两条二级评论
}

func (c *CommentModel) ConstructPb() *pb.CommentInfo {
	replys := make([]*pb.CommentInfo, 0)
	for _, r := range c.Replys {
		replys = append(replys, &pb.CommentInfo{
			Id:         r.ID,
			Oid:        r.Oid,
			IsRepost:   r.IsRepost,
			Cid:        r.Cid,
			Content:    r.Content,
			ImgId:      r.ImgEx,
			ImgEx:      r.ImgID,
			LikesCount: int32(r.LikesCount),
			ReplyCount: int32(r.ReplyCount),
			Author: &pb.UserBaseInfo{
				Uname:    r.Uname,
				Uid:      r.UID,
				AvatarId: r.AvatarID,
				AvatarEx: r.AvatarEx,
			},
			IsLiking:  r.IsLiking,
			CreatedAt: r.CreatedAt,
		})
	}
	return &pb.CommentInfo{
		Id:         c.ID,
		Oid:        c.Oid,
		IsRepost:   c.IsRepost,
		Cid:        c.Cid,
		Content:    c.Content,
		ImgId:      c.ImgEx,
		ImgEx:      c.ImgID,
		LikesCount: int32(c.LikesCount),
		ReplyCount: int32(c.ReplyCount),
		Author: &pb.UserBaseInfo{
			Uname:    c.Uname,
			Uid:      c.UID,
			AvatarId: c.AvatarID,
			AvatarEx: c.AvatarEx,
		},
		IsLiking:  c.IsLiking,
		Replys:    replys,
		CreatedAt: c.CreatedAt,
	}
}

func (c *CommentModel) DestructPb(in *pb.CommNewArgs) *CommentModel {
	c.ID = xid.New().String()
	c.Oid = in.Oid
	c.IsRepost = in.IsRepost
	c.Cid = in.Cid
	c.Uname = in.Author.Uname
	c.UID = in.Author.Uid
	c.AvatarID = in.Author.AvatarId
	c.AvatarEx = in.Author.AvatarEx
	c.Content = in.Content
	c.ImgID = in.ImgId
	c.ImgEx = in.ImgEx
	c.CreatedAt = utime.Millisec()
	return c
}

type CommentModels []CommentModel

func (c CommentModels) Len() int {
	return len(c)
}

func (c CommentModels) Less(i, j int) bool {
	return c[i].CreatedAt > c[j].CreatedAt
}

func (c CommentModels) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CommentLikeModel struct {
	ID        string `json:"id" bson:"_id,omitempty"` // Id = uid + cid
	Oid       string `json:"oid" bson:"oid"`
	UID       string `json:"uid" bson:"uid"`
	Cid       string `json:"cid" bson:"cid"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

func (c *CommentLikeModel) DestructPb(in *pb.CommLikeArgs) *CommentLikeModel {
	c.ID = in.Uid + in.Cid
	c.Oid = in.Oid
	c.UID = in.Uid
	c.Cid = in.Cid
	c.CreatedAt = utime.Millisec()
	return c
}

type CommentLikeModels []CommentLikeModel

type LikeModel struct {
	ID  string `json:"id" bson:"_id,omitempty"` // Id = uid + oid
	Oid string `json:"oid" bson:"oid"`

	Uname    string `json:"uname" bson:"uname"`
	UID      string `json:"uid" bson:"uid"`
	AvatarID string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
}

func (l *LikeModel) ConstructPb() *pb.LikeInfo {
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

func (l *LikeModel) DestructPb(in *pb.LikeNewArgs) *LikeModel {
	l.ID = in.Author.Uid + in.Oid
	l.Oid = in.Oid
	l.Uname = in.Author.Uname
	l.UID = in.Author.Uid
	l.AvatarID = in.Author.AvatarId
	l.AvatarEx = in.Author.AvatarEx
	l.CreatedAt = utime.Millisec()
	return l
}

type LikeModels []LikeModel

type SummaryModel struct {
	ID              string `json:"id" bson:"_id,omitempty"` // Id = oid
	CommsCount      int    `json:"comms_count" bson:"comms_count"`
	CommsFirstCount int    `json:"comms_first_count" bson:"comms_first_count"`
	LikesCount      int    `json:"likes_count" bson:"likes_count"`
	RepostCount     int    `json:"repost_count" bson:"repost_count"`
	IsLiking        bool   `json:"is_liking" bson:"is_liking,omitempty"`
}

func (s *SummaryModel) ConstructPb() *pb.BoardSummaryInfo {
	return &pb.BoardSummaryInfo{
		Id:              s.ID,
		CommsCount:      int32(s.CommsCount),
		CommsFirstCount: int32(s.CommsFirstCount),
		LikesCount:      int32(s.LikesCount),
		RepostCount:     int32(s.RepostCount),
		IsLiking:        s.IsLiking,
	}
}

type SummaryModels []SummaryModel
