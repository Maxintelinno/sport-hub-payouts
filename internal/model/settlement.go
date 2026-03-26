package model

import (
	"time"

	"github.com/google/uuid"
)

type SettlementStatus string

const (
	StatusPending    SettlementStatus = "pending"
	StatusAvailable  SettlementStatus = "available"
	StatusProcessing SettlementStatus = "processing"
	StatusPaid       SettlementStatus = "paid"
	StatusHold       SettlementStatus = "hold"
	StatusReversed   SettlementStatus = "reversed"
)

type OwnerSettlement struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	BookingID      uuid.UUID        `json:"booking_id" db:"booking_id"`
	OwnerID        uuid.UUID        `json:"owner_id" db:"owner_id"`
	GrossAmount    float64          `json:"gross_amount" db:"gross_amount"`
	PlatformFee    float64          `json:"platform_fee" db:"platform_fee"`
	DiscountAmount float64          `json:"discount_amount" db:"discount_amount"`
	NetAmount      float64          `json:"net_amount" db:"net_amount"`
	Status         SettlementStatus `json:"status" db:"status"`
	AvailableAt    *time.Time       `json:"available_at,omitempty" db:"available_at"`
	PaidAt         *time.Time       `json:"paid_at,omitempty" db:"paid_at"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

type WalletSummary struct {
	AvailableBalance  float64 `json:"available_balance"`
	ProcessingBalance float64 `json:"processing_balance"`
	PaidOutTotal      float64 `json:"paid_out_total"`
}
