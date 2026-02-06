package routes

import (
	"technical-test/src/handler"
	"technical-test/src/middleware"
	"technical-test/src/repository"
	"technical-test/src/usecase"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	stepRepo := repository.NewStepRepository(db)
	requestRepo := repository.NewRequestRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)
	workflowUsecase := usecase.NewWorkflowUsecase(workflowRepo)
	stepUsecase := usecase.NewStepUsecase(stepRepo, workflowRepo)
	requestUsecase := usecase.NewRequestUsecase(requestRepo, stepRepo, workflowRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	workflowHandler := handler.NewWorkflowHandler(workflowUsecase)
	stepHandler := handler.NewStepHandler(stepUsecase, workflowUsecase)
	requestHandler := handler.NewRequestHandler(requestUsecase, workflowUsecase)

	// Setup routes
	v1 := app.Group("/v1")

	// Auth routes (public)
	authGroup := v1.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// Protected routes
	protected := v1.Group("/", middleware.JWTProtected())

	// Workflow routes
	workflowGroup := protected.Group("/workflows")
	workflowGroup.Post("/", workflowHandler.CreateWorkflow)
	workflowGroup.Get("/", workflowHandler.FindAllWorkflows)
	workflowGroup.Get("/:workflowId", workflowHandler.GetWorkflowByID)

	// Step routes
	workflowGroup.Post("/:workflowId/steps", stepHandler.CreateStep)
	workflowGroup.Get("/:workflowId/steps", stepHandler.FindStepsByWorkflowID)

	// Request routes
	requestGroup := protected.Group("/requests")
	requestGroup.Post("/", requestHandler.CreateRequest)
	requestGroup.Get("/", requestHandler.FindAllRequests)
	requestGroup.Get("/:requestId", requestHandler.GetRequestByID)
	requestGroup.Post("/:requestId/approve", requestHandler.ApproveRequest)
	requestGroup.Post("/:requestId/reject", requestHandler.RejectRequest)
}
