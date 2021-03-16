package logentry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type HttpClientLogEntry interface {
	Start()
	End()
	Json() []byte
	Text() string

	SetReqUrl(reqUrl string)
	SetMethod(method string)
	SetStatusCode(statusCode int)
	SetLocalIP(localIP string)
	SetLocalApp(localApp string)
	SetRemoteIP(remoteIP string)
	SetRemoteApp(remoteApp string)
	SetReqSizeBytes(reqSizeBytes string)
	SetRespSizeBytes(respSizeBytes string)

	SetTraceID(traceID string)
}

func NewHttpClientLogEntry(ctx context.Context, options ...Option) HttpClientLogEntry {
	h := &httpClientLogEntry{
		ctx: ctx,
	}

	for _, o := range options {
		o.Apply(h)
	}

	return h
}

func (h *httpClientLogEntry) Start() {
	h.startTime = time.Now()
}

func (h *httpClientLogEntry) End() {
	h.Duration = time.Since(h.startTime)
	h.CostMS = fmt.Sprintf("%d", h.Duration.Milliseconds())
}

func (h *httpClientLogEntry) Json() []byte {
	traceID := ""
	if traceIDInterface := h.ctx.Value(h.traceIDCtxName); traceIDInterface != nil {
		value, ok := traceIDInterface.(string)
		if ok {
			traceID = value
		}
	}

	h.TraceID = traceID
	h.Started = h.startTime.Format("2006-01-02 15:04:05.000")
	buf, err := json.Marshal(h)
	if err != nil {

	}

	return buf
}

func (h *httpClientLogEntry) Text() string {
	return fmt.Sprintf("started=%s,cost_ms=%s,req_url=%s,method=%s,status_code=%s,req_size_bytes=%s,resp_size_bytes=%s,local_ip=%s,local_app=%s,remote_ip=%s,remote_app=%s,trace_id%s,span_id=%s",
		h.Started, h.CostMS, h.ReqUrl, h.Method, h.StatusCode, h.ReqSizeBytes, h.RespSizeBytes,
		h.LocalIP, h.LocalApp, h.RemoteIP, h.RemoteApp,
		h.TraceID, h.SpanID)
}

func (h *httpClientLogEntry) SetReqUrl(reqUrl string) {
	h.ReqUrl = reqUrl
}

func (h *httpClientLogEntry) SetMethod(method string) {
	h.Method = method
}

func (h *httpClientLogEntry) SetStatusCode(statusCode int) {
	h.StatusCode = fmt.Sprintf("%d", statusCode)
}

func (h *httpClientLogEntry) SetLocalIP(localIP string) {
	h.LocalIP = localIP
}

func (h *httpClientLogEntry) SetLocalApp(localApp string) {
	h.LocalApp = localApp
}

func (h *httpClientLogEntry) SetTraceID(traceID string) {
	h.TraceID = traceID
}

func (h *httpClientLogEntry) SetRemoteApp(remoteApp string) {
	h.RemoteApp = remoteApp
}

func (h *httpClientLogEntry) SetRemoteIP(remoteIP string) {
	h.RemoteIP = remoteIP
}

func (h *httpClientLogEntry) SetReqSizeBytes(reqSizeBytes string) {
	h.ReqSizeBytes = reqSizeBytes
}

func (h *httpClientLogEntry) SetRespSizeBytes(respSizeBytes string) {
	h.RespSizeBytes = respSizeBytes
}

type httpClientLogEntry struct {
	ctx            context.Context
	traceIDCtxName string
	spanIDCtxName  string

	//startTime 请求开始时间
	startTime time.Time `json:"-"`

	//Started 请求开始时间字符串格式2006-01-02 15:04:05.000
	Started string `json:"started"`

	Duration time.Duration `json:"-"`

	//CostMS 请求耗时
	CostMS string `json:"cost_ms"`

	//ReqUrl 请求地址
	ReqUrl string `json:"req_url"`

	//Method HTTP method
	Method string `json:"method"`

	//StatusCode http状态码
	StatusCode string `json:"status_code"`

	//LocalIP 当前应用服务器ip
	LocalIP string `json:"local_ip"`

	//LocalApp 当前应用名
	LocalApp string `json:"local_app"`

	// 请求应用名
	RemoteApp string `json:"remote_app"`

	// 请求服务ip
	RemoteIP string `json:"remote_ip"`

	//traceId TraceId
	TraceID string `json:"trace_id"`

	//spanId SpanId
	SpanID string `json:"span_id"`

	//SpanKind Span 类型

	//reqSizeBytes 请求大小
	ReqSizeBytes string `json:"req_size_bytes"`

	//respSizeBytes 响应大小
	RespSizeBytes string `json:"resp_size_bytes"`

	//sysBaggage 系统透传的 baggage 数据

	//bizBaggage 业务透传的 baggage 数据
}

type accessLogEntry struct {
	Method       string  `json:"method"`
	Host         string  `json:"host"`
	Path         string  `json:"path"`
	IP           string  `json:"ip"`
	ResponseTime float64 `json:"response_time"`

	encoded []byte
	err     error
}

// An Option configures a mutex.
type Option interface {
	Apply(*httpClientLogEntry)
}

// OptionFunc is a function that configures a mutex.
type OptionFunc func(*httpClientLogEntry)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *httpClientLogEntry) {
	f(mutex)
}

// WithExpiry can be used to set the expiry of a mutex to the given value.
func WithTraceIDCtxName(traceIDCtxName string) Option {
	return OptionFunc(func(m *httpClientLogEntry) {
		m.traceIDCtxName = traceIDCtxName
	})
}

// WithTries can be used to set the number of times lock acquire is attempted.
func WithSpanIDCtxName(spanIDCtxName string) Option {
	return OptionFunc(func(m *httpClientLogEntry) {
		m.spanIDCtxName = spanIDCtxName
	})
}
