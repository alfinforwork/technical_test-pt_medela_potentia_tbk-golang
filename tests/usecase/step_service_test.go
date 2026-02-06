package usecase

import (
	"technical-test/src/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type StepUsecaseTestSuite struct {
	BaseTestSuite
	stepUsecase     usecase.StepUsecase
	workflowUsecase usecase.WorkflowUsecase
}

func (suite *StepUsecaseTestSuite) SetupTest() {
	err := suite.InitializeDB("step_usecase")
	suite.NoError(err)

	suite.stepUsecase, suite.workflowUsecase = suite.CreateStepUsecaseWithDeps()
}

func (suite *StepUsecaseTestSuite) TestCreateStep_Valid() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step, err := suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), workflow.ID, step.WorkflowID)
	assert.Equal(suite.T(), "Manager", step.Actor)
	assert.Equal(suite.T(), uint(1), step.Level)
}

func (suite *StepUsecaseTestSuite) TestCreateStep_NonExistentWorkflow() {
	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	_, err := suite.stepUsecase.CreateStep(9999, "Manager", conditions)

	assert.Error(suite.T(), err)
}

func (suite *StepUsecaseTestSuite) TestCreateStep_IncrementalLevels() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	step1, err := suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), step1.Level)

	step2, err := suite.stepUsecase.CreateStep(int(workflow.ID), "Director", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(2), step2.Level)

	step3, err := suite.stepUsecase.CreateStep(int(workflow.ID), "CEO", conditions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(3), step3.Level)
}

func (suite *StepUsecaseTestSuite) TestGetNextLevelForWorkflow() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)
	suite.stepUsecase.CreateStep(int(workflow.ID), "Director", conditions)

	nextLevel, err := suite.stepUsecase.GetNextLevelForWorkflow(int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(3), nextLevel)
}

func (suite *StepUsecaseTestSuite) TestFindStepsByWorkflowID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))

	suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)
	suite.stepUsecase.CreateStep(int(workflow.ID), "Director", conditions)

	steps, err := suite.stepUsecase.FindStepsByWorkflowID(int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), steps, 2)
}

func (suite *StepUsecaseTestSuite) TestFindStepsByWorkflowID_Empty() {
	steps, err := suite.stepUsecase.FindStepsByWorkflowID(9999)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), steps, 0)
}

func (suite *StepUsecaseTestSuite) TestFindStepByLevelAndWorkflowID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)

	step, err := suite.stepUsecase.FindStepByLevelAndWorkflowID(1, int(workflow.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdStep.ID, step.ID)
	assert.Equal(suite.T(), "Manager", step.Actor)
}

func (suite *StepUsecaseTestSuite) TestFindStepByLevelAndWorkflowID_NotFound() {
	workflow := suite.CreateTestWorkflow()

	_, err := suite.stepUsecase.FindStepByLevelAndWorkflowID(99, int(workflow.ID))

	assert.Error(suite.T(), err)
}

func (suite *StepUsecaseTestSuite) TestGetStepByID() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)

	step, err := suite.stepUsecase.GetStepByID(int(createdStep.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdStep.ID, step.ID)
}

func (suite *StepUsecaseTestSuite) TestGetStepByID_NotFound() {
	_, err := suite.stepUsecase.GetStepByID(9999)

	assert.Error(suite.T(), err)
}

func (suite *StepUsecaseTestSuite) TestUpdateStep() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	createdStep, _ := suite.stepUsecase.CreateStep(int(workflow.ID), "Manager", conditions)

	newConditions := datatypes.JSON([]byte(`{"min_amount": 200}`))
	updatedStep, err := suite.stepUsecase.UpdateStep(int(createdStep.ID), 1, "Director", newConditions)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Director", updatedStep.Actor)
	assert.Equal(suite.T(), newConditions, updatedStep.Conditions)
}

func (suite *StepUsecaseTestSuite) TestUpdateStep_NotFound() {
	conditions := datatypes.JSON([]byte(`{"min_amount": 100}`))
	_, err := suite.stepUsecase.UpdateStep(9999, 1, "Manager", conditions)

	assert.Error(suite.T(), err)
}

func TestStepUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(StepUsecaseTestSuite))
}
