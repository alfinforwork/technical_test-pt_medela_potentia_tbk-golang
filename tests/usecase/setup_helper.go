package usecase

import (
	"fmt"
	"technical-test/src/model"
	"technical-test/src/repository"
	"technical-test/src/usecase"

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

func (suite *BaseTestSuite) CreateRequestUsecaseWithDeps() (usecase.RequestUsecase, usecase.WorkflowUsecase, usecase.StepUsecase) {
	workflowRepo := repository.NewWorkflowRepository(suite.DB)
	stepRepo := repository.NewStepRepository(suite.DB)
	requestRepo := repository.NewRequestRepository(suite.DB)

	workflowUsecase := usecase.NewWorkflowUsecase(workflowRepo)
	stepUsecase := usecase.NewStepUsecase(stepRepo, workflowRepo)
	requestUsecase := usecase.NewRequestUsecase(requestRepo, stepRepo, workflowRepo)
	return requestUsecase, workflowUsecase, stepUsecase
}

func (suite *BaseTestSuite) CreateStepUsecaseWithDeps() (usecase.StepUsecase, usecase.WorkflowUsecase) {
	workflowRepo := repository.NewWorkflowRepository(suite.DB)
	stepRepo := repository.NewStepRepository(suite.DB)

	workflowUsecase := usecase.NewWorkflowUsecase(workflowRepo)
	stepUsecase := usecase.NewStepUsecase(stepRepo, workflowRepo)
	return stepUsecase, workflowUsecase
}
