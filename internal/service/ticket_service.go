package service

import (
	"errors"
	"strings"

	"tickets/internal/models"
	"tickets/internal/repository"
)

type TicketService interface {
	CreateTicket(title, description, contact string) (*models.Ticket, error)
	GetAllTicket(params models.TicketQueryParams) ([]models.Ticket, error)
}

type ticketService struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &ticketService{repo: repo}
}

func (ts *ticketService) CreateTicket(title, description, contact string) (*models.Ticket, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(description) == "" || strings.TrimSpace(contact) == "" {
		return nil, errors.New("title, description, and contact are required")
	}
	ticket, err := ts.repo.CreateTicket(title, description, contact)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, errors.New("ticket already exists")
	}
	return ticket, err
}

func (ts *ticketService) GetAllTicket(params models.TicketQueryParams) ([]models.Ticket, error) {
	return ts.repo.GetAllTicket(params)
}

