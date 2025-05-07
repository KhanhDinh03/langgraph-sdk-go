package client

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"
	"strings"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type StoreClient struct {
	http *http.HttpClient
}

func NewStoreClient(httpClient *http.HttpClient) *StoreClient {
	return &StoreClient{http: httpClient}
}

func (c *StoreClient) PutItem(ctx context.Context, namespace []string, key string, value map[string]any, index any, ttl int, headers map[string]string) error {
	for _, label := range namespace {
		if containsDot(label) {
			return fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	payload := map[string]any{
		"namespace": namespace,
		"key":       key,
		"value":     value,
		"index":     index,
		"ttl":       ttl,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	_, err := c.http.Put(ctx, "/store/items", payload, &headers)
	return err
}

func (c *StoreClient) GetItem(ctx context.Context, namespace []string, key string, refreshTtl bool, headers map[string]string) (map[string]any, error) {
	for _, label := range namespace {
		if containsDot(label) {
			return nil, fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	params := url.Values{}
	params.Add("namespace", strings.Join(namespace, "."))
	params.Add("key", key)

	if refreshTtl {
		params.Add("refresh_ttl", "true")
	}

	resp, err := c.http.Get(ctx, "/store/items", params, &headers)
	if err != nil {
		return nil, err
	}

	var item map[string]any
	err = json.Unmarshal(resp.Body(), &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

//Delete an item.
//
// Args:
// 	key: The unique identifier for the item.
// 	namespace: Optional list of strings representing the namespace path.
// 	headers: Optional custom headers to include with the request.

// Returns:
// 	None

// Example Usage:

// await client.store.delete_item(
//
//	["documents", "user123"],
//	key="item456",
//
// )
func (c *StoreClient) DeleteItem(ctx context.Context, namespace []string, key string, headers map[string]string) error {
	for _, label := range namespace {
		if containsDot(label) {
			return fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	jsonData := map[string]any{
		"namespace": namespace,
		"key":       key,
	}

	err := c.http.Delete(ctx, "/store/items", jsonData, &headers)
	if err != nil {
		return err
	}

	return nil
}

func (c *StoreClient) SearchItems(
	namespace []string,
	filter map[string]any,
	limit int,
	offset int,
	query string,
	refreshTtl bool,
	headers map[string]string,
) (schema.SearchItemsResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"namespace":   namespace,
		"filter":      filter,
		"limit":       limit,
		"offset":      offset,
		"query":       query,
		"refresh_ttl": refreshTtl,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	ctx := context.Background()
	resp, err := c.http.Post(ctx, "/store/items/search", payload, &headers)
	if err != nil {
		return schema.SearchItemsResponse{}, err
	}

	var searchItemsResponse schema.SearchItemsResponse

	err = json.Unmarshal(resp.Body(), &searchItemsResponse)
	if err != nil {
		return schema.SearchItemsResponse{}, err
	}

	return searchItemsResponse, nil
}

func (c *StoreClient) ListNamespaces(ctx context.Context, prefix []string, suffix []string, maxDepth int, limit int, offset int, headers map[string]string) ([]schema.ListNamespaceResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"prefix":    prefix,
		"suffix":    suffix,
		"max_depth": maxDepth,
		"limit":     limit,
		"offset":    offset,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "/store/namespaces", payload, &headers)
	if err != nil {
		return []schema.ListNamespaceResponse{}, err
	}

	var namespaces []schema.ListNamespaceResponse
	err = json.Unmarshal(resp.Body(), &namespaces)
	if err != nil {
		return []schema.ListNamespaceResponse{}, err
	}

	return namespaces, nil
}
