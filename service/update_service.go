package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bagasadiii/maxcloud_vps/model"
	"github.com/bagasadiii/maxcloud_vps/repository"
	"github.com/bagasadiii/maxcloud_vps/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TransactionSchedulerServiceImpl interface {
	SchedulerWorkerService(ctx context.Context, worker int)
}
type TransactionSchedulerService struct {
	db     *pgxpool.Pool
	repo   repository.TransactionSchedulerRepoImpl
	logger *zap.Logger
}

func NewTransactionSchedulerService(db *pgxpool.Pool, repo repository.TransactionSchedulerRepoImpl, logger *zap.Logger) *TransactionSchedulerService {
	return &TransactionSchedulerService{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (hs *TransactionSchedulerService) SchedulerWorkerService(ctx context.Context, worker int) {
	jobs := make(chan *model.UpdateClient, 100)

	for i := 0; i < worker; i++ {
		go func(id int) {
			hs.logger.Info("Worker started", zap.Int("worker_id", id))
			for client := range jobs {
				hs.logger.Info("Worker processing transaction", zap.String("client_id", client.ClientID.String()))
				if err := hs.transactionService(ctx, client); err != nil {
					info := "failed to process transaction"
					hs.logger.Error(utils.ErrInternal.Error(), zap.String("error", info), zap.Error(err), zap.Any("client", client))
				}
				hs.logger.Info("Transaction success, balance deducted", zap.Any("client", client))
			}
		}(i)
	}

	hs.logger.Info("Worker pool started")
	for {
		select {
		case <-ctx.Done():
			hs.logger.Info("Stopping worker")
			return
		default:
			hs.schedulerService(ctx, jobs)
		}
	}
}

func (hs *TransactionSchedulerService) schedulerService(ctx context.Context, jobs chan<- *model.UpdateClient) {
	clients, err := hs.repo.GetActiveClient(ctx)
	if err != nil {
		info := "failed to get clients"
		hs.logger.Error(utils.ErrInternal.Error(), zap.String("error", info), zap.Error(err))
		return
	}
	for _, client := range clients {
		clientCopy := client
		now := time.Now()

		if now.Sub(client.UpdatedAt).Hours() >= 1 {
			jobs <- &clientCopy
		}
	}
}

func (hs *TransactionSchedulerService) transactionService(ctx context.Context, data *model.UpdateClient) error {
	tx, err := hs.db.Begin(ctx)
	if err != nil {
		info := "failed to begin transaction"
		hs.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}

	defer func() {
		if err != nil {
			hs.logger.Error(utils.ErrDatabase.Error(), zap.String("error", "rolling back transaction"), zap.Error(err))
			tx.Rollback(ctx)
		}
	}()

	threshold := int(float64(data.MonthlyFee) * 0.10)
	if data.Balance < threshold {
		hs.logger.Warn("Client have less than 10% of monthly fee", zap.Any("client", data))
	}

	newBalance := data.Balance - data.CostPerHour
	err = hs.repo.UpdateBalance(ctx, tx, data.ClientID, newBalance)
	if err != nil {
		return err
	}

	newFee := data.TotalFee + data.CostPerHour
	err = hs.repo.UpdateTotalFee(ctx, tx, data.BillingID, newFee)
	if err != nil {
		return err
	}

	err = hs.repo.UpdateClientInfo(ctx, tx, data.ClientID)
	if err != nil {
		return err
	}

	newUptime := data.Uptime + 1
	err = hs.repo.UpdateBillingInfo(ctx, tx, data.BillingID, newUptime)
	if err != nil {
		return err
	}
	if newBalance < 0 {
		err = hs.repo.SuspendClient(ctx, tx, data.ClientID)
		if err != nil {
			return err
		}
		hs.logger.Warn("Client suspended", zap.Any("client", data))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
