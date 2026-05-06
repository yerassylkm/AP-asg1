package domain

import "context"

type ProcessedMessageRepository interface {
	IsProcessed(ctx context.Context, messageID string) (bool, error)
	MarkProcessed(ctx context.Context, messageID string) error
}
