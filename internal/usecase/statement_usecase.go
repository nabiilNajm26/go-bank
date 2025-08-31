package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/repository"
)

type StatementUseCase struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewStatementUseCase(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) *StatementUseCase {
	return &StatementUseCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *StatementUseCase) GeneratePDFStatement(ctx context.Context, accountID uuid.UUID, fromDate, toDate time.Time) ([]byte, error) {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	filter := &domain.TransactionFilter{
		AccountID: accountID,
		FromDate:  fromDate,
		ToDate:    toDate,
		Limit:     1000,
	}

	transactions, err := uc.transactionRepo.GetByAccountID(ctx, accountID, filter)
	if err != nil {
		return nil, err
	}

	cfg := config.NewBuilder().Build()
	mrt := maroto.New(cfg)

	// Header
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New("BANK STATEMENT", props.Text{
					Size: 16,
				}),
			),
		),
	)

	// Account Info
	mrt.AddRows(
		row.New(5).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Account Number: %s", account.AccountNumber), props.Text{Size: 10}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Current Balance: %s %s", account.Balance.String(), account.Currency), props.Text{Size: 10}),
			),
		),
	)

	// Generate simple PDF content
	for _, tx := range transactions {
		description := "Transfer"
		if tx.Description != nil {
			description = *tx.Description
		}

		amount := tx.Amount.String()
		if tx.FromAccountID != nil && *tx.FromAccountID == accountID {
			amount = "-" + amount
		} else {
			amount = "+" + amount
		}

		mrt.AddRows(
			row.New(5).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("%s | %s | %s | %s", 
						tx.CreatedAt.Format("2006-01-02"), 
						string(tx.Type), 
						description, 
						amount), props.Text{Size: 8}),
				),
			),
		)
	}

	document, err := mrt.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

func (uc *StatementUseCase) GenerateCSVStatement(ctx context.Context, accountID uuid.UUID, fromDate, toDate time.Time) ([]byte, error) {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	filter := &domain.TransactionFilter{
		AccountID: accountID,
		FromDate:  fromDate,
		ToDate:    toDate,
		Limit:     1000,
	}

	transactions, err := uc.transactionRepo.GetByAccountID(ctx, accountID, filter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	headers := []string{"Date", "Type", "Description", "Amount", "Balance", "Reference"}
	writer.Write(headers)

	// Transactions
	for _, tx := range transactions {
		description := "Transfer"
		if tx.Description != nil {
			description = *tx.Description
		}

		amount := tx.Amount.String()
		if tx.FromAccountID != nil && *tx.FromAccountID == accountID {
			amount = "-" + amount
		} else {
			amount = "+" + amount
		}

		record := []string{
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
			string(tx.Type),
			description,
			amount,
			account.Balance.String(),
			tx.Reference,
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}