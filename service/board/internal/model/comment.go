package model

import (
	"github.com/rs/xid"

	"a.com/go-server/common/utime"
	"a.com/go-server/proto/pb"
)

type Comment struct {
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

	Replys []Comment `json:"replys" bson:"replys"` // 展示两条二级评论
}

func (c *Comment) ConstructPb() *pb.CommentInfo {
	replys := make([]*pb.CommentInfo, 0)
	for _, r := range c.Replys {
		replys = append(replys, &pb.CommentInfo{
			Id:         r.ID,
			IsRepost:   r.IsRepost,
			Oid:        r.Oid,
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
		IsRepost:   c.IsRepost,
		Oid:        c.Oid,
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

func (c *Comment) DestructPb(in *pb.CommNewArgs) *Comment {
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

type Comments []Comment

func (c Comments) Len() int {
	return len(c)
}

func (c Comments) Less(i, j int) bool {
	return c[i].CreatedAt > c[j].CreatedAt
}

func (c Comments) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
