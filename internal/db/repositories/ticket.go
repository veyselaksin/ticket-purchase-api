package repositories

import (
	"context"
	"gorm.io/gorm"
	"sync"
	"ticket-purchase/internal/db/models"
)

//go:generate mockgen -destination=../../mocks/repositories/ticket_repository_mock.go -package=repositories ticket-purchase/internal/db/repositories TicketRepository
type TicketRepository interface {
	FindAll(ctx context.Context) ([]models.Ticket, error)
	FindById(ctx context.Context, id string) (*models.Ticket, error)
	Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	Update(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
}

type ticketRepository struct {
	db        *gorm.DB
	dbMutex   sync.Mutex
	tableName string
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	var ticketModel models.Ticket
	return &ticketRepository{db: db, tableName: ticketModel.TableName()}
}

func (r *ticketRepository) FindAll(ctx context.Context) ([]models.Ticket, error) {
	var tickets []models.Ticket
	result := r.db.Table(r.tableName).WithContext(ctx).Find(&tickets)
	return tickets, result.Error
}

func (r *ticketRepository) FindById(ctx context.Context, id string) (*models.Ticket, error) {
	var ticket models.Ticket
	result := r.db.Table(r.tableName).WithContext(ctx).Where("id = ?", id).First(&ticket)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ticket, nil
}

func (r *ticketRepository) Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	result := r.db.Table(r.tableName).WithContext(ctx).Create(ticket).Scan(&ticket)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	return ticket, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	r.dbMutex.Lock()
	defer r.dbMutex.Unlock()

	// Transaction begins
	tx := r.db.Begin()
	defer tx.Commit()

	result := tx.Table(r.tableName).WithContext(ctx).Where("id = ?", ticket.Id).Updates(ticket).Scan(&ticket)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	return ticket, nil
}
