package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/usecase"
)

type TransactionHandler struct {
	transactionUseCase *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func (h *TransactionHandler) Transfer(c *fiber.Ctx) error {
	var req domain.TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	transaction, err := h.transactionUseCase.Transfer(c.Context(), &req)
	if err != nil {
		if err == usecase.ErrInsufficientBalance {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == usecase.ErrSameAccount || err == usecase.ErrInvalidAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process transfer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(transaction)
}

func (h *TransactionHandler) GetTransactionHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	
	// For simplicity, we'll get transactions for the user's first account
	// In production, you'd want to handle this differently
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "account_id is required",
		})
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	filter := &domain.TransactionFilter{
		AccountID: accountID,
		Limit:     50,
		Offset:    0,
	}

	transactions, err := h.transactionUseCase.GetTransactionHistory(c.Context(), accountID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
	})
}