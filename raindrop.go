package raindrop

import (
	"fmt"
	"strings"
	"time"
)

const (
	// Raindrop types
	RaindropTypeLink     = "link"
	RaindropTypeArticle  = "article"
	RaindropTypeImage    = "image"
	RaindropTypeVideo    = "video"
	RaindropTypeDocument = "document"
	RaindropTypeAudio    = "audio"
)

type Raindrop struct {
	ID         int       `json:"_id"`
	Collection IDRef     `json:"collection"`
	Cover      string    `json:"cover"`
	Created    time.Time `json:"created"`
	Domain     string    `json:"domain"`
	Excerpt    string    `json:"excerpt"`
	Note       string    `json:"note"`
	LastUpdate time.Time `json:"lastUpdate"`
	Link       string    `json:"link"`
	Media      []Media   `json:"media"`
	Tags       []string  `json:"tags"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	User       IDRef     `json:"user"`
}

type Media struct {
	Link string `json:"link"`
}

// https://developer.raindrop.io/v1/raindrops/single#get-raindrop
func (c *Client) GetRaindrop(id int) (*Raindrop, error) {
	var result struct {
		Result bool     `json:"result"`
		Item   Raindrop `json:"item"`
	}
	err := c.doRequest("GET", fmt.Sprintf("%s/raindrop/%d", baseURL, id), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/raindrops/single#create-raindrop
func (c *Client) CreateRaindrop(raindrop *Raindrop) (*Raindrop, error) {
	var result struct {
		Result bool     `json:"result"`
		Item   Raindrop `json:"item"`
	}
	err := c.doRequest("POST", fmt.Sprintf("%s/raindrop", baseURL), raindrop, &result)
	if err != nil {
		return nil, err
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/raindrops/single#update-raindrop
func (c *Client) UpdateRaindrop(id int, updates map[string]interface{}) (*Raindrop, error) {
	var result struct {
		Result bool     `json:"result"`
		Item   Raindrop `json:"item"`
	}
	err := c.doRequest("PUT", fmt.Sprintf("%s/raindrop/%d", baseURL, id), updates, &result)
	if err != nil {
		return nil, err
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/raindrops/single#remove-raindrop
func (c *Client) RemoveRaindrop(id int) error {
	return c.doRequest("DELETE", fmt.Sprintf("%s/raindrop/%d", baseURL, id), nil, nil)
}

// https://developer.raindrop.io/v1/raindrops/multiple#get-raindrops
func (c *Client) GetRaindrops(collectionID int, params map[string]string) ([]Raindrop, error) {
	url := fmt.Sprintf("%s/raindrops/%d", baseURL, collectionID)

	// Add query parameters
	if len(params) > 0 {
		query := make([]string, 0, len(params))
		for k, v := range params {
			query = append(query, fmt.Sprintf("%s=%s", k, v))
		}
		url += "?" + strings.Join(query, "&")
	}

	var result struct {
		Result       bool       `json:"result"`
		Items        []Raindrop `json:"items"`
		Count        int        `json:"count"`
		Collectionid int        `json:"collectionId"`
	}
	err := c.doRequest("GET", url, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// https://developer.raindrop.io/v1/raindrops/multiple#create-many-raindrops
func (c *Client) CreateRaindrops(raindrops []Raindrop) ([]Raindrop, error) {
	if len(raindrops) > 100 {
		return nil, fmt.Errorf("maximum number of raindrops to create is 100")
	}
	var result struct {
		Result bool       `json:"result"`
		Items  []Raindrop `json:"items"`
	}
	err := c.doRequest("POST", fmt.Sprintf("%s/raindrops", baseURL), raindrops, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// https://developer.raindrop.io/v1/raindrops/multiple#update-many-raindrops
func (c *Client) UpdateRaindrops(collectionID int, updates map[string]interface{}, ids []int) (int, error) {
	updates["ids"] = ids
	var result struct {
		Result   bool `json:"result"`
		Modified int  `json:"modified"`
	}
	err := c.doRequest("PUT", fmt.Sprintf("%s/raindrops/%d", baseURL, collectionID), updates, &result)
	if err != nil {
		return 0, err
	}
	return result.Modified, nil
}

// https://developer.raindrop.io/v1/raindrops/multiple#remove-many-raindrops
func (c *Client) RemoveRaindrops(collectionID int, ids []int) (int, error) {
	var result struct {
		Result   bool `json:"result"`
		Modified int  `json:"modified"`
	}
	body := map[string]interface{}{"ids": ids}
	err := c.doRequest("DELETE", fmt.Sprintf("%s/raindrops/%d", baseURL, collectionID), body, &result)
	if err != nil {
		return 0, err
	}
	return result.Modified, nil
}
