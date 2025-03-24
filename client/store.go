package client

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"
	"strings"

	"github.com/KhanhDinh03/langgraph-sdk-go/http"
	"github.com/KhanhDinh03/langgraph-sdk-go/schema"
)

type StoreClient struct {
	http *http.HttpClient
}

func NewStoreClient(httpClient *http.HttpClient) *StoreClient {
	return &StoreClient{http: httpClient}
}

func (c *StoreClient) PutItem(ctx context.Context, namespace []string, key string, value map[string]any, index any) error {
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
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	_, err := c.http.Put(ctx, "/store", payload)
	return err
}

func (c *StoreClient) GetItem(ctx context.Context, namespace []string, key string) (map[string]any, error) {
	for _, label := range namespace {
		if containsDot(label) {
			return nil, fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	params := url.Values{}
	params.Add("namespace", strings.Join(namespace, "."))
	params.Add("key", key)

	resp, err := c.http.Get(ctx, "/store", params)
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

func (c *StoreClient) DeleteItem(ctx context.Context, namespace []string, key string) error {
	for _, label := range namespace {
		if containsDot(label) {
			return fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	params := url.Values{}
	params.Add("namespace", strings.Join(namespace, "."))
	params.Add("key", key)

	err := c.http.Delete(ctx, "/store", params)
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
) (schema.SearchItemsResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"namespace": namespace,
		"filter":    filter,
		"limit":     limit,
		"offset":    offset,
		"query":     query,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	ctx := context.Background()
	resp, err := c.http.Post(ctx, "/store/items/search", payload)
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

func (c *StoreClient) ListNamespaces(ctx context.Context, prefix []string, suffix []string, maxDepth int, limit int, offset int) ([]schema.ListNamespaceResponse, error) {
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

	resp, err := c.http.Post(ctx, "/store/namespaces", payload)
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
