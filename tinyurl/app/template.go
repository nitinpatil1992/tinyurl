package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type appTemplate struct {
	t *template.Template
}

func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles("./app/templates/base.html"))
	path := filepath.Join("./app/templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) *appError {
	d := struct {
		Data        interface{}
		FormEnabled bool
	}{
		Data:        data,
		FormEnabled: true,
	}
	if err := tmpl.t.Execute(w, d); err != nil {
		return appErrorf(err, "could not write template: %v", err)
	}
	return nil
}
