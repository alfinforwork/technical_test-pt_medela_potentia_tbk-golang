package router

import (
	"technical-test/src/controller"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRequestRoutes(router fiber.Router, db *gorm.DB, ws *service.WorkflowService, ss *service.StepService, rs *service.RequestService) {
	requestGroup := router.Group("/requests")
	requestController := controller.NewRequestController(*rs, *ws)

	requestGroup.Post("/", requestController.CreateRequest)
	requestGroup.Get("/", requestController.FindAllRequests)
	requestGroup.Get("/:requestId", requestController.GetRequestByID)
	requestGroup.Post("/:requestId/approve", requestController.ApproveRequest)
	requestGroup.Post("/:requestId/reject", requestController.RejectRequest)
}
