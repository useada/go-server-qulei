用于网关服务并发请求Grpc服务

```
import (
    "a.com/go-server/common/concurrency"
)

type Result struct {
    A int `json:"a"`
    B int `json:"b"`
    C int `json:"c"`
}
resp := Result{}

wait := concurrency.WaitGroupWrapper{}

x := 1
wait.Wrap(func() {
    resp.A = FunctionA(x)
})

y := 2
wait.Wrap(func() {
    resp.B = FunctionB(y)
})

z := 3
wait.Wrap(func() {
    resp.C = FunctionC(z)
})

wait.Wait()

return resp
```
