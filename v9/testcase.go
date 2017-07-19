package main

type TestCase struct {
	TIdTC       int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Lens        string `json:"lens"`
	Author      string `json:"Author"`
}
type TestCases []TestCase
