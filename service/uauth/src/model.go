package main

import (
	"a.com/go-server/proto/model"
)

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
