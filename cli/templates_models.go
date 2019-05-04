package cli

import (
	"fmt"
	"os"

	"github.com/jinzhu/inflection"
)

type field struct {
	Name      string
	FieldName string
	Type      string
	Required  bool
	Tag       string
}

type model struct {
	Name       string
	StructName string
	Fields     []*field
	BelongsTo  []string `yaml:"belongs_to"`
	HasMany    []string `yaml:"has_many"`
}

func createModels() {
	cmdPath := config.SRCPath + "/models"
	err := os.MkdirAll(cmdPath, os.ModePerm)
	handleError(err)

	// Create models.
	for _, m := range config.Models {
		m.StructName = snakeToCamel(m.Name)

		// Build relations.
		if len(m.BelongsTo) > 0 {
			println("model belongs to: ")
			for _, r := range m.BelongsTo {
				println("relation: ", r)
				f := &field{
					Name:      fmt.Sprintf("%s_id", r),
					FieldName: fmt.Sprintf("%sID", snakeToCamel(r)),
					Type:      "goat.ID",
				}
				m.Fields = append([]*field{f}, m.Fields...)
			}
		}
		if len(m.HasMany) > 0 {
			println("model has many: ")
			for _, r := range m.HasMany {
				println(r)
				t := inflection.Singular(r)
				f := &field{
					Name:      r,
					FieldName: snakeToCamel(r),
					Type:      fmt.Sprintf("[]*%s", snakeToCamel(t)),
				}
				m.Fields = append(m.Fields, f)
			}
		}

		// Set field names and annotations.
		for _, f := range m.Fields {
			if f.FieldName == "" {
				f.FieldName = snakeToCamel(f.Name)
			}
			var extra string
			if f.Required {
				extra = ` binding:"required"`
			}
			f.Tag = fmt.Sprintf(`json:"%s"%s`, f.Name, extra)
		}
		generateModel(cmdPath, m.Name, modelTemplate, *m)
	}
}

var modelTemplate = `
package models

import "github.com/68696c6c/goat"

// swagger:model {{.StructName}}
{{ $tick := "` + "`" + `" }}
type {{.StructName}} struct {
	goat.Model
	{{- range $key, $value := .Fields }}
	{{ $value.FieldName }} {{ $value.Type }} {{ $tick }}{{ $value.Tag }}{{ $tick }}
	{{- end }}
}

`
