package websocket

import (
	"context"
	"log"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/internal/payload/outgoing"
	systemService "github.com/crusttech/crust/system/service"
)

func (s *Session) execCommand(ctx context.Context, c *incoming.ExecCommand) error {
	// @todo: check access / can we join this channel (should be done by service layer)

	log.Printf("Received command '%s(%v)", c.Command, c.Params)

	if c.Command == "echo" {
		if c.Input != "" {
			if user, err := systemService.User(ctx).FindByID(s.user.Identity()); err != nil {
				return err
			} else {
				return s.sendReply(&outgoing.Message{
					ID:        factory.Sonyflake.NextID(),
					User:      payload.User(user),
					CreatedAt: time.Now(),
					Type:      "hallucination",
					ChannelID: c.ChannelID,
					Message:   c.Input,
				})
			}
		}
	}
	return s.svc.command.With(ctx).Do(c.ChannelID, c.Command, c.Input)
}
