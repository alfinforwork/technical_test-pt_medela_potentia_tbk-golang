package service

import (
	svc "technical-test/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type StepServiceTestSuite struct {
	BaseTestSuite
	stepService     *svc.StepService
	workflowService *svc.WorkflowService
}

func (suite *StepServiceTestSuite) SetupTest() {
	err := suite.InitializeDB()
	suite.NoError(err)

	suite.stepService, suite.workflowService = suite.CreateStepServiceWithDeps()
}

func (suite *StepServiceTestSuite) TestCreateStep_Valid() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step, err := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), workflow.ID, step.WorkflowID)
	assert.Equal(suite.T(), "Manager", step.Actor)
	assert.Equal(suite.T(), uint(1), step.Level)
}

func (suite *StepServiceTestSuite) TestCreateStep_NonExistentWorkflow() {
	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	_, err := suite.stepService.CreateStep(9999, "Manager", conditions)

	assert.Error(suite.T(), err)
}

func (suite *StepServiceTestSuite) TestCreateStep_IncrementalLevels() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	step1, err := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), step1.Level)

	step2, err := suite.stepService.CreateStep(int(workflow.ID), "Director", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(2), step2.Level)

	step3, err := suite.stepService.CreateStep(int(workflow.ID), "CEO", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(3), step3.Level)
}

func (suite *StepServiceTestSuite) TestGetNextLevelForWorkflow() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)
	suite.stepService.CreateStep(int(workflow.ID), "Director", conditions)

	nextLevel, err := suite.stepService.GetNextLevelForWorkflow(int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(3), nextLevel)
}

func (suite *StepServiceTestSuite) TestFindStepsByWorkflowID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)
	suite.stepService.CreateStep(int(workflow.ID), "Director", conditions)

	steps, err := suite.stepService.FindStepsByWorkflowID(int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), steps, 2)
}

func (suite *StepServiceTestSuite) TestFindStepsByWorkflowID_Empty() {
	steps, err := suite.stepService.FindStepsByWorkflowID(9999)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), steps, 0)
}

func (suite *StepServiceTestSuite) TestFindStepByLevelAndWorkflowID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)

	step, err := suite.stepService.FindStepByLevelAndWorkflowID(1, int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdStep.ID, step.ID)
	assert.Equal(suite.T(), "Manager", step.Actor)
}

func (suite *StepServiceTestSuite) TestFindStepByLevelAndWorkflowID_NotFound() {
	workflow := suite.CreateTestWorkflow()

	_, err := suite.stepService.FindStepByLevelAndWorkflowID(99, int(workflow.ID))

	assert.Error(suite.T(), err)
}

func (suite *StepServiceTestSuite) TestGetStepByID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)

	step, err := suite.stepService.GetStepByID(int(createdStep.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdStep.ID, step.ID)
}

func (suite *StepServiceTestSuite) TestGetStepByID_NotFound() {
	_, err := suite.stepService.GetStepByID(9999)

	assert.Error(suite.T(), err)
}

func (suite *StepServiceTestSuite) TestUpdateStep() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)

	newConditions := datatypes.JSON([]byte(`{"min_amount": 200}`))
	updatedStep, err := suite.stepService.UpdateStep(int(createdStep.ID), 1, "Director", newConditions)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Director", updatedStep.Actor)
	assert.Equal(suite.T(), newConditions, updatedStep.Conditions)
}

func (suite *StepServiceTestSuite) TestUpdateStep_NotFound() {
	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	_, err := suite.stepService.UpdateStep(9999, 1, "Manager", conditions)

	assert.Error(suite.T(), err)
}

func (suite *StepServiceTestSuite) TestDeleteStep() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepService.CreateStep(int(workflow.ID), "Manager", conditions)

	err := suite.stepService.DeleteStep(int(createdStep.ID))

	assert.NoError(suite.T(), err)

	_, err = suite.stepService.GetStepByID(int(createdStep.ID))
	assert.Error(suite.T(), err)
}

func (suite *StepServiceTestSuite) TestDeleteStep_NonExistent() {
	err := suite.stepService.DeleteStep(9999)

	assert.NoError(suite.T(), err)
}

func TestStepServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StepServiceTestSuite))
}
