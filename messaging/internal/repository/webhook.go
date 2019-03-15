package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/messaging/types"
)

type (
	WebhookRepository interface {
		With(ctx context.Context, db *factory.DB) WebhookRepository

		Get(webhookID uint64) (*types.Webhook, error)
		GetWithToken(webhookID uint64, webhookToken string) (*types.Webhook, error)

		Find(filter *types.WebhookFilter) (types.WebhookSet, error)

		Delete(webhookID uint64) error
		DeleteWithToken(webhookID uint64, webhookToken string) error
	}

	webhook struct {
		webhook string

		*repository
	}
)

func Webhook(ctx context.Context, db *factory.DB) WebhookRepository {
	return (&webhook{}).With(ctx, db)
}

func (r *webhook) With(ctx context.Context, db *factory.DB) WebhookRepository {
	return &webhook{
		webhook:    "messaging_webhook",
		repository: r.repository.With(ctx, db),
	}
}

func (r *webhook) Get(webhookID uint64) (*types.Webhook, error) {
	hook := &types.Webhook{}
	if err := r.db().Get(&hook, "select * from "+r.webhook+" where id=?", webhookID); err != nil {
		return nil, err
	}
	return hook, nil
}

func (r *webhook) GetWithToken(webhookID uint64, webhookToken string) (*types.Webhook, error) {
	webhook, err := r.Get(webhookID)
	switch {
	case err != nil:
		return nil, err
	case webhook.AuthToken == webhookToken:
		return webhook, nil
	default:
		return nil, errors.New("Invalid Webhook Token")
	}
}

// Find webhooks based on filter
//
// If ChannelID > 0 it returns webhooks created on a specific channel
// If OwnerUserID > 0 it returns webhooks owned by a specific user
func (r *webhook) Find(filter *types.WebhookFilter) (types.WebhookSet, error) {
	params := make([]interface{}, 0)
	vv := types.WebhookSet{}
	sql := "select * from messaging_webhook where 1=1"

	if filter != nil {
		if filter.OwnerUserID > 0 {
			// scope: only channel we have access to
			sql += " AND rel_owner=?"
			params = append(params, filter.OwnerUserID)
		}
		if filter.ChannelID > 0 {
			// scope: only channel we have access to
			sql += " AND rel_channel=?"
			params = append(params, filter.ChannelID)
		}
	}

	return vv, r.db().Select(&vv, sql, params...)
}

func (r *webhook) Delete(webhookID uint64) error {
	_, err := r.Get(webhookID)
	if err != nil {
		return err
	}
	_, err = r.db().Exec("delete from "+r.webhook+" where id=?", webhookID)
	return err
}

func (r *webhook) DeleteWithToken(webhookID uint64, webhookToken string) error {
	_, err := r.GetWithToken(webhookID, webhookToken)
	if err != nil {
		return err
	}
	_, err = r.db().Exec("delete from "+r.webhook+" where id=?", webhookID)
	return err
}
