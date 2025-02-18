package models

import "time"

type Ticket struct {
	ID		  int       `json:"id"`
	Title    string    `json:"title"`
	Description string `json:"description"`
	Contact string    `json:"contact"`
	Status  string    `json:"status"`
	LastUpdatedBy string `json:"last_updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TicketQueryParams struct {
	Status string `form:"status"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

type User struct {
	ID int `json:"id"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Role string `json:"role"`
}