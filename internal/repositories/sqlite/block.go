package sqlite

import (
	"context"

	"github.com/sisu-network/interview/internal/models"
	"github.com/sisu-network/interview/internal/repositories"
)

type blockRepository struct {
	q *models.Queries
}

func NewBlockRepository(q *models.Queries) repositories.BlockRepository {
	return &blockRepository{q: q}
}

func (r *blockRepository) Create(ctx context.Context, block *models.Block) error {
	panic("not implemented") // TODO: Implement
}

func (r *blockRepository) GetByID(ctx context.Context, id string) (*models.Block, error) {
	panic("not implemented") // TODO: Implement
}

func (r *blockRepository) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	panic("not implemented") // TODO: Implement
}

func (r *blockRepository) GetAll(ctx context.Context) ([]*models.Block, error) {
	panic("not implemented") // TODO: Implement
}
