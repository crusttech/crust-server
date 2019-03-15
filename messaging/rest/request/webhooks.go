package request

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
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Webhooks webhookGet request parameters
type WebhooksWebhookGet struct {
	WebhookID uint64 `json:",string"`
}

func NewWebhooksWebhookGet() *WebhooksWebhookGet {
	return &WebhooksWebhookGet{}
}

func (wReq *WebhooksWebhookGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))

	return err
}

var _ RequestFiller = NewWebhooksWebhookGet()

// Webhooks webhookDelete request parameters
type WebhooksWebhookDelete struct {
	WebhookID uint64 `json:",string"`
}

func NewWebhooksWebhookDelete() *WebhooksWebhookDelete {
	return &WebhooksWebhookDelete{}
}

func (wReq *WebhooksWebhookDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))

	return err
}

var _ RequestFiller = NewWebhooksWebhookDelete()

// Webhooks webhookDeletePublic request parameters
type WebhooksWebhookDeletePublic struct {
	WebhookID    uint64 `json:",string"`
	WebhookToken string
}

func NewWebhooksWebhookDeletePublic() *WebhooksWebhookDeletePublic {
	return &WebhooksWebhookDeletePublic{}
}

func (wReq *WebhooksWebhookDeletePublic) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))
	wReq.WebhookToken = chi.URLParam(r, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksWebhookDeletePublic()

// Webhooks webhookPublish request parameters
type WebhooksWebhookPublish struct {
	Username     string
	AvatarURL    string
	Content      string
	WebhookID    uint64 `json:",string"`
	WebhookToken string
}

func NewWebhooksWebhookPublish() *WebhooksWebhookPublish {
	return &WebhooksWebhookPublish{}
}

func (wReq *WebhooksWebhookPublish) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["username"]; ok {

		wReq.Username = val
	}
	if val, ok := get["avatarURL"]; ok {

		wReq.AvatarURL = val
	}
	if val, ok := get["content"]; ok {

		wReq.Content = val
	}
	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))
	wReq.WebhookToken = chi.URLParam(r, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksWebhookPublish()
