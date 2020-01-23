package catalog

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
)

func renderHTML(examples interface{}, wr io.Writer) error {
	exampleTmpl := template.New("example")

	_, err := exampleTmpl.Parse(mustReadFile("pkg/catalog/catalog.html"))
	if err != nil {
		panic(err)
	}
	return exampleTmpl.Execute(wr, examples)
}

func renderJSON(examples interface{}, wr io.Writer) error {

	j, err := json.Marshal(examples)
	if err != nil {
		return err
	}
	_, err = wr.Write(j)

	return err
}

func mustReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
