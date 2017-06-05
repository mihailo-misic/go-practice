package main

import (
    "net/http"
    "log"
    "github.com/satori/go.uuid"
    "fmt"
)

func main() {
    http.HandleFunc("/", foo)
    http.Handle("/favicon.ico", http.NotFoundHandler())

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func foo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    c, err := r.Cookie("session")
    if err != nil {
        id := uuid.NewV4()
        c = &http.Cookie{
            Name:     "session",
            Value:    id.String(),
            HttpOnly: true,
        }
    }

    fmt.Fprintf(w, `
    <h2>Name: %s</h2>
    <h2>Value: %s</h2>
    `,c.Name, c.Value)
}
