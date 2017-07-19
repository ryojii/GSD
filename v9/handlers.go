package main

import (
	"encoding/json"
	"fmt"
	"github.com/yosssi/ace"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func ExecIndex(w http.ResponseWriter, r *http.Request) {
	template, err := ace.Load("execs", "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = template.Execute(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ExecsSearch(w http.ResponseWriter, r *http.Request) {
	reviewer :=	r.FormValue("reviewer")
	campaign := r.FormValue("campaign")
	status := r.FormValue("status")
	name := r.FormValue("name")
	fmt.Println("DEBUG: will look for '"+reviewer+"' '"+campaign+"' "+status)
	var execs Execs
	var searchReviewer string
	var searchStatus string
	var searchCampaign string
	var searchName string
	search := make([]string, 0)
	if campaign != "" {
		searchCampaign= "idcampaign LIKE '%"+ campaign +"%'"
		search = append(search, searchCampaign)
	}
	if reviewer != "" {
		searchReviewer = "reviewer LIKE '%"+ reviewer +"%'"
		search = append(search, searchReviewer)
	}
	if name != "" {
		searchName = "name LIKE '%" + name + "%'"
		search = append(search, searchName)
	}
	if status != "" {
		searchStatus = "status = \"" + status+ "\""
		search = append(search, searchStatus)
	}
	if len(search) > 0 {
		fmt.Println("slice length:" + strconv.Itoa(len(search)))
		var searchString string
		if len(search) > 1 {
			for i := 0 ; i < len(search) ; i++ {
				fmt.Println("current: "+ search[i])
				if i == 0 {
					searchString = search[i]
				} else {
					searchString += " AND "+search[i]
				}
			}
		} else {
			searchString = search[0]
		}
		execs = findExecs(searchString)
		if len(execs) > 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(execs); err != nil {
				panic(err)
			}
			return
		}
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

//view one execution
func ExecShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	path := strings.Split(r.URL.Path, "/")
	var id string = path[len(path)-1]
	if exec, err := readExec(id); err == nil {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(exec); err != nil {
			panic(err)
		}
	} else {
		// If we didn't find it, 404
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	}
	return
}

//view all executions
func ExecsShow(w http.ResponseWriter, r *http.Request) {
	execs := readExecs()
	if len(execs) > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(execs); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"idcampaign":"4.2.4","name":"New Execution"}' http://localhost:8080/execs

{"idexec":0,"idcampaign":"4.2.4","name":"New Execution","status":"","trace":"","fstatus":"","start":"0001-01-01T00:00:00Z","end":"0001-01-01T00:00:00Z"}


*/
func ExecCreate(w http.ResponseWriter, r *http.Request) {
	var exec Exec
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &exec); err != nil {
		fmt.Println("unable to unmarshall json : " + err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	t := addExec(exec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func ExecDel(w http.ResponseWriter, r *http.Request) {
	deleteId(r.FormValue("id"))
	w.WriteHeader(http.StatusOK)
}

// URL : 	"/update/{id}/reviewer/{name}",
func ExecUpdateReviewer(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(path[len(path)-3])
	name := path[len(path)-1]
	if err := updateId(id, "reviewer", name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

// URL : 	"/update/{id}/status/{name}",
func ExecUpdateStatus(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(path[len(path)-3])
	name := path[len(path)-1]
	updateId(id, "forcedStatus", name)
}
