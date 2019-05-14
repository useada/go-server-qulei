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

errgroup 遇到错误取消整个goroutine
```
import "golang.org/x/sync/errgroup"
......
eg, ctx := errgroup.WithContext(context.TODO())
for _, w := range work {
		w := w
		eg.Go(func() error {
			// Do something with w and
			// listen for ctx cancellation
		})
}
// If any of the goroutines returns an error ctx will be
// canceled and err will be non-nil.
if err := eg.Wait(); err != nil {
	return err
}
```
