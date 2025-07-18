package main

type locationAreaS struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []locationResult `json:"results"`
}

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
