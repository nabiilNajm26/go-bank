package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/usecase"
)

type AccountHandler struct {
	accountUseCase *usecase.AccountUseCase
}

func NewAccountHandler(accountUseCase *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
	}
}

// CreateAccount godoc
// @Summary Create new account
// @Description Create a new bank account for authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateAccountRequest true "Account creation request"
// @Success 201 {object} domain.Account
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req domain.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	account, err := h.accountUseCase.CreateAccount(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

func (h *AccountHandler) GetAccount(c *fiber.Ctx) error {
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	account, err := h.accountUseCase.GetAccount(c.Context(), accountID)
	if err != nil {
		if err == usecase.ErrAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get account",
		})
	}

	return c.JSON(account)
}

func (h *AccountHandler) GetUserAccounts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	accounts, err := h.accountUseCase.GetUserAccounts(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get accounts",
		})
	}

	return c.JSON(fiber.Map{
		"accounts": accounts,
	})
}

// UpdateAccount godoc
// @Summary Update account details
// @Description Update account name or other details for authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param request body domain.UpdateAccountRequest true "Account update request"
// @Success 200 {object} domain.Account
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	var req domain.UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	account, err := h.accountUseCase.UpdateAccount(c.Context(), userID, accountID, &req)
	if err != nil {
		if err == usecase.ErrAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == usecase.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only update your own accounts",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update account",
		})
	}

	return c.JSON(account)
}

// DeleteAccount godoc
// @Summary Delete account
// @Description Delete account if it has zero balance (safety measure)
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	err = h.accountUseCase.DeleteAccount(c.Context(), userID, accountID)
	if err != nil {
		if err == usecase.ErrAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == usecase.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only delete your own accounts",
			})
		}
		if err == usecase.ErrAccountNotEmpty {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Cannot delete account with non-zero balance",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}