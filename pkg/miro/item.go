package miro

import (
	"connectors/pkg/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeItem = "item"

type Items struct {
	Data []Item `json:"data"`
}

type Geometry struct {
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	Width    float64 `json:"width"`
}

type Position struct {
	Origin     string  `json:"origin"`
	RelativeTo string  `json:"relativeTo"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
}

type Item struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	CreatedBy Member `json:"createdBy"`
	Data      struct {
		Content string `json:"content"`
	} `json:"data"`
	Geometry   Geometry  `json:"geometry"`
	Links      Links     `json:"links"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy Member    `json:"modifiedBy"`
	Parent     struct {
		ID    int64 `json:"id"`
		Links Links `json:"links"`
	} `json:"parent"`
	Position Position `json:"position"`
	Type     string   `json:"type"`
}

func (c *Client) GetItems(ctx context.Context, boardId string) ([]Item, error) {
	url := fmt.Sprintf("%sboards/%s/items", c.url, boardId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items data for boardId %s: %w", boardId, err)
	}
	defer resp.Body.Close()

	var items Items
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response for boardId %s: %w", boardId, err)
	}

	return items.Data, nil
}

func (i Item) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         i.Data.Content,
		EntityUrl:    "",
		ExternalId:   i.Id,
		Type:         EntityTypeItem,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: i.ModifiedAt,
		Data:         i,
	}
}
