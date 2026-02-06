package usecase

import (
	"encoding/json"
	"errors"
	"technical-test/src/model"
	"technical-test/src/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RequestUsecase interface {
	CreateRequest(workflowID int, amount float64) (model.Request, error)
	GetRequestByID(id int) (model.Request, error)
	FindAllRequestsWithPagination(page, pageSize int, search, status string) ([]model.Request, int64, error)
	ApproveRequest(id int) (model.Request, error)
	RejectRequest(id int) (model.Request, error)
}

type requestUsecase struct {
	requestRepo  repository.RequestRepository
	stepRepo     repository.StepRepository
	workflowRepo repository.WorkflowRepository
}

type stepConditions struct {
	MinAmount    float64 `json:"min_amount"`
	ApprovalType string  `json:"approval_type"`
}

var (
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInvalidRequestState = errors.New("request is not in pending state")
	ErrAmountBelowMinimum  = errors.New("amount does not meet minimum requirement for this step")
)

func NewRequestUsecase(requestRepo repository.RequestRepository, stepRepo repository.StepRepository, workflowRepo repository.WorkflowRepository) RequestUsecase {
	return &requestUsecase{
		requestRepo:  requestRepo,
		stepRepo:     stepRepo,
		workflowRepo: workflowRepo,
	}
}

func (uc *requestUsecase) CreateRequest(workflowID int, amount float64) (model.Request, error) {
	if amount <= 0 {
		return model.Request{}, ErrInvalidAmount
	}

	_, err := uc.workflowRepo.FindByID(workflowID)
	if err != nil {
		return model.Request{}, err
	}

	_, err = uc.stepRepo.FindByLevelAndWorkflowID(1, workflowID)
	if err != nil {
		return model.Request{}, err
	}

	existingRequest, err := uc.requestRepo.FindPendingByWorkflowID(workflowID)
	if err == nil && existingRequest.ID != 0 {
		existingRequest.Amount += amount

		accumulatedMinAmount, err := uc.getAccumulatedMinAmount(workflowID, existingRequest.CurrentStep)
		if err != nil {
			return model.Request{}, err
		}

		if existingRequest.Amount >= accumulatedMinAmount {
			nextStep, err := uc.stepRepo.FindByLevelAndWorkflowID(existingRequest.CurrentStep+1, workflowID)
			if err == nil && nextStep.ID != 0 {
				existingRequest.CurrentStep += 1
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return model.Request{}, err
			} else {
				existingRequest.Status = "APPROVED"
			}
		}

		if err := uc.requestRepo.Update(&existingRequest); err != nil {
			return model.Request{}, err
		}
		return existingRequest, nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Request{}, err
	}

	request := model.Request{
		WorkflowID:  uint(workflowID),
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      amount,
	}

	accumulatedMinAmount, err := uc.getAccumulatedMinAmount(workflowID, 1)
	if err != nil {
		return model.Request{}, err
	}

	if amount >= accumulatedMinAmount {
		nextStep, err := uc.stepRepo.FindByLevelAndWorkflowID(2, workflowID)
		if err == nil && nextStep.ID != 0 {
			request.CurrentStep = 2
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Request{}, err
		} else {
			request.Status = "APPROVED"
		}
	}

	if err := uc.requestRepo.Create(&request); err != nil {
		return model.Request{}, err
	}

	return request, nil
}

func (uc *requestUsecase) GetRequestByID(id int) (model.Request, error) {
	return uc.requestRepo.FindByID(id)
}

func (uc *requestUsecase) FindAllRequestsWithPagination(page, pageSize int, search, status string) ([]model.Request, int64, error) {
	offset := (page - 1) * pageSize
	return uc.requestRepo.FindAllWithPagination(offset, pageSize, search, status)
}

func (uc *requestUsecase) ApproveRequest(id int) (model.Request, error) {
	var request model.Request

	tx := uc.requestRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	request, err := uc.requestRepo.FindByIDWithLock(tx, id)
	if err != nil {
		tx.Rollback()
		return request, err
	}

	if request.Status != "PENDING" {
		tx.Rollback()
		return request, ErrInvalidRequestState
	}

	step, err := uc.stepRepo.FindByLevelAndWorkflowIDTx(tx, request.CurrentStep, int(request.WorkflowID))
	if err != nil {
		tx.Rollback()
		return request, err
	}

	conditions, err := parseConditions(step.Conditions)
	if err != nil {
		tx.Rollback()
		return request, err
	}

	if conditions.ApprovalType == "API" {
		accumulatedMinAmount, err := uc.getAccumulatedMinAmountTx(tx, int(request.WorkflowID), request.CurrentStep)
		if err != nil {
			tx.Rollback()
			return request, err
		}

		if request.Amount < accumulatedMinAmount {
			tx.Rollback()
			return request, nil
		}
	}

	request.Status = "APPROVED"
	if err := uc.requestRepo.UpdateTx(tx, &request); err != nil {
		tx.Rollback()
		return request, err
	}

	if err := tx.Commit().Error; err != nil {
		return request, err
	}

	return request, nil
}

func (uc *requestUsecase) RejectRequest(id int) (model.Request, error) {
	request, err := uc.requestRepo.FindByID(id)
	if err != nil {
		return request, err
	}

	if request.Status != "PENDING" {
		return request, ErrInvalidRequestState
	}

	request.Status = "REJECTED"
	if err := uc.requestRepo.Update(&request); err != nil {
		return request, err
	}

	return request, nil
}

func (uc *requestUsecase) getAccumulatedMinAmount(workflowID int, currentLevel uint) (float64, error) {
	var total float64 = 0

	for level := uint(1); level <= currentLevel; level++ {
		step, err := uc.stepRepo.FindByLevelAndWorkflowID(level, workflowID)
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

func (uc *requestUsecase) getAccumulatedMinAmountTx(tx *gorm.DB, workflowID int, currentLevel uint) (float64, error) {
	var total float64 = 0

	for level := uint(1); level <= currentLevel; level++ {
		step, err := uc.stepRepo.FindByLevelAndWorkflowIDTx(tx, level, workflowID)
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
