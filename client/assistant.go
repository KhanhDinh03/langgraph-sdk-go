package client

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
	"github.com/go-resty/resty/v2"
)

// Client for managing assistants in LangGraph.
//
// This class provides methods to interact with assistants, which are versioned configurations of your graph.
//
// Example:
//
//	client := langgraph.GetClient()
//	assistant, err := client.Assistants.Get("assistant-id")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(assistant)
type AssistantsClient struct {
	http *http.HttpClient
}

func NewAssistantsClient(httpClient *http.HttpClient) *AssistantsClient {
	return &AssistantsClient{http: httpClient}
}

// Get retrieves an assistant by its ID.
//
// Args:
//
//	assistantID: The ID of the assistant to retrieve
//
// Returns:
//
//	schema.Assistant: The assistant object
//	error: Any error encountered during the API request
//
// Example:
//
//	 ```go
//	 ctx := context.Background()
//	 assistant, err := client.Assistants.Get(ctx, "my_assistant_id")
//	 if err != nil {
//	   fmt.Println(err)
//	 }
//	 fmt.Println(assistant)
//	 ```
//	```json
//	  {
//	       "assistant_id": "my_assistant_id",
//	       "graph_id": "agent",
//	       "created_at": "2024-06-25T17:10:33.109781+00:00",
//	       "updated_at": "2024-06-25T17:10:33.109781+00:00",
//	       "config": {},
//	       "metadata": {"created_by": "system"}
//	   }
//	 ```
func (c *AssistantsClient) Get(ctx context.Context, assistantID string, headers map[string]string) (schema.Assistant, error) {
	resp, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s", assistantID), nil, &headers)
	if err != nil {
		return schema.Assistant{}, err
	}
	var assistant schema.Assistant
	err = json.Unmarshal(resp.Body(), &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}
	return assistant, nil
}

// Get the graph of an assistant by ID.
//
// Args:
//
//	assistantID: The ID of the assistant to retrieve the graph for
//	xray: Whether to include xray data in the response
//
// Returns:
//
//	schema.Graph: The graph of the assistant
//	error: Any error encountered during the API request
//
// Example:
//
//	 ```go
//	 ctx := context.Background()
//	 graph, err := client.Assistants.GetGraph(ctx, "assistant-id", true)
//	 if err != nil {
//	   fmt.Println(err)
//	 }
//	 fmt.Println(graph)
//	 ```
//	```json
//	{
//	    "nodes":
//	        [
//	            {"id": "__start__", "type": "schema", "data": "__start__"},
//	            {"id": "__end__", "type": "schema", "data": "__end__"},
//	            {"id": "agent","type": "runnable","data": {"id": ["langgraph", "utils", "RunnableCallable"],"name": "agent"}},
//	        ],
//	    "edges":
//	        [
//	            {"source": "__start__", "target": "agent"},
//	            {"source": "agent","target": "__end__"}
//	        ]
//	}
//
//	```
func (c *AssistantsClient) GetGraph(ctx context.Context, assistantID string, xray any, headers map[string]string) (schema.Graph, error) {
	params := url.Values{}
	params.Add("xray", fmt.Sprintf("%v", xray))

	resp, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s/graph", assistantID), params, &headers)
	if err != nil {
		return schema.Graph{}, err
	}

	var graph schema.Graph
	err = json.Unmarshal(resp.Body(), &graph)
	if err != nil {
		return schema.Graph{}, err
	}

	return graph, nil
}

// Get the schemas of an assistant by ID.
//
// Args:
//
//	assistantID: The ID of the assistant to retrieve the schemas for
//
// Returns:
//
//	schema.GraphSchema: The schemas of the assistant
//	error: Any error encountered during the API request
//
// Example:
//
//	 ```go
//	 ctx := context.Background()
//	 schemas, err := client.Assistants.GetSchemas(ctx, "assistant-id")
//	 if err != nil {
//	   fmt.Println(err)
//	 }
//	 fmt.Println(schemas)
//	 ```
//	```json
//
//	{
//	                "graph_id": "agent",
//	                "state_schema":
//	                    {
//	                        "title": "LangGraphInput",
//	                        "$ref": "#/definitions/AgentState",
//	                        "definitions":
//	                            {
//	                                "BaseMessage":
//	                                    {
//	                                        "title": "BaseMessage",
//	                                        "description": "Base abstract Message class. Messages are the inputs and outputs of ChatModels.",
//	                                        "type": "object",
//	                                        "properties":
//	                                            {
//	                                             "content":
//	                                                {
//	                                                    "title": "Content",
//	                                                    "anyOf": [
//	                                                        {"type": "string"},
//	                                                        {"type": "array","items": {"anyOf": [{"type": "string"}, {"type": "object"}]}}
//	                                                    ]
//	                                                },
//	                                            "additional_kwargs":
//	                                                {
//	                                                    "title": "Additional Kwargs",
//	                                                    "type": "object"
//	                                                },
//	                                            "response_metadata":
//	                                                {
//	                                                    "title": "Response Metadata",
//	                                                    "type": "object"
//	                                                },
//	                                            "type":
//	                                                {
//	                                                    "title": "Type",
//	                                                    "type": "string"
//	                                                },
//	                                            "name":
//	                                                {
//	                                                    "title": "Name",
//	                                                    "type": "string"
//	                                                },
//	                                            "id":
//	                                                {
//	                                                    "title": "Id",
//	                                                    "type": "string"
//	                                                }
//	                                            },
//	                                        "required": ["content", "type"]
//	                                    },
//	                                "AgentState":
//	                                    {
//	                                        "title": "AgentState",
//	                                        "type": "object",
//	                                        "properties":
//	                                            {
//	                                                "messages":
//	                                                    {
//	                                                        "title": "Messages",
//	                                                        "type": "array",
//	                                                        "items": {"$ref": "#/definitions/BaseMessage"}
//	                                                    }
//	                                            },
//	                                        "required": ["messages"]
//	                                    }
//	                            }
//	                    },
//	                "config_schema":
//	                    {
//	                        "title": "Configurable",
//	                        "type": "object",
//	                        "properties":
//	                            {
//	                                "model_name":
//	                                    {
//	                                        "title": "Model Name",
//	                                        "enum": ["anthropic", "openai"],
//	                                        "type": "string"
//	                                    }
//	                            }
//	                    }
//	            }
//
//	```
func (c *AssistantsClient) GetSchemas(ctx context.Context, assistantID string, headers map[string]string) (schema.GraphSchema, error) {
	resp, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s/schemas", assistantID), nil, &headers)
	if err != nil {
		return schema.GraphSchema{}, err
	}

	var graphSchema schema.GraphSchema

	err = json.Unmarshal(resp.Body(), &graphSchema)
	if err != nil {
		return schema.GraphSchema{}, err
	}

	return graphSchema, nil
}

// Get the subgraphs of an assistant by ID.
//
// Args:
//
//	assistantID: The ID of the assistant to retrieve the subgraphs for
//
// Returns:
//
//	schema.Subgraphs: The subgraphs of the assistant
//	error: Any error encountered during the API request
func (c *AssistantsClient) GetSubgraphs(ctx context.Context, assistantID string, namespace string, recurse bool, headers map[string]string) (schema.Subgraphs, error) {
	// 	if namespace is not None:
	// 	return await self.http.get(
	// 		f"/assistants/{assistant_id}/subgraphs/{namespace}",
	// 		params={"recurse": recurse},
	// 		headers=headers,
	// 	)
	// else:
	// 	return await self.http.get(
	// 		f"/assistants/{assistant_id}/subgraphs",
	// 		params={"recurse": recurse},
	// 		headers=headers,
	// 	)
	var (
		resp *resty.Response
		err  error
	)

	if namespace != "" {
		resp, err = c.http.Get(ctx, fmt.Sprintf("/assistants/%s/subgraphs/%s", assistantID, namespace), nil, &headers)
	} else {
		resp, err = c.http.Get(ctx, fmt.Sprintf("/assistants/%s/subgraphs", assistantID), nil, &headers)
	}

	var subgraphs schema.Subgraphs
	err = json.Unmarshal(resp.Body(), &subgraphs)
	if err != nil {
		return schema.Subgraphs{}, err
	}

	return subgraphs, nil
}

// Create an new assistant.
//
// Useful when graph is configurable and you want to create different assistants based on different configurations.
//
// Args:
//
//		graphID: The ID of the graph the assistant should use. The graph ID is normally set in your langgraph.json configuration.
//		config: Configuration to use for the graph.
//		metadata: Metadata to add to assistant.
//		assistantID: Assistant ID to use, will default to a random UUID if not provided.
//		ifExists: How to handle duplicate creation. Defaults to "raise" under the hood.
//	       	  Must be either "raise" (raise error if duplicate), or "do_nothing" (return existing assistant).
//		name: The name of the assistant. Defaults to "Untitled" under the hood.
//
// Returns:
//
//	schema.Assistant: The assistant object
//	error: Any error encountered during the API request
//
// Example:
//
//		```go
//		ctx := context.Background()
//		assistant, err := client.Assistants.Create(
//				ctx,
//				"agent",
//				&schema.Config{"configurable": {"model_name": "openai"}},
//				{"number":1},
//				 "my-assistant-id",
//	 			string(schema.OnConflictBehaviorDoNothing),
//				"my-name")
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(assistant)
//		```
//		```json
//		{
//			"assistant_id": "my-assistant-id",
//			"graph_id": "agent",
//			"created_at": "2024-06-25T17:10:33.109781+00:00",
//			"updated_at": "2024-06-25T17:10:33.109781+00:00",
//			"config": {},
//			"metadata": {"number": 1}
//		}
//		```
func (c *AssistantsClient) Create(
	ctx context.Context,
	graphID string,
	config *schema.Config,
	metadata schema.Json,
	assistantID string,
	ifExists schema.OnConflictBehavior,
	name string,
	headers map[string]string,
	description string,
) (schema.Assistant, error) {
	payload := map[string]any{
		"graph_id": graphID,
	}
	if config != nil {
		payload["config"] = config
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if assistantID != "" {
		payload["assistant_id"] = assistantID
	}
	if ifExists != "" {
		payload["if_exists"] = ifExists
	}
	if name != "" {
		payload["name"] = name
	}
	if description != "" {
		payload["description"] = description
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "/assistants", payload, &headers)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(resp.Body(), &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}

// Update an assistant.
//
// Use this to point to a different graph, update the configuration, or change the metadata of an assistant.
//
// Args:
//
//	assistantID: The ID of the assistant to update
//	graphID: The ID of the graph the assistant should use.
//			 The graph ID is normally set in your langgraph.json configuration. If None, assistant will keep pointing to same graph.
//	config: Configuration to use for the graph.
//	metadata: Metadata to merge with existing assistant metadata.
//	name: The name of the assistant.
//
// Returns:
//
//	Assistant: The updated assistant.
//	error: Any error encountered during the API request.
//
// Example:
//
//	 ```go
//	 ctx := context.Background()
//	 assistant, err := client.Assistants.Update(
//		 ctx,
//		 "e280dad7-8618-443f-87f1-8e41841c180f",
//		 "other-graph",
//		 &schema.Config{"configurable": {"model_name": "openai"}},
//		 {"number":1},
//		 "")
//	 if err != nil {
//	   fmt.Println(err)
//	 }
//	 fmt.Println(assistant)
//	 ```
func (c *AssistantsClient) Update(
	ctx context.Context,
	assistantID string,
	graphID string,
	config *schema.Config,
	metadata schema.Json,
	name string,
	headers map[string]string,
	description string,
) (schema.Assistant, error) {
	payload := map[string]any{}
	if graphID != "" {
		payload["graph_id"] = graphID
	}
	if config != nil {
		payload["config"] = config
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if name != "" {
		payload["name"] = name
	}
	if description != "" {
		payload["description"] = description
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Patch(ctx, fmt.Sprintf("/assistants/%s", assistantID), payload, &headers)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(resp.Body(), &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}

// Delete an assistant by ID.
//
// Args:
//
//	assistantID: The ID of the assistant to delete
//
// Returns:
//
//	 ctx: The context to use for the request
//		error: Any error encountered during the API request
//
// Example:
//
//		```go
//	 ctx := context.Background()
//		err := client.Assistants.Delete(ctx, "assistant-id")
//		if err != nil {
//		  fmt.Println(err)
//		}
//		```
func (c *AssistantsClient) Delete(ctx context.Context, assistantID string, headers map[string]string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/assistants/%s", assistantID), nil, &headers)
	if err != nil {
		return err
	}

	return nil
}

// Search for assistants.
//
// Args:
//
//	metadata: Metadata to filter by. Exact match filter for each key-value pair.
//	graphID: The ID of the graph to filter by.
//			The graph ID is normally set in your langgraph.json configuration.
//	limit: The maximum number of assistants to return. Defaults to 10.
//	offset: The number of results to skip. Defaults to 0.
//
// Returns:
//
//	[]schema.Assistant: The list of assistants that match the search criteria.
//	error: Any error encountered during the API request.
//
// Example:
//
//	```go
//	ctx := context.Background()
//	assistants, err := client.Assistants.Search(
//		ctx,
//		{"created_by": "system"},
//		"agent",
//		10,
//		0)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(assistants)
//	```
func (c *AssistantsClient) Search(
	ctx context.Context,
	metadata schema.Json,
	graphID string,
	limit int,
	offset int,
	sortBy schema.AssistantSortBy,
	sortOrder schema.SortOrder,
	headers map[string]string,
) ([]schema.Assistant, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"limit":  limit,
		"offset": offset,
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if graphID != "" {
		payload["graph_id"] = graphID
	}
	if sortBy != "" {
		payload["sort_by"] = sortBy
	}
	if sortOrder != "" {
		payload["sort_order"] = sortOrder
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "/assistants/search", payload, &headers)
	if err != nil {
		return []schema.Assistant{}, err
	}

	var assistants []schema.Assistant

	err = json.Unmarshal(resp.Body(), &assistants)
	if err != nil {
		return []schema.Assistant{}, err
	}

	return assistants, nil
}

// List all versions of an assistant by ID.
//
// Args:
//
//	assistantID: The ID of the assistant to retrieve versions for
//	metadata: Metadata to filter by. Exact match filter for each key-value pair.
//	limit: The maximum number of versions to return. Defaults to 10.
//	offset: The number of results to skip. Defaults to 0.
//
// Returns:
//
//	[]schema.Assistant: The list of assistants that match the search criteria.
//	error: Any error encountered during the API request.
//
// Example:
//
//	```go
//	ctx := context.Background()
//	assistants, err := client.Assistants.GetVersions(ctx, "assistant-id", nil, 10, 0)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(assistants)
//	```
func (c *AssistantsClient) GetVersions(
	ctx context.Context,
	assistantID string,
	metadata schema.Json,
	limit int,
	offset int,
	headers map[string]string,
) ([]schema.Assistant, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"limit":  limit,
		"offset": offset,
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, fmt.Sprintf("/assistants/%s/versions", assistantID), payload, &headers)
	if err != nil {
		return []schema.Assistant{}, err
	}

	var assistants []schema.Assistant

	err = json.Unmarshal(resp.Body(), &assistants)
	if err != nil {
		return []schema.Assistant{}, err
	}

	return assistants, nil
}

// Change the latest version of an assistant.
//
// Args:
//
//	assistantID: The ID of the assistant to set the latest version for
//	version: The version number to set as the latest
//
// Returns:
//
//	Assistant: Assistant Object
//	error: Any error encountered during the API request
//
// Example:
//
//	```go
//	ctx := context.Background()
//	assistant, err := client.Assistants.SetLatest(ctx, "assistant-id", 1)
//	if err != nil {
//	  fmt.Println(err)
//	}
//	fmt.Println(assistant)
//	```
func (c *AssistantsClient) SetLatest(ctx context.Context, assistantID string, version int, headers map[string]string) (schema.Assistant, error) {

	payload := map[string]any{
		"version": version,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, fmt.Sprintf("/assistants/%s/versions/latest", assistantID), payload, &headers)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant

	err = json.Unmarshal(resp.Body(), &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}
