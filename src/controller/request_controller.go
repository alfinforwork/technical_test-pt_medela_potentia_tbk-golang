package controller

import (
	"strconv"
	"technical-test/src/response"
	"technical-test/src/service"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type RequestController struct {
	RequestService  service.RequestService
	WorkflowService service.WorkflowService
}

func NewRequestController(requestService service.RequestService, workflowService service.WorkflowService) *RequestController {
	return &RequestController{
		RequestService:  requestService,
		WorkflowService: workflowService,
	}
}

func (rc *RequestController) CreateRequest(c fiber.Ctx) error {
	var body struct {
		WorkflowID int     `json:"workflow_id" validate:"required"`
		Amount     float64 `json:"amount" validate:"required,gt=0"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	request, err := rc.RequestService.CreateRequest(body.WorkflowID, body.Amount)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request created successfully", request, nil)
}

func (rc *RequestController) FindAllRequests(c fiber.Ctx) error {
	params := utils.GetPaginationParams(c)
	status := c.Query("status")

	requests, total, err := rc.RequestService.FindAllRequestsWithPagination(params.Page, params.PageSize, params.Search, status)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return response.Error(c, "Failed to retrieve requests", nil)
	}

	totalPages := utils.CalculateTotalPages(total, params.PageSize)
	meta := utils.PaginationMeta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	data := fiber.Map{
		"requests":   requests,
		"pagination": meta,
	}

	return response.Success(c, "Requests retrieved successfully", data, nil)
}

func (rc *RequestController) GetRequestByID(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := rc.RequestService.GetRequestByID(requestId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Request not found", nil)
	}

	return response.Success(c, "Request retrieved successfully", request, nil)
}

func (rc *RequestController) ApproveRequest(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := rc.RequestService.ApproveRequest(requestId)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request approved successfully", request, nil)
}

func (rc *RequestController) RejectRequest(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := rc.RequestService.RejectRequest(requestId)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request rejected successfully", request, nil)
}
