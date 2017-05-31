//go:generate goagen bootstrap -d github.com/hairyhenderson/restdemo/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/hairyhenderson/restdemo/app"
)

func main() {
	// Create service
	service := goa.New("WorstSocialNetwork")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "health" controller
	c := NewHealthController(service)
	app.MountHealthController(service, c)
	// Mount "post" controller
	c2 := NewPostController(service)
	app.MountPostController(service, c2)
	// Mount "public" controller
	c3 := NewPublicController(service)
	app.MountPublicController(service, c3)

	// Start service
	if err := service.ListenAndServe(":8000"); err != nil {
		service.LogError("startup", "err", err)
	}
}
