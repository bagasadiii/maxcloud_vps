package model

const (
	CPUPrice       = 200
	RAMPrice       = 200
	StoragePrice   = 200
	BasicBilling   = "basic"
	NormalBilling  = "normal"
	PremiumBilling = "premium"
)

// Price per hour
var (
	Basic = Billing{
		CPU:         1,
		RAM:         1024,
		Storage:     8,
		DownPayment: 15000,
	}
	Normal = Billing{
		CPU:         2,
		RAM:         2048,
		Storage:     16,
		DownPayment: 25000,
	}
	Premium = Billing{
		CPU:         4,
		RAM:         4096,
		Storage:     16,
		DownPayment: 40000,
	}
)