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
				threshold := client.Balance
				if client.Balance < threshold {
					hs.logger.Info("This client's balance is less than 10%", zap.Any("client", client))
				}
				if client.Balance < 0 {
					client.Suspended = true
				}
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
			time.Sleep(5 * time.Second)
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

		nextPayment := time.Until(client.CreatedAt.Add(time.Hour))

		time.Sleep(nextPayment)

		jobs <- &clientCopy
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

	err = hs.repo.UpdateBillingInfo(ctx, tx, data.BillingID, data.Uptime+1)
	if err != nil {
		return err
	}

	if newBalance < 0 {
		err = hs.repo.SuspendClient(ctx, tx, data.BillingID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
