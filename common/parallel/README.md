用于网关服务并发请求Grpc服务

```
import (
    "a.com/go-server/common/parallel"
)

type Result struct {
    A int `json:"a"`
    B int `json:"b"`
    C int `json:"c"`
}

x := 1
y := 2
c := 3
resp := Result{}

wait := parallel.WaitGroupWrapper{}
wait.Wrap(func() {
    resp.A = FunctionA(x)
})
wait.Wrap(func() {
    resp.B = FunctionB(y)
})
wait.Wrap(func() {
    resp.C = FunctionC(z)
})
wait.Wait()

return resp
```
