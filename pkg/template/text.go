package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type TextTemplate struct {
	Template string
}

func (t *TextTemplate) Render(data interface{}) (string, error) {
	tmpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(t.Template)

	if err != nil {
		return "", fmt.Errorf("TextTemplate parse failed: %s", err)
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, data)

	if err != nil {
		return "", fmt.Errorf("TextTemplate render failed: %s", err)
	}

	return rendered.String(), nil
}
