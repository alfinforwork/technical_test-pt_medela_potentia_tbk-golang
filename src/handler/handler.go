package handler

import (
	"technical-test/src/middleware"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	ws := service.NewWorkflowService(db)
	ss := service.NewStepService(db, ws)
	rs := service.NewRequestService(db, ws, ss)
	as := service.NewAuthService(db)

	v1 := app.Group("/v1")
	SetupAuthHandlers(v1, db, as)

	protected := v1.Group("/", middleware.JWTProtected())
	SetupWorkflowHandlers(protected, db, ws)
	SetupStepHandlers(protected, db, ws, ss)
	SetupRequestHandlers(protected, db, ws, ss, rs)

}
