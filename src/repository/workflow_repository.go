package repository

import (
	"technical-test/src/model"

	"gorm.io/gorm"
)

type WorkflowRepository interface {
	FindByName(name string) (model.Workflow, error)
	Create(workflow *model.Workflow) error
	FindAll() ([]model.Workflow, error)
	FindAllWithPagination(offset, limit int, search string) ([]model.Workflow, int64, error)
	FindByID(id int) (model.Workflow, error)
}

type workflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) FindByName(name string) (model.Workflow, error) {
	var workflow model.Workflow
	err := r.db.Where("name = ?", name).First(&workflow).Error
	return workflow, err
}

func (r *workflowRepository) Create(workflow *model.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *workflowRepository) FindAll() ([]model.Workflow, error) {
	var workflows []model.Workflow
	err := r.db.Select("id", "name", "created_at").Find(&workflows).Error
	return workflows, err
}

func (r *workflowRepository) FindAllWithPagination(offset, limit int, search string) ([]model.Workflow, int64, error) {
	var workflows []model.Workflow
	var total int64

	query := r.db.Model(&model.Workflow{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return workflows, 0, err
	}

	err := query.Select("id", "name", "created_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&workflows).Error

	return workflows, total, err
}

func (r *workflowRepository) FindByID(id int) (model.Workflow, error) {
	var workflow model.Workflow
	err := r.db.First(&workflow, id).Error
	return workflow, err
}
