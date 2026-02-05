package router

import (
	"technical-test/src/controller"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupWorkflowRoutes(router fiber.Router, db *gorm.DB, ws *service.WorkflowService) {
	workflowGroup := router.Group("/workflows")
	workflowController := controller.NewWorkflowController(*ws)

	workflowGroup.Post("/", workflowController.CreateWorkflow)
	workflowGroup.Get("/", workflowController.FindAllWorkflows)
	workflowGroup.Get("/:workflowId", workflowController.GetWorkflowByID)
}
