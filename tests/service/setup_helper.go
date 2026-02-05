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

func (suite *BaseTestSuite) InitializeDB(suiteName string) error {
	if suite.DB != nil {
		suite.TestCounter++
		return nil
	}

	suite.TestCounter++

	if suiteName == "" {
		suiteName = "default"
	}

	// Use named shared in-memory database to avoid "no such table" across pooled connections.
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", suiteName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// Keep a single connection to ensure the in-memory schema is preserved.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

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
