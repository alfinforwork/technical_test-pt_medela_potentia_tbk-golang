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
	workflows, err := wc.WorkflowService.FindAllWorkflows()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to retrieve workflows", nil)
	}
	return response.Success(c, "Workflows retrieved successfully", workflows, nil)
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
