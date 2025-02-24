package internal

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"toYaml":   ToYaml,
		"fromYaml": FromYaml,
		"toJson":   ToJson,
		"fromJson": FromJson,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

// ExecuteTemplate runs a template
func ExecuteTemplate(text string, values interface{}) (result string, err error) {
	tmpl, err := template.New("tpl").Funcs(funcMap()).Parse(text)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, values)
	if err != nil {
		return
	}
	result = buf.String()
	return
}
