package main

import (
    "database/sql"
    "fmt"
    _"github.com/go-sql-driver/mysql"
    "net/http"
    "io"
)

var db *sql.DB
var err error

func main() {
    db, err = sql.Open("mysql", "root:secret@/practice?charset=utf8")
    check(err)
    defer db.Close()

    err = db.Ping()
    check(err)

    http.HandleFunc("/", index)
    http.Handle("/favicon.ico", http.NotFoundHandler())
    err := http.ListenAndServe(":8080", nil)
    check(err)
}

func index(w http.ResponseWriter, r *http.Request) {
    _, err = io.WriteString(w, "Successfully completed.")
    check(err)

    rows, err := db.Query(`SELECT firstname FROM users`)
    check(err)

    var s, name string

    for rows.Next() {
        s += "\n"
        err = rows.Scan(&name)
        check(err)
        s += name
    }
    io.WriteString(w, s)
}

func check(err error) {
    if err != nil {
        fmt.Println(err)
    }
}
