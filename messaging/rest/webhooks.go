package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/types"
)

var _ = errors.Wrap

type Webhooks struct {
	webhook service.WebhookService
}

func (Webhooks) New() *Webhooks {
	return &Webhooks{}
}

func (ctrl *Webhooks) WebhookGet(ctx context.Context, r *request.WebhooksWebhookGet) (interface{}, error) {
	return ctrl.webhook.With(ctx).Get(r.WebhookID)
}

func (ctrl *Webhooks) WebhookDelete(ctx context.Context, r *request.WebhooksWebhookDelete) (interface{}, error) {
	return nil, ctrl.webhook.With(ctx).Delete(r.WebhookID)
}

func (ctrl *Webhooks) WebhookList(ctx context.Context, r *request.WebhooksWebhookList) (interface{}, error) {
	return ctrl.webhook.With(ctx).Find(&types.WebhookFilter{
		ChannelID:   r.ChannelID,
		OwnerUserID: r.UserID,
	})
}

func (ctrl *Webhooks) WebhookCreate(ctx context.Context, r *request.WebhooksWebhookCreate) (interface{}, error) {
	// Webhook request parameters
	parameters := types.WebhookRequest{
		r.Username,
		r.Avatar,
		r.AvatarURL,
		r.Trigger,
		r.Url,
	}
	return ctrl.webhook.With(ctx).Create(r.Kind, r.ChannelID, parameters)
}

func (ctrl *Webhooks) WebhookUpdate(ctx context.Context, r *request.WebhooksWebhookUpdate) (interface{}, error) {
	// Webhook request parameters
	parameters := types.WebhookRequest{
		r.Username,
		r.Avatar,
		r.AvatarURL,
		r.Trigger,
		r.Url,
	}
	return ctrl.webhook.With(ctx).Update(r.WebhookID, r.Kind, r.ChannelID, parameters)
}
