package store

import (
	"a.com/go-server/service/esearch/model"
)

type Store interface {
	UsersByName(name string, offset, limit int) ([]model.SearchInfo, error)
	UsersByNear(lat, lon float64, gender, offset, limit int) ([]model.SearchInfo, error)
}
