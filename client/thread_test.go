package client

import (
	"context"
	"testing"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/stretchr/testify/assert"
)

func TestThreadsClient_Get(t *testing.T) {
	httpClient := http.NewHttpClient("http://localhost:2024", nil, 0, nil)
	client := NewThreadsClient(httpClient)

	threadID := "test-thread-id"
	_, err := client.Get(context.Background(), threadID, nil)

	assert.NoError(t, err, "Expected no error when fetching thread")
}
