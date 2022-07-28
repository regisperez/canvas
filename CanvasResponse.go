package main

type CanvasResponse struct {
	ID      string   `json:"id"`
	Drawing []string `json:"drawing"`
	CreationDate string `json:"creationDate"`
}

