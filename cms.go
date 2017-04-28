package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "github.com/artyom/autoflags"
    "flag"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/pat"
    "io/ioutil"
    "os"
    "strings"
)

var config = struct {
    Port    int `flag:"port,port number to listen on"`
    Conn    string `flag:"conn,connection string"`
    DocumentationPath string `flag:"path,path to documentation"`
}{

}

func main() {

    setupFlags()

    address := fmt.Sprintf(":%d", config.Port)
    start(address)

}

func setupFlags() () {
    autoflags.Define(&config)
    flag.Parse()
}

func start(address string) {
    r := pat.New()
    r.Get("/{id:[a-z_]+}", getPageHandler)
    http.Handle("/", r)
    error := http.ListenAndServe(address, nil)
    if (nil != error) {
        log.Fatal(error)
    }
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
    pageId := r.URL.Query().Get(":id")

    stat, err := os.Stat(config.DocumentationPath + "/page/" + pageId + ".md")
    if err == nil {

        dat, err := ioutil.ReadFile(config.DocumentationPath + "/page/" + pageId + ".md")
        if err == nil {

            p := strings.SplitAfterN(string(dat), "\n", 2)
            page := new(Page)
            page.PageId = pageId
            page.Title = strings.Trim(p[0], "\n# ")
            page.Content = strings.Trim(p[1], "\n")
            page.ContentType = "markdown"
            page.Modified = stat.ModTime().Format("2006-01-02 15:04:05")
            page.Created = stat.ModTime().Format("2006-01-02 15:04:05")

            js, err := json.Marshal(page)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type", "application/json")
            w.Write(js)
            return
        }
    }

    db, err := sql.Open("mysql", config.Conn)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    stmt, err := db.Prepare("select page_id, title, content, content_type, modified, created from page where page_id = ?")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    page := new(Page)

    err = stmt.QueryRow(pageId).Scan(&page.PageId, &page.Title, &page.Content, &page.ContentType, &page.Modified, &page.Created)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    js, err := json.Marshal(page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

type Page struct {
    PageId string `json:"page_id"`
    Title string `json:"title"`
    Content string `json:"content"`
    ContentType string `json:"content_type"`
    Created string `json:"created"`
    Modified string `json:"modified"`
}
