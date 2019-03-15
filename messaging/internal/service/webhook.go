package service

import (
	"context"

	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	webhook struct {
		db  db
		ctx context.Context

		users   systemService.UserService
		webhook repository.WebhookRepository
		message repository.MessageRepository
	}

	WebhookService interface {
		With(ctx context.Context) WebhookService

		Get(webhookID uint64) (*types.Webhook, error)

		ListByChannel(channelID uint64) (types.WebhookSet, error)

		Delete(webhookID uint64) error
		DeleteWithToken(webhookID uint64, webhookToken string) error

		PublishMessage(webhookID uint64, webhookToken string, username, avatar, message string) (interface{}, error)
	}
)

func Webhook() WebhookService {
	return (&webhook{}).With(context.Background())
}

func (svc *webhook) With(ctx context.Context) WebhookService {
	db := repository.DB(ctx)
	return &webhook{
		db:  db,
		ctx: ctx,

		users: systemService.User(ctx),

		webhook: repository.Webhook(ctx, db),
		message: repository.Message(ctx, db),
	}
}

func (svc *webhook) Get(webhookID uint64) (*types.Webhook, error) {
	return svc.webhook.Get(webhookID)
}

func (svc *webhook) ListByChannel(channelID uint64) (types.WebhookSet, error) {
	return svc.webhook.Find(&types.WebhookFilter{
		ChannelID: channelID,
	})
}

func (svc *webhook) Delete(webhookID uint64) error {
	return svc.webhook.Delete(webhookID)
}

func (svc *webhook) DeleteWithToken(webhookID uint64, webhookToken string) error {
	return svc.webhook.DeleteWithToken(webhookID, webhookToken)
}

func (svc *webhook) PublishMessage(webhookID uint64, webhookToken string, username, avatar, message string) (interface{}, error) {
	if webhook, err := svc.webhook.GetWithToken(webhookID, webhookToken); err != nil {
		return nil, err
	} else {
		message := &types.Message{
			ChannelID: webhook.ChannelID,
			UserID:    webhook.UserID,
			Message:   message,
		}
		return svc.message.CreateMessage(message)
	}
}
