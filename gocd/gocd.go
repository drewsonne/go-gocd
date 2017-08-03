package gocd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	libraryVersion = "1"
	userAgent      = "go-gocd/" + libraryVersion
	apiV1          = "application/vnd.go.cd.v1+json"
	apiV2          = "application/vnd.go.cd.v2+json"
	apiV3          = "application/vnd.go.cd.v3+json"
	apiV4          = "application/vnd.go.cd.v4+json"
)

type StringResponse struct {
	Message string `json:"message"`
}

type ClientInterface interface{}

type APIResponse struct {
	Http    *http.Response
	Body    string
	Request *APIRequest
}

type APIRequest struct {
	Http *http.Request
	Body string
}

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
	Auth      *Auth

	Agents            *AgentsService
	PipelineGroups    *PipelineGroupsService
	Stages            *StagesService
	Jobs              *JobsService
	PipelineTemplates *PipelineTemplatesService
	Pipelines         *PipelinesService
	PipelineConfigs   *PipelineConfigsService

	common service
	cookie string
}

type PaginationResponse struct {
	Offset   int64 `json:"offset"`
	Total    int64 `json:"total"`
	PageSize int64 `json:"page_size"`
}

type service struct {
	client *Client
}

type Auth struct {
	Username string
	Password string
}

type Configuration struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	SslCheck bool   `yaml:"ssl_check,omitempty"`
}

func (c *Configuration) HasAuth() bool {
	return (c.Username != "") && (c.Password != "")
}

func (c *Configuration) Client() *Client {
	return NewClient(c, nil)
}

func NewClient(cfg *Configuration, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if strings.HasPrefix(cfg.Server, "https") {
		if !cfg.SslCheck {
			httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	}

	baseURL, _ := url.Parse(cfg.Server)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.common.client = c
	c.Agents = (*AgentsService)(&c.common)
	c.PipelineGroups = (*PipelineGroupsService)(&c.common)
	c.Stages = (*StagesService)(&c.common)
	c.Jobs = (*JobsService)(&c.common)
	c.PipelineTemplates = (*PipelineTemplatesService)(&c.common)
	c.Pipelines = (*PipelinesService)(&c.common)
	c.PipelineConfigs = (*PipelineConfigsService)(&c.common)
	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}, apiVersion string) (*APIRequest, error) {
	rel, err := url.Parse("api/" + urlStr)

	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	request := &APIRequest{}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		bdy, _ := ioutil.ReadAll(buf)
		request.Body = string(bdy)

		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	request.Http = req
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if apiVersion != "" {
		req.Header.Set("Accept", apiVersion)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	if c.cookie == "" {
		if c.Auth != nil {
			req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
		}
	} else {
		req.Header.Set("Cookie", c.cookie)
	}

	return request, nil
}

func (c *Client) Do(ctx context.Context, req *APIRequest, v interface{}) (*APIResponse, error) {

	req.Http = req.Http.WithContext(ctx)

	response := &APIResponse{
		Request: req,
	}

	resp, err := c.client.Do(req.Http)
	if err != nil {
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}

	response.Http = resp
	err = CheckResponse(response.Http)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			bdy, err := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(bdy, v)
			response.Body = string(bdy)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}
func CheckResponse(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		bdy, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf(
			"Received HTTP Status '%s': '%s'",
			response.Status,
			bdy,
		)
	}
	return nil
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
//func addOptions(s string, opt interface{}) (string, error) {
func addOptions(s string) (string, error) {
	//v := reflect.ValueOf(opt)
	//if v.Kind() == reflect.Ptr && v.IsNil() {
	//	return s, nil
	//}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	return u.String(), nil
}

//type ResourceNotFound struct {
//	When     time.Time
//	Resource string
//}
//
//func (e ResourceNotFound) Error() string {
//	return fmt.Sprintf("Could not find '%s'.", e.Resource)
//	error.
//}