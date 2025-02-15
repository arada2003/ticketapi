package models

import "time"

type Ticket struct {
	ID		  int       `json:"id"`
	Title    string    `json:"title"`
	Description string `json:"description"`
	Contact string    `json:"contact"`
	Status  string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TicketQueryParams struct {
	Status string `form:"status"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}