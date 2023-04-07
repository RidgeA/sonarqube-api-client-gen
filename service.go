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
	serviceTemplateName     = "service.tpl"
	serviceTemplateFileName = "./tpl/service.tpl"
)

var (
	serviceTemplate = template.Must(template.New(serviceTemplateName).Funcs(templateHelpers).ParseFiles(serviceTemplateFileName))
)

func renderService(in io.Writer, data *webService) error {

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
