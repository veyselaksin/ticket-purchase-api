package ticket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"ticket-purchase/internal/dto"
	"ticket-purchase/internal/i18n"
	"ticket-purchase/internal/i18n/messages"
	"ticket-purchase/internal/services"
	"ticket-purchase/pkg/cresponse"
)

type Handler interface {
	CreateTicket(ctx *fiber.Ctx) error
	GetTicket(ctx *fiber.Ctx) error
	PurchaseTicket(ctx *fiber.Ctx) error
}

type handler struct {
	ticketService services.TicketService
}

func New(ticketService services.TicketService) Handler {
	return &handler{
		ticketService: ticketService,
	}
}

// TicketCreate godoc
// @Summary Create a new ticket
// @Description Create a new ticket
// @Tags Ticket
// @Accept application/json
// @Produce application/json
// @Param ticket body dto.TicketCreateRequest true "Ticket data"
// @Success 201 {object} dto.TicketResponse
// @Router /tickets [post]
func (h *handler) CreateTicket(ctx *fiber.Ctx) error {
	var request dto.TicketCreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}

	response, err := h.ticketService.Create(ctx.Context(), &request)
	if err != nil {
		var status int
		var message string
		if err.Error() == messages.ErrorTicketCreate {
			status = fiber.StatusInternalServerError
			message = i18n.CreateMsg(ctx, messages.ErrorTicketCreate)
		} else {
			status = fiber.StatusInternalServerError
			message = i18n.CreateMsg(ctx, messages.UnexpectedError)
		}

		log.Error("Error creating ticket: ", err)
		return cresponse.ErrorResponse(ctx, status, message)
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusCreated, response)
}

// TicketGet godoc
// @Summary Get ticket by ID
// @Description Get ticket by ID
// @Tags Ticket
// @Accept application/json
// @Produce application/json
// @Param id path string true "Ticket ID"
// @Success 200 {object} dto.TicketResponse
// @Router /tickets/{id} [get]
func (h *handler) GetTicket(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	response, err := h.ticketService.FindById(ctx.Context(), id)
	if err != nil {
		var status int
		var message string
		if err.Error() == messages.NotFound {
			status = fiber.StatusNotFound
			message = i18n.CreateMsg(ctx, messages.NotFound)
		} else {
			status = fiber.StatusInternalServerError
			message = i18n.CreateMsg(ctx, messages.UnexpectedError)
		}

		log.Error("Error getting ticket: ", err)
		return cresponse.ErrorResponse(ctx, status, message)
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}

// TicketPurchase godoc
// @Summary Purchase a ticket
// @Description Purchase a ticket
// @Tags Ticket
// @Accept application/json
// @Produce application/json
// @Param id path string true "Ticket ID"
// @Param purchase body dto.TicketPurchaseRequest true "Purchase data"
// @Success 200 {object} interface{}
// @Router /tickets/{id}/purchase [post]
func (h *handler) PurchaseTicket(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var request dto.TicketPurchaseRequest
	if err := ctx.BodyParser(&request); err != nil {
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}
	request.TicketId = id

	err := h.ticketService.TicketPurchase(ctx.Context(), &request)
	if err != nil {
		var status int
		var message string
		if err.Error() == messages.ErrorPurchase {
			status = fiber.StatusInternalServerError
			message = i18n.CreateMsg(ctx, messages.ErrorPurchase)
		} else if err.Error() == messages.NotFound {
			status = fiber.StatusNotFound
			message = i18n.CreateMsg(ctx, messages.NotFound)
		} else if err.Error() == messages.ErrorTicketAllocations {
			status = fiber.StatusBadRequest
			message = i18n.CreateMsg(ctx, messages.ErrorTicketAllocations)
		} else {
			status = fiber.StatusInternalServerError
			message = i18n.CreateMsg(ctx, messages.UnexpectedError)
		}

		log.Error("Error purchasing ticket: ", err)
		return cresponse.ErrorResponse(ctx, status, message)
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, nil)
}
