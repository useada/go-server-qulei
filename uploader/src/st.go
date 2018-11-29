package main

type FileInfoDB struct {
	Id     string `json:"id"`    // 文件id <primary key>
	Ex     string `json:"ex"`    // 文件扩展名
	Typef  int    `json:"typef"` // 文件类型 头像/图片/普通文件
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int64  `json:"size"`
}

func (fi *FileInfoDB) TableName() string {
	return "file_info"
}
