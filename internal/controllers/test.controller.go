package controllers

import (
	"engractice/internal/models"
	"engractice/internal/services"
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
	if request.Quantity == 0 || (request.Tags == "") {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	test, err := tc.testService.CreateTest(request.Quantity, request.Tags)
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
