package controllers

import (
	"engractice/internal/models"
	"engractice/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VocabularyController struct {
	vocabularyService *services.VocabularyService
}

func NewVocabularyController(svc *services.VocabularyService) *VocabularyController {
	return &VocabularyController{
		vocabularyService: svc,
	}
}

// GetAll godoc
// @Summary Get all vocabulary
// @Description Get all vocabulary
// @Tags vocabulary
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=[]models.Vocabulary}
// @Failure 500 {object} models.Response
// @Router /vocabulary [get]

func (vc *VocabularyController) GetAll(c *fiber.Ctx) error {
	// Default page = 1 if not provided
	pageParam := c.Query("page", "1")

	// Convert page to int64
	page, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	// Call service with context and page
	vocabulary, err := vc.vocabularyService.GetAll(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to get vocabulary",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    vocabulary,
	})
}
