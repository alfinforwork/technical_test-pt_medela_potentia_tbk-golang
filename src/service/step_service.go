package service

import (
	"technical-test/src/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StepService struct {
	db              *gorm.DB
	workflowService *WorkflowService
}

func NewStepService(db *gorm.DB, ws *WorkflowService) *StepService {
	return &StepService{
		db:              db,
		workflowService: ws,
	}
}

func (ss *StepService) CreateStep(workflowID int, Actor string, Conditions datatypes.JSON) (model.Step, error) {
	workflow, err := ss.workflowService.GetWorkflowByID(workflowID)
	if err != nil {
		return model.Step{}, err
	}

	nextLevel, err := ss.GetNextLevelForWorkflow(int(workflow.ID))
	if err != nil {
		return model.Step{}, err
	}

	step := model.Step{
		WorkflowID: uint(workflowID),
		Level:      nextLevel,
		Actor:      Actor,
		Conditions: Conditions,
	}
	result := ss.db.Create(&step)
	return step, result.Error
}

func (ss *StepService) GetNextLevelForWorkflow(workflowID int) (uint, error) {
	var maxLevel uint
	ss.db.Model(&model.Step{}).
		Where("workflow_id = ?", workflowID).
		Select("COALESCE(MAX(level), 0)").
		Row().
		Scan(&maxLevel)
	return maxLevel + 1, nil
}

func (ss *StepService) FindStepsByWorkflowID(workflowID int) ([]model.Step, error) {
	var steps []model.Step
	result := ss.db.Where("workflow_id = ?", workflowID).Find(&steps)
	return steps, result.Error
}

func (ss *StepService) FindStepsByWorkflowIDWithPagination(workflowID int, page, pageSize int, search string) ([]model.Step, int64, error) {
	var steps []model.Step
	var total int64

	query := ss.db.Where("workflow_id = ?", workflowID)
	if search != "" {
		query = query.Where("actor LIKE ?", "%"+search+"%")
	}

	if err := query.Model(&model.Step{}).Count(&total).Error; err != nil {
		return steps, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.
		Order("level ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&steps)

	return steps, total, result.Error
}

func (ss *StepService) FindStepByLevelAndWorkflowID(level uint, workflowID int) (model.Step, error) {
	var step model.Step
	result := ss.db.Where("level = ? AND workflow_id = ?", level, workflowID).First(&step)
	return step, result.Error
}

func (ss *StepService) GetStepByID(id int) (model.Step, error) {
	var step model.Step
	result := ss.db.First(&step, id)
	return step, result.Error
}

func (ss *StepService) UpdateStep(id int, Level uint, Actor string, Conditions datatypes.JSON) (model.Step, error) {
	var step model.Step
	if err := ss.db.First(&step, id).Error; err != nil {
		return step, err
	}
	step.Level = Level
	step.Actor = Actor
	step.Conditions = Conditions
	result := ss.db.Save(&step)
	return step, result.Error
}

func (ss *StepService) DeleteStep(id int) error {
	result := ss.db.Delete(&model.Step{}, id)
	return result.Error
}
