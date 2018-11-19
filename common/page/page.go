package page

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// PageToken 分页口令, 包含下一分页的参数
type PageToken struct {
	Offset int64 `json:"offset"` // created_at / updated_at / price
	Limit  int   `json:"limit"`  //
}

// 默认PageToken
func DefaultPageToken(limit int) (string, error) {
	bytes, err := json.Marshal(PageToken{
		Offset: 0,
		Limit:  limit,
	})
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
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
		return errors.New("empty page token")
	}

	bytes, err := base64.StdEncoding.DecodeString(tok)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, p)
}
