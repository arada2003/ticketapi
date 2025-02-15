package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"tickets/internal/config"
	"tickets/internal/handler"
	"tickets/internal/middleware"
	"tickets/internal/repository"
	"tickets/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	
	db, err := repository.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	ticketRepo := repository.NewTicketRepository(db)
	ticketServ := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketServ)

	r := gin.Default()

	// Health Check API
	r.GET("/health", func(c *gin.Context) {
		if err := repository.CheckDBConnection(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
	})

	authReq := r.Group("/v1", middleware.BearerAuth(cfg.APIToken))
	{
		authReq.POST("/tickets", ticketHandler.CreateTicket)
		authReq.GET("/tickets", ticketHandler.GetAllTicket)
	}
	
	r.Run(":80")
}