package viewmodel

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type VM interface {
	FS() fs.FS
	Data() VM
}

func New[T VM](name string, fs fs.FS, data T) Root {
	return &raw{inner: data, fs: fs, name: name}
}

// If the viewmodel has NO values it is basic
type baseModel struct {
	fs     fs.FS
	Values any
}

func Basic(fs fs.FS, values any) *baseModel { return &baseModel{fs: fs, Values: values} }
func (vm *baseModel) Data() VM              { return nil }
func (vm *baseModel) FS() fs.FS             { return vm.fs }

type raw struct {
	name  string
	fs    fs.FS
	inner VM
}

type Root interface{ Execute(w http.ResponseWriter) }

func (raw *raw) Execute(w http.ResponseWriter) {
	fs, paths := Merge(allFSs(raw.inner))
	templ, err := template.
		New(raw.name).
		Funcs(template.FuncMap{
			"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		}).
		ParseFS(fs, paths...)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := templ.Execute(&buf, &raw.inner); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := minifyHTML(w, &buf); err != nil {
		log.Println("HTML minify error:", err)
		w.Write(buf.Bytes()) // fallback
	}
}
