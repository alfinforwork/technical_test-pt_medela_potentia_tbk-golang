package usecase

import (
	"errors"
	"technical-test/src/model"
	"technical-test/src/repository"

	"gorm.io/gorm"
)

type WorkflowUsecase interface {
	CreateWorkflow(name string) (model.Workflow, error)
	FindAllWorkflows() ([]model.Workflow, error)
	FindAllWorkflowsWithPagination(page, pageSize int, search string) ([]model.Workflow, int64, error)
	GetWorkflowByID(id int) (model.Workflow, error)
}

type workflowUsecase struct {
	workflowRepo repository.WorkflowRepository
}

var ErrWorkflowNameExists = errors.New("workflow name already exists")

func NewWorkflowUsecase(workflowRepo repository.WorkflowRepository) WorkflowUsecase {
	return &workflowUsecase{
		workflowRepo: workflowRepo,
	}
}

func (uc *workflowUsecase) CreateWorkflow(name string) (model.Workflow, error) {
	existing, err := uc.workflowRepo.FindByName(name)
	if err == nil && existing.ID != 0 {
		return model.Workflow{}, ErrWorkflowNameExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Workflow{}, err
	}

	workflow := model.Workflow{Name: name}
	if err := uc.workflowRepo.Create(&workflow); err != nil {
		return model.Workflow{}, err
	}

	return workflow, nil
}

func (uc *workflowUsecase) FindAllWorkflows() ([]model.Workflow, error) {
	return uc.workflowRepo.FindAll()
}

func (uc *workflowUsecase) FindAllWorkflowsWithPagination(page, pageSize int, search string) ([]model.Workflow, int64, error) {
	offset := (page - 1) * pageSize
	return uc.workflowRepo.FindAllWithPagination(offset, pageSize, search)
}

func (uc *workflowUsecase) GetWorkflowByID(id int) (model.Workflow, error) {
	return uc.workflowRepo.FindByID(id)
}
