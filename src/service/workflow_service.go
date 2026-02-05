package service

import (
	"errors"
	"technical-test/src/model"

	"gorm.io/gorm"
)

type WorkflowService struct {
	db *gorm.DB
}

var ErrWorkflowNameExists = errors.New("workflow name already exists")

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

func (ws *WorkflowService) CreateWorkflow(Name string) (model.Workflow, error) {
	var existing model.Workflow
	err := ws.db.Where("name = ?", Name).First(&existing).Error
	if err == nil {
		return model.Workflow{}, errors.New("Workflow name already exists")
	}

	workflow := model.Workflow{Name: Name}
	result := ws.db.Create(&workflow)
	return workflow, result.Error
}

func (ws *WorkflowService) FindAllWorkflows() ([]model.Workflow, error) {
	var workflows []model.Workflow
	result := ws.db.Select("id", "name", "created_at").Find(&workflows)
	return workflows, result.Error
}

func (ws *WorkflowService) GetWorkflowByID(id int) (model.Workflow, error) {
	var workflow model.Workflow
	result := ws.db.First(&workflow, id)
	return workflow, result.Error
}
