package main

const (
	NORMAL_STATUS = 0
	DELETE_STATUS = 1
)

type CommRecord struct {
	Id    string `json:"id" bson:"_id,omitempty"`
	Oid   string `json:"oid" bson:"oid"`
	Cid   string `json:"cid" bson:"cid"`
	Level int    `json:"level" bson:"level"`

	Uid      string `json:"uid" bson:"uid"`
	Uname    string `json:"uname" bson:"uname"`
	AvatarId string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	IsRepost bool   `json:"is_repost" bson:"is_repost"`
	Content  string `json:"content" bson:"content"`
	ImgId    string `json:"img_id" bson:"img_id"`
	ImgEx    string `json:"img_ex" bson:"img_ex"`

	LikesCount int64 `json:"likes_count" bson:"likes_count"`
	ReplyCount int64 `json:"reply_count" bson:"reply_count"`
	CreatedAt  int64 `json:"created_at" bson:"created_at"`

	LatestReplys []CommRecord `json:"latest_replys" bson:"latest_replys"`
}

type CommRecordList []CommRecord

func (c CommRecordList) Len() int {
	return len(c)
}

func (c CommRecordList) Less(i, j int) bool {
	return c[i].CreatedAt < c[j].CreatedAt
}

func (c CommRecordList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CommLikeRecord struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Oid       string `json:"oid" bson:"oid"`
	Uid       string `json:"uid" bson:"uid"`
	CommentId string `json:"comment_id" bson:"comment_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

type LikeRecord struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Oid       string `json:"oid" bson:"oid"`
	Uid       string `json:"uid" bson:"uid"`
	Uname     string `json:"uname" bson:"uname"`
	AvatarId  string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx  string `json:"avatar_ex" bson:"avatar_ex"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

type SumRecord struct {
	Id          string `json:"id" bson:"_id,omitempty"`
	CommsTotal  int    `json:"comms_total" bson:"comms_total"`
	CommsTop    int    `json:"comms_top" bson:"comms_top"`
	LikesTotal  int    `json:"likes_total" bson:"likes_total"`
	RepostTotal int    `json:"repost_total" bson:"repost_total"`
}
