package services

import (
	"context"
	"errors"
	"ticket-purchase/internal/db/models"
	"ticket-purchase/internal/db/repositories"
	"ticket-purchase/internal/dto"
	"ticket-purchase/internal/i18n/messages"
	"time"
)

var timeNow = time.Now

type TicketService interface {
	// Create creates a new ticket
	Create(ctx context.Context, request *dto.TicketCreateRequest) (*dto.TicketResponse, error)
	FindById(ctx context.Context, id string) (*dto.TicketResponse, error)
	TicketPurchase(ctx context.Context, request *dto.TicketPurchaseRequest) error
}

type ticketService struct {
	ticketRepo   repositories.TicketRepository
	purchaseRepo repositories.PurchaseRepository
}

func NewTicketService(
	ticketRepo repositories.TicketRepository,
	purchaseRepo repositories.PurchaseRepository,
) TicketService {
	return &ticketService{
		ticketRepo:   ticketRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (s *ticketService) Create(ctx context.Context, request *dto.TicketCreateRequest) (*dto.TicketResponse, error) {
	ticket := models.Ticket{
		Name:        request.Name,
		Description: request.Description,
		Allocation:  request.Allocation,
	}

	data, err := s.ticketRepo.Create(ctx, &ticket)
	if err != nil {
		return nil, errors.New(messages.ErrorTicketCreate)
	}

	response := dto.TicketResponse{
		Id:          data.Id,
		Name:        data.Name,
		Description: data.Description,
		Allocation:  data.Allocation,
	}

	return &response, nil
}

func (s *ticketService) FindById(ctx context.Context, id string) (*dto.TicketResponse, error) {
	data, err := s.ticketRepo.FindById(ctx, id)
	if err != nil && err.Error() == "record not found" {
		return nil, errors.New(messages.NotFound)
	}

	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	response := dto.TicketResponse{
		Id:          data.Id,
		Name:        data.Name,
		Description: data.Description,
		Allocation:  data.Allocation,
	}

	return &response, nil
}

func (s *ticketService) TicketPurchase(ctx context.Context, request *dto.TicketPurchaseRequest) error {
	ticketPurchase := models.Purchase{
		TicketId:  request.TicketId,
		UserId:    request.UserId,
		Quantity:  request.Quantity,
		CreatedBy: request.UserId,
		UpdatedBy: request.UserId,
		CreatedAt: timeNow(),
		UpdatedAt: timeNow(),
	}

	err := s.purchaseRepo.Create(ctx, &ticketPurchase)
	if err != nil {
		return errors.New(messages.ErrorPurchase)
	}

	// Update ticket allocation
	ticket, err := s.ticketRepo.FindById(ctx, request.TicketId)
	if err != nil && err.Error() == "record not found" {
		return errors.New(messages.NotFound)
	}

	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	if ticket.Allocation < request.Quantity {
		return errors.New(messages.ErrorTicketAllocations)
	}

	ticket.Allocation -= request.Quantity
	_, err = s.ticketRepo.Update(ctx, ticket)
	if err != nil {
		return errors.New(messages.ErrorTicketUpdate)
	}

	return nil
}
