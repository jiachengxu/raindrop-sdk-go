package raindrop

import (
	"fmt"
	"time"
)

const (
	// Access levels
	AccessLevelPublicReadOnly        = 1 // read only access (equal to public=true)
	AccessLevelCollaboratorReadOnly  = 2 // collaborator with read only access
	AccessLevelCollaboratorWriteOnly = 3 // collaborator with write only access
	AccessLevelOwner                 = 4 // owner

	// View types
	ViewList    = "list"
	ViewSimple  = "simple"
	ViewGrid    = "grid"
	ViewMasonry = "masonry" // Pinterest like grid

	// System collections
	CollectionUnsorted = -1
	CollectionTrash    = -99
)

type Collection struct {
	ID            int            `json:"_id"`
	Access        Access         `json:"access"`
	Collaborators []Collaborator `json:"collaborators,omitempty"`
	Color         string         `json:"color"`
	Count         int            `json:"count"`
	Cover         []string       `json:"cover"`
	Created       time.Time      `json:"created"`
	Expanded      bool           `json:"expanded"`
	LastUpdate    time.Time      `json:"lastUpdate"`
	Parent        IDRef          `json:"parent"`
	Public        bool           `json:"public"`
	Sort          int            `json:"sort"`
	Title         string         `json:"title"`
	View          string         `json:"view"`
}

type Access struct {
	Level     int  `json:"level"`
	Draggable bool `json:"draggable"`
}

type Collaborator struct {
	ID         int    `json:"_id"`
	Email      string `json:"email,omitempty"`
	EmailMD5   string `json:"email_MD5"`
	FullName   string `json:"fullName"`
	Registered string `json:"registered"`
	Role       string `json:"role"`
}

type IDRef struct {
	ID int `json:"$id"`
}

// https://developer.raindrop.io/v1/collections/methods#get-root-collections
func (c *Client) GetRootCollections() ([]Collection, error) {
	var result struct {
		Result bool         `json:"result"`
		Items  []Collection `json:"items"`
	}
	err := c.doRequest("GET", fmt.Sprintf("%s/collections", baseURL), nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// https://developer.raindrop.io/v1/collections/methods#get-child-collections
func (c *Client) GetChildCollections(parentCollectionID int) ([]Collection, error) {
	var result struct {
		Result bool         `json:"result"`
		Items  []Collection `json:"items"`
	}
	err := c.doRequest("GET", fmt.Sprintf("%s/collections/childrens", baseURL), nil, &result)
	if err != nil {
		return nil, err
	}
	collections := make([]Collection, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Parent.ID == parentCollectionID {
			collections = append(collections, item)
		}
	}
	return collections, nil
}

// https://developer.raindrop.io/v1/collections/methods#get-collection
func (c *Client) GetCollection(id int) (*Collection, error) {
	var result struct {
		Result bool       `json:"result"`
		Item   Collection `json:"item"`
	}
	err := c.doRequest("GET", fmt.Sprintf("%s/collection/%d", baseURL, id), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/collections/methods#create-collection
func (c *Client) CreateCollection(title string, parentCollectionID int) (*Collection, error) {
	body := map[string]interface{}{
		"title": title,
		"parent": map[string]interface{}{
			"$id": parentCollectionID,
		},
	}
	var result struct {
		Result bool       `json:"result"`
		Item   Collection `json:"item"`
	}
	err := c.doRequest("POST", fmt.Sprintf("%s/collection", baseURL), body, &result)
	if err != nil {
		return nil, err
	}
	if !result.Result {
		return nil, fmt.Errorf("failed to create collection")
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/collections/methods#update-collection
func (c *Client) UpdateCollection(id int, updates map[string]interface{}) (*Collection, error) {
	var result struct {
		Result bool       `json:"result"`
		Item   Collection `json:"item"`
	}
	err := c.doRequest("PUT", fmt.Sprintf("%s/collection/%d", baseURL, id), updates, &result)
	if err != nil {
		return nil, err
	}
	return &result.Item, nil
}

// https://developer.raindrop.io/v1/collections/methods#remove-collection
func (c *Client) RemoveCollection(id int) error {
	return c.doRequest("DELETE", fmt.Sprintf("%s/collection/%d", baseURL, id), nil, nil)
}

// https://developer.raindrop.io/v1/collections/methods#remove-multiple-collections
func (c *Client) RemoveMultipleCollections(ids []int) error {
	body := map[string]interface{}{
		"ids": ids,
	}
	return c.doRequest("DELETE", fmt.Sprintf("%s/collections", baseURL), body, nil)
}

// https://developer.raindrop.io/v1/collections/methods#reorder-all-collections
func (c *Client) ReorderAllCollections(sort string) error {
	body := map[string]interface{}{
		"sort": sort,
	}
	return c.doRequest("PUT", fmt.Sprintf("%s/collections", baseURL), body, nil)
}

// https://developer.raindrop.io/v1/collections/methods#expand-collapse-all-collections
func (c *Client) ExpandCollapseAllCollections(expanded bool) error {
	body := map[string]interface{}{
		"expanded": expanded,
	}
	return c.doRequest("PUT", fmt.Sprintf("%s/collections", baseURL), body, nil)
}

// https://developer.raindrop.io/v1/collections/methods#merge-collections
func (c *Client) MergeCollections(to int, ids []int) error {
	body := map[string]interface{}{
		"to":  to,
		"ids": ids,
	}
	return c.doRequest("PUT", fmt.Sprintf("%s/collections/merge", baseURL), body, nil)
}

// https://developer.raindrop.io/v1/collections/methods#remove-all-empty-collections
func (c *Client) RemoveAllEmptyCollections() (int, error) {
	var result struct {
		Result bool `json:"result"`
		Count  int  `json:"count"`
	}
	err := c.doRequest("PUT", fmt.Sprintf("%s/collections/clean", baseURL), nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// https://developer.raindrop.io/v1/collections/methods#empty-trash
func (c *Client) EmptyTrash() error {
	return c.doRequest("DELETE", fmt.Sprintf("%s/collection/-99", baseURL), nil, nil)
}
