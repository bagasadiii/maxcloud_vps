package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bagasadiii/maxcloud_vps/model"
	"github.com/bagasadiii/maxcloud_vps/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TransactionSchedulerRepoImpl interface {
	GetActiveClient(ctx context.Context) ([]model.UpdateClient, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, clientID uuid.UUID, newBalance int) error
	UpdateTotalFee(ctx context.Context, tx pgx.Tx, billingID uuid.UUID, newFee int) error
	UpdateClientInfo(ctx context.Context, tx pgx.Tx, clientID uuid.UUID) error
	UpdateBillingInfo(ctx context.Context, tx pgx.Tx, billingID uuid.UUID, newUptime int) error
	SuspendClient(ctx context.Context, tx pgx.Tx, billingID uuid.UUID) error
}
type TransactionSchedulerRepo struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewTransactionSchedulerRepo(db *pgxpool.Pool, logger *zap.Logger) *TransactionSchedulerRepo {
	return &TransactionSchedulerRepo{
		db:     db,
		logger: logger,
	}
}

func (hr *TransactionSchedulerRepo) GetActiveClient(ctx context.Context) ([]model.UpdateClient, error) {
	rows, err := hr.db.Query(ctx, `
    SELECT c.client_id, c.suspended, c.balance, c.created_at, b.monthly_fee, b.cost_per_hour, b.total_fee, b.uptime, b.billing_id
    FROM clients c
    JOIN billings b ON c.client_id = b.client_id
    WHERE c.suspended = false
    FOR UPDATE
    `)
	if err != nil {
		info := "failed to get client info"
		hr.logger.Error(utils.ErrBadRequest.Error(), zap.String("error", info), zap.Error(err))
		return nil, fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	defer rows.Close()
	var clients []model.UpdateClient
	for rows.Next() {
		var client model.UpdateClient
		err := rows.Scan(
			&client.ClientID,
			&client.Suspended,
			&client.Balance,
			&client.CreatedAt,
			&client.MonthlyFee,
			&client.CostPerHour,
			&client.TotalFee,
			&client.Uptime,
			&client.BillingID,
		)
		if err != nil {
			info := "failed while scanning client info"
			hr.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
			return nil, fmt.Errorf("%s: %w", info, utils.ErrDatabase)
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (hr *TransactionSchedulerRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, clientID uuid.UUID, newBalance int) error {
	_, err := tx.Exec(ctx, `
		UPDATE clients SET balance = $1 WHERE client_id = $2
	`, newBalance, clientID)
	if err != nil {
		info := "failed to update balance"
		hr.logger.Error(utils.ErrDatabase.Error(),
			zap.String("error", info),
			zap.String("client_id", clientID.String()),
			zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return nil
}

func (hr *TransactionSchedulerRepo) UpdateTotalFee(ctx context.Context, tx pgx.Tx, billingID uuid.UUID, newFee int) error {
	_, err := tx.Exec(ctx, `
		UPDATE billings SET total_fee = $1 WHERE billing_id = $2
	`, newFee, billingID)
	if err != nil {
		info := "failed to update total fee"
		hr.logger.Error(utils.ErrDatabase.Error(),
			zap.String("error", info),
			zap.String("billing_id", billingID.String()),
			zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return nil
}

func (hr *TransactionSchedulerRepo) UpdateClientInfo(ctx context.Context, tx pgx.Tx, clientID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		UPDATE clients SET updated_at = $1 WHERE client_id = $2
	`, time.Now(), clientID)
	if err != nil {
		info := "failed to update client info"
		hr.logger.Error(utils.ErrDatabase.Error(),
			zap.String("error", info),
			zap.String("client_id", clientID.String()),
			zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return nil
}

func (hr *TransactionSchedulerRepo) UpdateBillingInfo(ctx context.Context, tx pgx.Tx, billingID uuid.UUID, newUptime int) error {
	_, err := tx.Exec(ctx, `
		UPDATE billings SET uptime = $1, updated_at = $2 WHERE billing_id = $3
	`, newUptime, time.Now(), billingID)
	if err != nil {
		info := "failed to update billing info"
		hr.logger.Error(utils.ErrDatabase.Error(),
			zap.String("error", info),
			zap.String("billing_id", billingID.String()),
			zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return nil
}

func (hr *TransactionSchedulerRepo) SuspendClient(ctx context.Context, tx pgx.Tx, billingID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE clients SET suspended = true WHERE billing_id = $1`, billingID)
	if err != nil {
		info := "failed to suspend clients"
		hr.logger.Error(utils.ErrDatabase.Error(),
			zap.String("error", info),
			zap.String("billing_id", billingID.String()),
			zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return nil
}
