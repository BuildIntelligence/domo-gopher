package domo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// User is the object for a Domo User
type User struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Role           string `json:"role,omitempty"`
	Title          string `json:"title,omitempty"`
	AlternateEmail string `json:"alternateEmail,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Location       string `json:"location,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	ImageURI       string `json:"image,omitempty"`
	EmployeeNumber int    `json:"employeeNumber,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	// Groups []DomoGroup `json:"groups,omitempty"`
}

// CreateNewUser creates a new Domo User and sends an email invite if set to send invite
func (c *Client) CreateNewUser(user User, sendInvite bool) (*User, error) {
	domoURL := fmt.Sprintf("%s/v1/users?sendInvite=%t", c.baseURL, sendInvite)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(user)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result User
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserInfo retrieves a given Domo User by userID.
func (c *Client) GetUserInfo(userID int) (*User, error) {
	domoURL := fmt.Sprintf("%s/v1/users/%d", c.baseURL, userID)

	var d *User

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// UpdateUser Updates an existing User of the given userID.
func (c *Client) UpdateUser(userID int, updatedFields User) error {
	domoURL := fmt.Sprintf("%s/v1/users/%d", c.baseURL, userID)

	buf := new(bytes.Buffer) // I think this could be buf, err := json.Marshal(dataset) instead
	err := json.NewEncoder(buf).Encode(updatedFields)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", domoURL, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}

// DeleteUser deletes a Domo User for the given userID.
func (c *Client) DeleteUser(userID int) error {
	domoURL := fmt.Sprintf("%s/v1/users/%d", c.baseURL, userID)
	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}

// ListUsers returns a list of users based on the given limit and start offset.
func (c *Client) ListUsers(limit int, offset int) ([]*User, error) {
	domoURL := fmt.Sprintf("%s/v1/users?limit=%d&offset=%d", c.baseURL, limit, offset)

	var d []*User

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
