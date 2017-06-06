package main

import (
    "net/http"
    "log"
    "github.com/julienschmidt/httprouter"
    "html/template"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
    "time"
    "net/url"
    "strings"
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
    router.POST("/login-submit", ProcessLogin)
    router.GET("/register", Register)
    router.POST("/register-submit", ProcessRegister)
    router.GET("/logout", Logout)
    router.GET("/me", Me)
    router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

    log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if !isLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    old_data := make(map[string]string, 2)
    // Get old data from cookie
    cookie, err := r.Cookie("old_data")
    if err == nil {
        parsed, err := url.ParseQuery(cookie.Value)
        check(err, w)
        for k, v := range parsed {
            old_data[k] = strings.Join(v, "")
        }
    }

    tpl.ExecuteTemplate(w, "login.gohtml", old_data)
}

func ProcessLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // If already logged in just redirect to home
    if isLoggedIn(r) {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if r.Method == http.MethodPost {
        err := r.ParseForm()
        check(err, w)

        // Get email and password from form
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Check if the required fields are set
        if len(email) == 0 || len(password) == 0 {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
        }

        o := orm.NewOrm()

        u := User{
            Email:    email,
            Password: password,
        }

        // Get user from database
        found := o.Read(&u, "Email", "Password")
        if found != nil {
            expires := time.Now().Add(time.Second * 3)
            cookie := http.Cookie{Name: "old_data", Value: r.PostForm.Encode(), Expires: expires}
            http.SetCookie(w, &cookie)

            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        // Set user session
        session := http.Cookie{Name: "session", Value: u.Email}
        http.SetCookie(w, &session)

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func ProcessRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // Delete session
    c := &http.Cookie{
        Name:   "session",
        Value:  "",
        MaxAge: -1,
    }
    http.SetCookie(w, c)

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

func isLoggedIn(r *http.Request) bool {
    _, err := r.Cookie("session")
    if err != nil {
        return false
    }
    return true
}
