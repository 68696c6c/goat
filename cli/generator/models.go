package generator

import (
	"fmt"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
)

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

type field struct {
	Name      string
	FieldName string
	Type      string
	Required  bool
	Tag       string
}

type Model struct {
	Name       string
	StructName string
	Fields     []*field
	BelongsTo  []string `yaml:"belongs_to"`
	HasMany    []string `yaml:"has_many"`
}

func CreateModels(config *ProjectConfig) error {
	err := CreateDir(config.ModelsPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create models directory '%s'", config.ModelsPath)
	}

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

		err = GenerateFile(config.ModelsPath, m.Name, modelTemplate, *m)
		if err != nil {
			return errors.Wrap(err, "failed to generate model")
		}
	}

	return nil
}
