package main

import "time"

type Exec struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Status	  string    `json:"status"`
	StartDate time.Time `json:"start"`
	EndDate   time.Time `json:"end"`
}

type Execs []Exec
