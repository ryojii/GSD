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

// DB Management (crud)
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
	sqlStmt := "CREATE TABLE Execution (idcampaign VARCHAR(100), name VARCHAR(255), status VARCHAR(50), reviewer VARCHAR(70) DEFAULT 'null',  trace BLOB DEFAULT 'null', forcedStatus varchar(50) DEFAULT 'null', start DATETIME, end DATETIME);"
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

// Execution management

// Create
func addExec(exec Exec) error {
	fmt.Println("DEBUG: exec content:" + exec.IdCampaign + " " + exec.Name + " " + exec.Status + " " + exec.Trace)
	stmt, err := db.Prepare("INSERT INTO Execution (idcampaign, name, status, trace, start, end) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(exec.IdCampaign, exec.Name, exec.Status, exec.Trace, exec.StartDate, exec.EndDate)
	if err2 != nil {
		log.Fatal(err2)
	}
	return err2
}

// Read

func readExec(id string) (Exec, error) {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT rowid, idcampaign, name, status, forcedStatus, reviewer, trace, start, end FROM Execution WHERE rowid = " + id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exec Exec = Exec{}
	rows.Next()
	err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.ForcedStatus, &exec.Reviewer, &exec.Trace, &exec.StartDate, &exec.EndDate)
	if err2 != nil {
		log.Fatal(err2)
	}
	return exec, err2
}

func readExecs() Execs {
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	fmt.Println("DEBUG: SELECT rowid, idcampaign, name, status, reviewer, trace, start, end FROM Execution")
	rows, err := db.Query("SELECT rowid, idcampaign, name, status, forcedStatus, reviewer, trace, start, end FROM Execution")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var execs Execs
	for rows.Next() {
		exec := Exec{}
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.ForcedStatus, &exec.Reviewer, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
		execs = append(execs, exec)
	}
	return execs
}

// Delete

func deleteId(id string) {
	stmt, err := db.Prepare("DELETE FROM Execution WHERE rowid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}

// Update
func updateId(id int, field string, value string) error {
	fmt.Println("UPDATE Execution SET " + field + " = " + value + " WHERE rowid= " + strconv.Itoa(id))
	stmt, err := db.Prepare("UPDATE Execution SET " + field + " = ? WHERE rowid= ?")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, id)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Search

func FindExecById(id int) Exec {
	return findExec("rowid = " + strconv.Itoa(id))
}

func findExec(search string) Exec {
	fmt.Println("DEBUG: search...")
	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	rows, err := db.Query("SELECT rowid, idcampaign, name, status, forcedStatus,, reviewer, trace, start, end FROM Execution WHERE " + search)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exec Exec
	for rows.Next() {
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.ForcedStatus, &exec.Reviewer, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return exec
}

// also add search by date

func FindExecsByDate(search string) Execs {
	return findExecs("date contain(\"" + search + "\")")
}

func FindExecsByTrace(search string) Execs {
	return findExecs("trace contain(\"" + search + "\")")
}

func findExecs(search string) Execs {

	if db == nil {
		fmt.Println("pointeur de DB null")
		log.Fatal()
	}
	fmt.Println("SELECT rowid, idcampaign, name, status, forcedStatus, reviewer, trace, start, end FROM Execution WHERE " + search)
	rows, err := db.Query("SELECT rowid, idcampaign, name, status, forcedStatus, reviewer, trace, start, end FROM Execution WHERE " + search)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var execs Execs
	for rows.Next() {
		exec := Exec{}
		err2 := rows.Scan(&exec.IdExec, &exec.IdCampaign, &exec.Name, &exec.Status, &exec.ForcedStatus, &exec.Reviewer, &exec.Trace, &exec.StartDate, &exec.EndDate)
		if err2 != nil {
			log.Fatal(err2)
		}
		execs = append(execs, exec)
	}
	for _, execution := range execs{
		fmt.Println("find: " + execution.Name)
	}
	return execs
}
