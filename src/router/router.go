package router

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
	SetupAuthRoutes(v1, db, as)

	protected := v1.Group("/", middleware.JWTProtected())
	SetupWorkflowRoutes(protected, db, ws)
	SetupStepRoutes(protected, db, ws, ss)
	SetupRequestRoutes(protected, db, ws, ss, rs)

}
