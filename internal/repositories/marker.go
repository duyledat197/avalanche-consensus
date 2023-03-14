package repositories

import (
	"context"

	"github.com/sisu-network/interview/internal/models"
)

type MarkerRepository interface {
	MarkBlock(ctx context.Context, blockID string) error
	GetByBlockID(ctx context.Context, blockID string) (*models.Marker, error)
}
