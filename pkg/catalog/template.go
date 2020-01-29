package catalog

const catalogTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <title>Example</title>
  </head>
  <body>
    <table>
      	<tr>
          <td>Service Name</td>
          <td>Team</td>
          <td>Owner</td>
          <td>Email</td>
          <td>Filename</td>
          <td>URL</td>
          <td>Type</td>
          <td>Link</td>
    	</tr>
      	{{ with . }}
			{{ range . }}
      			<tr>
                  <td>{{.ServiceName}}</td>
                  <td>{{.Team}}</td>
                  <td>{{.Owner}}</td>
                  <td>{{.Email}}</td>
              		<td>{{.Filename }}</td>
                  <td>{{.URL}}</td>
                  <td>{{.Type}}</td>
                  <td><a href=".{{.Link}}/">here</a></td>
      			</tr>
			{{ end }} 
      	{{ end }}
    </table>
  </body>
</html>`
