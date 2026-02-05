package controller

import (
	"encoding/json"
	"strconv"
	"technical-test/src/response"
	"technical-test/src/service"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/datatypes"
)

type StepController struct {
	StepService     service.StepService
	WorkflowService service.WorkflowService
}

func NewStepController(workflowService service.WorkflowService, stepService service.StepService) *StepController {
	return &StepController{
		StepService:     stepService,
		WorkflowService: workflowService,
	}
}

func (sc *StepController) CreateStep(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}
	_, err = sc.WorkflowService.GetWorkflowByID(workflowId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Workflow not found", nil)
	}

	var body struct {
		Actor      string          `json:"actor" validate:"required"`
		Conditions json.RawMessage `json:"conditions"`
	}
	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	var conditionsJSON datatypes.JSON
	if len(body.Conditions) > 0 {
		if !json.Valid(body.Conditions) {
			c.Status(fiber.StatusBadRequest)
			return response.Error(c, "Invalid conditions format", nil)
		}
		conditionsJSON = datatypes.JSON(body.Conditions)
	}

	step, err := sc.StepService.CreateStep(workflowId, body.Actor, conditionsJSON)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Step created successfully", step, nil)
}

func (sc *StepController) FindStepsByWorkflowID(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}
	_, err = sc.WorkflowService.GetWorkflowByID(workflowId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Workflow not found", nil)
	}

	params := utils.GetPaginationParams(c)

	steps, total, err := sc.StepService.FindStepsByWorkflowIDWithPagination(workflowId, params.Page, params.PageSize, params.Search)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to retrieve steps", nil)
	}

	totalPages := utils.CalculateTotalPages(total, params.PageSize)
	meta := utils.PaginationMeta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	data := fiber.Map{
		"steps":      steps,
		"pagination": meta,
	}

	return response.Success(c, "Steps retrieved successfully", data, nil)
}

func (sc *StepController) FindStepByID(c fiber.Ctx) error {
	stepId, err := strconv.Atoi(c.Params("stepId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid step ID", nil)
	}

	step, err := sc.StepService.GetStepByID(stepId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Step not found", nil)
	}

	return response.Success(c, "Step retrieved successfully", step, nil)
}

func (sc *StepController) UpdateStep(c fiber.Ctx) error {
	stepId, err := strconv.Atoi(c.Params("stepId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid step ID", nil)
	}

	var body struct {
		Level      uint            `json:"level"`
		Actor      string          `json:"actor"`
		Conditions json.RawMessage `json:"conditions"`
	}
	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request body", nil)
	}

	var conditionsJSON datatypes.JSON
	if len(body.Conditions) > 0 {
		if !json.Valid(body.Conditions) {
			c.Status(fiber.StatusBadRequest)
			return response.Error(c, "Invalid conditions format", nil)
		}
		conditionsJSON = datatypes.JSON(body.Conditions)
	}

	step, err := sc.StepService.UpdateStep(stepId, body.Level, body.Actor, conditionsJSON)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to update step", nil)
	}

	return response.Success(c, "Step updated successfully", step, nil)
}

func (sc *StepController) DeleteStep(c fiber.Ctx) error {
	stepId, err := strconv.Atoi(c.Params("stepId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid step ID", nil)
	}

	err = sc.StepService.DeleteStep(stepId)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to delete step", nil)
	}

	return response.Success(c, "Step deleted successfully", nil, nil)
}
