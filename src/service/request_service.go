package service

import (
	"encoding/json"
	"errors"
	"technical-test/src/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RequestService struct {
	db              *gorm.DB
	workflowService *WorkflowService
	stepService     *StepService
}

type stepConditions struct {
	MinAmount    float64 `json:"min_amount"`
	ApprovalType string  `json:"approval_type"` // "API" or "MANUAL"
}

var (
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInvalidRequestState = errors.New("request is not in pending state")
	ErrAmountBelowMinimum  = errors.New("amount does not meet minimum requirement for this step")
)

func NewRequestService(db *gorm.DB, ws *WorkflowService, ss *StepService) *RequestService {
	return &RequestService{
		db:              db,
		workflowService: ws,
		stepService:     ss,
	}
}

func (rs *RequestService) CreateRequest(workflowID int, amount float64) (model.Request, error) {
	if amount <= 0 {
		return model.Request{}, ErrInvalidAmount
	}
	_, err := rs.workflowService.GetWorkflowByID(workflowID)
	if err != nil {
		return model.Request{}, err
	}
	_, err = rs.stepService.FindStepByLevelAndWorkflowID(1, workflowID)
	if err != nil {
		return model.Request{}, err
	}

	var existingRequest model.Request
	err = rs.db.Where("workflow_id = ? AND status = ?", workflowID, "PENDING").First(&existingRequest).Error
	if err == nil {
		existingRequest.Amount += amount

		accumulatedMinAmount, err := rs.getAccumulatedMinAmount(workflowID, existingRequest.CurrentStep)
		if err != nil {
			return model.Request{}, err
		}

		if existingRequest.Amount >= accumulatedMinAmount {
			nextStep, err := rs.stepService.FindStepByLevelAndWorkflowID(existingRequest.CurrentStep+1, workflowID)
			if err == nil && nextStep.ID != 0 {
				existingRequest.CurrentStep += 1
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return model.Request{}, err
			} else {
				existingRequest.Status = "APPROVED"
			}
		}

		return existingRequest, rs.db.Save(&existingRequest).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Request{}, err
	}

	request := model.Request{
		WorkflowID:  uint(workflowID),
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      amount,
	}

	accumulatedMinAmount, err := rs.getAccumulatedMinAmount(workflowID, 1)
	if err != nil {
		return model.Request{}, err
	}

	if amount >= accumulatedMinAmount {
		nextStep, err := rs.stepService.FindStepByLevelAndWorkflowID(2, workflowID)
		if err == nil && nextStep.ID != 0 {
			request.CurrentStep = 2
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Request{}, err
		} else {
			request.Status = "APPROVED"
		}
	}

	result := rs.db.Create(&request)
	return request, result.Error
}

func (rs *RequestService) GetRequestByID(id int) (model.Request, error) {
	var request model.Request
	result := rs.db.First(&request, id)
	return request, result.Error
}

func (rs *RequestService) ApproveRequest(id int) (model.Request, error) {
	var request model.Request
	if err := rs.db.First(&request, id).Error; err != nil {
		return request, err
	}

	if request.Status != "PENDING" {
		return request, ErrInvalidRequestState
	}

	step, err := rs.stepService.FindStepByLevelAndWorkflowID(request.CurrentStep, int(request.WorkflowID))
	if err != nil {
		return request, err
	}

	conditions, err := parseConditions(step.Conditions)
	if err != nil {
		return request, err
	}

	if conditions.ApprovalType == "API" {
		accumulatedMinAmount, err := rs.getAccumulatedMinAmount(int(request.WorkflowID), request.CurrentStep)
		if err != nil {
			return request, err
		}

		if request.Amount >= accumulatedMinAmount {
			request.Status = "APPROVED"
			return request, rs.db.Save(&request).Error
		}

		return request, nil
	}
	request.Status = "APPROVED"
	return request, rs.db.Save(&request).Error

}

func (rs *RequestService) RejectRequest(id int) (model.Request, error) {
	var request model.Request
	if err := rs.db.First(&request, id).Error; err != nil {
		return request, err
	}

	if request.Status != "PENDING" {
		return request, ErrInvalidRequestState
	}

	request.Status = "REJECTED"
	return request, rs.db.Save(&request).Error
}

func parseMinAmount(conditions datatypes.JSON) (float64, error) {
	if len(conditions) == 0 {
		return 0, nil
	}

	var cond stepConditions
	if err := json.Unmarshal(conditions, &cond); err != nil {
		return 0, err
	}

	return cond.MinAmount, nil
}

func parseConditions(conditions datatypes.JSON) (stepConditions, error) {
	var cond stepConditions
	if len(conditions) == 0 {
		return cond, nil
	}

	if err := json.Unmarshal(conditions, &cond); err != nil {
		return cond, err
	}

	return cond, nil
}

func (rs *RequestService) getAccumulatedMinAmount(workflowID int, currentLevel uint) (float64, error) {
	var total float64 = 0

	for level := uint(1); level <= currentLevel; level++ {
		step, err := rs.stepService.FindStepByLevelAndWorkflowID(level, workflowID)
		if err != nil {
			return 0, err
		}

		minAmount, err := parseMinAmount(step.Conditions)
		if err != nil {
			return 0, err
		}

		total += minAmount
	}

	return total, nil
}
