package controllers

import (
	"engractice/internal/models"
	"engractice/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TestController struct {
	testService *services.TestService
}

func NewTestController(svc *services.TestService) *TestController {
	return &TestController{
		testService: svc,
	}
}

// Test godoc
// @Summary Create new test
// @Schemes http https
// @Description Create a new test
// @Tags test
// @Accept json
// @Produce json
// @Param {object} body models.TestRequest true "Create new test"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/test [post]
func (tc *TestController) CreateTest(c *fiber.Ctx) error {
	var request models.TestRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}
	if request.Quantity == 0 || (request.Tags == "" && !request.NearestMode) {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	test, err := tc.testService.CreateTest(&request)
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
		Data:    test,
	})

}

// Vocabulary godoc
// @Summary Get all tests
// @Schemes http https
// @Description Get all tests
// @Tags test
// @Accept json
// @Produce json
// @Param date  query string false "query using date"
// @Param tags query string false "query using tags"
// @Param nearestMode query bool false "query using nearest mode"
// @Param quantity query int false "query using quantity"
// @Param page query int  false "page response"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/test [get]
func (tc *TestController) GetAllTest(c *fiber.Ctx) error {
	date := c.Query("date", "")
	tags := c.Query("tags", "")
	nearestMode := c.Query("nearestMode", "")
	quantity := c.Query("quantity", "")
	page := c.Query("page", "")

	tests, err := tc.testService.GetAllTest(date, tags, nearestMode, quantity, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to get tests",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    tests,
	})
}

// Vocabulary godoc
// @Summary Upload finished test
// @Schemes http https
// @Description upload finished test
// @Tags test
// @Accept json
// @Produce json
// @Param {object} body models.Test true "Finish a test"
// @Success  200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/test/finish [post]
func (tc *TestController) FinishTest(c *fiber.Ctx) error {
	var test models.Test

	if err := c.BodyParser(&test); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}
	test.Date = time.Now()
	_, err := tc.testService.FinishTest(&test)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  false,
			Message: "Failed to upload the test result",
			Data:    nil,
		})
	}

	return c.JSON(models.Response{
		Status:  true,
		Message: "Success",
		Data:    nil,
	})
}
