package links

import "net/url"

type Service interface {
	SetUrl(key string, value *url.URL)
	AddBaseUrlPath(key, path string)
	GetUrl(key ...string) *url.URL
	GetBaseUrl() *url.URL
	// GetUrlPath(key, path string) *url.URL
}

// type Url struct {
// 	*url.URL
// }
//
// func (u *Url) SetQueryParams(params url.Values) {
// 	q := u.Query()
// 	for key, value := range params {
// 		for _, v := range value {
// 			q.Add(key, v)
// 		}
// 	}
// }

func copyUrl(input *url.URL) *url.URL {
	result, _ := url.Parse(input.String())
	return result
}

func NewService(baseUrl *url.URL) Service {
	return &service{
		baseUrl: baseUrl,
		urls:    make(urlMap),
	}
}

type urlMap map[string]*url.URL

type service struct {
	baseUrl *url.URL
	urls    urlMap
}

func (s *service) SetUrl(key string, value *url.URL) {
	s.urls[key] = value
}

func (s *service) GetUrl(key ...string) *url.URL {
	if len(key) > 0 {
		u, ok := s.urls[key[0]]
		if ok {
			return copyUrl(u)
		}
	}
	return s.GetBaseUrl()
}

func (s *service) GetBaseUrl() *url.URL {
	return copyUrl(s.baseUrl)
}

// func (s *service) GetUrlPath(key, path string) *url.URL {
// 	println("getting url path for '" + key + "' with path '" + path + "'")
// 	base, ok := s.urls[key]
// 	if ok {
// 		return base.JoinPath(path)
// 	}
// 	// TODO: test this?
// 	return s.baseUrl.JoinPath(path)
// }

func (s *service) AddBaseUrlPath(key, path string) {
	s.urls[key] = s.baseUrl.JoinPath(path)
}
