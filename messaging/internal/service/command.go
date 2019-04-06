package service

import (
	"context"

	"github.com/crusttech/crust/messaging/types"
)

type (
	command struct {
		ctx context.Context
	}

	CommandService interface {
		With(context.Context) CommandService

		Do(channelID uint64, command, input string) error
	}
)

func Command(ctx context.Context) CommandService {
	return (&command{}).With(ctx)
}

func (svc *command) With(ctx context.Context) CommandService {
	return &command{
		ctx: ctx,
	}
}

func (svc *command) Do(channelID uint64, command, input string) error {
	switch command {
	case "tableflip":
		fallthrough
	case "unflip":
		fallthrough
	case "shrug":
		messages := map[string]string{
			"tableflip": "(╯°□°）╯︵ ┻━┻",
			"unflip":    "┬─┬ ノ( ゜-゜ノ)",
			"shrug":     "¯\\_(ツ)_/¯",
		}
		msg := &types.Message{
			ChannelID: channelID,
			Message:   messages[command],
		}

		if input != "" {
			msg.Message = input + " " + msg.Message
		}
		_, err := DefaultMessage.With(svc.ctx).Create(msg)
		return err
	default:
		webhookSvc := DefaultWebhook.With(svc.ctx)
		webhooks, err := webhookSvc.Find(&types.WebhookFilter{
			ChannelID: channelID,
			OutgoingTrigger: command,
		})
		if err != nil || len(webhooks) == 0 {
			return err
		}
		_, err = webhookSvc.Do(webhooks[0], input)
		return err
	}
	return nil
}
