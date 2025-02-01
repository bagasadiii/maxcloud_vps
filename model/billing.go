package model

import (
	"time"

	"github.com/google/uuid"
)
type Billing struct {
	BillingID   uuid.UUID `json:"billing_id"`
	CPU         int       `json:"cpu"`
	RAM         int       `json:"ram"`
	Storage     int       `json:"storage"`
	DownPayment int       `json:"down_payment"`
	MonthlyFee  int       `json:"monthly_fee"`
	CostPerHour int       `json:"cost_per_hour"`
	TotalFee    int       `json:"total_fee"`
	Uptime      int       `json:"uptime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClientID    uuid.UUID `json:"client_id"`
}

func (b *Billing) CalculateCostPerHour() int {
	cost := (b.CPU * CPUPrice) + (b.RAM / 1024 * RAMPrice) + (b.Storage * StoragePrice)
	return cost
}

func (b *Billing) CalculateMonthlyFee() int {
	hourlyCost := b.CalculateCostPerHour()
	monthlyCost := hourlyCost * 24 * 30
	return monthlyCost
}

