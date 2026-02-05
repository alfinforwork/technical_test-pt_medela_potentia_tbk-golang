package service

import (
	"technical-test/src/model"
	svc "technical-test/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type RequestServiceTestSuite struct {
	BaseTestSuite
	requestService  *svc.RequestService
	workflowService *svc.WorkflowService
	stepService     *svc.StepService
}

func (suite *RequestServiceTestSuite) SetupTest() {
	err := suite.InitializeDB()
	suite.NoError(err)

	suite.requestService, suite.workflowService, suite.stepService = suite.CreateRequestServiceWithDeps()
}

// Test CreateRequest with valid amount
func (suite *RequestServiceTestSuite) TestCreateRequest_ValidAmount() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions,
	}
	suite.DB.Create(&step)

	// Create request with valid amount
	request, err := suite.requestService.CreateRequest(int(workflow.ID), 150)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), workflow.ID, request.WorkflowID)
	assert.Equal(suite.T(), 150.0, request.Amount)
	assert.Equal(suite.T(), "APPROVED", request.Status)
}

// Test CreateRequest with invalid amount
func (suite *RequestServiceTestSuite) TestCreateRequest_InvalidAmount() {
	workflow := suite.CreateTestWorkflow()

	// Create request with invalid amount
	_, err := suite.requestService.CreateRequest(int(workflow.ID), -50)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), svc.ErrInvalidAmount, err)
}

// Test CreateRequest with zero amount
func (suite *RequestServiceTestSuite) TestCreateRequest_ZeroAmount() {
	workflow := suite.CreateTestWorkflow()

	_, err := suite.requestService.CreateRequest(int(workflow.ID), 0)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), svc.ErrInvalidAmount, err)
}

// Test CreateRequest with non-existent workflow
func (suite *RequestServiceTestSuite) TestCreateRequest_NonExistentWorkflow() {
	_, err := suite.requestService.CreateRequest(9999, 100)

	assert.Error(suite.T(), err)
}

// Test CreateRequest with amount below minimum requirement
func (suite *RequestServiceTestSuite) TestCreateRequest_BelowMinimum() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions,
	}
	suite.DB.Create(&step)

	// Create request with amount below minimum
	request, err := suite.requestService.CreateRequest(int(workflow.ID), 50)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), request.CurrentStep)
	assert.Equal(suite.T(), "PENDING", request.Status)
	assert.Equal(suite.T(), 50.0, request.Amount)
}

// Test CreateRequest with multi-level workflow
func (suite *RequestServiceTestSuite) TestCreateRequest_MultiLevelWorkflow() {
	workflow := suite.CreateTestWorkflow()

	// Create step 1
	conditions1 := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step1 := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions1,
	}
	suite.DB.Create(&step1)

	// Create step 2
	conditions2 := datatypes.JSON([]byte(`{"min_amount": 200, "approval_type": "MANUAL"}`))
	step2 := model.Step{
		WorkflowID: workflow.ID,
		Level:      2,
		Actor:      "Director",
		Conditions: conditions2,
	}
	suite.DB.Create(&step2)

	// Create request with amount meeting step 1 requirement but not step 2
	request, err := suite.requestService.CreateRequest(int(workflow.ID), 150)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(2), request.CurrentStep)
	assert.Equal(suite.T(), "PENDING", request.Status)
}

// Test ApproveRequest with API approval type
func (suite *RequestServiceTestSuite) TestApproveRequest_APIApprovalType() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions,
	}
	suite.DB.Create(&step)

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      150,
	}
	suite.DB.Create(&request)

	// Approve request
	approvedRequest, err := suite.requestService.ApproveRequest(int(request.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "APPROVED", approvedRequest.Status)
}

// Test ApproveRequest with MANUAL approval type
func (suite *RequestServiceTestSuite) TestApproveRequest_ManualApprovalType() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "MANUAL"}`))
	step := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions,
	}
	suite.DB.Create(&step)

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      50,
	}
	suite.DB.Create(&request)

	// Approve request (should approve regardless of amount for MANUAL type)
	approvedRequest, err := suite.requestService.ApproveRequest(int(request.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "APPROVED", approvedRequest.Status)
}

// Test ApproveRequest with invalid request state
func (suite *RequestServiceTestSuite) TestApproveRequest_InvalidState() {
	workflow := suite.CreateTestWorkflow()

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "APPROVED",
		Amount:      100,
	}
	suite.DB.Create(&request)

	// Try to approve already approved request
	_, err := suite.requestService.ApproveRequest(int(request.ID))

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), svc.ErrInvalidRequestState, err)
}

// Test RejectRequest
func (suite *RequestServiceTestSuite) TestRejectRequest() {
	workflow := suite.CreateTestWorkflow()

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      100,
	}
	suite.DB.Create(&request)

	// Reject request
	rejectedRequest, err := suite.requestService.RejectRequest(int(request.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "REJECTED", rejectedRequest.Status)
}

// Test RejectRequest with invalid state
func (suite *RequestServiceTestSuite) TestRejectRequest_InvalidState() {
	workflow := suite.CreateTestWorkflow()

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "REJECTED",
		Amount:      100,
	}
	suite.DB.Create(&request)

	// Try to reject already rejected request
	_, err := suite.requestService.RejectRequest(int(request.ID))

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), svc.ErrInvalidRequestState, err)
}

// Test GetRequestByID
func (suite *RequestServiceTestSuite) TestGetRequestByID() {
	workflow := suite.CreateTestWorkflow()

	request := model.Request{
		WorkflowID:  workflow.ID,
		CurrentStep: 1,
		Status:      "PENDING",
		Amount:      100,
	}
	suite.DB.Create(&request)

	// Get request
	fetchedRequest, err := suite.requestService.GetRequestByID(int(request.ID))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), request.ID, fetchedRequest.ID)
	assert.Equal(suite.T(), request.Amount, fetchedRequest.Amount)
}

// Test GetRequestByID with non-existent ID
func (suite *RequestServiceTestSuite) TestGetRequestByID_NotFound() {
	_, err := suite.requestService.GetRequestByID(9999)

	assert.Error(suite.T(), err)
}

// Test request accumulation
func (suite *RequestServiceTestSuite) TestRequestAccumulation() {
	workflow := suite.CreateTestWorkflow()

	conditions := datatypes.JSON([]byte(`{"min_amount": 100, "approval_type": "API"}`))
	step := model.Step{
		WorkflowID: workflow.ID,
		Level:      1,
		Actor:      "Manager",
		Conditions: conditions,
	}
	suite.DB.Create(&step)

	// Create first request
	request1, err := suite.requestService.CreateRequest(int(workflow.ID), 60)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "PENDING", request1.Status)

	// Create second request (should accumulate)
	request2, err := suite.requestService.CreateRequest(int(workflow.ID), 50)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "APPROVED", request2.Status)
	assert.Equal(suite.T(), 110.0, request2.Amount)
}

// Run the test suite
func TestRequestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RequestServiceTestSuite))
}
