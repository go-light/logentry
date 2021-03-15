package logentry

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"unsafe"
)

func TestHttpClientLogEntry_Json(t *testing.T) {

	buf1 := []byte("sadsadas")
	fmt.Printf("n2 的类型 %T n2占中的字节数是 %d", buf1, unsafe.Sizeof(buf1))

	s := "postgres://user:pass@host.com:5432/path?k=v#f"
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		panic(err)

	}
	fmt.Println(u.Scheme, u.Host, u.Path, u.RawQuery)

	ctx := context.WithValue(context.Background(), "trace-id", 1)
	httpClientLogEntry := NewHttpClientLogEntry(ctx, WithTraceIDCtxName("trace-id"))
	httpClientLogEntry.Start()
	httpClientLogEntry.SetReqUrl(s)
	httpClientLogEntry.SetStatusCode(200)
	httpClientLogEntry.End()
	buf := httpClientLogEntry.Json()

	fmt.Println(string(buf))
	fmt.Println(httpClientLogEntry.Text())

}
