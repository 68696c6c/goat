package goat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type HandlerTest struct {
	router     Router
	remoteAddr string
	request    handlerTestRequest
}

type HandlerTestResponse struct {
	*httptest.ResponseRecorder
	BodyString string
}

// Map a handler test response to the provided struct.
func (r *HandlerTestResponse) Map(m any) error {
	err := json.Unmarshal(r.Body.Bytes(), m)
	return err
}

func NewHandlerTest(r Router) *HandlerTest {
	return &HandlerTest{
		router:     r,
		remoteAddr: "127.0.0.1",
	}
}

// Send returns the result of sending the current request.
// Panics if the request creation fails.
func (h *HandlerTest) Send() *HandlerTestResponse {
	rr := httptest.NewRecorder()
	h.router.ServeHTTP(rr, h.request.Build())
	return &HandlerTestResponse{rr, rr.Body.String()}
}

// Headers sets the headers for the current request.
func (h *HandlerTest) Headers(headers map[string]string) *HandlerTest {
	h.request.Headers = headers
	return h
}

// Body sets the body for the current request.
func (h *HandlerTest) Body(data *map[string]any) *HandlerTest {
	h.request.SetBody(data)
	return h
}

// SetRemoteAddr overwrites the default remote address to use when sending requests.
func (h *HandlerTest) SetRemoteAddr(ip string) {
	h.remoteAddr = ip
}

func (h *HandlerTest) Request(method, url string) *HandlerTest {
	h.setRequest(method, url)
	return h
}

func (h *HandlerTest) RequestF(method, urlf string, a ...any) *HandlerTest {
	h.setRequest(method, fmt.Sprintf(urlf, a...))
	return h
}

func (h *HandlerTest) Get(url string) *HandlerTest {
	h.setRequest("GET", url)
	return h
}

func (h *HandlerTest) GetF(urlf string, a ...any) *HandlerTest {
	h.setRequest("GET", fmt.Sprintf(urlf, a...))
	return h
}

func (h *HandlerTest) Post(url string) *HandlerTest {
	h.setRequest("POST", url)
	return h
}

func (h *HandlerTest) PostF(urlf string, a ...any) *HandlerTest {
	h.setRequest("POST", fmt.Sprintf(urlf, a...))
	return h
}

func (h *HandlerTest) Put(url string) *HandlerTest {
	h.setRequest("PUT", url)
	return h
}

func (h *HandlerTest) PutF(urlf string, a ...any) *HandlerTest {
	h.setRequest("PUT", fmt.Sprintf(urlf, a...))
	return h
}

func (h *HandlerTest) Delete(url string) *HandlerTest {
	h.setRequest("DELETE", url)
	return h
}

func (h *HandlerTest) DeleteF(urlf string, a ...any) *HandlerTest {
	h.setRequest("DELETE", fmt.Sprintf(urlf, a...))
	return h
}

func (h *HandlerTest) setRequest(method, url string) {
	h.request = handlerTestRequest{
		RemoteAddr: h.remoteAddr,
		Method:     method,
		URL:        url,
	}
}

type handlerTestRequest struct {
	RemoteAddr string
	Method     string
	URL        string
	Headers    map[string]string
	body       *bytes.Reader
	hasBody    bool
}

// Build returns an http.Request struct built from a handlerTestRequest.
// Panics if the http.Request cannot be built.
func (r *handlerTestRequest) Build() *http.Request {
	req, err := http.NewRequest(r.Method, r.URL, r.GetBody())
	if err != nil {
		panic(err)
	}
	req.RemoteAddr = r.RemoteAddr
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	return req
}

func (r *handlerTestRequest) SetBody(data *map[string]any) {
	if data == nil {
		r.body = bytes.NewReader([]byte{})
		r.hasBody = true
		return
	}
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	r.body = bytes.NewReader(body)
	r.hasBody = true
	return
}

func (r *handlerTestRequest) GetBody() io.Reader {
	if r.hasBody {
		return r.body
	}
	return nil
}
