package handler

import (
	"net/http"

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
