package repository

import (
	"context"
	"notification-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) domain.ProcessedMessageRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) IsProcessed(ctx context.Context, messageID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM processed_messages WHERE message_id = $1)", messageID)
	return exists, err
}

func (r *postgresRepo) MarkProcessed(ctx context.Context, messageID string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO processed_messages (message_id) VALUES ($1) ON CONFLICT DO NOTHING", messageID)
	return err
}
