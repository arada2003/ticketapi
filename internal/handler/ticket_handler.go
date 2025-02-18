package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tickets/internal/models"
	"tickets/internal/service"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketServ service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketServ}
}

func (th *TicketHandler) CreateTicket(c *gin.Context) {
	var req_ticket struct {
		Title	   string `json:"title"`
		Description string `json:"description"`
		Contact    string `json:"contact"`
	}
	if err := c.ShouldBindJSON(&req_ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticket, err := th.ticketService.CreateTicket(req_ticket.Title, req_ticket.Description, req_ticket.Contact)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

func (th *TicketHandler) GetAllTicket(c *gin.Context) {
	params := models.TicketQueryParams {
		Status: c.Query("status"),
		Sort:	c.Query("sort"),
		Order:  c.Query("order"),
	}

	tickets, err := th.ticketService.GetAllTicket(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (th *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
        return
    }
	
	var req_status struct {
		Status string `json:"status"`
		Email  string `json:"email"`
		Description	string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req_status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	ticket, err := th.ticketService.UpdateTicketStatus(id, req_status.Status, req_status.Email, req_status.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (th *TicketHandler) GetTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := th.ticketService.GetTicketByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (th *TicketHandler) GetUsers(c *gin.Context) {
	users, err := th.ticketService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (th *TicketHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := th.ticketService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (th *TicketHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := th.ticketService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (th *TicketHandler) GetUserByFirstname(c *gin.Context) {
	firstname := c.Param("firstname")

	user, err := th.ticketService.GetUserByFirstname(firstname)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
