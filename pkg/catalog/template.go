package catalog

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
        <td>Service Name</td>
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
      <td>{{.ServiceName}}</td>
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
