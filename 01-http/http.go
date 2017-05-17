package main

import (
    "net/http"
    "html/template"
    "os"
    "io"
    "fmt"
    "strings"
)

type Person struct {
    FirstName string
    LastName  string
}

func main() {
    tmp, err := template.ParseFiles("tmp.gohtml")
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {

        if req.Method == "POST" {
            // Making the dir if it doesn't exist
            if _, err := os.Stat("./uploads"); err != nil {
                err := os.Mkdir("./uploads", 0755)
                if err != nil {
                    http.Error(resp, err.Error(), 500)
                    return
                } else {
                    fmt.Println("Made it!")
                }
            }

            // Fetch the filename from the form
            fileName := req.FormValue("file_name")

            // Fetch the file from the form
            file, fileInfo, err := req.FormFile("my_file")
            if err != nil {
                http.Error(resp, err.Error(), 500)
                return
            }
            defer file.Close()

            // Get extension
            file_name := fileInfo.Filename
            _split := strings.Split(file_name, ".")
            ext := _split[len(_split)-1]

            dst, err := os.Create("./uploads/" + fileName + "." + ext)
            if err != nil {
                http.Error(resp, err.Error(), 500)
                return
            }
            defer dst.Close()

            io.Copy(dst, file)
        }

        resp.Header().Set("Content-Type", "text/html")
        tmp.Execute(resp, nil)
    })

    http.ListenAndServe(":9000", nil)
}
