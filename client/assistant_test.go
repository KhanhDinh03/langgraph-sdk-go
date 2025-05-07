package client

import (
	"context"
	"testing"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/stretchr/testify/assert"
)

func TestAssistantsClient_Get(t *testing.T) {
	httpClient := http.NewHttpClient("http://localhost:2024", nil, 0, nil)
	client := NewAssistantsClient(httpClient)

	assistantID := "test-assistant-id"
	_, err := client.Get(context.Background(), assistantID, nil)

	assert.NoError(t, err, "Expected no error when fetching assistant")
}
