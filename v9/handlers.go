package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
    "github.com/yosssi/ace"
    "fmt"
    "strconv"
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
	var searchMethod string
	var search string
	path := strings.Split(r.URL.Path, "/")
	searchMethod = path[len(path) -2]
	search = path[len(path)-1]
	var execs Execs
	switch searchMethod  {
		case "status" :
			execs = FindExecsByStatus(search)
		case "testName" :
			execs = FindExecsByMatchingName(search)
		case "containTestName" :
			execs = FindExecsBySimilarMatchingName(search)
		case "campaignName" :
			execs = FindExecsByCampaign(search)
		case "containCampaignName" :
			execs = FindExecsBySimilarCampaign(search)
		case "trace" :
			execs = FindExecsByTrace(search)
		case "date" :
			execs = FindExecsByDate(search)
	}
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
	id, _ := strconv.Atoi(path[len(path) -3])
	name := path[len(path)-1]
    updateId(id, "reviewer", name)
}

// URL : 	"/update/{id}/status/{name}",
func ExecUpdateStatus(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(path[len(path) -3])
	name := path[len(path)-1]
    updateId(id, "forcedStatus", name)
}
