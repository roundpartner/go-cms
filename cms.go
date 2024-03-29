package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/artyom/autoflags"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/pat"
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var ServiceName = "go-cms"

var config = struct {
	Port              int    `flag:"port,port number to listen on"`
	Conn              string `flag:"conn,connection string"`
	DocumentationPath string `flag:"path,path to documentation"`
	Debug             bool   `flag:"debug,notify if pulling from database"`
}{
	Port: 7335,
}

func main() {
	log.SetOutput(os.Stdout)
	setupFlags()

	address := fmt.Sprintf(":%d", config.Port)
	start(address)

}

func setupFlags() {
	autoflags.Define(&config)
	flag.Parse()
}

func start(address string) {
	r := pat.New()
	r.Get("/check", check)
	r.Get("/{slug:[a-z_]+}/{id:[a-z_-]+}", getPageWithSlugHandler)
	r.Get("/{id:[a-z_]+}", getPageHandler)
	http.Handle("/", r)
	server := &http.Server{Addr: address}
	ShutdownGracefully(server)
	log.Printf("[INFO] [%s] Server starting on address %s", ServiceName, address)
	err := server.ListenAndServe()
	if nil != err {
		log.Printf("[INFO] [%s] %s", ServiceName, err.Error())
	}
}

func check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Use of deprectated endpoint detected: %s", r.URL.Query().Get(":id"))
	getPage(w, r, "page")
}

func getPageWithSlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	getPage(w, r, slug)
}

func getPage(w http.ResponseWriter, r *http.Request, slug string) {
	pageId := r.URL.Query().Get(":id")

	stat, err := os.Stat(config.DocumentationPath + "/" + slug + "/" + pageId + ".md")
	if err == nil {

		dat, err := ioutil.ReadFile(config.DocumentationPath + "/" + slug + "/" + pageId + ".md")
		if err == nil {

			p := strings.SplitAfterN(string(dat), "\n", 2)
			page := new(Page)
			page.PageId = pageId
			page.Title = strings.Trim(p[0], "\n# ")
			page.Content = strings.Trim(p[1], "\n")
			page.Content = string(github_flavored_markdown.Markdown([]byte(page.Content)))
			page.ContentType = "html"
			page.Modified = stat.ModTime().Format("2006-01-02 15:04:05")
			page.Created = stat.ModTime().Format("2006-01-02 15:04:05")

			js, err := json.Marshal(page)
			if err != nil {
				log.Fatalf("[ERROR] Error marshalling page: %s %s\n", pageId, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("[INFO] Returning page: %s/%s.md", slug, pageId)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
		log.Printf("[INFO] Unable to read: %s", err.Error())
	}
	log.Printf("[INFO] Not found: %s", config.DocumentationPath+"/"+slug+"/"+pageId+".md")

	if "page" != slug {
		log.Printf("[INFO] Page not found: %s\n", pageId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page := new(Page)
	page, err = getPageFromDatabase(pageId)
	if err != nil {
		log.Printf("[ERROR] Error getting page: %s %s\n", pageId, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page == nil {
		log.Printf("[INFO] Page not found: %s\n", pageId)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("[INFO] Found page in database: %s\n", pageId)

	if config.Debug {
		page.Content = "___Page was pulled from the database___\r\n" + page.Content
	}

	js, err := json.Marshal(page)
	if err != nil {
		log.Printf("[ERROR] Error marshalling page: %s %s\n", pageId, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getPageFromDatabase(pageId string) (*Page, error) {
	page := new(Page)

	db, err := sql.Open("mysql", config.Conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select page_id, title, content, content_type, modified, created from page where page_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(pageId).Scan(&page.PageId, &page.Title, &page.Content, &page.ContentType, &page.Modified, &page.Created)
	if err != nil {
		return nil, nil
	}

	return page, nil
}

type Page struct {
	PageId      string `json:"page_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Created     string `json:"created"`
	Modified    string `json:"modified"`
}
