package json

type BaseIterface interface {
	GetId() string
	GetState() int
	GetCreatedAt() int64
	GetUpdatedAt() int64
}

type Base struct {
	Id        string `json: "id"`
	State     int    `json: "state"`
	CreatedAt int64  `json: "created_at"`
	UpdatedAt int64  `json: "updated_at"`
}

func (b *Base) GetId() string {
	return b.Id
}

func (b *Base) GetState() int {
	return b.State
}

func (b *Base) GetCreatedAt() int64 {
	return b.CreatedAt
}

func (b *Base) GetUpdatedAt() int64 {
	return b.UpdatedAt
}

type Geo struct {
	Lon  float64 `json:"lon" `
	Lat  float64 `json:"lat" `
	Addr string  `json:"addr" binding:"lte=255"`
}

type Image struct {
	ImgId string `json:"img_id"`
	ImgEx string `json:"img_ex"`
}
