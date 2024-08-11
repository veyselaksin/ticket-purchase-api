package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"testing"
	"ticket-purchase/internal/db/models"
	"ticket-purchase/internal/dto"
	"ticket-purchase/internal/i18n"
	"ticket-purchase/internal/mocks/repositories"
	"time"
)

var mockTicketData = []models.Ticket{
	{
		Id:          "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		Name:        "Ticket 1",
		Description: "Description 1",
		Allocation:  100,
	},
	{
		Id:          "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		Name:        "Ticket 2",
		Description: "Description 2",
		Allocation:  200,
	},
}

var mockPurchaseData = []models.Purchase{
	{
		Id:        "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		TicketId:  "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		UserId:    "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		Quantity:  1,
		CreatedBy: "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		UpdatedBy: "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		CreatedAt: time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
	},
	{
		Id:        "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		TicketId:  "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		UserId:    "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		Quantity:  2,
		CreatedBy: "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		UpdatedBy: "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4c",
		CreatedAt: time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
	},
}

var fiberCtx *fiber.Ctx
var s TicketService
var ticketRepo *repositories.MockTicketRepository
var purchaseRepo *repositories.MockPurchaseRepository

func setupTicketTest(t *testing.T) func() {
	ct := gomock.NewController(t)
	defer ct.Finish()

	app := fiber.New()
	fiberCtx = app.AcquireCtx(&fasthttp.RequestCtx{})

	// Assign language to fiber context header
	fiberCtx.Request().Header.Set("Accept-Language", "en")

	i18n.InitBundle("./../i18n/languages")
	ticketRepo = repositories.NewMockTicketRepository(ct)
	purchaseRepo = repositories.NewMockPurchaseRepository(ct)

	s = NewTicketService(ticketRepo, purchaseRepo)
	return func() {
		s = nil
		defer ct.Finish()
	}
}

func TestTicketService_Create_Success(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	request := dto.TicketCreateRequest{
		Name:        "Ticket 3",
		Description: "Description 3",
		Allocation:  100,
	}

	ticket := models.Ticket{
		Name:        request.Name,
		Description: request.Description,
		Allocation:  request.Allocation,
	}

	ticketRepo.EXPECT().Create(fiberCtx.Context(), &ticket).Return(&ticket, nil)

	response, err := s.Create(fiberCtx.Context(), &request)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, request.Name, response.Name)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.Allocation, response.Allocation)
}

func TestTicketService_Create_Failure(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	request := dto.TicketCreateRequest{
		Name:        "Ticket 3",
		Description: "Description 3",
		Allocation:  100,
	}

	ticket := models.Ticket{
		Name:        request.Name,
		Description: request.Description,
		Allocation:  request.Allocation,
	}

	ticketRepo.EXPECT().Create(fiberCtx.Context(), &ticket).Return(nil, assert.AnError)

	response, err := s.Create(fiberCtx.Context(), &request)
	if err == nil {
		t.Fatalf("Expected error to be not nil, got nil")
	}

	assert.Nil(t, response)
}

func TestTicketService_FindById_Success(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	id := "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b"
	ticket := mockTicketData[0]

	ticketRepo.EXPECT().FindById(fiberCtx.Context(), id).Return(&ticket, nil)

	response, err := s.FindById(fiberCtx.Context(), id)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, ticket.Name, response.Name)
	assert.Equal(t, ticket.Description, response.Description)
	assert.Equal(t, ticket.Allocation, response.Allocation)
}

func TestTicketService_FindById_Record_Not_Found(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	id := "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b"

	ticketRepo.EXPECT().FindById(fiberCtx.Context(), id).Return(&models.Ticket{}, assert.AnError)

	response, err := s.FindById(fiberCtx.Context(), id)
	if err == nil {
		t.Fatalf("Expected error to be not nil, got nil")

	}

	assert.Nil(t, response)
}

func TestTicketService_FindById_Failure(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	id := "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b"

	ticketRepo.EXPECT().FindById(fiberCtx.Context(), id).Return(nil, assert.AnError)

	response, err := s.FindById(fiberCtx.Context(), id)
	if err == nil {
		t.Fatalf("Expected error to be not nil, got nil")
	}

	assert.Nil(t, response)
}

func TestTicketService_TicketPurchase_Success(t *testing.T) {
	teardown := setupTicketTest(t)
	defer teardown()

	request := dto.TicketPurchaseRequest{
		TicketId: "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		UserId:   "4a4b3b3b-1b4b-4b3b-8b3b-3b4b3b4b3b4b",
		Quantity: 1,
	}

	mockTime := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time {
		return mockTime
	}
	defer func() { timeNow = time.Now }()

	purchase := models.Purchase{
		TicketId:  request.TicketId,
		UserId:    request.UserId,
		Quantity:  request.Quantity,
		CreatedBy: request.UserId,
		UpdatedBy: request.UserId,
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
	}

	purchaseRepo.EXPECT().Create(fiberCtx.Context(), &purchase).Return(nil)
	ticketRepo.EXPECT().FindById(fiberCtx.Context(), request.TicketId).Return(&mockTicketData[0], nil)
	ticketRepo.EXPECT().Update(fiberCtx.Context(), &mockTicketData[0]).Return(&mockTicketData[0], nil)

	err := s.TicketPurchase(fiberCtx.Context(), &request)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err)
	}

	assert.Equal(t, 99, mockTicketData[0].Allocation)
	assert.Equal(t, mockTime, mockTicketData[0].UpdatedAt)
	assert.Equal(t, 1, len(mockPurchaseData))
}
