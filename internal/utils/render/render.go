package render

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/antonmedv/expr"
)

// Render returns the rendered string of the supplied key with the supplied variables and any error encountered.
func Render(key string, aux map[string]any) (string, error) {
	vars := Variables
	for k, v := range aux {
		vars[k] = v
	}
	ct, err := expr.Compile(key, expr.Env(vars))
	if err != nil {
		return "", fmt.Errorf("rendering check: compiling `%s`: %v", key, err)
	}
	t, err := expr.Run(ct, vars)
	if err != nil {
		return "", fmt.Errorf("rendering check: executing `%s`: %v", key, err)
	}
	return fmt.Sprint(t), nil
}

// TextTemplate returns the corresponding rendered string from a supplied object and recipe template.
//
// It also returns any errors encountered.
func TextTemplate(object any, recipe string) (string, error) {
	buf := &bytes.Buffer{}
	template, err := text.New("text").Parse(recipe)
	if err != nil {
		return "", err
	}
	err = template.Execute(buf, object)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// HtmlTemplate returns the corresponding html string from a supplied object and recipe template.
//
// It also returns any errors encountered.
func HtmlTemplate(object any, recipe string) (any, error) {
	buf := &bytes.Buffer{}
	template, err := html.New("html").Parse(recipe)
	if err != nil {
		return "", err
	}
	err = template.Execute(buf, object)
	if err != nil {
		return "", err
	}
	return html.HTML(buf.String()), nil
}
