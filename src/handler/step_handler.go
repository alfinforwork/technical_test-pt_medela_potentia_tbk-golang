package handler

import (
	"technical-test/src/controller"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupStepHandlers(router fiber.Router, db *gorm.DB, ws *service.WorkflowService, ss *service.StepService) {
	stepGroup := router.Group("/workflows/:workflowId/steps")
	stepController := controller.NewStepController(*ws, *ss)

	stepGroup.Post("/", stepController.CreateStep)
	stepGroup.Get("/", stepController.FindStepsByWorkflowID)
	stepGroup.Get("/:stepId", stepController.FindStepByID)
	stepGroup.Put("/:stepId", stepController.UpdateStep)
	stepGroup.Delete("/:stepId", stepController.DeleteStep)
}
