package eth

import (
	"context"
	"encoding/json"
	"github.com/TobiaszCudnik/infura-interview/internal/utils"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	rpc *jrpc2.Client
}

func New(addr string) *Client {
	ch := jhttp.NewChannel(addr, &jhttp.ChannelOptions{
		Client: &http.Client{
			Timeout: utils.RpcTimeout,
		},
	})
	return &Client{
		rpc: jrpc2.NewClient(ch, nil),
	}
}

func (c *Client) Start(ctx context.Context, ready chan<- bool) {
	// init
	ready <- true
	<-ctx.Done()
	// teardown
}

func (c *Client) GetBlock(ctx context.Context, blockNum string) (*string, error) {
	var res json.RawMessage
	err := c.rpc.CallResult(ctx, "eth_getBlockByNumber", []any{blockNum, false}, &res)
	if err != nil {
		log.Println("[ERR eth.GetBlock]" + err.Error())
		return nil, err
	}

	// test helper
	testDelay := os.Getenv("ETH_REQ_DELAY_SEC")
	if testDelay != "" {
		td, err := strconv.ParseFloat(testDelay, 10)
		if err != nil {
			log.Println("[ERR eth.GetBlock]" + err.Error())
			return nil, err
		}
		log.Printf("[eth.GetBlock] sleeping for %f \n", td)
		time.Sleep(time.Second * time.Duration(td))
	}

	// cast
	ret := string(res[:])
	return &ret, nil
}
