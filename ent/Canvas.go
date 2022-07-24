package ent

import "time"

type Canvas struct {
	ID           string    `json:"id"`
	Drawing      string    `json:"Drawing"`
	CreationDate time.Time `json:"CreationDate"`
}
