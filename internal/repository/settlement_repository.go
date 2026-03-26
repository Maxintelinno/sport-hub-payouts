package repository

import (
	"context"

	"github.com/sport-hub/sport-hub-payouts/internal/model"
	"gorm.io/gorm"
)

type SettlementRepository interface {
	GetWalletSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error)
}

type settlementRepository struct {
	db *gorm.DB
}

func NewSettlementRepository(db *gorm.DB) SettlementRepository {
	return &settlementRepository{db: db}
}

func (r *settlementRepository) GetWalletSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error) {
	var summary model.WalletSummary

	err := r.db.WithContext(ctx).Table("owner_settlements").
		Select(`
			COALESCE(SUM(CASE WHEN status = 'available' THEN net_amount ELSE 0 END), 0) as available_balance,
			COALESCE(SUM(CASE WHEN status = 'processing' THEN net_amount ELSE 0 END), 0) as processing_balance,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN net_amount ELSE 0 END), 0) as paid_out_total
		`).
		Where("owner_id = ?", ownerID).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
