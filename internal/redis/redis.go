package redis

import (
	"context"
	"github.com/TobiaszCudnik/infura-interview/internal/utils"
	"github.com/go-redis/redis/v8"
	"log"
)

type Client struct {
	addr string
	rdb  *redis.Client
	ps   *redis.PubSub
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Start(ctx context.Context, ready chan<- bool) {
	// init

	c.rdb = redis.NewClient(&redis.Options{
		Addr: c.addr, Password: "", DB: 0,
	})

	c.ps = c.rdb.Subscribe(ctx, "block-completed", "tx-completed")

	ready <- true
	<-ctx.Done()
	// teardown

	err := c.rdb.Close()
	if err != nil {
		log.Println("[ERR Start] " + err.Error())
	}
	err = c.ps.Close()
	if err != nil {
		log.Println("[ERR Start] " + err.Error())
	}
}

func (c *Client) ReceiveCompleted() (ch <-chan *redis.Message) {
	return c.ps.Channel()
}

func (c *Client) GetBlock(ctx context.Context, blockNum string) (*string, error) {
	k := "block-" + blockNum
	val, err := c.rdb.Get(ctx, k).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Println("[ERR redis.GetBlock] " + err.Error())
		return nil, err
	}
	return &val, nil
}

func (c *Client) AcquireBlockLock(ctx context.Context, blockNum string) error {
	k := "block-lock-" + blockNum
	err := c.rdb.Set(ctx, k, true, utils.LockTimeout).Err()
	if err != nil {
		log.Println("[ERR redis.AcquireBlockLock] " + err.Error())
		return err
	}
	return nil
}

func (c *Client) ReleaseBlockLock(ctx context.Context, blockNum string) error {
	k := "block-lock-" + blockNum
	err := c.rdb.Set(ctx, k, nil, 0).Err()
	if err != nil {
		log.Println("[ERR redis.ReleaseBlockLock] " + err.Error())
		return err
	}
	return nil
}

func (c *Client) CheckBlockLock(ctx context.Context, blockNum string) (bool, error) {
	k := "block-lock-" + blockNum
	err := c.rdb.Get(ctx, k).Err()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Println("[ERR redis.CheckBlockLock] " + err.Error())
		return false, err
	}
	return true, nil
}

func (c *Client) SetBlock(ctx context.Context, blockNum string, content *string) error {
	k := "block-" + blockNum
	err := c.rdb.Set(ctx, k, *content, 0).Err()
	if err != nil {
		log.Println("[ERR redis.SetBlock] " + err.Error())
		return err
	}
	return nil
}

func (c *Client) PublishBlock(ctx context.Context, blockNum string) error {
	ch := "block-completed"
	err := c.rdb.Publish(ctx, ch, blockNum).Err()
	if err != nil {
		log.Println("[ERR redis.PublishBlock] " + err.Error())
		return err
	}
	return nil
}
