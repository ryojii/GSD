package main

import "time"

type Exec struct {
	IdExec       int       `json:"idexec"`
	IdCampaign   int       `json:"idcampaign"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Trace        string    `json:"trace"`
	ForcedStatus string    `json:"fstatus"`
	StartDate    time.Time `json:"start"`
	EndDate      time.Time `json:"end"`
}

type Execs []Exec
