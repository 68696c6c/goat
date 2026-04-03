package goat

import (
	"bytes"
	html "html/template"
	text "text/template"

	"github.com/pkg/errors"
)

func ParseTextTemplateToString(name, templateString string, data any) (string, error) {
	var tpl bytes.Buffer
	t := text.Must(text.New(name).Parse(templateString))
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", errors.Wrapf(err, "failed parse text template '%s'", name)
	}
	return tpl.String(), nil
}

func ParseHTMLTemplateToString(name, templateString string, data any) (string, error) {
	var tpl bytes.Buffer
	t := html.Must(html.New(name).Parse(templateString))
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", errors.Wrapf(err, "failed parse html template '%s'", name)
	}
	return tpl.String(), nil
}
