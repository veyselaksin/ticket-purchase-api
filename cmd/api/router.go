package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	"ticket-purchase/cmd/api/handlers/v1/ticket"
	"ticket-purchase/internal/db/repositories"
	"ticket-purchase/internal/services"
)

// HealthCheck godoc
// @Summary Health Check API
// @Description Health Check for the API
// @Tags Health Check
// @Accept application/json
// @Produce application/json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func health(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

func InitializeRouters(app *fiber.App, connection *gorm.DB) {

	// Repositories
	ticketRepository := repositories.NewTicketRepository(connection)
	purchaseRepository := repositories.NewPurchaseRepository(connection)

	// Services
	ticketService := services.NewTicketService(ticketRepository, purchaseRepository)

	// Handlers
	ticketHandler := ticket.New(ticketService)

	// Initialize the routes for the application here
	v1 := app.Group("/v1")

	// Swagger documentation
	v1.Get("/docs/*", swagger.HandlerDefault)

	// Health check
	v1.Get("/health", health)

	// Initialize the routes for the application here
	ticketRouter := v1.Group("/tickets")
	ticketRouter.Post("/", ticketHandler.CreateTicket)
	ticketRouter.Get("/:id", ticketHandler.GetTicket)
	ticketRouter.Post("/:id/purchase", ticketHandler.PurchaseTicket)
}
