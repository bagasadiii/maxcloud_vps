package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bagasadiii/maxcloud_vps/model"
	"github.com/bagasadiii/maxcloud_vps/model/req"
	"github.com/bagasadiii/maxcloud_vps/model/res"
	"github.com/bagasadiii/maxcloud_vps/repository"
	"github.com/bagasadiii/maxcloud_vps/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClientServiceImpl interface {
	CreateClientService(ctx context.Context, req *req.NewClient) error
	GetClientInfoService(ctx context.Context, clientID uuid.UUID) (*res.ClientInfo, error)
}
type ClientService struct {
	repo   repository.ClientRepoImpl
	logger *zap.Logger
}

func NewClientService(repo repository.ClientRepoImpl, logger *zap.Logger) *ClientService {
	return &ClientService{
		repo:   repo,
		logger: logger,
	}
}

func (cs *ClientService) CreateClientService(ctx context.Context, req *req.NewClient) error {
	billing, err := selectBilling(req.Plan)
	if err != nil {
		cs.logger.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return fmt.Errorf("%v: %w",err, utils.ErrBadRequest)
	}
	remainingBalance := req.Balance - billing.DownPayment
	if remainingBalance < billing.CalculateMonthlyFee() {
		info := fmt.Sprintf("remaining balance: %d, Monthly fee: %d", remainingBalance, billing.CalculateMonthlyFee())
		cs.logger.Error(utils.ErrBadRequest.Error(), zap.String("insufficient balance", info))
		return fmt.Errorf("insufficient fund: %s: %w", info, utils.ErrBadRequest)
	}
	client := &model.Client{
		ClientID:  uuid.New(),
		Email:     req.Email,
		Suspended: false,
		Balance:   remainingBalance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	clientBilling := &model.Billing{
		BillingID:   uuid.New(),
		ClientID:    client.ClientID,
		CPU:         billing.CPU,
		RAM:         billing.RAM,
		Storage:     billing.Storage,
		MonthlyFee:  billing.CalculateMonthlyFee(),
		CostPerHour: billing.CalculateCostPerHour(),
		TotalFee:    0,
		Uptime:      0,
	}

	return cs.repo.CreateClientRepo(ctx, client, clientBilling)
}

func (cs *ClientService) GetClientInfoService(ctx context.Context, clientID uuid.UUID) (*res.ClientInfo, error) {
	return cs.repo.GetClientInfoRepo(ctx, clientID)
}

func selectBilling(plan string) (*model.Billing, error) {
	switch plan {
	case model.BasicBilling:
		return &model.Billing{
			CPU:         model.Basic.CPU,
			RAM:         model.Basic.RAM,
			Storage:     model.Basic.Storage,
			DownPayment: model.Basic.DownPayment,
		}, nil
	case model.NormalBilling:
		return &model.Billing{
			CPU:         model.Normal.CPU,
			RAM:         model.Normal.RAM,
			Storage:     model.Normal.Storage,
			DownPayment: model.Normal.DownPayment,
		}, nil
	case model.PremiumBilling:
		return &model.Billing{
			CPU:         model.Premium.CPU,
			RAM:         model.Premium.RAM,
			Storage:     model.Premium.Storage,
			DownPayment: model.Premium.DownPayment,
		}, nil
	default:
		return nil, fmt.Errorf("plan '%s' is not recognized", plan)
	}
}
