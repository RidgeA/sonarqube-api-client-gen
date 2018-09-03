package main

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"text/template"

	"github.com/pkg/errors"
)

const (
	serviceTemplateName = "service"
)

var (
	serviceTemplateText = `package {{.PackageName}}

import (
	"context"
	"net/http"
	
	"github.com/pkg/errors"
)

// {{.Description | formatDescription }}
{{- if .Since }}
// Since : {{.Since}}
{{- end}}
{{- if .Deprecated }}
// Deprecated 
{{- end}}
{{- if .Internal }}
// Internal  
{{- end}}
type {{.ServiceName}} struct {
	client *Client
	url string
}


{{- if .Deprecated }}
// Deprecated 
{{- end}}
func New{{.ServiceName}} (client *Client) *{{.ServiceName}}{
	s := &{{.ServiceName}}{
		client: client,
		url: "{{.Path}}",
	}
	return s
}

{{- range $index, $element := .Actions}}
{{template "action" $element}}
{{- end}}
` + actionTemplateText + requestTemplateText + responseTemplateText

	actionTemplateText = `{{- define "action"}}
// {{.Description | formatDescription }}
{{- if .Since}}
// Since {{.Since}}
{{- end}}
{{- if len .Changelog }}
//
// Changelog:
	{{- range .Changelog }}
// {{.String}}
	{{- end}}
{{- end}}
{{- if .Deprecated }}
//
// Deprecated since {{.DeprecatedSince}}
{{- end}}
func (s *{{.ServiceName}}) {{.MethodName}} (ctx context.Context{{- if .Params}}, request *{{.RequestTypeName}}{{- end}}) (*{{.ResponseTypeName}}, error) {
	resp, err := s.client.invoke(ctx, {{.Post}}, s.url + "/" + "{{.Key}}", {{- if .Params}} request {{- else}} nil {{- end}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call {{.ServiceName}}.{{.MethodName}}")
	}
	return &{{.ResponseTypeName}}{
		Response: resp,
	}, nil
}
{{- if .Params }}
{{ template "request" .}}
{{- end}}

{{ template "response" .}}

{{- end}}`

	requestTemplateText = `{{- define "request"}}
type {{.RequestTypeName}} struct {
{{- /* see https://github.com/golang/go/issues/18221#issuecomment-394255883 */}}
{{- range .Params}}
	// {{.Description | formatDescription }}
	{{- if .Since }}
	// Since {{.Since}}
	{{- end}}
	{{- if .Required}}
	// Required
	{{- end }}
	{{- if .Internal}}
	// Internal
	{{- end }}
	{{- if .DefaultValue }}
	// Default: {{.DefaultValue}}
	{{- end}}
	{{- if .ExampleValue}}
	// Example: {{.ExampleValue}}
	{{- end }}
	{{- if .PossibleValues }}
	// Possible values: {{- range .PossibleValues}} "{{.}}", {{- end}}
	{{- end}}
	{{- if .Deprecated}}
	// Deprecated since {{.DeprecatedSince.String}}
	{{- end }}
	{{.ParamName}} *string {{tick}}json:"{{.Key}},omitempty" url:"{{.Key}},omitempty"{{tick}}
{{- end}}
}
{{- end}}
`
	responseTemplateText = `{{- define "response"}}
type {{.ResponseTypeName}} struct {
	*http.Response
} 
{{- end}}
`
	serviceTemplate = template.Must(template.New(serviceTemplateName).Funcs(templateHelpers).Parse(serviceTemplateText))
)

func renderService(in io.WriteCloser, data *webService) error {
	defer in.Close()

	buff := bytes.NewBuffer([]byte{})

	if err := serviceTemplate.Execute(buff, data); err != nil {
		return errors.Wrapf(err, "failed to render service %s", data.ServiceName())
	}

	src := buff.Bytes()

	formatted, err := format.Source(src)
	if err != nil {
		log.Printf("failed to format source of %s, %s", data.ServiceName(), err.Error())
		formatted = src
	}

	_, err = in.Write(formatted)

	return err
}
