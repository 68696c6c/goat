package goat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	HeaderContentType string = "Content-Type"
	ContentTypeJson   string = "application/json"
	ContentTypeForm   string = "application/x-www-form-urlencoded"
)

type Request struct {
	baseURL         *url.URL
	query           url.Values
	formData        url.Values
	jsonBody        any
	contentType     string
	redirectHandler RedirectHandler
	Method          string
	URL             string
	Headers         map[string]string
}

func NewRequest(baseUrl string) (*Request, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	query := make(url.Values)
	if u.RawQuery != "" {
		query, err = url.ParseQuery(u.RawQuery)
		if err != nil {
			return nil, err
		}
	}
	return &Request{
		baseURL:         u,
		query:           query,
		formData:        make(url.Values),
		jsonBody:        nil,
		redirectHandler: nil,
		Method:          "",
		URL:             "",
		Headers:         make(map[string]string),
	}, nil
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) AddQueryParam(key, value string) *Request {
	r.query.Add(key, value)
	return r
}

func (r *Request) SetQueryParam(key, value string) *Request {
	r.query.Set(key, value)
	return r
}

func (r *Request) AddFormParam(key, value string) *Request {
	r.formData.Add(key, value)
	return r
}

func (r *Request) SetFormParam(key, value string) *Request {
	r.formData.Set(key, value)
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.Headers[key] = value
	if key == HeaderContentType {
		r.contentType = value
	}
	return r
}

func (r *Request) SetContentTypeJSON() *Request {
	r.SetHeader(HeaderContentType, ContentTypeJson)
	return r
}

func (r *Request) SetContentTypeForm() *Request {
	r.SetHeader(HeaderContentType, ContentTypeForm)
	return r
}

func (r *Request) SetBodyJSON(body any) *Request {
	r.SetContentTypeJSON()
	r.jsonBody = body
	return r
}

func (r *Request) GetBodyJSON() (io.Reader, error) {
	if r.jsonBody == nil {
		return nil, nil
	}
	result, err := json.Marshal(&r.jsonBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json body")
	}
	return bytes.NewBuffer(result), nil
}

func (r *Request) SetBodyForm(data url.Values) *Request {
	r.SetContentTypeForm()
	r.formData = data
	return r
}

func (r *Request) GetBodyForm() (io.Reader, error) {
	if r.formData == nil {
		return nil, nil
	}
	return strings.NewReader(r.formData.Encode()), nil
}

type RedirectHandler func(req *http.Request, via []*http.Request) error

func (r *Request) SetRedirectHandler(handler RedirectHandler) *Request {
	r.redirectHandler = handler
	return r
}

func (r *Request) Build() (*http.Request, error) {
	r.baseURL.RawQuery = r.query.Encode()
	reqURL := r.baseURL.String()

	var err error
	var body io.Reader

	switch r.contentType {
	case ContentTypeJson:
		body, err = r.GetBodyJSON()
		if err != nil {
			return nil, errors.Wrap(err, "failed to build request body json")
		}
	case ContentTypeForm:
		body, err = r.GetBodyForm()
		if err != nil {
			return nil, errors.Wrap(err, "failed to build request body form")
		}
	default:
		body = nil
	}

	req, err := http.NewRequest(r.Method, reqURL, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (r *Request) Send() (*http.Response, error) {
	req, err := r.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}

	client := http.DefaultClient

	if r.redirectHandler != nil {
		client.CheckRedirect = r.redirectHandler
	}

	response, err := client.Do(req)
	if err != nil {
		return response, errors.Wrap(err, "failed to send request")
	}

	if response.StatusCode > 299 {
		return response, errors.Wrapf(err, "received error response: %s", response.Status)
	}

	return response, nil
}

func (r *Request) SendAndRead(output any) (*http.Response, error) {
	response, err := r.Send()
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response, errors.Wrap(err, "failed to read response body")
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return response, errors.Wrap(err, "failed to unmarshal response body")
	}

	return response, nil
}

// RequestData is an alias of url.Values that supports method chaining.
type RequestData url.Values

func MakeRequestData() RequestData {
	return make(RequestData)
}

func (v RequestData) Set(key, value string) RequestData {
	v[key] = []string{value}
	return v
}

func (v RequestData) Add(key, value string) RequestData {
	v[key] = append(v[key], value)
	return v
}

func (v RequestData) Del(key string) RequestData {
	delete(v, key)
	return v
}

func (v RequestData) Has(key string) bool {
	_, ok := v[key]
	return ok
}

func (v RequestData) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v RequestData) Values() url.Values {
	return url.Values(v)
}
