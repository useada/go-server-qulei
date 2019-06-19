```
package main

import (
	"context"
	"fmt"
	"time"
)

func init() {
	httpclient.Init(httpclient.Config{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90,
	})
}

var HttpHost = "https://api.something.com/logic/v1/list"

func main() {
	{
		t1 := time.Now()
		path := fmt.Sprintf("?uid=%d", 802799)

		r, err := httpclient.NewRequester("GET", HttpHost+path, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		r.SetHeaders(map[string]string{"Content-Type": "application/json"})

		data, err := r.Do(context.TODO(), 500*time.Millisecond)
		elapsed := time.Since(t1)
		fmt.Println(data, err, elapsed)
	}

	{
		t1 := time.Now()
		path := fmt.Sprintf("?uid=%d", 802799)

		r, err := httpclient.NewRequester("GET", HttpHost+path, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		r.SetHeaders(map[string]string{"Content-Type": "application/json"})

		data, err := r.Do(context.TODO(), 100*time.Millisecond)
		elapsed := time.Since(t1)
		fmt.Println(data, err, elapsed)
	}
}
```
