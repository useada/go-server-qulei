### Board 接口

#### 1. 评论列表

---

##### 必须Auth认证
```
否
```

---

##### 请求参数

| | 必选 | 类型 | 说明 |
|---------|------|--------------|--------------|
| oid | y | string | 对象(feed) id|
| cid | n | string | 父级评论id(如果想要拉取某评论的二级评论列表)|
| page_token| n | string | 用于分页的token 初始为空串|

---

##### 请求方法
```
GET
```

---

##### 调用样例
```
192.168.200.120:8091/v1/board/comm/list?oid=test_oid_00001&cid=&page_token=
```
---

##### 返回结果

```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "items": [
            {
                "id": "bg2frooscsgmt2ta60gg",
                "oid": "test_oid_00001",
                "content": "Hello, this is NO 41 comment.",
                "reply_count": 3,
                "created_at": 1543831011790,
                "author": {
                    "uid": "testuid"
                },
                "replys": [
                    {
                        "id": "bg2h5dgscsgg39kdhjl0",
                        "oid": "test_oid_00001",
                        "cid": "bg2frooscsgmt2ta60gg",
                        "content": "Hello, this is NO 41 - 04 comment.",
                        "created_at": 1543836342689,
                        "author": {
                            "uid": "testuid"
                        }
                    },
                    {
                        "id": "bg2h5c8scsgg39kdhjkg",
                        "oid": "test_oid_00001",
                        "cid": "bg2frooscsgmt2ta60gg",
                        "content": "Hello, this is NO 41 - 03 comment.",
                        "created_at": 1543836337567,
                        "author": {
                            "uid": "testuid"
                        }
                    }
                ]
            },
            .
            .
            .
            {
                "id": "bg2fmooscsgmt2ta6070",
                "oid": "test_oid_00001",
                "content": "Hello, this is NO 22 comment.",
                "created_at": 1543830371100,
                "author": {
                    "uid": "testuid"
                }
            }
        ],
        "page_token": "eyJvZmZzZXQiOjE1NDM4MzAzNzExMDAsImxpbWl0IjoyMH0="
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 2. 获取评论

---

##### 必须Auth认证
```
否
```

---

##### 请求参数

| | 必选 | 类型 | 说明 |
|---------|----|-------------|-------------|
| id  | y | string |评论id |
| oid | y | string |对象(feed) id |

---

##### 请求方法
```
GET
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/comm/get?id=bg2ctb8scsgjultl73k0&oid=test_oid_00001
```

---

##### 返回结果

```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "bg2ctb8scsgjultl73k0",
        "oid": "test_oid_00001",
        "content": "Hello, this is NO 21 comment.",
        "created_at": 1543818925008,
        "author": {
            "uid": "testuid"
        }
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 3. 创建评论

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001",
  "cid": "bg2frooscsgmt2ta60gg",                // cid不为空，则当前评论是 bg2frooscsgmt2ta60gg 的二级评论
  "is_repost": false,                           // 该条评论是否是一条feed转发
  "content": "Hello, this is NO 01 comment.",
  "img_id": "",                                 // 已经上传到文件服务的图片id
  "img_ex": ""                                  // jpg / png
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/comm/new
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "bg2h5dgscsgg39kdhjl0"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 4. 删除评论

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001",
  "id": "bg2csq8scsgjultl73eg"
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/comm/del
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "bg2foroscsgmt2ta60d0"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 5. 评论点赞

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001",
  "cid": "bg2fpmoscsgmt2ta60fg"
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/comm/like
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "testuidbg2fpmoscsgmt2ta60fg"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 6. 评论取消点赞

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001",
  "cid": "bg2fq70scsgmt2ta60g0"
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/comm/unlike
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "testuidbg2fpmoscsgmt2ta60fg"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 7. 点赞列表

---

##### 必须Auth认证
```
否
```

---

##### 请求参数

| | 必选 | 类型 | 说明 |
|---------|----|-------------|-------------|
| oid | y | string |对象(feed) id |
| page_token  | n | string | 分页token|

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/like/list?oid=test_oid_00001&page_token=
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "items": [
            {
                "id": "testuidtest_oid_00001",
                "oid": "test_oid_00001",
                "author": {
                    "uid": "testuid"
                }
            }
        ]
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 8. 创建点赞

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001"
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/like/new
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "testuidtest_oid_00001"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 9. 取消点赞

---

##### 必须Auth认证
```
是
```

---

##### 请求参数

```
{
  "oid": "test_oid_00001"
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/like/del
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "id": "testuidtest_oid_00001"
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---

#### 10. 评论点赞汇总

---

##### 必须Auth认证
```
否
```

---

##### 请求参数

```
{
  "oids": ["test_oid_00001"]
}
```

---

##### 请求方法
```
POST
```

---

##### 调用样例

```
192.168.200.120:8091/v1/board/summary/mget
```

---

##### 返回结果


```
正确返回值
{
    "code": "2000",
    "msg": "OK",
    "data": {
        "items": [
            {
                "id": "test_oid_00001",
                "comms_count": 45,
                "comms_first_count": 41
            }
        ]
    }
}
```

```
错误返回值
{
    "code": "4000",
    "msg": "Key: '.Oid' Error:Field validation for 'Oid' failed on the 'required' tag",
    "data": null
}
```
---