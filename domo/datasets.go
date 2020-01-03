package domo

import (
	"context"
	"fmt"
	"net/http"
)

// DatasetsService handles communication with the dataset
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo/com/
type DatasetsService service

// List the datasets. Limit should be between 1 and 50.
func (s *DatasetsService) List(ctx context.Context, limit, offset int) ([]*DatasetDetails, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 50 {
		return nil, nil, fmt.Errorf("limit must be 50 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/datasets?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var datasets []*DatasetDetails
	resp, err := s.client.Do(ctx, req, &datasets)
	if err != nil {
		return nil, resp, err
	}

	return datasets, resp, nil
}
