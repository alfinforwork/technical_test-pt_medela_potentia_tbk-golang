package service

import (
	"fmt"
	"technical-test/src/model"
	svc "technical-test/src/service"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BaseTestSuite struct {
	suite.Suite
	DB          *gorm.DB
	TestCounter int
}

func (suite *BaseTestSuite) InitializeDB() error {
	if suite.DB != nil {
		suite.TestCounter++
		return nil
	}

	suite.TestCounter++
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Workflow{},
		&model.Step{},
		&model.Request{},
	)
	if err != nil {
		return err
	}

	suite.DB = db
	return nil
}

func (suite *BaseTestSuite) CreateTestWorkflow() model.Workflow {
	workflow := model.Workflow{Name: fmt.Sprintf("Test Workflow %d", suite.TestCounter)}
	suite.DB.Create(&workflow)
	return workflow
}

func (suite *BaseTestSuite) CreateRequestServiceWithDeps() (*svc.RequestService, *svc.WorkflowService, *svc.StepService) {
	workflowService := svc.NewWorkflowService(suite.DB)
	stepService := svc.NewStepService(suite.DB, workflowService)
	requestService := svc.NewRequestService(suite.DB, workflowService, stepService)
	return requestService, workflowService, stepService
}

func (suite *BaseTestSuite) CreateStepServiceWithDeps() (*svc.StepService, *svc.WorkflowService) {
	workflowService := svc.NewWorkflowService(suite.DB)
	stepService := svc.NewStepService(suite.DB, workflowService)
	return stepService, workflowService
}
