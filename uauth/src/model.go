package main

import (
	"a.com/go-server/proto/st"
)

type AuthInfoModel struct {
	st.Base
	Phone  string `json:"phone"`
	Wechat string `json:"wechat"`
	Qicq   string `json:"qicq"`
	Passwd string `json:"passwd"`
	Salt   string `json:"salt"`
}

func (a *AuthInfoModel) TableName() string {
	return "auth_info"
}
