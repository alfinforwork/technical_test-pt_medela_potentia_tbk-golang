# Unit Tests Documentation

## Overview
This project includes comprehensive unit tests for the approval workflow system, covering:
- **Approval Logic Tests** - Testing the request approval workflow
- **Validation Tests** - Testing input validation and edge cases

## Test Files

### 1. `src/service/request_service_test.go`
Tests for the request approval service, including:

#### Approval Logic Tests
- `TestApproveRequest_APIApprovalType` - Tests automatic approval when amount meets API approval threshold
- `TestApproveRequest_ManualApprovalType` - Tests manual approval type (approves regardless of amount)
- `TestRejectRequest` - Tests request rejection workflow
- `TestRequestAccumulation` - Tests accumulation of requests until minimum amount is reached

#### Validation Tests
- `TestCreateRequest_ValidAmount` - Tests valid request creation
- `TestCreateRequest_InvalidAmount` - Tests rejection of negative amounts
- `TestCreateRequest_ZeroAmount` - Tests rejection of zero amounts
- `TestCreateRequest_BelowMinimum` - Tests request pending state when below minimum
- `TestCreateRequest_MultiLevelWorkflow` - Tests multi-step approval workflows
- `TestCreateRequest_NonExistentWorkflow` - Tests error handling for invalid workflow
- `TestApproveRequest_InvalidState` - Tests that already-approved requests cannot be re-approved
- `TestRejectRequest_InvalidState` - Tests that already-rejected requests cannot be re-rejected
- `TestGetRequestByID` - Tests retrieving request by ID
- `TestGetRequestByID_NotFound` - Tests error handling for non-existent requests
- `TestAccumulatedMinAmount` - Tests calculation of accumulated minimum amounts across approval levels

### 2. `src/service/step_service_test.go`
Tests for the step/approval level service:

#### CRUD Tests
- `TestCreateStep_Valid` - Tests creating a step with valid data
- `TestCreateStep_IncrementalLevels` - Tests that step levels increment properly
- `TestGetStepByID` - Tests retrieving step by ID
- `TestUpdateStep` - Tests updating step details
- `TestDeleteStep` - Tests deleting a step
- `TestFindStepsByWorkflowID` - Tests finding all steps for a workflow
- `TestFindStepByLevelAndWorkflowID` - Tests finding specific step by level

#### Validation Tests
- `TestCreateStep_NonExistentWorkflow` - Tests error handling for invalid workflow
- `TestGetStepByID_NotFound` - Tests error handling for non-existent steps
- `TestUpdateStep_NotFound` - Tests error handling when updating non-existent step
- `TestDeleteStep_NonExistent` - Tests handling of deletion for non-existent steps
- `TestFindStepByLevelAndWorkflowID_NotFound` - Tests handling of non-existent step level
- `TestFindStepsByWorkflowID_Empty` - Tests handling of workflow with no steps
- `TestGetNextLevelForWorkflow` - Tests calculation of next step level

## Running Tests

### Run all service tests:
```bash
go test ./src/service -v
```

### Run with coverage:
```bash
go test ./src/service -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run specific test:
```bash
go test ./src/service -v -run TestCreateRequest_ValidAmount
```

### Run with timeout:
```bash
go test ./src/service -v -timeout 30s
```

## Test Setup
All tests use:
- **In-memory SQLite database** for isolation and speed
- **Test suite framework** (testify/suite) for organized setup/teardown
- **Unique workflow names** per test to avoid constraint violations
- **Assertion library** (testify/assert) for readable assertions

## Key Test Patterns

### 1. Approval Logic Testing
```go
// Test that API approval type validates amount
func (suite *RequestServiceTestSuite) TestApproveRequest_APIApprovalType() {
    // Setup workflow and steps
    // Create request
    // Approve and verify status
}
```

### 2. Validation Testing
```go
// Test that invalid amounts are rejected
func (suite *RequestServiceTestSuite) TestCreateRequest_InvalidAmount() {
    workflow := suite.createTestWorkflow()
    _, err := suite.requestService.CreateRequest(int(workflow.ID), -50)
    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), ErrInvalidAmount, err)
}
```

### 3. Multi-step Workflow Testing
```go
// Test request progression through multiple approval levels
func (suite *RequestServiceTestSuite) TestCreateRequest_MultiLevelWorkflow() {
    // Create workflow with multiple steps
    // Test that request moves through correct levels
}
```

## Test Coverage

- **RequestService**: 15 test methods covering all major functions
- **StepService**: 15 test methods covering CRUD and validation
- **Total**: 30+ test cases

## Dependencies
- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/suite` - Test suites
- `gorm.io/driver/sqlite` - SQLite driver for tests
