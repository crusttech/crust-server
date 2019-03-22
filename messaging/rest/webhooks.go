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
	return ctrl.webhook.Get(r.WebhookID)
}

func (ctrl *Webhooks) WebhookDelete(ctx context.Context, r *request.WebhooksWebhookDelete) (interface{}, error) {
	return nil, ctrl.webhook.Delete(r.WebhookID)
}

func (ctrl *Webhooks) WebhookDeletePublic(ctx context.Context, r *request.WebhooksWebhookDeletePublic) (interface{}, error) {
	return nil, ctrl.webhook.DeleteByToken(r.WebhookID, r.WebhookToken)
}

func (ctrl *Webhooks) WebhookMessageCreate(ctx context.Context, r *request.WebhooksWebhookMessageCreate) (interface{}, error) {
	return ctrl.webhook.Message(r.WebhookID, r.WebhookToken, r.Username, r.AvatarURL, r.Content)
}

func (ctrl *Webhooks) WebhookList(ctx context.Context, r *request.WebhooksWebhookList) (interface{}, error) {
	return ctrl.webhook.With(ctx).Find(&types.WebhookFilter{
		ChannelID:   r.ChannelID,
		OwnerUserID: r.UserID,
	})
}
func (ctrl *Webhooks) WebhookCreate(ctx context.Context, r *request.WebhooksWebhookCreate) (interface{}, error) {
	// Webhooks webhookCreate request parameters
	/*
	   ChannelID uint64 `json:",string"`
	   Kind      types.WebhookKind
	   Trigger   string
	   Url       string
	   Username  string
	   Avatar    *multipart.FileHeader
	*/

	// @todo: process r.Avatar file upload for webhook

	switch {
	case r.Kind == types.IncomingWebhook:
		return ctrl.webhook.With(ctx).CreateIncoming(r.ChannelID, r.Username, r.Avatar)
	case r.Kind == types.OutgoingWebhook:
		return ctrl.webhook.With(ctx).CreateOutgoing(r.ChannelID, r.Username, r.Avatar, r.Trigger, r.Url)
	}
	return nil, errors.New("Unknown webhook type")
}
