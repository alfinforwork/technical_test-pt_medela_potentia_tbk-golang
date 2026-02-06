package usecase

import (
	"technical-test/src/model"
	"technical-test/src/repository"

	"gorm.io/datatypes"
)

type StepUsecase interface {
	CreateStep(workflowID int, actor string, conditions datatypes.JSON) (model.Step, error)
	GetNextLevelForWorkflow(workflowID int) (uint, error)
	FindStepsByWorkflowID(workflowID int) ([]model.Step, error)
	FindStepsByWorkflowIDWithPagination(workflowID int, page, pageSize int, search string) ([]model.Step, int64, error)
	FindStepByLevelAndWorkflowID(level uint, workflowID int) (model.Step, error)
	GetStepByID(id int) (model.Step, error)
	UpdateStep(id int, level uint, actor string, conditions datatypes.JSON) (model.Step, error)
}

type stepUsecase struct {
	stepRepo     repository.StepRepository
	workflowRepo repository.WorkflowRepository
}

func NewStepUsecase(stepRepo repository.StepRepository, workflowRepo repository.WorkflowRepository) StepUsecase {
	return &stepUsecase{
		stepRepo:     stepRepo,
		workflowRepo: workflowRepo,
	}
}

func (uc *stepUsecase) CreateStep(workflowID int, actor string, conditions datatypes.JSON) (model.Step, error) {
	workflow, err := uc.workflowRepo.FindByID(workflowID)
	if err != nil {
		return model.Step{}, err
	}

	nextLevel, err := uc.GetNextLevelForWorkflow(int(workflow.ID))
	if err != nil {
		return model.Step{}, err
	}

	step := model.Step{
		WorkflowID: uint(workflowID),
		Level:      nextLevel,
		Actor:      actor,
		Conditions: conditions,
	}

	if err := uc.stepRepo.Create(&step); err != nil {
		return model.Step{}, err
	}

	return step, nil
}

func (uc *stepUsecase) GetNextLevelForWorkflow(workflowID int) (uint, error) {
	maxLevel, err := uc.stepRepo.GetMaxLevel(workflowID)
	if err != nil {
		return 0, err
	}
	return maxLevel + 1, nil
}

func (uc *stepUsecase) FindStepsByWorkflowID(workflowID int) ([]model.Step, error) {
	return uc.stepRepo.FindByWorkflowID(workflowID)
}

func (uc *stepUsecase) FindStepsByWorkflowIDWithPagination(workflowID int, page, pageSize int, search string) ([]model.Step, int64, error) {
	offset := (page - 1) * pageSize
	return uc.stepRepo.FindByWorkflowIDWithPagination(workflowID, offset, pageSize, search)
}

func (uc *stepUsecase) FindStepByLevelAndWorkflowID(level uint, workflowID int) (model.Step, error) {
	return uc.stepRepo.FindByLevelAndWorkflowID(level, workflowID)
}

func (uc *stepUsecase) GetStepByID(id int) (model.Step, error) {
	return uc.stepRepo.FindByID(id)
}

func (uc *stepUsecase) UpdateStep(id int, level uint, actor string, conditions datatypes.JSON) (model.Step, error) {
	step, err := uc.stepRepo.FindByID(id)
	if err != nil {
		return step, err
	}

	step.Level = level
	step.Actor = actor
	step.Conditions = conditions

	if err := uc.stepRepo.Update(&step); err != nil {
		return step, err
	}

	return step, nil
}
