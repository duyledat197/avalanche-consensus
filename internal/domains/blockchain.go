package domains

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sisu-network/interview/configs"
	"github.com/sisu-network/interview/internal/models"
	"github.com/sisu-network/interview/internal/repositories"
)

type BlockchainDomain interface {
	Validate(ctx context.Context, data []int) error
	SnowBall(ctx context.Context, blockID string, data []int) error
	PingNeighbourNodes(ctx context.Context, blockID string)
}

type blockchainDomain struct {
	nodeRepo   repositories.NodeRepository
	blockRepo  repositories.BlockRepository
	markerRepo repositories.MarkerRepository
	configs    *configs.Config
}

func NewBlockchainDomain(
	nodeRepo repositories.NodeRepository,
	blockRepo repositories.BlockRepository,
	markerRepo repositories.MarkerRepository,
	configs *configs.Config,
) BlockchainDomain {
	return &blockchainDomain{nodeRepo, blockRepo, markerRepo, configs}
}

func (d *blockchainDomain) Validate(ctx context.Context, data []int) error {
	block, err := d.blockRepo.GetLatestBlock(ctx)
	if err != nil {
		return err
	}

	m := make(map[int]int)
	for pos, num := range block.Data {
		m[int(num)] = pos
	}
	pointer := -1
	for _, num := range data {
		if _, ok := m[num]; ok {
			if m[num] < pointer {
				return fmt.Errorf("wrong position")
			}
			pointer = m[num]
		}
	}
	return nil
}

func (d *blockchainDomain) PingNeighbourNodes(ctx context.Context, blockID string) {
	_, err := d.markerRepo.GetByBlockID(ctx, blockID)
	if err != nil {
		log.Printf("unable to get block: %v", err)
		return
	}
	if err := d.markerRepo.MarkBlock(ctx, blockID); err != nil {
		log.Printf("unable to mark block: %v", err)
		return
	}
	nodes, err := d.nodeRepo.GetAll(ctx)
	if err != nil {
		log.Printf("unable to get nodes: %v", err)
		return
	}
	for _, node := range nodes {
		go func() {
			conn, err := net.Dial("tcp", node.Address)
			if err != nil {
				log.Printf("unable to connect %s: %v", node.Address, err)
				return
			}
			defer conn.Close()
			var (
				req models.Request
			)
			b, _ := json.Marshal(&req)
			if _, err := conn.Write(b); err != nil {
				log.Printf("unable to write data : %v", err)
				return
			}
		}()
	}
}

func (d *blockchainDomain) SnowBall(ctx context.Context, blockID string, data []int) error {
	preference := false
	consecutiveSuccesses := 0
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	startTime := time.Now()

	for {
		if time.Since(startTime) > 10*time.Second {
			break
		}
		chosenNodes, err := d.nodeRepo.GetRandom(ctx, d.configs.SampleSize)
		if err != nil {
			log.Printf("error choosing nodes: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		totalAccept := 0
		totalDecline := 0
		for _, node := range chosenNodes {
			wg.Add(1)

			go func(ctx context.Context, addr string) {
				defer wg.Done()
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					log.Printf("unable to connect %s: %w", addr, err)
					return
				}
				defer conn.Close()
				var (
					req  models.Request
					resp models.Response
				)
				b, _ := json.Marshal(&req)
				if _, err := conn.Write(b); err != nil {
					log.Printf("unable to write data : %w", err)
					return
				}

				data := make([]byte, 1024)
				n, err := conn.Read(data)
				if err != nil {
					log.Printf("unable to reading :%v\n", err)
					return
				}

				json.Unmarshal(data[:n], &resp)
				if resp.Err != nil {
					log.Printf("unable to unmarshal :%v\n", err)
				}

				mutex.Lock()
				if resp.IsAccept {
					totalAccept += 1
				} else {
					totalDecline += 1
				}
				mutex.Unlock()
			}(ctx, node.Address)
		}
		wg.Wait()
		if totalAccept < d.configs.QuorumSize && totalDecline < d.configs.QuorumSize {
			consecutiveSuccesses = 0
			continue
		}
		currenntPref := totalAccept > totalDecline
		if currenntPref != preference {
			preference = currenntPref
			consecutiveSuccesses = 1
		} else {
			consecutiveSuccesses++
		}

		if consecutiveSuccesses >= d.configs.ThreshHold {
			d.decide(ctx, blockID, data)
			break
		}
	}
	return nil
}

func (d *blockchainDomain) decide(ctx context.Context, blockID string, data []int) {
	var nums []int32
	for _, num := range data {
		nums = append(nums, int32(num))
	}
	d.blockRepo.Create(ctx, &models.Block{
		ID:   blockID,
		Data: nums,
	})
}
