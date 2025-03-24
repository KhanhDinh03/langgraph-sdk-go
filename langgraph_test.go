package langgraph_sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	client := GetClient("http://localhost:2024", "test-api-key", nil)

	assert.NotNil(t, client, "Expected a valid LangGraphClient instance")
}
