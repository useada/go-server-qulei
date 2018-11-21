package main

const (
	COMMENT_LEVEL_FIRST = 1
	COMMENT_LEVEL_CHILD = 2
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

	Replys []CommRecord `json:"replys" bson:"replys"` // 露出两条二级评论
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
	Id        string `json:"id" bson:"_id,omitempty"` // Id = uid + cid
	Oid       string `json:"oid" bson:"oid"`
	Uid       string `json:"uid" bson:"uid"`
	Cid       string `json:"cid" bson:"cid"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
type CommLikeRecordList []CommLikeRecord

type LikeRecord struct {
	Id        string `json:"id" bson:"_id,omitempty"` // Id = uid + oid
	Oid       string `json:"oid" bson:"oid"`
	Uid       string `json:"uid" bson:"uid"`
	Uname     string `json:"uname" bson:"uname"`
	AvatarId  string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx  string `json:"avatar_ex" bson:"avatar_ex"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
type LikeRecordList []LikeRecord

type SumRecord struct {
	Id              string `json:"id" bson:"_id,omitempty"` // Id = oid
	CommsCount      int    `json:"comms_count" bson:"comms_count"`
	CommsFirstCount int    `json:"comms_first_count" bson:"comms_first_count"`
	LikesCount      int    `json:"likes_count" bson:"likes_count"`
	RepostCount     int    `json:"repost_count" bson:"repost_count"`
}
type SumRecordList []SumRecord
