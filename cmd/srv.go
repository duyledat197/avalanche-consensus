package cmd

import (
	"context"

	"github.com/sisu-network/interview/configs"

	"go.uber.org/zap"
)

type server struct {

	//* config
	config *configs.Config

	//* logger
	logger *zap.Logger

	processors []processor
	factories  []factory
}

type processor interface {
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type factory interface {
	Connect(ctx context.Context) error
	Stop(ctx context.Context) error
}

func (s *server) loadDatabaseClients(ctx context.Context) error {

}

func (s *server) loadLogger() error {
	// s.logger = logger.NewZapLogger("INFO", true)
	return nil
}
func (s *server) loadRepositories() error {

	return nil
}

func (s *server) loadServices() error {

	return nil
}

func (s *server) loadDeliveries() error {

	return nil
}

func (s *server) loadConfig(ctx context.Context) error {
	return nil
}

func (s *server) loadClients(ctx context.Context) error {

	return nil
}

func (s *server) loadServers(ctx context.Context) error {

	return nil
}
