package domo

import "fmt"

// DatasetDetails contains basic data about a domo Dataset.
type DatasetDetails struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Columns       int    `json:"columns"`
	Rows          int    `json:"rows"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	DataCurrentAt string `json:"dataCurrentAt"`
}

// GetDatasets gets Domo Datasets Lists given limit and offset.
func (c *Client) GetDatasets(limit int, offset int) ([]*DatasetDetails, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets?limit=%d&offset=%d", c.baseURL, limit, offset)

	var d []*DatasetDetails

	// Note that this is using lowercase `get` which is a wrapper aroung `Get` defined in domo.go
	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
