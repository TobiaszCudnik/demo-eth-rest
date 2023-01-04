package service

import (
	"context"
	"github.com/TobiaszCudnik/infura-interview/internal/eth"
	rdb "github.com/TobiaszCudnik/infura-interview/internal/redis"
	"github.com/TobiaszCudnik/infura-interview/internal/utils"
	"log"
	"time"
)

type Service struct {
	redis *rdb.Client
	eth   *eth.Client
	// map of lists of channels to notify when "block-completed" triggered for block `k`
	waitBlock map[string]([]chan bool)
}

func New(rdb *rdb.Client, eth *eth.Client) *Service {
	return &Service{
		redis:     rdb,
		eth:       eth,
		waitBlock: map[string]([]chan bool){},
	}
}

func (s *Service) Start(ctx context.Context, ready chan<- bool) {
	// init
	go s.dispatchCompleted()
	ready <- true
	<-ctx.Done()
	// teardown
}

// GetBlock returns a block using one of the following:
// - from cache
// - by awaiting for block-completed
// - from the rpc source
func (s *Service) GetBlock(ctx context.Context, blockNum string) (*string, error) {
	log.Println("[service.getBlock] get " + blockNum)
	var content *string
	var err error

	for {
		log.Println("[service.getBlock] check cache for " + blockNum)
		content, err = s.redis.GetBlock(ctx, blockNum)
		if err != nil {
			return nil, err
		}
		if content != nil {
			// cached
			log.Println("[service.getBlock] using cache for " + blockNum)
			return content, nil
		}
		lock, err := s.redis.CheckBlockLock(ctx, blockNum)
		if err != nil {
			return nil, err
		}
		if lock {
			// already requested
			log.Println("[service.getBlock] locked, waiting for " + blockNum)
			select {
			case <-s.waitForBlock(blockNum):
				log.Println("[service.GetBlock] block-completed received for " + blockNum)
				content, err := s.redis.GetBlock(ctx, blockNum)
				if err != nil {
					log.Println("[ERR service.GetBlock] cache missing after wait")
					return nil, err
				}
				return content, nil
			case <-time.After(utils.LockTimeout):
				log.Println("[service.GetBlock] timeout for " + blockNum)
				// pass
			}
			// TODO sub to completions
			// timeout
		} else {
			// request content (no cache, no lock)
			break
		}
	}
	content, err = s.requestBlock(ctx, blockNum)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// requestBlock requests a block from the RPC source and returns it, after acquiring the lock.
// Then triggers a background cache save and notifies awaiting dup requests.
func (s *Service) requestBlock(ctx context.Context, blockNum string) (*string, error) {
	log.Println("[service.requestBlock] request " + blockNum)
	// lock
	if err := s.redis.AcquireBlockLock(ctx, blockNum); err != nil {
		return nil, err
	}
	// request
	content, err := s.eth.GetBlock(ctx, blockNum)
	if err != nil {
		return nil, err
	}
	// cache & notify in the background
	go func() {
		ctx := context.Background()
		// cache
		err = s.redis.SetBlock(ctx, blockNum, content)
		if err != nil {
			return
		}
		// notify
		err = s.redis.PublishBlock(ctx, blockNum)
		if err != nil {
			return
		}
	}()
	return content, nil
}

// waitForBlock returns a channel which reads when blockNum becomes available in the cache.
func (s *Service) waitForBlock(blockNum string) chan bool {
	ch := make(chan bool, 1)
	if _, ok := s.waitBlock[blockNum]; !ok {
		s.waitBlock[blockNum] = []chan bool{}
	}
	s.waitBlock[blockNum] = append(s.waitBlock[blockNum], ch)
	return ch
}

// dispatchCompleted reacts to block-completed on pubsub and notifies all the channels awaiting it,
// then disposes the listeners.
func (s *Service) dispatchCompleted() {
	ch := s.redis.ReceiveCompleted()
	for msg := range ch {
		if msg.Channel == "block-completed" {
			blockNum := msg.Payload
			subs, ok := s.waitBlock[blockNum]
			if !ok {
				continue
			}
			log.Println("[service.dispatchCompleted] sending block-completed for " + blockNum)
			// notify all pending requests
			for i := range subs {
				subs[i] <- true
			}
			// dispose
			delete(s.waitBlock, blockNum)
		}
	}
}
