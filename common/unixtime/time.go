package unixtime

import "time"

func UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
