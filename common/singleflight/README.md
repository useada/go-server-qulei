```
import (
    "sync"
    "sync/atomic"
    "time"

    "a.com/go-server/common/singleflight"
)

// 使用示例
var g singleflight.Group

// 在一个时间段内只被调用一次
fn := func() (interface{}, error) {
    return db.find().where(), nil
}

v, err := g.Do("key", fn)
if err != nil {
    fmt.Println("Do error: %v", err)
}
```

https://github.com/golang/groupcache/tree/master/singleflight
