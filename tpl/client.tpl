package {{.PackageName}}

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"

	"github.com/pkg/errors"
)

type Client struct {
	host string
	username string
	password string
	transport *http.Client
{{- range .WebServices}}
	{{.Variable}} *{{.ServiceName}}
{{- end }}
}

type httpErrorResponse struct {
	Errors []*httpErrorResponseMsg {{tick}}json:"errors"{{tick}}
}

func (er *httpErrorResponse) String() string {
	msgA := make([]string, len(er.Errors))
	for i, em := range er.Errors {
		msgA[i] = em.String()
	}
	return strings.Join(msgA, ", ")
}

type httpErrorResponseMsg struct {
	Msg string {{tick}}json:"msg"{{tick}}
}

func (em *httpErrorResponseMsg) String() string {
	return em.Msg
}

type HttpError struct {
	status   int
	parsed   *httpErrorResponse
	response *http.Response
	err      error
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("http code - %d, msg: %s", he.status, he.parsed.String())
}

func (he *HttpError) Response() *http.Response {
	return he.response
}

func checkHttpErrors(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil
	default:
		var err error
		bytesData, err := ioutil.ReadAll(resp.Body)

		errorResponse := &httpErrorResponse{}
		if err == nil {
			err = json.Unmarshal(bytesData, errorResponse)
		}

		if err != nil {
			msg := &httpErrorResponseMsg{
				Msg: "Unknown error: " + err.Error(),
			}
			errorResponse = &httpErrorResponse{
				Errors: []*httpErrorResponseMsg{msg},
			}
		}

		result := &HttpError{
			status:   resp.StatusCode,
			response: resp,
			parsed:   errorResponse,
		}
		return result
	}
}

func NewClient(httpClient *http.Client, host, username, password string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		host: host,
		transport: httpClient,
		username: username,
		password: password,
	}

{{- range .WebServices}}
	c.{{.Variable}} = New{{.ServiceName}}(c)
{{- end }}

	return c
}

func (c *Client) invoke(ctx context.Context, post bool, url string, payload interface{}) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url = c.host + "/" + url

	method := http.MethodGet
	if post {
		method = http.MethodPost
	}

	var err error
	values, err := query.Values(payload)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse payload")
	}

	var req *http.Request
	var body io.Reader
	if method == http.MethodGet {
		url = url + "?" + values.Encode()
	} else {
		body = strings.NewReader(values.Encode())
	}

	req, err = http.NewRequest(method, url, body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	if len(c.username) != 0 && len(c.password) != 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	req.WithContext(ctx)

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkHttpErrors(resp); err != nil {
		return nil, errors.Wrapf(err, "got error response (url: %s)", url)
	}
	return resp, nil
}

{{- range .WebServices}}
{{- template "getter" .}}
{{- end}}

{{- define "getter"}}
// {{.Description}}
{{- if .Since }}
// Since : {{.Since}}
{{- end}}
{{- if .Deprecated }}
// Deprecated
{{- end}}
{{- if .Internal }}
// Internal
{{- end}}
func (c *Client) {{.Getter}}() *{{.ServiceName}} {
	return c.{{.Variable}}
}
{{- end}}

// Helper function to convert string to pointer to string
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
