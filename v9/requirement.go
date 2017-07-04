package main

type Requirement struct {
    Id              int     `json:"id"`
    Name            int     `json:"name"`
    Description     string  `json:"description"`
    Lens            string  `json:"lens"`
    Category        string  `json:"category"`
}
