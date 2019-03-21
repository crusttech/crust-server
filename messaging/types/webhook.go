package types

import (
	"time"
)

type (
	Webhook struct {
		ID uint64 `json:"id" db:"id"`

		Kind      WebhookKind `json:"kind" db:"webhook_kind"`
		AuthToken string      `json:"-" db:"webhook_token"`

		OwnerUserID uint64 `json:"userId" db:"rel_owner"`

		// Created bot User ID
		UserID    uint64 `json:"userId" db:"rel_user"`
		ChannelID uint64 `json:"channelId" db:"rel_channel"`

		// Outgoing webhook details
		OutgoingTrigger string `json:"trigger" db:"outgoing_trigger"`
		OutgoingURL     string `json:"url" db:"outgoing_url"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	WebhookFilter struct {
		ChannelID   uint64
		OwnerUserID uint64
	}

	WebhookKind string
)

const (
	IncomingWebhook WebhookKind = "incoming"
	OutgoingWebhook             = "outgoing"
)
