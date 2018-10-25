syntax = "proto3";                                                                                                                       

package pb;
 
service MsgBoard {
    // 评论
    rpc GetComment(GetCommRequest) returns (CommentResponse) {}   // 获取评论
    rpc NewComment(NewCommRequest) returns (CommentResponse) {}   // 创建评论
    rpc DelComment(DelCommRequest) returns (CommentResponse) {}   // 删除评论
    rpc ListComments(ListCommRequest) returns () {} 
    rpc LikeComment() returns () {}
    rpc UnLikeComment() returns () {}

    // 点赞
    rpc GetLike() returns () {}
    rpc NewLike() returns () {}
    rpc DelLike() returns () {} 
    rpc ListLikes() returns () {}

    // 每个对象的评论点赞汇总信息
    rpc GetStatistics() returns () {}
    rpc MGetStatistics() returns () {}
}

// 评论Request

message GetCommRequest {
    string object_id   = 1; // 对象id
    string id          = 2; // 评论id
}

message NewCommRequest {
    string object_id      = 1; 
    string father_id      = 2; // 父级评论id

    string self_uid       = 3;
    string self_uname     = 4;
    string self_avatar_id = 5;
    string self_avatar_ex = 6;

    string to_uid         = 7;
    string to_uname       = 8;

    string content        = 9;
    string photo_id       = 10;
    string photo_ex       = 11;
}

message DelCommRequest {
    string object_id   = 1;
    string father_id   = 2;
    string id          = 3;
}

message ListCommRequest {
    string object_id  = 1; 
    string page_token = 2;
    string direct     = 3;
}

// 评论Response

message CommentResponse {
    string object_id      = 1;
    string father_id      = 2;
    string id             = 3;

    string self_uid       = 4;
    string self_uname     = 5;
    string self_avatar_id = 6;
    string self_avatar_ex = 7;

    string to_uid         = 8;
    string to_uname       = 9;

    string content        = 10;
    string photo_id       = 11;
    string photo_ex       = 12;

    bool   liking         = 13;
    int32  likes_count    = 14;
    int32  reply_count    = 15
    int64  created_at     = 16;
}

