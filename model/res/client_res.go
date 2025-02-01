package res

import (
	"time"

	"github.com/google/uuid"
)

type ClientInfo struct {
	ClientID       uuid.UUID
	Email          string
	Suspended      bool
	Plan           string
	Balance        int
	ClientCreated  time.Time
	ClientUpdated  time.Time
	BillingID      uuid.UUID
	CPU            int
	RAM            int
	Storage        int
	MonthlyFee     int
	CostPerHour    int
	TotalFee       int
	Uptime         int
	BillingCreated time.Time
	BillingUpdated time.Time
}
