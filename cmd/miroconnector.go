package main

import (
	config "connectors/internal"
	"connectors/pkg/idstorage"
	"connectors/pkg/miro"
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	ctx := context.Background()

	userIds, err := idstorage.Load(idstorage.FromFile(cfg.IdStoragePath))
	if err != nil {
		log.Fatal(err)
	}

	var connectors []*miro.Client
	for _, id := range userIds {
		connectors = append(connectors, miro.NewClient(cfg, miro.ApiUrl, id.Owner, id.ApiKey))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(connectors))
	for _, c := range connectors {
		c := c
		go func() {
			defer wg.Done()
			c.GetEntities(ctx)
		}()
	}
	wg.Wait()

	ticker := time.NewTicker(time.Duration(cfg.SyncPeriod) * time.Second)
	for {
		<-ticker.C
		for _, connector := range connectors {
			go connector.GetEntities(ctx)
		}
	}
}
