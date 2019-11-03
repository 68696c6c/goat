package generator

import "github.com/pkg/errors"

const packageHTTP = "http"

const httpTemplate = `
package models

import "github.com/68696c6c/goat"

{{ $tick := "` + "`" + `" }}
// swagger:model {{ .StructName }}
type {{ .StructName }} struct {
	goat.Model
	{{- range $key, $value := .Fields }}
	{{ $value.FieldName }} {{ $value.Type }} {{ $tick }}{{ $value.Tag }}{{ $tick }}
	{{- end }}
}

`

const controllerTemplate = `
`

type middleware struct {
	name string
}

type handler struct {
	name         string
	dependencies string
	middlewares  []middleware
}

type controller struct {
	name       string
	structName string
	handlers   []handler
}

type HTTP struct {
	controllers []controller
	handlers    []handler
}

func CreateHTTP(config *ProjectConfig) error {
	err := CreateDir(config.Paths.HTTP)
	if err != nil {
		return errors.Wrapf(err, "failed to create http directory '%s'", config.Paths.Models)
	}

	createHandlers := func(handlers []handler) []handler {
		var result []handler
		for _, h := range handlers {
			result = append(result, handler{
				name:         h.name,
				dependencies: h.dependencies,
				middlewares:  h.middlewares,
			})
		}
		return result
	}

	// Check for controllers first.
	for _, c := range config.HTTP.controllers {
		c.structName = snakeToCamel(c.name)

		c.handlers = createHandlers(c.handlers)

		err = GenerateFile(config.Paths.Models, c.name, controllerTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate model")
		}
	}

	return nil
}
