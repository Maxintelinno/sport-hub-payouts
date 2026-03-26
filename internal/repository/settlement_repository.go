package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sport-hub/sport-hub-payouts/internal/model"
)

type SettlementRepository interface {
	GetWalletSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error)
}

type settlementRepository struct {
	db *sqlx.DB
}

func NewSettlementRepository(db *sqlx.DB) SettlementRepository {
	return &settlementRepository{db: db}
}

func (r *settlementRepository) GetWalletSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error) {
	var summary model.WalletSummary

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN status = 'available' THEN net_amount ELSE 0 END), 0) as available_balance,
			COALESCE(SUM(CASE WHEN status = 'processing' THEN net_amount ELSE 0 END), 0) as processing_balance,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN net_amount ELSE 0 END), 0) as paid_out_total
		FROM owner_settlements
		WHERE owner_id = $1
	`

	err := r.db.GetContext(ctx, &summary, query, ownerID)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
