package service

import (
	"context"

	"github.com/pkg/errors"

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

		Find(*types.WebhookFilter) (types.WebhookSet, error)

		Delete(webhookID uint64) error
		DeleteByToken(webhookID uint64, webhookToken string) error

		Message(webhookID uint64, webhookToken string, username, avatar, message string) (interface{}, error)

		CreateIncoming(channelID uint64, username string, avatar interface{}) (*types.Webhook, error)
		CreateOutgoing(channelID uint64, username string, avatar interface{}, trigger string, url string) (*types.Webhook, error)
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

func (svc *webhook) CreateIncoming(channelID uint64, username string, avatar interface{}) (*types.Webhook, error) {
	// @todo: create bot user, create webhook in db
	return nil, errors.New("Not implemented")
}
func (svc *webhook) CreateOutgoing(channelID uint64, username string, avatar interface{}, trigger string, url string) (*types.Webhook, error) {
	// @todo: create bot user, create webhook in db, check triggers in message(s)
	return nil, errors.New("Not implemented")
}

func (svc *webhook) Get(webhookID uint64) (*types.Webhook, error) {
	return svc.webhook.Get(webhookID)
}

func (svc *webhook) Find(filter *types.WebhookFilter) (types.WebhookSet, error) {
	return svc.webhook.Find(filter)
}

func (svc *webhook) Delete(webhookID uint64) error {
	return svc.webhook.Delete(webhookID)
}

func (svc *webhook) DeleteByToken(webhookID uint64, webhookToken string) error {
	return svc.webhook.DeleteByToken(webhookID, webhookToken)
}

func (svc *webhook) Message(webhookID uint64, webhookToken string, username, avatar, message string) (interface{}, error) {
	if webhook, err := svc.webhook.GetByToken(webhookID, webhookToken); err != nil {
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
