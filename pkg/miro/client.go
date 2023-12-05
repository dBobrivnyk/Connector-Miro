package miro

import (
	config "connectors/internal"
	"connectors/pkg/sink"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ApiUrl = "https://api.miro.com/v2/"

// tests ??
// for http.Client you can create proper interface
// in test add mock
// and then you can create proper tests
// use this or any other option to mock http call
// https://github.com/stretchr/testify#mock-package
type Client struct {
	cfg     *config.Config
	client  *http.Client
	url     string
	ownerId string
	apiKey  string
	sink    *sink.Sink
}

func NewClient(cfg *config.Config, url string, ownerId, apiKey string) *Client {
	return &Client{
		cfg:     cfg,
		client:  http.DefaultClient,
		url:     url,
		ownerId: ownerId,
		apiKey:  apiKey,
	}
}

func (c *Client) GetEntities(ctx context.Context) {
	c.sink = sink.New(c.cfg.BufferSize, c.ownerId)
	defer c.sink.Dump()
	defer c.sink.Close()

	log.Println("Searching among boards")
	boards, err := c.GetBoards(ctx)
	if err != nil {
		log.Print(err)
		return
	}

	log.Println("Getting items")
	items := c.getItemsForBoards(ctx, boards)
	log.Println(items)
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve workspaces: %s", string(responseBytes))
	}

	return resp, nil
}

func (c *Client) getItemsForBoards(ctx context.Context, boards []Board) []Item {
	wg := sync.WaitGroup{}
	wg.Add(len(boards))

	var items []Item
	for _, b := range boards {
		b := b
		go func() {
			defer wg.Done()
			c.sink.Push(b.ToEntity(c.ownerId))
			newItems, err := c.GetItems(ctx, b.Id)
			if err != nil {
				log.Println(err)
				return
			}
			for _, i := range newItems {
				i := i
				c.sink.Push(i.ToEntity(b.Id))
			}
			items = append(items, newItems...)
		}()
	}
	wg.Wait()

	return items
}
