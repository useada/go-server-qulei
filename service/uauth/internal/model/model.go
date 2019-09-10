package model

import (
	"a.com/go-server/proto/model"
)

type Record struct {
	Uid        string `json:"uid"`
	DeviceType string `json:"device_type"`
	DeviceId   string `json:"device_id"`
	LoginTime  int64  `json:"login_time"`
}

type Account struct {
	model.Base
	Phone  string `json:"phone"`
	Wechat string `json:"wechat"`
	Qicq   string `json:"qicq"`
	Uname  string `json:"uname"`
	Passwd string `json:"passwd"`
	Salt   string `json:"salt"`
}

func (a *Account) TableName() string {
	return "accounts"
}
