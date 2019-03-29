package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks_public.go`, `webhooks_public.util.go` or `webhooks_public_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks_public.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type WebhooksPublicAPI interface {
	WebhookDelete(context.Context, *request.WebhooksPublicWebhookDelete) (interface{}, error)
	WebhookMessageCreate(context.Context, *request.WebhooksPublicWebhookMessageCreate) (interface{}, error)
}

// HTTP API interface
type WebhooksPublic struct {
	WebhookDelete        func(http.ResponseWriter, *http.Request)
	WebhookMessageCreate func(http.ResponseWriter, *http.Request)
}

func NewWebhooksPublic(wh WebhooksPublicAPI) *WebhooksPublic {
	return &WebhooksPublic{
		WebhookDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksPublicWebhookDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookDelete(r.Context(), params)
			})
		},
		WebhookMessageCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksPublicWebhookMessageCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookMessageCreate(r.Context(), params)
			})
		},
	}
}

func (wh *WebhooksPublic) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/webhooks", func(r chi.Router) {
			r.Delete("/{webhookID}/{webhookToken}", wh.WebhookDelete)
			r.Post("/{webhookID}/{webhookToken}", wh.WebhookMessageCreate)
		})
	})
}
