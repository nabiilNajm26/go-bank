package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/usecase"
)

type StatementHandler struct {
	statementUseCase *usecase.StatementUseCase
}

func NewStatementHandler(statementUseCase *usecase.StatementUseCase) *StatementHandler {
	return &StatementHandler{
		statementUseCase: statementUseCase,
	}
}

func (h *StatementHandler) GeneratePDFStatement(c *fiber.Ctx) error {
	accountID, err := uuid.Parse(c.Params("account_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	fromDateStr := c.Query("from_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	toDateStr := c.Query("to_date", time.Now().Format("2006-01-02"))

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid from_date format. Use YYYY-MM-DD",
		})
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid to_date format. Use YYYY-MM-DD",
		})
	}

	pdfBytes, err := h.statementUseCase.GeneratePDFStatement(c.Context(), accountID, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF statement",
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=statement.pdf")
	return c.Send(pdfBytes)
}

func (h *StatementHandler) GenerateCSVStatement(c *fiber.Ctx) error {
	accountID, err := uuid.Parse(c.Params("account_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	fromDateStr := c.Query("from_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	toDateStr := c.Query("to_date", time.Now().Format("2006-01-02"))

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid from_date format. Use YYYY-MM-DD",
		})
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid to_date format. Use YYYY-MM-DD",
		})
	}

	csvBytes, err := h.statementUseCase.GenerateCSVStatement(c.Context(), accountID, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate CSV statement",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=statement.csv")
	return c.Send(csvBytes)
}