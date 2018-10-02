package main

import "html/template"

var header = template.Must(template.New("header").Parse(`<!DOCTYPE html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>{{.Title}}</title>
</head>
<body>
{{if not .Index}}<p><a href="/">Back to issue list</a></p>{{end}}`))

const footer = `
</body>
</html>`

var issueList = template.Must(template.New("issuelist").Parse(`<h1>{{len .}} issues</h1>
<table>
<tr style="text-align: left">
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .}}
<tr>
  <td><a href="issue/{{.Number}}">{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href="user/{{.User.ID}}">{{.User.Login}}</a></td>
  <td><a href="issue/{{.Number}}">{{.Title}}</a></td>
</tr>
{{end}}
</table>` + footer))

var issueTempl = template.Must(template.New("issue").Parse(`<h1>#{{.Number}}: {{.Title}}</h1>
<p>by <a href="/user/{{.User.ID}}">{{.User.Login}}</a> &ndash; {{.State}}</p>
{{if .Milestone -}}
<p>Milestone: <a href="/milestone/{{.Milestone.Number}}">{{.Milestone.Title}}</a></p>
{{- end}}
<p>Created at {{.CreatedAt}}</p>
<pre>{{.Body}}</pre>` + footer))

var userTempl = template.Must(template.New("user").Parse(`<h1><a href="{{.HTMLURL}}">{{.Login}}</a></h1>
<img src="{{.AvatarURL}}" alt="Avatar"/>` + footer))

var milestoneTempl = template.Must(template.New("milestone").Parse(`<h1>Milestone #{{.Number}}: {{.Title}}</h1>
<p><a href="{{.HTMLURL}}">Link</a> &ndash; {{.State}}</p>
<pre>{{.Description}}</pre>` + footer))
