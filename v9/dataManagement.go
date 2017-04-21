package main

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"os"
)

var currentId int

var execs Execs

var DATABASE = "exec.db"
var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3",DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(DATABASE); err != nil {
		createDB()
	}
}
func createDB() {
	sqlStmt := "CREATE TABLE Execution (id INTEGER, name VARCHAR(255), status VARCHAR(50), trace BLOB, forcedStatus varchar(50), start DATETIME, end DATETIME);"
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

func RepoFindExec(id int) Exec {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, name, status, trace, start, end FROM Execution WHERE id = " + strconv.Itoa(id) )
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exec Exec
	for rows.Next() {
		err2 := rows.Scan(&exec.IdExec, &exec.Name, &exec.Status, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return exec
}

func addExec(exec Exec) Exec{
	stmt, err := db.Prepare("INSERT INTO Execution (id, name, status, trace, start, end) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(exec.IdExec, exec.Name, exec.Status, exec.Trace,  exec.StartDate, exec.EndDate)
	if err2 != nil {
		log.Fatal(err2)
	}
	return exec
}

func readExecs () Execs {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, name, status, trace, start, end FROM Execution" )
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var execs Execs
	for rows.Next() {
		exec := Exec{}
		err2 := rows.Scan(&exec.IdExec, &exec.Name, &exec.Status, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
		execs = append(execs, exec)
	}
	return execs
}





//this is bad, I don't think it passes race condtions
func RepoCreateExec(t Exec) Exec {
	currentId += 1
	t.IdExec = currentId
	execs = append(execs, t)
	return t
}

func RepoDestroyExec(id int) error {
	for i, t := range execs {
		if t.IdExec == id {
			execs = append(execs[:i], execs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Exec with id of %d to delete", id)
}
