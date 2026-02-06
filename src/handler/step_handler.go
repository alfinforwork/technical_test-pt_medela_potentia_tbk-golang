package handler

import (
	"encoding/json"
	"strconv"
	"technical-test/src/response"
	"technical-test/src/usecase"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/datatypes"
)

type StepHandler struct {
	stepUsecase     usecase.StepUsecase
	workflowUsecase usecase.WorkflowUsecase
}

func NewStepHandler(stepUsecase usecase.StepUsecase, workflowUsecase usecase.WorkflowUsecase) *StepHandler {
	return &StepHandler{
		stepUsecase:     stepUsecase,
		workflowUsecase: workflowUsecase,
	}
}

// CreateStep godoc
// @Summary Create a new step in a workflow
// @Description Add a new step to an existing workflow with actor and optional conditions
// @Tags Steps
// @Security Bearer
// @Accept json
// @Produce json
// @Param workflowId path int true "Workflow ID"
// @Param body body object{actor=string,conditions=object} true "Create Step Request"
// @Success 200 {object} response.ResponseSuccess "Step created successfully"
// @Failure 400 {object} response.ResponseError "Validation error"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Workflow not found"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /v1/workflows/{workflowId}/steps [post]
func (h *StepHandler) CreateStep(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}
	_, err = h.workflowUsecase.GetWorkflowByID(workflowId)
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

	step, err := h.stepUsecase.CreateStep(workflowId, body.Actor, conditionsJSON)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Step created successfully", step, nil)
}

// FindStepsByWorkflowID godoc
// @Summary List all steps in a workflow
// @Description Get all steps for a specific workflow with pagination support
// @Tags Steps
// @Security Bearer
// @Accept json
// @Produce json
// @Param workflowId path int true "Workflow ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search by step actor"
// @Success 200 {object} response.ResponseSuccess "Steps retrieved successfully"
// @Failure 400 {object} response.ResponseError "Invalid workflow ID"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Workflow not found"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /v1/workflows/{workflowId}/steps [get]
func (h *StepHandler) FindStepsByWorkflowID(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}
	_, err = h.workflowUsecase.GetWorkflowByID(workflowId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Workflow not found", nil)
	}

	params := utils.GetPaginationParams(c)

	steps, total, err := h.stepUsecase.FindStepsByWorkflowIDWithPagination(workflowId, params.Page, params.PageSize, params.Search)
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

func (h *StepHandler) FindStepByID(c fiber.Ctx) error {
	stepId, err := strconv.Atoi(c.Params("stepId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid step ID", nil)
	}

	step, err := h.stepUsecase.GetStepByID(stepId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Step not found", nil)
	}

	return response.Success(c, "Step retrieved successfully", step, nil)
}

func (h *StepHandler) UpdateStep(c fiber.Ctx) error {
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

	step, err := h.stepUsecase.UpdateStep(stepId, body.Level, body.Actor, conditionsJSON)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to update step", nil)
	}

	return response.Success(c, "Step updated successfully", step, nil)
}
