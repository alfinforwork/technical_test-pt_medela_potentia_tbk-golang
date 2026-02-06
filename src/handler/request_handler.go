package handler

import (
	"strconv"
	"technical-test/src/response"
	"technical-test/src/usecase"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type RequestHandler struct {
	requestUsecase  usecase.RequestUsecase
	workflowUsecase usecase.WorkflowUsecase
}

func NewRequestHandler(requestUsecase usecase.RequestUsecase, workflowUsecase usecase.WorkflowUsecase) *RequestHandler {
	return &RequestHandler{
		requestUsecase:  requestUsecase,
		workflowUsecase: workflowUsecase,
	}
}

// CreateRequest godoc
// @Summary Create a new request
// @Description Submit a new request for a workflow
// @Tags Requests
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body object{workflow_id=int,amount=number} true "Create Request"
// @Success 200 {object} response.ResponseSuccess "Request created successfully"
// @Failure 400 {object} response.ResponseError "Validation error"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Router /v1/requests [post]
func (h *RequestHandler) CreateRequest(c fiber.Ctx) error {
	var body struct {
		WorkflowID int     `json:"workflow_id" validate:"required"`
		Amount     float64 `json:"amount" validate:"required,gt=0"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	request, err := h.requestUsecase.CreateRequest(body.WorkflowID, body.Amount)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request created successfully", request, nil)
}

// FindAllRequests godoc
// @Summary List all requests
// @Description Get all requests with pagination and optional status filtering
// @Tags Requests
// @Security Bearer
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search by request ID"
// @Param status query string false "Filter by status (pending, approved, rejected)"
// @Success 200 {object} response.ResponseSuccess "Requests retrieved successfully"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /v1/requests [get]
func (h *RequestHandler) FindAllRequests(c fiber.Ctx) error {
	params := utils.GetPaginationParams(c)
	status := c.Query("status")

	requests, total, err := h.requestUsecase.FindAllRequestsWithPagination(params.Page, params.PageSize, params.Search, status)
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

// GetRequestByID godoc
// @Summary Get request by ID
// @Description Retrieve a specific request by its ID
// @Tags Requests
// @Security Bearer
// @Accept json
// @Produce json
// @Param requestId path int true "Request ID"
// @Success 200 {object} response.ResponseSuccess "Request retrieved successfully"
// @Failure 400 {object} response.ResponseError "Invalid request ID"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Request not found"
// @Router /v1/requests/{requestId} [get]
func (h *RequestHandler) GetRequestByID(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := h.requestUsecase.GetRequestByID(requestId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return response.Error(c, "Request not found", nil)
	}

	return response.Success(c, "Request retrieved successfully", request, nil)
}

// ApproveRequest godoc
// @Summary Approve a request
// @Description Approve a pending request
// @Tags Requests
// @Security Bearer
// @Accept json
// @Produce json
// @Param requestId path int true "Request ID"
// @Success 200 {object} response.ResponseSuccess "Request approved successfully"
// @Failure 400 {object} response.ResponseError "Invalid request ID"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Request not found"
// @Router /v1/requests/{requestId}/approve [post]
func (h *RequestHandler) ApproveRequest(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := h.requestUsecase.ApproveRequest(requestId)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request approved successfully", request, nil)
}

// RejectRequest godoc
// @Summary Reject a request
// @Description Reject a pending request
// @Tags Requests
// @Security Bearer
// @Accept json
// @Produce json
// @Param requestId path int true "Request ID"
// @Success 200 {object} response.ResponseSuccess "Request rejected successfully"
// @Failure 400 {object} response.ResponseError "Invalid request ID"
// @Failure 401 {object} response.ResponseError "Unauthorized"
// @Failure 404 {object} response.ResponseError "Request not found"
// @Router /v1/requests/{requestId}/reject [post]
func (h *RequestHandler) RejectRequest(c fiber.Ctx) error {
	requestId, err := strconv.Atoi(c.Params("requestId"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, "Invalid request ID", nil)
	}

	request, err := h.requestUsecase.RejectRequest(requestId)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Request rejected successfully", request, nil)
}
