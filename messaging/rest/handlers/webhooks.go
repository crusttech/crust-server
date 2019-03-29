package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks.go`, `webhooks.util.go` or `webhooks_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks.go`
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
type WebhooksAPI interface {
	WebhookList(context.Context, *request.WebhooksWebhookList) (interface{}, error)
	WebhookCreate(context.Context, *request.WebhooksWebhookCreate) (interface{}, error)
	WebhookUpdate(context.Context, *request.WebhooksWebhookUpdate) (interface{}, error)
	WebhookGet(context.Context, *request.WebhooksWebhookGet) (interface{}, error)
	WebhookDelete(context.Context, *request.WebhooksWebhookDelete) (interface{}, error)
}

// HTTP API interface
type Webhooks struct {
	WebhookList   func(http.ResponseWriter, *http.Request)
	WebhookCreate func(http.ResponseWriter, *http.Request)
	WebhookUpdate func(http.ResponseWriter, *http.Request)
	WebhookGet    func(http.ResponseWriter, *http.Request)
	WebhookDelete func(http.ResponseWriter, *http.Request)
}

func NewWebhooks(wh WebhooksAPI) *Webhooks {
	return &Webhooks{
		WebhookList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookList(r.Context(), params)
			})
		},
		WebhookCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookCreate(r.Context(), params)
			})
		},
		WebhookUpdate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookUpdate(r.Context(), params)
			})
		},
		WebhookGet: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookGet(r.Context(), params)
			})
		},
		WebhookDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookDelete(r.Context(), params)
			})
		},
	}
}

func (wh *Webhooks) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/webhooks", func(r chi.Router) {
			r.Get("/", wh.WebhookList)
			r.Post("/", wh.WebhookCreate)
			r.Post("/{webhookID}", wh.WebhookUpdate)
			r.Get("/{webhookID}", wh.WebhookGet)
			r.Delete("/{webhookID}", wh.WebhookDelete)
		})
	})
}
