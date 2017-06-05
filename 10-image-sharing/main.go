package main

import (
    "net/http"
    "log"
    "github.com/julienschmidt/httprouter"
    "html/template"
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
)

var tpl *template.Template

type User struct {
    Id        int `orm:"auto"`
    Email     string `orm:"size(100)"`
    Password  string `orm:"size(100)"`
    FirstName string `orm:"size(100)"`
    LastName  string `orm:"size(100)"`
}

func (u *User) TableName() string {
    return "users"
}

func init() {
    tpl = template.Must(template.ParseGlob("./templates/*"))

    orm.RegisterModel(new(User))
    orm.RegisterDataBase("default", "mysql", "root:secret@/practice?charset=utf8", 30)
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/login", Login)
    router.POST("/login", ProcessLogin)
    router.GET("/register", Register)
    router.GET("/logout", Logout)
    router.GET("/me", Me)
    router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

    log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func ProcessLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // TODO Check if already logged in
    if r.Method == http.MethodPost {
        err := r.ParseForm()
        check(err, w)

        email := r.FormValue("email")
        password := r.FormValue("password")

        if len(email) == 0 || len(password) == 0 {
            http.Redirect(w, r, "/login", http.StatusUnprocessableEntity)
        }
        fmt.Fprintln(w, email, password)

        //o := orm.NewOrm()
        //
        //u := User{
        //    Email:    string(form["email"]),
        //    Password: string(form["password"]),
        //}
        //
        //o.Read(&u)
        //fmt.Fprintln(w, u)
    }
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // TODO Erase state logic
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Me(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tpl.ExecuteTemplate(w, "me.gohtml", nil)
}

// Check for errors and if they exist handle them with http
func check(err error, w http.ResponseWriter) {
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
