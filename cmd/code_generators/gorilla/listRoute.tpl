{{.listRoute}}

    uri := "/api/"+{{.APIVersion}} + "/" +
           "namespace"+{{.namespace}} +"/"+
           {{.endPointName}} +"/"+ "{{endPointName}}LIST"
    a.Router.HandleFunc(uri, a.{{.metohd}}{{.endPointName}}.Methods("GET")
    log.Println("{{.method}}: ", uri)

{{end}}
