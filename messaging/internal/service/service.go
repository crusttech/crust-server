package service

import (
	"time"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultAttachment  AttachmentService
	DefaultChannel     ChannelService
	DefaultMessage     MessageService
	DefaultPubSub      *pubSub
	DefaultEvent       EventService
	DefaultPermissions PermissionsService
	DefaultWebhook     WebhookService
)

func Init() error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	client, err := http.New(&config.HTTPClient{
		Timeout: 10,
	})
	if err != nil {
		return err
	}

	DefaultPermissions = Permissions()
	DefaultEvent = Event()
	DefaultAttachment = Attachment(fs)
	DefaultMessage = Message()
	DefaultChannel = Channel()
	DefaultPubSub = PubSub()
	DefaultWebhook = Webhook(client)

	return nil
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
