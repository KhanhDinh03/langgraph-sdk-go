package client

import (
	"context"
	"testing"

	"github.com/KhanhDinh03/langgraph-sdk-go/http"
	"github.com/stretchr/testify/assert"
)

func TestStoreClient_PutItem(t *testing.T) {
	httpClient := http.NewHttpClient("http://localhost:2024", nil, 0, nil)
	client := NewStoreClient(httpClient)

	err := client.PutItem(context.Background(), []string{"namespace"}, "key", map[string]any{"value": "test"}, nil)

	assert.NoError(t, err, "Expected no error when putting an item")
}
