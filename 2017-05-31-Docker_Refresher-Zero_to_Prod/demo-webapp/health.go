package main

import (
	"os"

	"github.com/goadesign/goa"
	"github.com/hairyhenderson/restdemo/app"
)

// HealthController implements the health resource.
type HealthController struct {
	*goa.Controller
}

// NewHealthController creates a health controller.
func NewHealthController(service *goa.Service) *HealthController {
	return &HealthController{Controller: service.NewController("HealthController")}
}

// Show runs the show action.
func (c *HealthController) Show(ctx *app.ShowHealthContext) error {
	// HealthController_Show: start_implement

	// Put your logic here
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	ver := os.Getenv("VER")
	res := &app.Health{hostname, ver}

	// HealthController_Show: end_implement
	return ctx.OK(res)
}
