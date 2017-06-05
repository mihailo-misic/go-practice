package main

import (
    "net/http"
    "log"
    "fmt"
    "strconv"
)

var counter int = 1

func main() {
    http.HandleFunc("/", set)
    http.HandleFunc("/forget", forget)
    http.Handle("/favicon.ico", http.NotFoundHandler())

    log.Fatal(http.ListenAndServe(":8080", nil))
}

// SETTING AND CHANGING THE COOKIE
func set(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    http.SetCookie(w, &http.Cookie{
        Name:  "count",
        Value: strconv.Itoa(counter),
    })
    counter ++

    c, err := r.Cookie("count")
    if err != nil {
        fmt.Fprintln(w, "You visited ", 0, " times.")
    } else {
        fmt.Fprintln(w, "You visited ", c.Value, " times.")
    }
    fmt.Fprintln(w, `<h1><a href="/forget">forget me</a></h1>`)
}

// DELETING COOKIE
func forget(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("count")
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    counter = 1
    c.MaxAge = -1
    http.SetCookie(w, c)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
