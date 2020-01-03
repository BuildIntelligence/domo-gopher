package domo

import (
	"bytes"
	"encoding/json"
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

// CreateNewGroup creates a new domo group with the given name
func (c *Client) CreateNewGroup(name string) (*Group, error) {
	domoURL := fmt.Sprintf("%s/v1/groups", c.baseURL)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result Group
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGroupInfo returns a domo group for a given groupID
func (c *Client) GetGroupInfo(groupID int) (*Group, error) {
	domoURL := fmt.Sprintf("%s/v1/groups/%d", c.baseURL, groupID)
	// 200 OK success
	var d *Group

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// UpdateGroup updates the name and/or active status of the given domo Group ID
func (c *Client) UpdateGroup(groupID int, updatedGroup Group) error {
	domoURL := fmt.Sprintf("%s/v1/groups/%d", c.baseURL, groupID)
	type fields struct {
		Name   string `json:"name,omitempty"`
		Active bool   `json:"active,omitempty"`
	}
	body := fields{Name: updatedGroup.Name, Active: updatedGroup.Active}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", domoURL, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 200)
	return err
}

// DeleteGroup deletes a domo group with the given ID
func (c *Client) DeleteGroup(groupID int) error {
	domoURL := fmt.Sprintf("%s/v1/groups/%d", c.baseURL, groupID)
	// 204 No Content
	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}

// ListGroups lists domo groups for a range based off given lim and offset
func (c *Client) ListGroups(lim, offset int) ([]*Group, error) {

	// 200 ok
	domoURL := fmt.Sprintf("%s/v1/groups?offset=%d&limit=%d", c.baseURL, offset, lim)

	var d []*Group

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// AddUserToGroup adds a given user to a specified group
func (c *Client) AddUserToGroup(groupID, userID int) error {
	// 204 No Content
	domoURL := fmt.Sprintf("%s/v1/groups/%d/users/%d", c.baseURL, groupID, userID)

	req, err := http.NewRequest("PUT", domoURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}

// ListUsersInGroup returns a list of userIDs for users in the group
func (c *Client) ListUsersInGroup(groupID int) ([]*int, error) {

	// 200 OK [ id, id ]
	domoURL := fmt.Sprintf("%s/v1/groups/%d/users", c.baseURL, groupID)
	var d []*int

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// RemoveUserFromGroup removes a given user from specified group
func (c *Client) RemoveUserFromGroup(groupID, userID int) error {
	// 204 No Content
	domoURL := fmt.Sprintf("%s/v1/groups/%d/users/%d", c.baseURL, groupID, userID)

	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}
