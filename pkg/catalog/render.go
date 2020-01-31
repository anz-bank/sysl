package catalog

import (
	"bytes"
	"encoding/json"
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

func renderJSON(examples interface{}) ([]byte, error) {
	j, err := json.Marshal(examples)
	if err != nil {
		return nil, err
	}
	return j, err
}
