package json

type Base struct {
	ID        string `json: "id"`
	State     int    `json: "state"`
	CreatedAt int64  `json: "created_at"`
	UpdatedAt int64  `json: "updated_at"`
}

type Geo struct {
	Lon  float64 `json:"lon" `
	Lat  float64 `json:"lat" `
	Addr string  `json:"addr" binding:"lte=255"`
}

type Image struct {
	ImgID  string `json:"img_id"`
	ImgExt string `json:"img_ext"`
}
