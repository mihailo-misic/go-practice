package main

import (
    "net/http"
    "io"
    "html/template"
    "log"
)

func main() {
    http.HandleFunc("/", foo)
    http.HandleFunc("/dog/", dog)
    http.Handle("/favicon.ico", http.NotFoundHandler())
    http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
    http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    io.WriteString(w, "foo ran")
}

func dog(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    tmp, err := template.ParseFiles("dog.gohtml")
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "text/html")
    tmp.Execute(w, nil)
}
