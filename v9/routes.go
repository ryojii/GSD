package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ExecIndex,
	},
	Route{
		"ExecIndex",
		"PUT",
		"/exec/{execId}",
		ExecIndex,
	},
	Route{
		"ExecCreate",
		"POST",
		"/exec",
		ExecCreate,
	},
	Route{
		"ExecShow",
		"GET",
		"/exec/{execId}",
		ExecShow,
	},
	Route{
		"ExecsShow",
		"GET",
		"/execs",
		ExecsShow,
	},
	Route{
		"ExecsSearch",
		"GET",
		"/search/{searchMethod}/{search}",
		ExecsSearch,
	},
}
