package main

import (
	"github.com/goadesign/goa"
	"github.com/hairyhenderson/restdemo/app"
)

// PostController implements the post resource.
type PostController struct {
	*goa.Controller
}

// NewPostController creates a post controller.
func NewPostController(service *goa.Service) *PostController {
	return &PostController{Controller: service.NewController("PostController")}
}

// Show runs the show action.
func (c *PostController) Show(ctx *app.ShowPostContext) error {
	// PostController_Show: start_implement
	if ctx.PostID == 0 {
		return ctx.NotFound()
	}

	res := &app.WsnPost{
		ID:     ctx.PostID,
		Title:  "test!",
		Body:   "test test",
		UserID: 1,
	}
	return ctx.OK(res)
	// PostController_Show: end_implement
}
