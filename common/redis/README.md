```
// 全局初始化调用一次
redis.InitRedis(conf)

// 使用示例
func GetData(key string) (string, error) {
    data, err := redis.GetBytes(key)
    if err != nil {
        return "", err
    }
    return data, nil
}
```
