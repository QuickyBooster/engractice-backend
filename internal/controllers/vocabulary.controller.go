package controllers

import (
	"engractice/internal/models"
	"engractice/internal/services"
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

// Test godoc
// @Summary Get all vocabulary
// @Schemes http https
// @Description Get all vocabulary
// @Tags vocabulary
// @Accept json
// @Produce json
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary [get]
func (vc *VocabularyController) GetVocabulary(c *fiber.Ctx) error {

	vocabs, err := vc.vocabularyService.GetAllWords()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to create test",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    vocabs,
	})

}
