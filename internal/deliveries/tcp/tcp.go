package tcp

import (
	"context"

	"github.com/sisu-network/interview/internal/domains"
	"github.com/sisu-network/interview/internal/models"
)

type BlockchainDelivery interface {
	RetrievePingEvent(ctx context.Context, req *models.Request) (*models.Response, error)
	ValidateData(ctx context.Context, req *models.Request) (*models.Response, error)
}

type blockchainDelivery struct {
	blockchainDomain domains.BlockchainDomain
}

func (d *blockchainDelivery) RetrievePingEvent(ctx context.Context, req *models.Request) (*models.Response, error) {
	if err := d.blockchainDomain.SnowBall(ctx, req.BlockID, req.Data); err != nil {
		return nil, err
	}
	return nil, nil
}
func (d *blockchainDelivery) ValidateData(ctx context.Context, req *models.Request) (*models.Response, error) {
	err := d.blockchainDomain.Validate(ctx, req.Data)
	if err != nil {
		return nil, err
	}
	return &models.Response{
		IsAccept: true,
	}, nil
}
