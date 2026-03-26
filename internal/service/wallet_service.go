package service

import (
	"context"

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
	return s.repo.GetWalletSummary(ctx, ownerID)
}
