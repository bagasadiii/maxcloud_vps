package repository

import (
	"context"
	"fmt"

	"github.com/bagasadiii/maxcloud_vps/model"
	"github.com/bagasadiii/maxcloud_vps/model/res"
	"github.com/bagasadiii/maxcloud_vps/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ClientRepoImpl interface {
	CreateClientRepo(ctx context.Context, client *model.Client, billing *model.Billing) error
	GetClientInfoRepo(ctx context.Context, clientID uuid.UUID) (*res.ClientInfo, error)
}

type ClientRepo struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewClientRepo(db *pgxpool.Pool, logger *zap.Logger) *ClientRepo {
	return &ClientRepo{
		db:     db,
		logger: logger,
	}
}

func (cr *ClientRepo) CreateClientRepo(ctx context.Context, client *model.Client, billing *model.Billing) error {
	// Register a client for using the VPS service
	var exists bool
	err := cr.db.QueryRow(ctx, `
    SELECT EXISTS (SELECT 1 FROM clients WHERE email = $1)
    `, client.Email).Scan(&exists)
	if err != nil {
		info := "error while checking client"
		cr.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	if exists {
		info := "email exists"
		cr.logger.Warn(utils.ErrExists.Error(), zap.String("warn", info), zap.String("email", client.Email))
		return fmt.Errorf("%s: %w", info, utils.ErrExists)
	}
	tx, err := cr.db.Begin(ctx)
	if err != nil {
		info := "failed to begin transaction"
		cr.logger.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	_, err = tx.Exec(ctx, `
    INSERT INTO clients
    (client_id, email, suspended, balance, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
		`, client.ClientID, client.Email, client.Suspended, client.Balance, client.CreatedAt, client.UpdatedAt)
	if err != nil {
		info := "failed to add client"
		cr.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO billings
		(billing_id, cpu, ram, storage, monthly_fee, cost_per_hour, total_fee, uptime, updated_at, client_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, billing.BillingID, billing.CPU, billing.RAM, billing.Storage, billing.MonthlyFee,
		billing.CostPerHour, billing.TotalFee, billing.Uptime, billing.UpdatedAt, client.ClientID)
	if err != nil {
		info := "failed to add billing"
		cr.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
		return fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	cr.logger.Info("user and billing created", zap.String("email", client.Email),
		zap.String("client_id", client.ClientID.String()))
	return nil
}

func (cr *ClientRepo) GetClientInfoRepo(ctx context.Context, clientID uuid.UUID) (*res.ClientInfo, error) {
	var clientInfo res.ClientInfo
	err := cr.db.QueryRow(ctx, `
  SELECT
	  c.client_id, c.email, c.suspended, c.balance, c.created_at, c.updated_at,
	  b.billing_id, b.cpu, b.ram, b.storage, b.monthly_fee,
	  b.cost_per_hour, b.total_fee, b.uptime, b.created_at AS billing_created_at, b.updated_at AS billing_updated_at
	FROM clients c
	LEFT JOIN billings b ON c.client_id = b.client_id
	WHERE c.client_id = $1
  `, clientID).Scan(
		&clientInfo.ClientID, &clientInfo.Email, &clientInfo.Suspended, &clientInfo.Balance,
		&clientInfo.ClientCreated, &clientInfo.ClientUpdated,
		&clientInfo.BillingID, &clientInfo.CPU, &clientInfo.RAM, &clientInfo.Storage,
		&clientInfo.MonthlyFee, &clientInfo.CostPerHour, &clientInfo.TotalFee, &clientInfo.Uptime,
		&clientInfo.BillingCreated, &clientInfo.BillingUpdated,
	)
	if err == pgx.ErrNoRows {
		info := "client id not found"
		cr.logger.Warn(utils.ErrNotFound.Error(), zap.String("warn", info), zap.Error(err))

		return nil, fmt.Errorf("%s: %w", info, utils.ErrNotFound)
	} else if err != nil {
		info := "failed while scanning client data"
		cr.logger.Error(utils.ErrDatabase.Error(), zap.String("error", info), zap.Error(err))
		return nil, fmt.Errorf("%s: %w", info, utils.ErrDatabase)
	}
	return &clientInfo, nil
}
