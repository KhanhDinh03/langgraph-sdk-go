package client

import (
	"context"
	"testing"

	"github.com/KhanhDinh03/langgraph-sdk-go/http"
	"github.com/KhanhDinh03/langgraph-sdk-go/schema"
	"github.com/stretchr/testify/assert"
)

func TestRunsClient_Create(t *testing.T) {
	httpClient := http.NewHttpClient("http://localhost:2024", nil, 0, nil)
	client := NewRunsClient(httpClient)
	_, err := client.Create(
		context.Background(),
		"",
		"test-assistant-id",
		map[string]any{},
		schema.Command{},
		"",
		false,
		map[string]any{},
		schema.Config{},
		schema.Checkpoint{},
		"",
		"",
		"",
		"",
		schema.MultitaskStrategy(""),
		schema.IfNotExists(""),
		schema.OnCompletionBehavior(""),
		0,
	)
	assert.NoError(t, err, "Expected no error when creating a run")
}
