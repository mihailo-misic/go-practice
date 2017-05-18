package main

import (
    "log"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "html/template"
    "strconv"
)

type Task struct {
    Completed bool
    Name      string
    Deleted   bool
}

type Tasks []Task

var tasks Tasks = make([]Task, 0)

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/new-task", NewTask)
    router.GET("/toggle-task/:id", CompleteTask)
    router.GET("/delete-task/:id", DeleteTask)
    router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

    log.Fatal(http.ListenAndServe("localhost:8000", router))
}

// View all tasks
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tmp, err := template.ParseFiles("index.gohtml")
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "text/html")
    tmp.Execute(w, tasks)
}

// Handle adding new tasks
func NewTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if r.Method == "POST" {
        tasks = append(tasks, Task{false, r.FormValue("task"), false})
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
    http.Error(w, "WRONG METHOD \n Got: "+r.Method+"\n Expected: POST", 500)
}

// Handle completing a tasks
func CompleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    if r.Method == "GET" {
        id, err := strconv.Atoi(p.ByName("id"))
        if err != nil {
            log.Fatal(err)
        }
        tasks[id].Completed = !tasks[id].Completed

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
    http.Error(w, "WRONG METHOD \nGot: "+r.Method+"\nExpected: GET", 500)
}

// Handle removing a tasks
func DeleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    if r.Method == "GET" {
        id, err := strconv.Atoi(p.ByName("id"))
        if err != nil {
            log.Fatal(err)
        }
        tasks[id].Deleted = true

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
    http.Error(w, "WRONG METHOD \n Got: "+r.Method+"\n Expected: GET", 500)
}
