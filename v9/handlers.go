package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"html/template"
)

var templates = template.Must(template.ParseFiles("exec.html", "execs.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func ExecIndex(w http.ResponseWriter, r *http.Request) {
	execs := readExecs()
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
		renderTemplateExecs(w, "execs", &execs )
}

func ExecsSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var searchMethod string
	var search string
	searchMethod = vars["method"]
	search = vars["search"]
	//la methode find a toujours la même signature, je dois pouvoir en faire un pointeur
	execs := searchMethod(search)
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

func ExecShow(w http.ResponseWriter, r *http.Request) {
	execs := readExecs()
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		renderTemplateExecs(w, "execs", &execs )
		return
}

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

curl -H "Content-Type: application/json" -d '{"name":"New Execution"}' http://localhost:8080/execs

*/
func ExecCreate(w http.ResponseWriter, r *http.Request) {
	var exec Exec
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &exec); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := addExec(exec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
 
func renderTemplate( w http.ResponseWriter, template string, exec *Exec) {
	err := templates.ExecuteTemplate(w, template+".html", exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateExecs( w http.ResponseWriter, template string, execs *Execs) {
	err := templates.ExecuteTemplate(w, template+".html", execs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
