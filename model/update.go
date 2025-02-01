package model

import (
	"time"

	"github.com/google/uuid"
)


type UpdateClient struct {
	ClientID    uuid.UUID
	BillingID		uuid.UUID
	Suspended   bool    
	Balance     int    
	MonthlyFee  int    
	CostPerHour int   
	TotalFee    int  
	Uptime      int
	CreatedAt		time.Time
}