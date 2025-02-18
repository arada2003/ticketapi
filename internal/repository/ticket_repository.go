package repository

import (
	"database/sql"
	"errors"
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
	UpdateStatus(id int, status string, email string, description string) (*models.Ticket, error)
	GetTicketByID(id int) (*models.Ticket, error)
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByFirstname(firstname string) (*models.User, error)
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
		"INSERT INTO tickets (title, description, contact) VALUES ($1, $2, $3) RETURNING id, title, description, contact, status, last_updated_by, created_at, updated_at",
		title, description, contact,
	).Scan(&tk.ID, &tk.Title, &tk.Description, &tk.Contact, &tk.Status, &tk.LastUpdatedBy, &tk.CreatedAt, &tk.UpdatedAt)

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

	sortField := map[string]string{
		"status":     "status",
		"updated_at": "updated_at",
	}

	orderDirection := map[string]string{
		"asc":  "ASC",
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
		if err := rows.Scan(&tk.ID, &tk.Title, &tk.Description, &tk.Contact, &tk.Status, &tk.LastUpdatedBy, &tk.CreatedAt, &tk.UpdatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, tk)
	}
	return tickets, nil
}

func (r *ticketRepository) GetTicketByID(id int) (*models.Ticket, error) {
	var t models.Ticket

	err := r.db.QueryRow(
		"SELECT id, title, description, contact, status, last_updated_by, created_at, updated_at FROM tickets WHERE id = $1",
		id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Contact, &t.Status, &t.LastUpdatedBy, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &models.Ticket{}, fmt.Errorf("ticket not found")
		}
		return &models.Ticket{}, fmt.Errorf("failed to get ticket: %v", err)
	}

	return &t, nil
}

func (r *ticketRepository) UpdateStatus(id int, status string, email string, description string) (*models.Ticket, error) {
    var role string
    err := r.db.QueryRow(
        "SELECT role FROM users WHERE email = $1",
        email,
    ).Scan(&role)

    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    } else if err != nil {
        return nil, err
    }

    if role != "admin" {
        return nil, errors.New("update not allowed: user is not an admin")
    }

    var t models.Ticket
    err = r.db.QueryRow(
        "UPDATE tickets SET status = $1,last_updated_by = $2, description = $3, updated_at=now() WHERE id = $4 RETURNING id, title, description, contact, status, last_updated_by, created_at, updated_at",
        status, email, description, id,
    ).Scan(&t.ID, &t.Title, &t.Description, &t.Contact, &t.Status, &t.LastUpdatedBy, &t.CreatedAt, &t.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, errors.New("ticket not found")
    } else if err != nil {
        return nil, err
    }
    return &t, nil
}

func (r *ticketRepository) GetUsers() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, firstname, lastname, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *ticketRepository) GetUserByID(id int) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, firstname, lastname, email, role FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return &models.User{}, fmt.Errorf("user not found")
		}
		return &models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return &u, nil
}

func (r *ticketRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, firstname, lastname, email, role FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return &models.User{}, fmt.Errorf("user not found")
		}
		return &models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return &u, nil
}

func (r *ticketRepository) GetUserByFirstname(firstname string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, firstname, lastname, email, role FROM users WHERE firstname = $1",
		firstname,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return &models.User{}, fmt.Errorf("user not found")
		}
		return &models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return &u, nil
}
