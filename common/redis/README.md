```
import (
    "a.com/go-server/common/redis"
)

// 全局初始化调用一次
func init() {
    redis.Init(conf)
}

// 使用示例
func DemoGetRedisData(key string) (string, error) {
    data, err := redis.GetBytes(key)
    if err != nil {
        return "", err
    }
    return data, nil
}
```
