package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"log"
	"fmt"
)

// Global constant
var DATABASE = "gsd.db"

type taskItem struct {
	id string
	description string
	category string
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3",DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createDB(db *sql.DB) {
	sqlStmt := "CREATE TABLE Task (id integer primary key, description text, category integer);"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}
	return
}

func addTask(db *sql.DB, items []taskItem) {
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

func readTask (db *sql.DB) []taskItem {
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

func main() {
	os.Remove(DATABASE)
	testItems := []taskItem{
		taskItem{"0", "desc0", "1"},
		taskItem{"1", "desc1", "2"},
	}
	db := initDB()
	defer db.Close()
	createDB(db)

	addTask(db, testItems)
	taskItems := readTask(db);
	for _, item := range taskItems {
		fmt.Println(item.id, item.description, item.category)
	}
}
