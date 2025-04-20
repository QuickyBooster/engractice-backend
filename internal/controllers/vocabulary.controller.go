package controllers

import (
	"engractice/internal/models"
	"engractice/internal/services"
	"strconv"
	"time"

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

// Vocabulary godoc
// @Summary Get all vocabulary
// @Schemes http https
// @Description Get all vocabulary
// @Tags vocabulary
// @Accept json
// @Produce json
// @Success  200 {object} models.Vocabulary
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary [get]
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

// Vocabulary godoc
// @Summary Get word by ID
// @Schemes http https
// @Description Get word by ID
// @Tags vocabulary
// @Accept json
// @Produce json
// @Param id path string true "Word ID"
// @Success  200 {object} models.Vocabulary
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary/{id} [get]
func (vc *VocabularyController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	vocabulary, err := vc.vocabularyService.GetByID(id)
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

// Vocabulary godoc
// @Summary Create a new word
// @Schemes http https
// @Description Create a new word
// @Tags vocabulary
// @Produce json
// @Accept json
// @Param {object} body models.VocabularyDTO true "Create a new word"
// @Success  201 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary [post]
func (vc *VocabularyController) Create(c *fiber.Ctx) error {
	var vocabulary models.Vocabulary

	if err := c.BodyParser(&vocabulary); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}
	vocabulary.CreatedAt = time.Now()

	if _, err := vc.vocabularyService.Create(&vocabulary); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to create vocabulary",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    vocabulary,
	})
}

// Vocabulary godoc
// @Summary Edit a word
// @Schemes http https
// @Description Edit a word
// @Tags vocabulary
// @Accept json
// @Produce json
// @Param {object} body models.VocabularyDTO true "Edit a word"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary/{id} [put]
func (vc *VocabularyController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var vocabulary models.Vocabulary

	if err := c.BodyParser(&vocabulary); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}

	if _, err := vc.vocabularyService.Update(id, &vocabulary); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to update vocabulary",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    vocabulary,
	})
}

// Vocabulary godoc
// @Summary Delete a word
// @Schemes http https
// @Description Delete a word
// @Tags vocabulary
// @Accept json
// @Produce json
// @Param id path string true "Word ID"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary/{id} [delete]
func (vc *VocabularyController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := vc.vocabularyService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to delete vocabulary",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    nil,
	})
}

// Vocabulary godoc
// @Summary Search a word
// @Schemes http https
// @Description Search a word
// @Tags vocabulary
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/vocabulary/search [get]
func (vc *VocabularyController) Search(c *fiber.Ctx) error {
	query := c.Query("query")
	pageParam := c.Query("page", "1")
	// Convert page to int64
	page, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	vocabulary, err := vc.vocabularyService.Search(&query, &page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to search vocabulary",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    vocabulary,
	})
}
