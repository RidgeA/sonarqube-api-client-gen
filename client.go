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
	clientTemplateName     = "client.tpl"
	clientTemplateFileName = "./tpl/client.tpl"
	clientFileName         = clientTemplateName + ".go"
)

var (
	clientTemplate = template.Must(template.New(clientTemplateName).Funcs(templateHelpers).ParseFiles(clientTemplateFileName))
)

func renderClient(in io.Writer, data *apiDefinition) error {

	buff := bytes.NewBuffer([]byte{})

	if err := clientTemplate.Execute(buff, data); err != nil {
		return errors.Wrapf(err, "failed to render client")
	}

	src := buff.Bytes()

	formatted, err := format.Source(src)
	if err != nil {
		log.Printf("failed to format source of %s/client.go: err:%s", data.PackageName, err.Error())
		formatted = src
	}

	_, err = in.Write(formatted)
	return err
}
