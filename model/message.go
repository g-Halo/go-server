package model

import "time"

type Message struct {
	Sender 			string 		`json:"sender"`
	Recipient	 	string 		`json:"recipient"`
	Body 			string		`json:"body"`
	CreatedAt 		time.Time 	`json:"created_at"`
	Status			string		`json:"status"`
}