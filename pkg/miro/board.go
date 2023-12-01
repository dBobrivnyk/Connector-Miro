package miro

import (
	"connectors/pkg/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeBoard = "board"

type Boards struct {
	Data []Board `json:"data"`
}

type Member struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type Membership struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type SharingPolicy struct {
	Access                            string `json:"access"`
	InviteToAccountAndBoardLinkAccess string `json:"inviteToAccountAndBoardLinkAccess"`
	OrganizationAccess                string `json:"organizationAccess"`
	TeamAccess                        string `json:"teamAccess"`
}

type PermissionsPolicy struct {
	CollaborationToolsStartAccess string `json:"collaborationToolsStartAccess"`
	CopyAccess                    string `json:"copyAccess"`
	CopyAccessLevel               string `json:"copyAccessLevel"`
	SharingAccess                 string `json:"sharingAccess"`
}

type Policy struct {
	PermissionsPolicy PermissionsPolicy `json:"permissionsPolicy"`
	SharingPolicy     SharingPolicy     `json:"sharingPolicy"`
}

type Links struct {
	Self    string `json:"self"`
	Related string `json:"related,omitempty"`
}

type Board struct {
	Id                    string            `json:"id"`
	Type                  string            `json:"type"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	Links                 Links             `json:"links"`
	CreatedAt             time.Time         `json:"createdAt"`
	CreatedBy             Member            `json:"createdBy"`
	CurrentUserMembership Membership        `json:"currentUserMembership"`
	ModifiedAt            time.Time         `json:"modifiedAt"`
	ModifiedBy            Member            `json:"modifiedBy"`
	Owner                 Member            `json:"owner"`
	PermissionsPolicy     PermissionsPolicy `json:"permissionsPolicy"`
	Policy                Policy            `json:"policy"`
	SharingPolicy         SharingPolicy     `json:"sharingPolicy"`
	Team                  Member            `json:"team"`
	ViewLink              string            `json:"viewLink"`
}

func (c *Client) GetBoards(ctx context.Context) ([]Board, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"boards", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve boards data: %w", err)
	}
	defer resp.Body.Close()

	var boards Boards
	err = json.NewDecoder(resp.Body).Decode(&boards)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return boards.Data, nil
}

func (b Board) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         b.Name,
		EntityUrl:    "",
		ExternalId:   b.Id,
		Type:         EntityTypeBoard,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         b,
	}
}
