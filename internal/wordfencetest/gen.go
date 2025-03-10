//go:build ignore

package main

import (
	"os"
	"text/template"
)

const fileTemplateRaw = `// Code generated by {{ .GeneratedBy }}; DO NOT EDIT.
//
// This file was generated using data from
{{- range .Sources }}
//   - {{ . }}
{{- end }}
package main

var {{ .VariableName }} = []string{
{{- range .Versions }}
	{{ printf "%q" . }},
{{- end }}
}
`

var fileTemplate = template.Must(template.New("").Parse(fileTemplateRaw))

type fileData struct {
	GeneratedBy  string
	Sources      []string
	VariableName string
	Versions     []string
}

func gen(out string, data fileData) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	return fileTemplate.Execute(f, data)
}
