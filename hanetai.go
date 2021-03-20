package hanetai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	tokenSource oauth2.TokenSource

	client HttpClient // HTTP client used to communicate with the API.

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent used when communicating with the Hanet AI API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Places  *PlaceService
	Persons *PersonService
}

type service struct {
	client *Client
}

const (
	defaultBaseURL = "https://partner.hanet.ai/"
	userAgent      = "hanetai-sdk"
)

func NewClient(httpClient HttpClient, ts oauth2.TokenSource) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			// > Cần thiết lập timeout cho API này từ 10 - 30s
			// Thiết lập timeout dùng cho 2 API employee/register và person/updateByFaceImage.
			Timeout: 30 * time.Second,
		}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:      httpClient,
		BaseURL:     baseURL,
		UserAgent:   userAgent,
		tokenSource: ts,
	}

	c.common.client = c

	c.Places = (*PlaceService)(&c.common)
	c.Persons = (*PersonService)(&c.common)

	return c
}

type requestBodyFn = func(token string) (io.Reader, string, error)

func urlencodeBody(body interface{}) requestBodyFn {
	return func(token string) (io.Reader, string, error) {
		v, err := query.Values(body)
		if err != nil {
			return nil, "", err
		}
		v.Add("token", token)
		body := strings.NewReader(v.Encode())

		return body, "application/x-www-form-urlencoded", nil
	}
}

func multipartBody(file io.Reader, fn func(m *multipart.Writer) error) requestBodyFn {
	return func(token string) (io.Reader, string, error) {
		body := bytes.NewBuffer(nil)

		w := multipart.NewWriter(body)
		err := fn(w)
		if err != nil {
			return nil, "", err
		}

		w.WriteField("token", token)
		f, err := w.CreateFormFile("file", fileName)
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return nil, "", err
		}

		err = w.Close()
		if err != nil {
			return nil, "", err
		}

		return body, w.FormDataContentType(), nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(urlStr string, fn requestBodyFn) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	body, contentType, err := fn(token.AccessToken)
	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	var env envelope
	err = json.NewDecoder(resp.Body).Decode(&env)
	if err != nil {
		return resp, err
	}

	if env.ReturnCode != 1 {
		return resp, &ServerError{
			Code:    env.ReturnCode,
			Message: env.ReturnMessage,
		}
	}

	if v != nil {
		err = json.Unmarshal(env.Data, v)
	}

	return resp, err
}

type envelope struct {
	StatusCode    int             `json:"statusCode"`
	ReturnCode    int             `json:"returnCode"`
	ReturnMessage string          `json:"returnMessage"`
	Data          json.RawMessage `json:"data"`
}
