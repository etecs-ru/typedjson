package main

import (
	"io"
	"strings"
	"text/template"
)

func generateCode(genArgs *GeneratorArgs, out io.Writer) error {
	tmpl := template.Must(template.New("gen").
		Funcs(template.FuncMap{
			"join":      strings.Join,
			"isPointer": isPointer,
			"trimStar":  trimStar,
		}).
		Parse(tmplRaw))
	err := tmpl.Execute(out, genArgs)
	if err != nil {
		return err
	}
	return nil
}

func isPointer(name string) bool {
	return len(name) > 0 && name[0] == '*'
}

func trimStar(name string) string {
	return strings.TrimLeft(name, "*")
}

var tmplRaw = `package {{.Package}}

// Code generated by "{{join .AllArgs " "}}"; DO NOT EDIT.

import (
	"encoding/json"
	"errors"

	{{ range .Imports -}} 
		"{{ . }}" 
	{{ end }}
) 

type {{.Typed}} struct {
	{{.Interface}}
} 

func (t {{.Typed}}) MarshalJSON() ([]byte, error) {
	if t.{{.Interface}} == nil {
		return nil, errors.New("nil interface in {{.Typed}}.{{.Interface}}")
	}
	typedString := t.{{.Interface}}.TypedJSON(nil)
	wrapper := struct {
		T string
		V {{.Interface}}
	}{
		T: typedString,
		V: t.{{.Interface}},
	}
	return json.Marshal(&wrapper)
} 

func (t *{{.Typed}}) UnmarshalJSON(src []byte) error {
	var wrapper struct {
		T string
		V json.RawMessage
	}
	err := json.Unmarshal(src, &wrapper)
	if err != nil {
		return err
	}
	data, err := GetEmpty{{.Interface}}(wrapper.T)
	if err != nil {
		return err
	}
	t.{{.Interface}} = data
	if err := json.Unmarshal(wrapper.V, t.{{.Interface}}); err != nil {
		return err
	}
	return nil
}

func GetEmpty{{.Interface}}(typedString string) ({{.Interface}}, error) {
	switch typedString {
	{{- range .Structs }}
	{{- if isPointer .Type -}}
	case "{{.Alias}}":
		return &{{trimStar .Type }}{}, nil
	{{else -}}
	case "{{.Alias}}":
		return {{trimStar .Type }}{}, nil
	{{- end }}{{ end }}
	default:
		return nil, errors.New("unknown type")
	}
}

{{ range .Structs }}
func (s {{.Type}}) TypedJSON(*{{$.Typed}}) string {
	return "{{.Alias}}"
}
{{ end }}
`
