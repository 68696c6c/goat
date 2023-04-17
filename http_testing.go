package goat

import (
	"net/http/httptest"
	"net/url"
)

type HandlerTest struct {
	router     Router
	remoteAddr string
	baseURL    *url.URL
	request    *Request
}

func NewHandlerTest(router Router) *HandlerTest {
	host := "127.0.0.1"
	u := new(url.URL)
	u.Scheme = "http"
	u.Host = host
	return &HandlerTest{
		router:     router,
		remoteAddr: host,
		baseURL:    u,
	}
}

// NewRequest creates a new Request using the HandlerTest remoteAddr as the base url.
func (h *HandlerTest) NewRequest(path string) (*Request, error) {
	request, err := NewRequest(h.baseURL.JoinPath(path).String())
	if err != nil {
		return nil, err
	}
	return request, nil
}

// SetRemoteAddr overwrites the default remote address to use when sending requests.
func (h *HandlerTest) SetRemoteAddr(ip string) *HandlerTest {
	h.remoteAddr = ip
	h.baseURL.Host = ip
	return h
}

func (h *HandlerTest) SetRequest(request *Request) *HandlerTest {
	h.request = request
	return h
}

func (h *HandlerTest) Send() (*httptest.ResponseRecorder, error) {
	rr := httptest.NewRecorder()
	request, err := h.request.Build()
	if err != nil {
		return nil, err
	}
	h.router.ServeHTTP(rr, request)
	return rr, nil
}

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// )
//
// type HandlerTest struct {
// 	router     Router
// 	remoteAddr string
// 	request    handlerTestRequest
// 	// request Request
// }
//
// type HandlerTestResponse struct {
// 	*httptest.ResponseRecorder
// 	BodyString string
// }
//
// // UnmarshalBody a handler test response to the provided struct.
// func (r *HandlerTestResponse) UnmarshalBody(m any) error {
// 	err := json.Unmarshal(r.Body.Bytes(), m)
// 	return err
// }
//
// // Map a handler test response to the provided struct.
// func (r *HandlerTestResponse) Map(m any) error {
// 	err := json.Unmarshal(r.Body.Bytes(), m)
// 	return err
// }
//
// func NewHandlerTest(r Router) *HandlerTest {
// 	return &HandlerTest{
// 		router:     r,
// 		remoteAddr: "127.0.0.1",
// 	}
// }
//
// // Send returns the result of sending the current request.
// // Panics if the request creation fails.
// func (h *HandlerTest) Send() *HandlerTestResponse {
// 	rr := httptest.NewRecorder()
// 	h.router.ServeHTTP(rr, h.request.Build())
// 	return &HandlerTestResponse{rr, rr.Body.String()}
// }
//
// // Headers sets the headers for the current request.
// func (h *HandlerTest) Headers(headers map[string]string) *HandlerTest {
// 	h.request.Headers = headers
// 	return h
// }
//
// // Body sets the body for the current request.
// func (h *HandlerTest) Body(data *map[string]any) *HandlerTest {
// 	h.request.SetBody(data)
// 	return h
// }
//
// // // SetFormData sets the form data for the current request.
// // func (h *HandlerTest) SetFormData(data url.Values) *HandlerTest {
// // 	h.request.SetFormData(data)
// // 	return h
// // }
//
// // SetRemoteAddr overwrites the default remote address to use when sending requests.
// func (h *HandlerTest) SetRemoteAddr(ip string) {
// 	h.remoteAddr = ip
// }
//
// func (h *HandlerTest) Request(method, url string) *HandlerTest {
// 	h.setRequest(method, url)
// 	return h
// }
//
// func (h *HandlerTest) RequestF(method, urlf string, a ...any) *HandlerTest {
// 	h.setRequest(method, fmt.Sprintf(urlf, a...))
// 	return h
// }
//
// func (h *HandlerTest) Get(url string) *HandlerTest {
// 	h.setRequest("GET", url)
// 	return h
// }
//
// func (h *HandlerTest) GetF(urlf string, a ...any) *HandlerTest {
// 	h.setRequest("GET", fmt.Sprintf(urlf, a...))
// 	return h
// }
//
// func (h *HandlerTest) Post(url string) *HandlerTest {
// 	h.setRequest("POST", url)
// 	return h
// }
//
// func (h *HandlerTest) PostF(urlf string, a ...any) *HandlerTest {
// 	h.setRequest("POST", fmt.Sprintf(urlf, a...))
// 	return h
// }
//
// func (h *HandlerTest) Put(url string) *HandlerTest {
// 	h.setRequest("PUT", url)
// 	return h
// }
//
// func (h *HandlerTest) PutF(urlf string, a ...any) *HandlerTest {
// 	h.setRequest("PUT", fmt.Sprintf(urlf, a...))
// 	return h
// }
//
// func (h *HandlerTest) Delete(url string) *HandlerTest {
// 	h.setRequest("DELETE", url)
// 	return h
// }
//
// func (h *HandlerTest) DeleteF(urlf string, a ...any) *HandlerTest {
// 	h.setRequest("DELETE", fmt.Sprintf(urlf, a...))
// 	return h
// }
//
// func (h *HandlerTest) setRequest(method, url string) {
// 	h.request = handlerTestRequest{
// 		RemoteAddr: h.remoteAddr,
// 		Method:     method,
// 		URL:        url,
// 	}
// }
//
// type handlerTestRequest struct {
// 	RemoteAddr  string
// 	Method      string
// 	URL         string
// 	Headers     map[string]string
// 	body        *bytes.Reader
// 	hasBody     bool
// 	formData    url.Values
// 	hasFormData bool
// }
//
// // Build returns an http.Request struct built from a handlerTestRequest.
// // Panics if the http.Request cannot be built.
// func (r *handlerTestRequest) Build() *http.Request {
// 	GetLogger().Info("http_test build request")
// 	req, err := http.NewRequest(r.Method, r.URL, r.GetBody())
// 	if err != nil {
// 		panic(err)
// 	}
// 	GetLogger().Info("http_test set remote addr")
// 	req.RemoteAddr = r.RemoteAddr
// 	GetLogger().Info("http_test set headers")
// 	for k, v := range r.Headers {
// 		req.Header.Set(k, v)
// 	}
// 	GetLogger().Info("http_test request build finished")
// 	return req
// }
//
// func (r *handlerTestRequest) SetBody(data *map[string]any) {
// 	if data == nil {
// 		r.body = bytes.NewReader([]byte{})
// 		r.hasBody = true
// 		return
// 	}
// 	body, err := json.Marshal(data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	r.body = bytes.NewReader(body)
// 	r.hasBody = true
// 	return
// }
//
// func (r *handlerTestRequest) GetBody() io.Reader {
// 	if r.hasBody {
// 		return r.body
// 	}
// 	return nil
// }
//
// func (r *handlerTestRequest) SetFormData(data url.Values) {
// 	r.formData = data
// 	r.hasFormData = true
// }
