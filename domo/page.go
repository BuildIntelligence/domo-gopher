package domo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Page models the Domo Page object
type Page struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	ParentID      int    `json:"parentId,omitempty"`
	OwnerID       int    `json:"ownerId,omitempty"`
	Locked        bool   `json:"locked,omitempty"`
	CollectionIDs []*int `json:"collectionIds,omitempty"`
	CardIDs       []*int `json:"cardIds,omitempty"`
	//Children
	//Visibility
	UserIDs  []*int `json:"userIds,omitempty"`
	GroupIDs []*int `json:"groupIds,omitempty"`
}

// PageCollection models the Domo Page Collection Object
type PageCollection struct {
	PageID           int    `json:"pageId,omitempty"`
	PageCollectionID int    `json:"page_collection_id,omitempty"`
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	CardIDs          []*int `json:"cardIds,omitempty"`
}

// GetPageInfo retrieves the info for a Domo Page for the given Domo pageID.
func (c *Client) GetPageInfo(pageID int) (*Page, error) {
	domoURL := fmt.Sprintf("%s/v1/pages/%d", c.baseURL, pageID)

	var p *Page

	err := c.get(domoURL, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// CreatePage
func (c *Client) CreatePage(page Page) (*Page, error) {
	// name, parentID, locked, cardIds, visibility, userIds, groupIds
	domoURL := fmt.Sprintf("%s/v1/pages", c.baseURL)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(page)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result Page
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdatePage
func (c *Client) UpdatePage(page Page) error {
	domoURL := fmt.Sprintf("%s/v1/pages/%d", c.baseURL, page.ID)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(page)
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

// DeletePage deletes a Domo Page given it's pageID.
func (c *Client) DeletePage(pageID int) error {
	domoURL := fmt.Sprintf("%s/v1/pages/%d", c.baseURL, pageID)
	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}

// ListPages retrieves an array of Domo Page info for a given limit and offset.
func (c *Client) ListPages(lim, offset int) ([]*Page, error) {
	domoURL := fmt.Sprintf("%s/v1/pages?limit=%d&offset=%d", c.baseURL, lim, offset)

	var p []*Page

	err := c.get(domoURL, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetPageCollection retrieves Domo Page collections for a given Domo pageID
func (c *Client) GetPageCollection(pageID int) ([]*PageCollection, error) {
	domoURL := fmt.Sprintf("%s/v1/pages/%d/collections", c.baseURL, pageID)

	var p []*PageCollection

	err := c.get(domoURL, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// CreatePageCollection creates a Domo Page Collection on the given Domo Page specified by pageID.
func (c *Client) CreatePageCollection(pageID int, collection PageCollection) error {
	domoURL := fmt.Sprintf("%s/v1/pages/%d/collections", c.baseURL, pageID)

	// create new struct to strip out fields except Title, Description, CardIDs
	nCollection := PageCollection{Title: collection.Title, Description: collection.Description, CardIDs: collection.CardIDs}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(nCollection)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 200)
	return err
}

// UpdatePageCollection
func (c *Client) UpdatePageCollection(pageID int, collection PageCollection) error {
	domoURL := fmt.Sprintf("%s/v1/pages/%d/collections/%d", c.baseURL, pageID, collection.PageCollectionID)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(collection)
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

// DeletePageCollection deletes a Domo Page Collection from the Domo Page specified via pageID.
func (c *Client) DeletePageCollection(pageID, pageCollectionID int) error {
	domoURL := fmt.Sprintf("%s/v1/pages/%d/collections/%d", c.baseURL, pageID, pageCollectionID)
	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	return err
}
