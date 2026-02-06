package handler

import (
	"strconv"
	"technical-test/src/response"
	"technical-test/src/usecase"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type WorkflowHandler struct {
	workflowUsecase usecase.WorkflowUsecase
}

func NewWorkflowHandler(workflowUsecase usecase.WorkflowUsecase) *WorkflowHandler {
	return &WorkflowHandler{
		workflowUsecase: workflowUsecase,
	}
}

// CreateWorkflow godoc
// @Summary Create a new workflow
// @Description Create a new workflow with a given name
// @Tags Workflows
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body object{name=string} true "Create Workflow Request"
// @Success 200 {object} response.ResponseSuccess "Workflow created successfully"
// @Failure 400 {object} response.ResponseError "Validation error"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Router /v1/workflows [post]
func (h *WorkflowHandler) CreateWorkflow(c fiber.Ctx) error {
	type Body struct {
		Name string `json:"name" form:"name" query:"name" validate:"required"`
	}
	body := new(Body)
	if err := c.Bind().Body(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}
	w, err := h.workflowUsecase.CreateWorkflow(body.Name)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}
	return response.Success(c, "Workflow created successfully", w, nil)
}

// FindAllWorkflows godoc
// @Summary List all workflows
// @Description Get all workflows with pagination support
// @Tags Workflows
// @Security Bearer
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search by workflow name"
// @Success 200 {object} response.ResponseSuccess "Workflows retrieved successfully"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /v1/workflows [get]
func (h *WorkflowHandler) FindAllWorkflows(c fiber.Ctx) error {
	params := utils.GetPaginationParams(c)

	workflows, total, err := h.workflowUsecase.FindAllWorkflowsWithPagination(params.Page, params.PageSize, params.Search)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to retrieve workflows", nil)
	}

	totalPages := utils.CalculateTotalPages(total, params.PageSize)
	meta := utils.PaginationMeta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	data := fiber.Map{
		"workflows":  workflows,
		"pagination": meta,
	}

	return response.Success(c, "Workflows retrieved successfully", data, nil)
}

// GetWorkflowByID godoc
// @Summary Get workflow by ID
// @Description Retrieve a specific workflow by its ID
// @Tags Workflows
// @Security Bearer
// @Accept json
// @Produce json
// @Param workflowId path int true "Workflow ID"
// @Success 200 {object} response.ResponseSuccess "Workflow retrieved successfully"
// @Failure 400 {object} response.ResponseError "Invalid workflow ID"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Workflow not found"
// @Router /v1/workflows/{workflowId} [get]
func (h *WorkflowHandler) GetWorkflowByID(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}

	workflow, err := h.workflowUsecase.GetWorkflowByID(workflowId)
	if err != nil {
		return response.Error(c, "Workflow not found", nil)
	}
	return response.Success(c, "Workflow retrieved successfully", workflow, nil)
}
