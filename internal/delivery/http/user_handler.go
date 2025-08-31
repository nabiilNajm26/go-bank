package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/s3"
	"github.com/nabiilNajm26/go-bank/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	s3Service   *s3.S3Service
}

func NewUserHandler(userUseCase *usecase.UserUseCase, s3Service *s3.S3Service) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		s3Service:   s3Service,
	}
}

// UploadProfileImage godoc
// @Summary Upload profile image
// @Description Upload profile image to S3 with safety limits (5MB max, PNG/JPG only)
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param image formData file true "Profile image file (PNG, JPG, JPEG, max 5MB)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile/image [post]
func (h *UserHandler) UploadProfileImage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Upload to S3
	result, err := h.s3Service.UploadProfileImage(c.Context(), userID, src, file)
	if err != nil {
		if strings.Contains(err.Error(), "exceeds maximum") || 
		   strings.Contains(err.Error(), "invalid file type") ||
		   strings.Contains(err.Error(), "maximum file limit") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image",
		})
	}

	// Update user profile image URL
	err = h.userUseCase.UpdateProfileImage(c.Context(), userID, result.URL)
	if err != nil {
		// If user update fails, cleanup S3 file
		h.s3Service.DeleteFile(c.Context(), result.Key)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user profile",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Profile image uploaded successfully",
		"image_url":    result.URL,
		"file_size":    result.Size,
		"content_type": result.MimeType,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	user, err := h.userUseCase.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user profile",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile information (name, email, phone)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateUserRequest true "User update request"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userUseCase.UpdateUser(c.Context(), userID, &req)
	if err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == usecase.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(user)
}

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Delete user account and all associated data (irreversible)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [delete]
func (h *UserHandler) DeleteProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.userUseCase.DeleteUser(c.Context(), userID)
	if err != nil {
		if err == usecase.ErrUserHasActiveAccounts {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Cannot delete user with active accounts",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile deleted successfully",
	})
}