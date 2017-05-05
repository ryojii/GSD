package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

var currentId int

var execs Execs

var DATABASE = "exec.db"
var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(DATABASE); err != nil {
		createDB()
	}
}
func createDB() {
	sqlStmt := "CREATE TABLE Execution (id INTEGER, idcampaign VARCHAR(100), name VARCHAR(255), status VARCHAR(50), trace BLOB, forcedStatus varchar(50), start DATETIME, end DATETIME);"
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

func FindExecById(id int) Exec {
	return findExec("id = " + strconv.Itoa(id))
}

func FindExcByName(search string) Exec {
	return findExec("name = \"" + search + "\"")
}

func findExec(search string) Exec {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, idcampaign, name, status, trace, start, end FROM Execution WHERE " + search)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exec Exec
	for rows.Next() {
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return exec
}

// also add search by date

func FindExecsByDate(search string) Execs {
	return findExecs("date contain(\""+ search +"\")")
}

func FindExecsByTrace(search string) Execs {
	return findExecs("trace contain(\""+ search +"\")")
}

func FindExecsByMatchingName(search string) Execs {
	//I should find another to do this, it's a sort of ... ugly
	return findExecs("name LIKE '%" + search +"%'" )
}

func FindExecsByStatus(status string) Execs {
	return findExecs("status = \"" + status + "\"")
}

func FindExecsByCampaign(campaign string) Execs {
	return findExecs("idcampaign = \"" + campaign + "\"")
}

func findExecs(search string) Execs {

	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, idcampaign, name, status, trace, start, end FROM Execution WHERE " + search)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var execs Execs
	for rows.Next() {
		exec := Exec{}
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
		execs = append(execs, exec)
	}
	return execs
}
func addExec(exec Exec) Exec {
	stmt, err := db.Prepare("INSERT INTO Execution (id, idcampaign, name, status, trace, start, end) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(exec.IdExec, exec.IdCampaign, exec.Name, exec.Status, exec.Trace, exec.StartDate, exec.EndDate)
	if err2 != nil {
		log.Fatal(err2)
	}
	return exec
}

func readExecs() Execs {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT id, idcampaign, name, status, trace, start, end FROM Execution")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var execs Execs
	for rows.Next() {
		exec := Exec{}
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
		execs = append(execs, exec)
	}
	return execs
}
