package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/sisu-network/interview/configs"
	"github.com/sisu-network/interview/internal/deliveries/tcp"
	"github.com/sisu-network/interview/internal/domains"
	"github.com/sisu-network/interview/internal/models"
	"github.com/sisu-network/interview/internal/repositories"
	"github.com/sisu-network/interview/internal/repositories/sqlite"
	"github.com/sisu-network/interview/pkg/tcp_server"

	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
)

type server struct {
	//* server
	tcpServer *tcp_server.TcpServer

	//* deliveries
	blockchainDelivery tcp.BlockchainDelivery

	//* domains
	blockchainDomain domains.BlockchainDomain

	//* repositories
	blockRepo  repositories.BlockRepository
	nodeRepo   repositories.NodeRepository
	markerRepo repositories.MarkerRepository

	db *sql.DB

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
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *server) loadLogger() error {
	// s.logger = logger.NewZapLogger("INFO", true)
	return nil
}
func (s *server) loadRepositories() error {

	q := models.New(s.db)
	s.blockRepo = sqlite.NewBlockRepository(q)
	s.markerRepo = sqlite.NewMarkerRepository(q)
	s.nodeRepo = sqlite.NewNodeRepository(q)
	return nil
}

func (s *server) loadDomains() error {
	s.blockchainDomain = domains.NewBlockchainDomain(s.nodeRepo, s.blockRepo, s.markerRepo, s.config)
	return nil
}

func (s *server) loadDeliveries() error {
	s.blockchainDelivery = tcp.NewBlockchainDelivery(s.blockchainDomain)
	return nil
}

func (s *server) loadConfig(ctx context.Context) error {
	s.config = &configs.Config{}
	s.config.SampleSize, _ = strconv.Atoi(os.Getenv("SAMPLE_SIZE"))
	s.config.QuorumSize, _ = strconv.Atoi(os.Getenv("QUORUM_SIZE"))
	s.config.DecisionThreshHold, _ = strconv.Atoi(os.Getenv("DECISION_THRESHOLD"))
	s.config.Tcp.Port = os.Getenv("PORT")
	return nil
}

func (s *server) loadServers(ctx context.Context) error {
	s.tcpServer = &tcp_server.TcpServer{
		Addr: configs.ConnectionAddr{
			Port: s.config.Tcp.Port,
		},
		Handlers: map[models.Event]func(ctx context.Context, req *models.Request) (*models.Response, error){
			models.PingEvent:     s.blockchainDelivery.RetrievePingEvent,
			models.ValidateEvent: s.blockchainDelivery.ValidateData,
		},
	}
	s.processors = append(s.processors, s.tcpServer)
	return nil
}
