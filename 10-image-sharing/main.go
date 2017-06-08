package main

import (
    "net/http"
    "log"
    "github.com/julienschmidt/httprouter"
    "html/template"
    _ "github.com/go-sql-driver/mysql"
    "time"
    "net/url"
    "strings"
    "golang.org/x/crypto/bcrypt"
    "errors"
    "github.com/jinzhu/gorm"
    "fmt"
    "crypto/sha1"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "io/ioutil"
)

var tpl *template.Template

type User struct {
    ID        uint
    Email     string `gorm:"size:100"`
    Password  string `gorm:"size:255"`
    FirstName string `gorm:"size:100"`
    LastName  string `gorm:"size:100"`
}

type Data struct {
    User   User
    Images []string
}

var db *gorm.DB

func init() {
    tpl = template.Must(template.ParseGlob("./templates/*"))
}

func main() {
    var err error
    db, err = gorm.Open("mysql", "root:secret@/practice?charset=utf8&parseTime=True")
    if err != nil {
        log.Fatal("Could not connect to database")
    }
    defer db.Close()

    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/login", Login)
    router.POST("/login-submit", ProcessLogin)
    router.GET("/register", Register)
    router.POST("/register-submit", ProcessRegister)
    router.GET("/logout", Logout)
    router.GET("/me", Me)
    router.POST("/me-submit", ProcessMe)
    router.GET("/images", Images)
    router.POST("/images-upload", ProcessUpload)
    router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
    router.ServeFiles("/storage/*filepath", http.Dir("storage/"))

    log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    u, err := getUser(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        Logout(w, r, nil)
        return
    }

    var data Data

    data.User = u

    // Get the images for the user
    img_dir := "/storage/uploads/" + strconv.Itoa(int(u.ID)) + "/"
    files, _ := ioutil.ReadDir("." + img_dir)
    // Fill the data.Images with the paths to images
    for _, f := range files {
        data.Images = append(data.Images, img_dir+f.Name())
    }

    tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // If already logged in just redirect to home
    if isLoggedIn(r) {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

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
        // Get email and password from form
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Check if the required fields are set
        if len(email) == 0 || len(password) == 0 {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
        }

        u := User{}

        q := db.First(&u, "email=?", email)
        if q.Error != nil {
            expires := time.Now().Add(time.Second * 3)
            cookie := http.Cookie{Name: "old_data", Value: r.PostForm.Encode(), Expires: expires}
            http.SetCookie(w, &cookie)

            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        // Check if passwords match
        err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
        if err != nil {
            fmt.Println(err)
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
    // Get registry data from form
    first_name := r.FormValue("first_name")
    last_name := r.FormValue("last_name")
    email := r.FormValue("email")
    password := r.FormValue("password")
    confirm_password := r.FormValue("confirm_password")

    // Check for required and if the password is confirmed
    if len(first_name) == 0 || len(last_name) == 0 || len(email) == 0 || len(password) == 0 || len(confirm_password) == 0 {
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    } else if password != confirm_password {
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Hashing the password
    pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    check(err, w)

    // Building the user model
    u := &User{
        FirstName: strings.Title(first_name),
        LastName:  strings.Title(last_name),
        Email:     email,
        Password:  string(pass),
    }

    exists := db.First(&u, "email=?", email)
    if exists.RowsAffected != 0 {
        http.Error(w, "Email is taken", http.StatusUnauthorized)
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Inserting the user into the DB
    q := db.Save(u)
    if q.Error != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
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
    u, err := getUser(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        Logout(w, r, nil)
        return
    }

    tpl.ExecuteTemplate(w, "me.gohtml", u)
}

func ProcessMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    r.ParseForm()
    form := r.PostForm

    // Check for empty fields
    for field := range form {
        value := form.Get(field)

        if len(value) == 0 {
            http.Error(w, field+" is required", http.StatusInternalServerError)
            return
        }
    }

    // Get User from database
    u, err := getUser(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if Old password is correct
    err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(form.Get("old_password")))
    if err != nil {
        http.Error(w, "The password is incorrect!", http.StatusInternalServerError)
        return
    }

    // Check if the password is different
    err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(form.Get("new_password")))
    if err == nil {
        http.Error(w, "New password is the same as the old one!", http.StatusInternalServerError)
        return
    }

    // Check if the email is taken by somebody else
    id := u.ID
    q := db.Find(&User{}, "id != ? AND email = ?", id, form.Get("email"))
    if q.RowsAffected > 0 {
        http.Error(w, "Email is taken!", http.StatusInternalServerError)
        return
    }

    password, err := bcrypt.GenerateFromPassword([]byte(form.Get("new_password")), bcrypt.DefaultCost)
    check(err, w)

    // Update the user
    u.FirstName = strings.Title(form.Get("first_name"))
    u.LastName = strings.Title(form.Get("last_name"))
    u.Email = form.Get("email")
    u.Password = string(password)

    db.Save(&u)

    // Update the session cookie
    cookie := http.Cookie{Name: "session", Value: u.Email}
    http.SetCookie(w, &cookie)

    http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func Images(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if !isLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    tpl.ExecuteTemplate(w, "images.gohtml", nil)
}

func ProcessUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    u, err := getUser(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        Logout(w, r, nil)
        return
    }

    mf, fh, err := r.FormFile("image")
    check(err, w)
    defer mf.Close()

    ext := strings.Split(fh.Filename, ".")[1]
    h := sha1.New()
    io.Copy(h, mf)
    fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
    // create new file
    wd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }
    // Create the dir if it does not exist
    dirpath := filepath.Join(wd, "storage", "uploads", strconv.Itoa(int(u.ID)))
    os.MkdirAll(dirpath, 0777)
    path := filepath.Join(dirpath, fname)
    nf, err := os.Create(path)
    if err != nil {
        fmt.Println(err)
    }
    defer nf.Close()
    // copy
    mf.Seek(0, 0)
    io.Copy(nf, mf)

    http.Redirect(w, r, "/", http.StatusSeeOther)
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

func getUser(w http.ResponseWriter, r *http.Request) (User, error) {
    if !isLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return User{}, nil
    }

    c, err := r.Cookie("session")
    check(err, w)

    u := User{}

    q := db.Find(&u, "email=?", c.Value)
    if q.Error != nil {
        expires := time.Now().Add(time.Second * 3)
        cookie := http.Cookie{Name: "old_data", Value: r.PostForm.Encode(), Expires: expires}
        http.SetCookie(w, &cookie)

        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return User{}, errors.New("User not found")
    }

    return u, nil
}

func (u User) FullName() string {
    return strings.Title(u.FirstName + " " + u.LastName)
}
