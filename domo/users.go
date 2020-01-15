package domo

import (
	"context"
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

// UsersService handles communication with the users
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2
type UsersService service

// List the users. Limit should be between 1 and 500.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2#List%20users
func (s *UsersService) List(ctx context.Context, limit, offset int) ([]*User, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/users?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var users []*User
	resp, err := s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// Info for the user for the given user id.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2#Retrieve%20a%20user
func (s *UsersService) Info(ctx context.Context, userID int) (*User, *http.Response, error) {
	u := fmt.Sprintf("v1/users/%d", userID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var usr *User
	resp, err := s.client.Do(ctx, req, &usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, nil
}

// Delete a domo user with the given user id.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2#Delete%20a%20user
func (s *UsersService) Delete(ctx context.Context, userID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/users/%d", userID)
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

// Create a new Domo User.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2#Create%20a%20user
func (s *UsersService) Create(ctx context.Context, user User, sendInvite bool) (*User, *http.Response, error) {
	u := fmt.Sprintf("v1/users?sendInvite=%t", sendInvite)
	req, err := s.client.NewRequest("POST", u, user)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var newUsr *User
	resp, err := s.client.Do(ctx, req, &newUsr)
	if err != nil {
		return nil, resp, err
	}

	return newUsr, resp, nil

}

// Update a Domo User.
//
// Domo API Docs: https://developer.domo.com/docs/users-api-reference/users-2#Update%20a%20user
// Updates the specified user by providing values to parameters passed.
// Any parameter left out of the request will cause the specific user’s
// attribute to remain unchanged.
// KNOWN LIMITATION: Currently all user fields are required.
func (s *UsersService) Update(ctx context.Context, user User) (*http.Response, error) {
	u := fmt.Sprintf("v1/users/%d", user.ID)
	req, err := s.client.NewRequest("PUT", u, user)
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
