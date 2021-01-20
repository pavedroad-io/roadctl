{{.keyedRoute}}

    uri := "/api/"+{{.APIVersion}} + "/" +
           "namespace"+{{.namespace}} +"/"+
           {{.endPointName}} +"/"+ "{key}"
    a.Router.HandleFunc(uri, a.{{.metohd}}{{.endPointName}}.Methods("{{.method}}")
    log.Println("{{.method}}: ", uri)

{{end}}
