package domo

import (
	"context"
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

// PagesService handles communication with the Page
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page
type PagesService service


// List the pages. Limit should be between 1 and 500.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#List%20pages
func (s *PagesService) List(ctx context.Context, limit, offset int) ([]*Page, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/pages?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var pages []*Page
	resp, err := s.client.Do(ctx, req, &pages)
	if err != nil {
		return nil, resp, err
	}

	return pages, resp, nil
}

// Info for the page for the given page id.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Retrieve%20a%20page
func (s *PagesService) Info(ctx context.Context, pageID int) (*Page, *http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d", pageID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var page *Page
	resp, err := s.client.Do(ctx, req, &page)
	if err != nil {
		return nil, resp, err
	}

	return page, resp, nil
}

// Delete a domo page with the given page id.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Delete%20a%20page
func (s *PagesService) Delete(ctx context.Context, pageID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d", pageID)
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

// Create a new Domo Page.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Create%20a%20page
func (s *PagesService) Create(ctx context.Context, page Page) (*Page, *http.Response, error) {
	u := "v1/pages"
	req, err := s.client.NewRequest("POST", u, page)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var newPage *Page
	resp, err := s.client.Do(ctx, req, &newPage)
	if err != nil {
		return nil, resp, err
	}

	return newPage, resp, nil
}

// Update a Domo Page.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Update%20a%20page
// Updates the specified page by providing values to parameters passed.
// Any parameter left out of the request will cause the specific page’s
// attribute to remain unchanged.
// Also, collections cannot be added or removed via this endpoint,
// only reordered. Giving access to a user or group will also cause
// that user or group to have access to the parent page
// (if the page is a subpage). Moving a page by updating the parentId
// will also cause everyone with access to the page to have access to the new parent page.
func (s *PagesService) Update(ctx context.Context, page Page) (*http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d", page.ID)
	req, err := s.client.NewRequest("PUT", u, page)
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

// CreateCollection on a Domo Page.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Create%20a%20page%20collection
func (s *PagesService) CreateCollection(ctx context.Context, pageID int, collection PageCollection) (*http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d/collections", pageID)
	req, err := s.client.NewRequest("POST", u, collection)
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

// UpdateCollection on a Domo Page.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Update%20a%20page%20collection
func (s *PagesService) UpdateCollection(ctx context.Context, pageID int, collection PageCollection) (*http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d/collections/%d", pageID, collection.PageCollectionID)
	req, err := s.client.NewRequest("PUT", u, collection)
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

// RemoveCollection from a Domo Page.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Delete%20a%20page%20collection
func (s *PagesService) RemoveCollection(ctx context.Context, pageID, collectionID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d/collections/%d", pageID, collectionID)
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

// Collections for the page for the given page id.
//
// Domo API Docs: https://developer.domo.com/docs/page-api-reference/page#Retrieve%20a%20page%20collection
func (s *PagesService) Collections(ctx context.Context, pageID int) ([]*PageCollection, *http.Response, error) {
	u := fmt.Sprintf("v1/pages/%d/collections", pageID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var collections []*PageCollection
	resp, err := s.client.Do(ctx, req, &collections)
	if err != nil {
		return nil, resp, err
	}

	return collections, resp, nil
}
