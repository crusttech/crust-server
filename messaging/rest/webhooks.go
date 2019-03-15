package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
)

var _ = errors.Wrap

type Webhooks struct {
	webhook service.WebhookService
}

func (Webhooks) New() *Webhooks {
	return &Webhooks{}
}

func (ctrl *Webhooks) WebhookGet(ctx context.Context, r *request.WebhooksWebhookGet) (interface{}, error) {
	return ctrl.webhook.Get(r.WebhookID)
}

func (ctrl *Webhooks) WebhookDelete(ctx context.Context, r *request.WebhooksWebhookDelete) (interface{}, error) {
	return nil, ctrl.webhook.Delete(r.WebhookID)
}

func (ctrl *Webhooks) WebhookDeletePublic(ctx context.Context, r *request.WebhooksWebhookDeletePublic) (interface{}, error) {
	return nil, ctrl.webhook.DeleteWithToken(r.WebhookID, r.WebhookToken)
}

func (ctrl *Webhooks) WebhookPublish(ctx context.Context, r *request.WebhooksWebhookPublish) (interface{}, error) {
	return ctrl.webhook.PublishMessage(r.WebhookID, r.WebhookToken, r.Username, r.AvatarURL, r.Content)
}
