package page

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// Token 分页口令, 包含下一分页的参数
type Token struct {
	Offset int64  `json:"offset"` // created_at / updated_at / price
	Limit  int    `json:"limit"`  // 20 / 50 /100
	Order  string `json:"order"`  // "id desc" / "updated_at asc"
}

// Default Token
func Default(offset int64, limit int) string {
	bytes, _ := json.Marshal(Token{
		Offset: offset,
		Limit:  limit,
		Order:  "",
	})
	return base64.StdEncoding.EncodeToString(bytes)
}

// Encode 返回一个 base64 字符串
func (t *Token) Encode() string {
	bytes, _ := json.Marshal(*t)
	return base64.StdEncoding.EncodeToString(bytes)
}

// Decode ...
func (t *Token) Decode(tok string) error {
	if len(tok) == 0 {
		return errors.New("empty page token")
	}

	bytes, err := base64.StdEncoding.DecodeString(tok)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, t)
}
