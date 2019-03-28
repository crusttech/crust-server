// +build integration,external

package service

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestOutgoingWebhook(t *testing.T) {
	var user = &systemTypes.User{ID: 1}
	var channel = &types.Channel{ID: 1}

	ctx := context.WithValue(context.Background(), "testing", true)
	ctx = auth.SetIdentityToContext(ctx, user)

	client, err := http.New(&config.HTTPClient{
		Timeout: 10,
	})
	test.Assert(t, err == nil, "Error creating HTTP client: %+v", err)

	/* create outgoing webhook */
	svc := Webhook(ctx, client)
	webhook, err := svc.CreateOutgoing(channel.ID, "test-webhook", "", "fortune", "https://api.scene-si.org/fortune.php")
	test.Assert(t, err == nil, "Error when creating webhook: %+v", err)

	/* find outgoing webhook */
	webhooks, err := svc.Find(&types.WebhookFilter{
		OutgoingTrigger: webhook.OutgoingTrigger,
	})
	test.Assert(t, err == nil, "Error when finding webhook: %+v", err)
	test.Assert(t, len(webhooks) == 1, "Expected to find 1 webhook, got %d", len(webhooks))

	/* trigger outgoing webhook */
	message, err := svc.Do(webhooks[0], "")
	test.Assert(t, err == nil, "Error when triggering webhook: %+v", err)
	test.Assert(t, strings.Contains(message.Message, "BOFH"), "Unexpected webhook output: %s", message.Message)
}
