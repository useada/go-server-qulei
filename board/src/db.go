package main

type CommentRecord struct {
	Id    string `json:"id" bson:"_id,omitempty"`
	ObjId string `json:"obj_id" bson:"obj_id"`

	Uid      string `json:"uid" bson:"uid"`
	Uname    string `json:"uname" bson:"uname"`
	AvatarId string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx string `json:"avatar_ex" bson:"avatar_ex"`

	IsRepost     bool            `json:"is_repost" bson:"is_repost"`
	Level        int             `json:"level" bson:"level"`
	TopId        string          `json:"top_id" bson:"top_id"`
	LatestReplys []CommentRecord `json:"latest_replys" bson:"latest_replys"`

	Content string `json:"content" bson:"content"`
	ImgId   string `json:"img_id" bson:"img_id"`
	ImgEx   string `json:"img_ex" bson:"img_ex"`

	LikesCount int64 `json:"likes_count" bson:"likes_count"`
	ReplyCount int64 `json:"reply_count" bson:"reply_count"`
	CreatedAt  int64 `json:"created_at" bson:"created_at"`
}

type CommLikeRecord struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	ObjId     string `json:"obj_id" bson:"obj_id"`
	Uid       string `json:"uid" bson:"uid"`
	CommentId string `json:"comment_id" bson:"comment_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

type LikeRecord struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	ObjId     string `json:"obj_id" bson:"obj_id"`
	Uid       string `json:"uid" bson:"uid"`
	Uname     string `json:"uname" bson:"uname"`
	AvatarId  string `json:"avatar_id" bson:"avatar_id"`
	AvatarEx  string `json:"avatar_ex" bson:"avatar_ex"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

type BoardSumRecord struct {
	ObjId       string `json:"obj_id" bson:"_id,omitempty"`
	CommsTotal  int    `json:"comms_total" bson:"comms_total"`
	CommsTop    int    `json:"comms_top" bson:"comms_top"`
	LikesTotal  int    `json:"likes_total" bson:"likes_total"`
	RepostTotal int    `json:"repost_total" bson:"repost_total"`
}
