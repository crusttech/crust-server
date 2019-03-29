package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
)

var _ = errors.Wrap

type WebhooksPublic struct {
	webhook service.WebhookService
}

func (WebhooksPublic) New() *WebhooksPublic {
	return &WebhooksPublic{}
}

func (ctrl *WebhooksPublic) WebhookDelete(ctx context.Context, r *request.WebhooksPublicWebhookDelete) (interface{}, error) {
	return nil, ctrl.webhook.With(ctx).DeleteByToken(r.WebhookID, r.WebhookToken)
}

func (ctrl *WebhooksPublic) WebhookMessageCreate(ctx context.Context, r *request.WebhooksPublicWebhookMessageCreate) (interface{}, error) {
	return ctrl.webhook.With(ctx).Message(r.WebhookID, r.WebhookToken, r.Username, r.AvatarURL, r.Content)
}
