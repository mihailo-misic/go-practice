package main

import (
    "html/template"
    "net/http"
    "log"
    "io/ioutil"
    "os"
    "path/filepath"
)

var tpl *template.Template

func init() {
    tpl = template.Must(template.ParseGlob("./templates/*"))
}

func main() {
    http.HandleFunc("/", foo)
    http.Handle("/favicon.ico", http.NotFoundHandler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func foo(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        file, h, err := r.FormFile("q")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        bs, err := ioutil.ReadAll(file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        dst, err := os.Create(filepath.Join("./uploads/", h.Filename))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        dst.Write(bs)
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
