package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
)

type Response struct{}

func (r *Response) html(w http.ResponseWriter, view string, data interface{}) {
	tpl := template.New("")
	funcs := make(map[string]interface{})
	funcs["ToString"] = ToString
	tpl.Funcs(funcs)

	base := path.Join("templates", "base.html")
	view = path.Join("templates", view)
	tpl, err := tpl.ParseFiles(base, view)
	if err != nil {
		panic(err)
		// r.html(w, path.Join(t.Path, "error.html"), "View file not found")
	}
	err = tpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		panic(err)
	}
}

func (r *Response) json(w http.ResponseWriter, out interface{}) {
	w.Header().Add("Content-Type", "application/json")
	res, _ := json.Marshal(out)
	w.Write(res)
}
