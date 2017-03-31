package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"log"
	"fmt"
	"net/http"
	"html/template"
	"regexp"
)

// Global constant
var DATABASE = "gsd.db"
var db *sql.DB 

type Page struct {
	Title string
	Body []string 
}

var templates = template.Must(template.ParseFiles("view.html", "insert.html"))
var validPath = regexp.MustCompile("^/(view|insert)/([a-zA-Z0-9])")

type taskItem struct {
	id string
	description string
	category string
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3",DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		fmt.Println("DB is nil")
	}
}

func createDB(db *sql.DB) {
	sqlStmt := "CREATE TABLE Task (id integer, description text, category integer);"
	if db == nil {
		fmt.Println("createDB : DB is nil")
		log.Fatal()
	}
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}
	return
}

func addTask(items []taskItem) {
	stmt, err := db.Prepare("INSERT INTO Task (id, description, category) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.id, item.description, item.category)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

func readTask () []taskItem {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, description, category FROM Task" )
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []taskItem
	for rows.Next() {
		item := taskItem{}
		err2 := rows.Scan(&item.id, &item.description, &item.category)
		if err2 != nil {
			log.Fatal(err2)
		}
		items = append(items, item)
	}
	return items
}

func makeHandler( fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println("m:"+m[2])
		if m == nil {
			http.NotFound(w,r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return 
	}
	renderTemplate(w, "view", p)
}

func insertHandler(w http.ResponseWriter, r *http.Request, title string) {
	query := r.URL.Query()
	fmt.Println(query["id"])
	item := []taskItem{taskItem{"", "desc3", "4"}}
	addTask(item)
	p, err := loadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "insert", p)
}

func loadPage(title string) (*Page, error) {
	taskItems := readTask();
	var body []string
	for _, item := range taskItems {
		body = append(body,item.id + " "+ item.description +" "+ item.category +"\n")
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	os.Remove(DATABASE)
	testItems := []taskItem{
		taskItem{"0", "desc0", "1"},
		taskItem{"1", "desc1", "2"},
	}
	initDB()
	if db == nil {
		fmt.Println("main: DB is nil")
	}
	createDB(db)

	addTask(testItems)
	taskItems := readTask();
	for _, item := range taskItems {
		fmt.Println(item.id, item.description, item.category)
	}
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/insert/", makeHandler(insertHandler))
	http.ListenAndServe(":8080", nil)
}
