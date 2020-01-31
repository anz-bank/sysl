package catalog

import (
	"bytes"
	"html/template"
)

func renderHTML(examples interface{}) ([]byte, error) {
	exampleTmpl := template.New("example")
	var contents bytes.Buffer
	_, err := exampleTmpl.Parse(catalogTemplate)
	if err != nil {
		return nil, err
	}
	err = exampleTmpl.Execute(&contents, examples)
	return contents.Bytes(), err
}
