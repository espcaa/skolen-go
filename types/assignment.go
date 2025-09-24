package types

import "time"

type Assignment struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	HTML              string    `json:"html"`
	Done              bool      `json:"done"`
	DueDateTime       time.Time `json:"dueDateTime"`
	DeliverWorkOnline bool      `json:"deliverWorkOnline"`
	OnlineDeliverURL  string    `json:"onlineDeliverUrl"`
	Subject           Subject   `json:"subject"`
}
