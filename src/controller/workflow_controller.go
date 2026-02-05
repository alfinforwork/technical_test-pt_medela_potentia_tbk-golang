package controller

import (
	"strconv"
	"technical-test/src/response"
	"technical-test/src/service"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type WorkflowController struct {
	WorkflowService service.WorkflowService
}

func NewWorkflowController(workflowService service.WorkflowService) *WorkflowController {
	return &WorkflowController{
		WorkflowService: workflowService,
	}
}

func (wc *WorkflowController) CreateWorkflow(c fiber.Ctx) error {
	type Body struct {
		Name string `json:"name" form:"name" query:"name"  validate:"required"`
	}
	body := new(Body)
	if err := c.Bind().Body(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}
	w, err := wc.WorkflowService.CreateWorkflow(body.Name)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}
	return response.Success(c, "Workflow created successfully", w, nil)
}

func (wc *WorkflowController) FindAllWorkflows(c fiber.Ctx) error {
	params := utils.GetPaginationParams(c)

	workflows, total, err := wc.WorkflowService.FindAllWorkflowsWithPagination(params.Page, params.PageSize, params.Search)
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

func (wc *WorkflowController) GetWorkflowByID(c fiber.Ctx) error {
	workflowId, err := strconv.Atoi(c.Params("workflowId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid workflow ID", nil)
	}

	workflow, err := wc.WorkflowService.GetWorkflowByID(workflowId)
	if err != nil {
		return response.Error(c, "Workflow not found", nil)
	}
	return response.Success(c, "Workflow retrieved successfully", workflow, nil)
}
