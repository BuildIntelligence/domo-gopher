package domo

import (
	"context"
	"fmt"
	"net/http"
)

// Group defines all the domo group related data in the Public API
type Group struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Default     bool   `json:"default,omitempty"`
	Active      bool   `json:"active,omitempty"`
	CreatorID   string `json:"creatorId,omitempty"`
	MemberCount int    `json:"memberCount,omitempty"`
	UserIDs     []int  `json:"userIds,omitempty"`
}

// GroupsService handles communication with the groups
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2
type GroupsService service

// List the groups. Limit should be between 1 and 500.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#List%20groups
func (s *GroupsService) List(ctx context.Context, limit, offset int) ([]*Group, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/groups?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var groups []*Group
	resp, err := s.client.Do(ctx, req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// Info for the group for the given group id.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Retrieve%20a%20group
func (s *GroupsService) Info(ctx context.Context, groupID int) (*Group, *http.Response, error) {
	u := fmt.Sprintf("v1/groups/%d", groupID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var group *Group
	resp, err := s.client.Do(ctx, req, &group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// Delete a domo group with the given group id.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Delete%20a%20group
func (s *GroupsService) Delete(ctx context.Context, groupID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/groups/%d", groupID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Create a new Domo Group.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Create%20a%20group
func (s *GroupsService) Create(ctx context.Context, group Group) (*Group, *http.Response, error) {
	u := "v1/groups"
	req, err := s.client.NewRequest("POST", u, group)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var newGroup *Group
	resp, err := s.client.Do(ctx, req, &newGroup)
	if err != nil {
		return nil, resp, err
	}

	return newGroup, resp, nil

}

// Update a Domo Group.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Update%20a%20group
// Updates the specified group by providing values to parameters passed.
// Any parameter left out of the request will cause the specific group’s
// attribute to remain unchanged.
func (s *GroupsService) Update(ctx context.Context, group Group) (*http.Response, error) {
	u := fmt.Sprintf("v1/groups/%d", group.ID)
	req, err := s.client.NewRequest("PUT", u, group)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil

}

// AddUser to a Domo Group.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Add%20a%20user%20to%20a%20group
func (s *GroupsService) AddUser(ctx context.Context, groupID, userID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/groups/%d/users/%d", groupID, userID)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
// RemoveUser to a Domo Group.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#Remove%20a%20user%20from%20a%20group
func (s *GroupsService) RemoveUser(ctx context.Context, groupID, userID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/groups/%d/users/%d", groupID, userID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UserIDs in the group. Limit should be between 1 and 500.
//
// Domo API Docs: https://developer.domo.com/docs/groups-api-reference/groups-2#List%20groups
func (s *GroupsService) UserIDs(ctx context.Context, groupID, limit, offset int) ([]int, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/groups/%d/users?limit=%d&offset=%d", groupID, limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var userIDs []int
	resp, err := s.client.Do(ctx, req, &userIDs)
	if err != nil {
		return nil, resp, err
	}

	return userIDs, resp, nil
}
