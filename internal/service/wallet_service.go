package service

import (
	"context"
	"fmt"

	"github.com/sport-hub/sport-hub-payouts/internal/model"
	"github.com/sport-hub/sport-hub-payouts/internal/repository"
)

type WalletService interface {
	GetSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error)
}

type walletService struct {
	repo repository.SettlementRepository
}

func NewWalletService(repo repository.SettlementRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetSummary(ctx context.Context, ownerID string) (*model.WalletSummary, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository is not initialized (possibly due to database connection failure)")
	}
	return s.repo.GetWalletSummary(ctx, ownerID)
}
