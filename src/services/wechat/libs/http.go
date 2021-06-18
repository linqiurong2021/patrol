package libs

import "github.com/go-resty/resty/v2"

type Http struct {
	//
	client *resty.Client
}
//
func NewHttp() *Http {
	client := resty.New()
	return &Http{
		client: client,
	}
}
// Get
func (h *Http) Get(url string) (result *resty.Response, err error) {
	result, err = h.client.R().Get(url)
	return
}
// Post
func (h *Http) Post(url string,body string) (result *resty.Response, err error) {
	result, err = h.client.R().SetBody(body).Post(url)
	return
}