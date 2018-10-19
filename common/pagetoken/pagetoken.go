package pagetoken

import (
	"encoding/base64"
	"encoding/json"
)

// PageToken 分页口令, 包含下一分页的参数
type PageToken struct {
	Offset int64 `json:"offset"` // created_at / updated_at / price
	Limit  int   `json:"limit"`  //
}

// Encode 返回一个 base64 字符串
func (p *PageToken) Encode() (string, error) {
	bytes, err := json.Marshal(*p)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Decode ...
func (p *PageToken) Decode(tok string) error {
	if len(tok) == 0 {
		return nil
	}

	bytes, err := base64.StdEncoding.DecodeString(tok)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, p)
}
