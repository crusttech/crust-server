package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	webhook struct {
		db     db
		ctx    context.Context
		client *http.Client

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

		Do(webhook *types.Webhook, message string) (*types.Message, error)
	}
)

func Webhook(client *http.Client) WebhookService {
	return (&webhook{
		client: client,
	}).With(context.Background())
}

func (svc *webhook) With(ctx context.Context) WebhookService {
	db := repository.DB(ctx)
	return &webhook{
		db:     db,
		ctx:    ctx,
		client: svc.client,

		users: systemService.User(ctx),

		webhook: repository.Webhook(ctx, db),
		message: repository.Message(ctx, db),
	}
}

func (svc *webhook) CreateIncoming(channelID uint64, username string, avatar interface{}) (*types.Webhook, error) {
	var userID = repository.Identity(svc.ctx)
	// @todo: avatar
	webhook := &types.Webhook{
		Kind:        types.IncomingWebhook,
		AuthToken:   "123", // @todo: JWT
		OwnerUserID: userID,
		UserID:      userID, // @todo: create bot user
		ChannelID:   channelID,
		CreatedAt:   time.Now(),
	}
	return svc.webhook.Create(webhook)
}
func (svc *webhook) CreateOutgoing(channelID uint64, username string, avatar interface{}, trigger string, url string) (*types.Webhook, error) {
	var userID = repository.Identity(svc.ctx)
	// @todo: avatar
	webhook := &types.Webhook{
		Kind:            types.OutgoingWebhook,
		OwnerUserID:     userID,
		UserID:          userID, // @todo: create bot user
		ChannelID:       channelID,
		OutgoingTrigger: trigger,
		OutgoingURL:     url,
	}
	return svc.webhook.Create(webhook)
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

func (svc *webhook) Do(webhook *types.Webhook, message string) (*types.Message, error) {
	if webhook.Kind != types.OutgoingWebhook {
		return nil, errors.Errorf("Unsupported webhook type: %s", webhook.Kind)
	}

	// replace url query %s with message
	url := strings.Replace(webhook.OutgoingURL, "%s", url.QueryEscape(message), -1)

	// post body contains only `text`
	requestBody := types.WebhookBody{message}
	req, err := svc.client.Post(url, requestBody)
	if err != nil {
		return nil, err
	}

	// execute outgoing webhook
	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response body
	responseBody := types.WebhookBody{}
	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "text/plain"):
		// keep plain/text as-is
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, errors.WithStack(err)
		} else {
			responseBody.Text = string(b)
		}
	default:
		switch resp.StatusCode {
		case 200:
			// assume the response is an expected json structure
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
				return nil, errors.WithStack(err)
			}
			if responseBody.Text == "" {
				return nil, errors.New("Empty webhook response")
			}
		default:
			return nil, http.ToError(resp)
		}
	}
	return &types.Message{
		ID:        factory.Sonyflake.NextID(),
		UserID:    webhook.UserID,
		ChannelID: webhook.ChannelID,
		CreatedAt: time.Now(),
		Message:   responseBody.Text,
	}, nil
}
