package catalog

import (
	"bytes"
	"html/template"
)

const catalogTemplate = `
{{$header := true}}
<!DOCTYPE html>
<html>
  <head>
    <title>Example</title>
  </head>
  <body>
    <table border = "1">
        <tr>
        <td>App Name</td>
        {{ with . }}
			  {{ range . }}
        {{if $header}}
        {{range .Fields }}
        <td>{{.}}</td>
        {{end}}
        {{$header = false}}
        <td>Link</td>
        {{end}}
        
    	</tr>
      <tr>
      <td>{{.AppName}}</td>
      {{$Attrs := .Attrs}}
      {{range .Fields }}
      <td>{{index $Attrs .}}</td>
      {{end}}
      <td> <a href="./{{.SwaggerUILink}}">{{.SwaggerUILink}}</a></td>
      </tr>
			{{ end }} 
      {{ end }}
    </table>
  </body>
</html>`

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
