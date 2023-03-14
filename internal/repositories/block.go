package repositories

import (
	"context"

	"github.com/sisu-network/interview/internal/models"
)

type BlockRepository interface {
	Create(ctx context.Context, block *models.Block) error
	GetByID(ctx context.Context, id string) (*models.Block, error)
	GetLatestBlock(ctx context.Context) (*models.Block, error)
}
