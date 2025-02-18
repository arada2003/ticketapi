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
	UpdateTicketStatus(id int, status string, email string) (*models.Ticket, error)
	GetTicketByID(id int) (*models.Ticket, error)
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByFirstname(firstname string) (*models.User, error)
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

func (ts *ticketService) UpdateTicketStatus(id int, status string, email string) (*models.Ticket, error) {
	ticket, err := ts.repo.UpdateStatus(id, status, email)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, errors.New("ticket not found")
	}
	return ticket, err
}

func (ts *ticketService) GetTicketByID(id int) (*models.Ticket, error) {
	ticket, err := ts.repo.GetTicketByID(id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, errors.New("ticket not found")
	}
	return ticket, err
}

func (s *ticketService) GetUsers() ([]models.User, error) {
	return s.repo.GetUsers()
}

func (s *ticketService) GetUserByID(id int) (*models.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (s *ticketService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (s *ticketService) GetUserByFirstname(firstname string) (*models.User, error) {
	user, err := s.repo.GetUserByFirstname(firstname)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, errors.New("user not found")
	}
	return user, err
}
