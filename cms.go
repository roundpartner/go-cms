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
)

var config = struct {
    Port    int `flag:"port,port number to listen on"`
    Conn    string `flag:"conn,connection string"`
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
    http.HandleFunc("/", handler)
    error := http.ListenAndServe(address, nil)
    if (nil != error) {
        log.Fatal(error)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {

        pageId := r.URL.Query().Get("page")

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

        err = stmt.QueryRow(pageId).Scan(&page.Page_id, &page.Title, &page.Content, &page.ContentType, &page.Modified, &page.Created)
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
}

type Page struct {
    Page_id string
    Title string
    Content string
    ContentType string
    Created string
    Modified string
}
