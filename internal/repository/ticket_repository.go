package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tickets/internal/config"
	"tickets/internal/models"

	_ "github.com/lib/pq"
)

type TicketRepository interface {
	CreateTicket(title, description, contact string) (*models.Ticket, error)
	GetAllTicket(params models.TicketQueryParams) ([]models.Ticket, error)
}

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func CheckDBConnection(db *sql.DB) error {
	return db.Ping()
}

func (r *ticketRepository) CreateTicket(title, description, contact string) (*models.Ticket, error) {
	var tk models.Ticket
	err := r.db.QueryRow(
		"INSERT INTO tickets (title, description, contact) VALUES ($1, $2, $3) RETURNING id, title, description, contact, status, created_at, updated_at", 
		title, description, contact,
	).Scan(&tk.ID, &tk.Title, &tk.Description, &tk.Contact, &tk.Status, &tk.CreatedAt, &tk.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &tk, nil
}

func (r *ticketRepository) GetAllTicket(params models.TicketQueryParams) ([]models.Ticket, error) {
	query := "SELECT * FROM tickets"

	args := []interface{}{}
	placeholderCount := 1

	if params.Status != "" {
		query += fmt.Sprintf(" WHERE status = $%d", placeholderCount)
		args = append(args, params.Status)
		placeholderCount++
	}

	sortField := map[string]string {
		"status": "status",
		"updated_at": "updated_at",
	}

	orderDirection := map[string]string {
		"asc": "ASC",
		"desc": "DESC",
	}

	sortColumn, ok := sortField[params.Sort]
	if !ok {
		sortColumn = "updated_at"
	}

	orderDir, ok := orderDirection[strings.ToLower(params.Order)]
	if !ok {
		orderDir = "ASC"
	}
	
	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, orderDir)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var tk models.Ticket
		if err := rows.Scan(&tk.ID, &tk.Title, &tk.Description, &tk.Contact, &tk.Status, &tk.CreatedAt, &tk.UpdatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, tk)
	}
	return tickets, nil
}