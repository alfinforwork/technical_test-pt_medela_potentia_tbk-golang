package repository

import (
	"technical-test/src/model"

	"gorm.io/gorm"
)

type StepRepository interface {
	Create(step *model.Step) error
	FindByWorkflowID(workflowID int) ([]model.Step, error)
	FindByWorkflowIDWithPagination(workflowID int, offset, limit int, search string) ([]model.Step, int64, error)
	FindByLevelAndWorkflowID(level uint, workflowID int) (model.Step, error)
	FindByLevelAndWorkflowIDTx(tx *gorm.DB, level uint, workflowID int) (model.Step, error)
	FindByID(id int) (model.Step, error)
	GetMaxLevel(workflowID int) (uint, error)
	Update(step *model.Step) error
	Delete(id int) error
}

type stepRepository struct {
	db *gorm.DB
}

func NewStepRepository(db *gorm.DB) StepRepository {
	return &stepRepository{db: db}
}

func (r *stepRepository) Create(step *model.Step) error {
	return r.db.Create(step).Error
}

func (r *stepRepository) FindByWorkflowID(workflowID int) ([]model.Step, error) {
	var steps []model.Step
	err := r.db.Where("workflow_id = ?", workflowID).Find(&steps).Error
	return steps, err
}

func (r *stepRepository) FindByWorkflowIDWithPagination(workflowID int, offset, limit int, search string) ([]model.Step, int64, error) {
	var steps []model.Step
	var total int64

	query := r.db.Model(&model.Step{}).Where("workflow_id = ?", workflowID)
	if search != "" {
		query = query.Where("actor LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return steps, 0, err
	}

	err := query.
		Order("level ASC").
		Offset(offset).
		Limit(limit).
		Find(&steps).Error

	return steps, total, err
}

func (r *stepRepository) FindByLevelAndWorkflowID(level uint, workflowID int) (model.Step, error) {
	var step model.Step
	err := r.db.Where("level = ? AND workflow_id = ?", level, workflowID).First(&step).Error
	return step, err
}

func (r *stepRepository) FindByLevelAndWorkflowIDTx(tx *gorm.DB, level uint, workflowID int) (model.Step, error) {
	var step model.Step
	err := tx.Where("level = ? AND workflow_id = ?", level, workflowID).First(&step).Error
	return step, err
}

func (r *stepRepository) FindByID(id int) (model.Step, error) {
	var step model.Step
	err := r.db.First(&step, id).Error
	return step, err
}

func (r *stepRepository) GetMaxLevel(workflowID int) (uint, error) {
	var maxLevel uint
	r.db.Model(&model.Step{}).
		Where("workflow_id = ?", workflowID).
		Select("COALESCE(MAX(level), 0)").
		Row().
		Scan(&maxLevel)
	return maxLevel, nil
}

func (r *stepRepository) Update(step *model.Step) error {
	return r.db.Save(step).Error
}

func (r *stepRepository) Delete(id int) error {
	return r.db.Delete(&model.Step{}, id).Error
}
