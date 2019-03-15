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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type WebhooksAPI interface {
	WebhookGet(context.Context, *request.WebhooksWebhookGet) (interface{}, error)
	WebhookDelete(context.Context, *request.WebhooksWebhookDelete) (interface{}, error)
	WebhookDeletePublic(context.Context, *request.WebhooksWebhookDeletePublic) (interface{}, error)
	WebhookPublish(context.Context, *request.WebhooksWebhookPublish) (interface{}, error)
}

// HTTP API interface
type Webhooks struct {
	WebhookGet          func(http.ResponseWriter, *http.Request)
	WebhookDelete       func(http.ResponseWriter, *http.Request)
	WebhookDeletePublic func(http.ResponseWriter, *http.Request)
	WebhookPublish      func(http.ResponseWriter, *http.Request)
}

func NewWebhooks(wh WebhooksAPI) *Webhooks {
	return &Webhooks{
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
		WebhookDeletePublic: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookDeletePublic()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookDeletePublic(r.Context(), params)
			})
		},
		WebhookPublish: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksWebhookPublish()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.WebhookPublish(r.Context(), params)
			})
		},
	}
}

func (wh *Webhooks) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/webhooks", func(r chi.Router) {
			r.Get("/webhook/{webhookID}", wh.WebhookGet)
			r.Delete("/webhook/{webhookID}", wh.WebhookDelete)
			r.Delete("/webhook/{webhookID}/{webhookToken}", wh.WebhookDeletePublic)
			r.Post("/webhook/{webhookID}/{webhookToken}", wh.WebhookPublish)
		})
	})
}
