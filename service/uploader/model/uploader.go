package model

import "a.com/go-server/proto/pb"

type FileInfo struct {
	ID     string `json:"id"`    // 文件id <primary key>
	Ex     string `json:"ex"`    // 文件扩展名
	Type   int    `json:"xtype"` // 文件类型 头像/图片/普通文件
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int64  `json:"size"`
}

func (f *FileInfo) ConstructPb() *pb.FileInfo {
	return &pb.FileInfo{
		Id:     f.ID,
		Ex:     f.Ex,
		Type:   pb.TYPE(f.Type),
		Width:  int32(f.Width),
		Height: int32(f.Height),
		Size_:  f.Size,
	}
}

func (f *FileInfo) DestructPb(in *pb.FileUploadArgs) *FileInfo {
	f.ID = in.Id
	f.Ex = in.Ex
	f.Type = int(in.Type)
	f.Height = int(in.Height)
	f.Width = int(in.Width)
	f.Size = int64(len(in.Data))
	return f
}

func (f *FileInfo) TableName() string {
	return "file_info"
}
