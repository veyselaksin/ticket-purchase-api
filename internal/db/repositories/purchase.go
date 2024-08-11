package repositories

import (
	"context"
	"gorm.io/gorm"
	"sync"
	"ticket-purchase/internal/db/models"
)

//go:generate mockgen -destination=../../mocks/repositories/purchase_repository_mock.go -package=repositories ticket-purchase/internal/db/repositories PurchaseRepository
type PurchaseRepository interface {
	Create(ctx context.Context, purchase *models.Purchase) error
}

type purchaseRepository struct {
	db        *gorm.DB
	dbMutex   sync.Mutex
	tableName string
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	var purchaseModel models.Purchase
	return &purchaseRepository{
		db:        db,
		tableName: purchaseModel.TableName(),
	}
}

func (r *purchaseRepository) Create(ctx context.Context, purchase *models.Purchase) error {
	tx := r.db.Begin()
	defer tx.Commit()

	result := tx.Table(r.tableName).WithContext(ctx).Create(purchase)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	return nil
}
