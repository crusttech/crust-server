package types

import (
	"time"
)

type (
	Webhook struct {
		ID uint64 `json:"id" db:"id"`

		AuthToken string `json:"-" db:"webhook_token"`

		OwnerUserID uint64 `json:"userId" db:"rel_owner"`

		// Created bot User ID
		UserID    uint64 `json:"userId" db:"rel_user"`
		ChannelID uint64 `json:"channelId" db:"rel_channel"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	WebhookFilter struct {
		ChannelID   uint64
		OwnerUserID uint64
	}
)
