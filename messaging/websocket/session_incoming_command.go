package websocket

import (
	"context"
	"log"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

func (s *Session) execCommand(ctx context.Context, c *incoming.ExecCommand) (*outgoing.Message, error) {
	// @todo: check access / can we join this channel (should be done by service layer)

	log.Printf("Received command '%s(%v)", c.Command, c.Params)

	switch c.Command {
	case "echo":
		if c.Input != "" {
			if user, err := systemService.User(ctx).FindByID(s.user.Identity()); err != nil {
				return nil, err
			} else {
				return &outgoing.Message{
					ID:        factory.Sonyflake.NextID(),
					User:      payload.User(user),
					CreatedAt: time.Now(),
					Type:      "hallucination",
					ChannelID: c.ChannelID,
					Message:   c.Input,
				}, nil
			}
		}
	case "shrug":
		msg := &types.Message{
			ChannelID: payload.ParseUInt64(c.ChannelID),
			Message:   "¯\\_(ツ)_/¯",
		}

		if c.Input != "" {
			msg.Message = c.Input + " " + msg.Message
		}
		_, err := s.svc.msg.With(ctx).Create(msg)
		return nil, err
	default:
		user, err := systemService.User(ctx).FindByID(s.user.Identity())
		if err != nil {
			return nil, err
		}

		webhooks, err := s.svc.webhook.Find(&types.WebhookFilter{
			OutgoingTrigger: c.Command,
		})
		if err != nil || len(webhooks) == 0 {
			// @todo: list available commands, webhook triggers?
			return nil, nil
		}
		message, err := s.svc.webhook.Do(webhooks[0], c.Input)
		if err != nil {
			return &outgoing.Message{
				ID:        factory.Sonyflake.NextID(),
				User:      payload.User(user),
				CreatedAt: time.Now(),
				Type:      "hallucination",
				ChannelID: c.ChannelID,
				Message:   "Error running webhook: " + err.Error(),
			}, nil
		}
		if message != nil {
			_, err := s.svc.msg.With(ctx).Create(message)
			return nil, err
		}
	}
	return nil, nil
}
