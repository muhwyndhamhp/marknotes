package internal

import (
	"context"
	"time"
)

type ModerationStatus int

const (
	ModerationUnverified ModerationStatus = iota
	ModerationOK
	ModerationWarning
	ModerationDangerous
	ModerationNothing
)

var moderationStatuses = map[ModerationStatus]string{
	ModerationUnverified: "UNVERIFIED",
	ModerationOK:         "OK",
	ModerationWarning:    "WARNING",
	ModerationDangerous:  "DANGEROUS",
	ModerationNothing:    "x",
}

func (r ModerationStatus) String() string {
	return moderationStatuses[r]
}

type Moderation struct {
	LastModeratedAt  *time.Time
	ModerationStatus ModerationStatus `json:"moderation_status"`
	ModerationReason string           `json:"moderation_reason"`
}

type LLM interface {
	ModerateReplies(ctx context.Context, replies []Reply) ([]Reply, error)
}
