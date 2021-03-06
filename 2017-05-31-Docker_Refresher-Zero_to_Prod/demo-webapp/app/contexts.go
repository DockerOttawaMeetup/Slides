// Code generated by goagen v1.2.0, DO NOT EDIT.
//
// API "WorstSocialNetwork": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/hairyhenderson/restdemo/design
// --out=$(GOPATH)/src/github.com/hairyhenderson/restdemo
// --version=v1.2.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
	"strconv"
)

// ShowHealthContext provides the health show action context.
type ShowHealthContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewShowHealthContext parses the incoming request URL and body, performs validations and creates the
// context used by the health controller show action.
func NewShowHealthContext(ctx context.Context, r *http.Request, service *goa.Service) (*ShowHealthContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := ShowHealthContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowHealthContext) OK(r *Health) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.health+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// ShowPostContext provides the post show action context.
type ShowPostContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	PostID int
}

// NewShowPostContext parses the incoming request URL and body, performs validations and creates the
// context used by the post controller show action.
func NewShowPostContext(ctx context.Context, r *http.Request, service *goa.Service) (*ShowPostContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := ShowPostContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramPostID := req.Params["postID"]
	if len(paramPostID) > 0 {
		rawPostID := paramPostID[0]
		if postID, err2 := strconv.Atoi(rawPostID); err2 == nil {
			rctx.PostID = postID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("postID", rawPostID, "integer"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowPostContext) OK(r *WsnPost) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.wsn.post+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowPostContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}
