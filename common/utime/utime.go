package utime

import "time"

func Millisec() int64 {
	return time.Now().UnixNano() / 1e6
}

// UtcZero 当天utc零点的时间戳
func UtcZero() (int64, error) {
	timess := time.Now().Format("2006-01-02")
	t, err := time.Parse("2006-01-02", timess)
	return t.UnixNano() / 1e6, err
}
