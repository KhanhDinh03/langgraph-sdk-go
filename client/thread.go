package client

import (
	"context"
	"encoding/json"
	"fmt"
	"langgraph-sdk/http"
	"langgraph-sdk/schema"
)

// Client for managing threads in LangGraph.
//
// A thread maintains the state of a graph across multiple interactions/invocations (aka runs).
// It accumulates and persists the graph"s state, allowing for continuity between separate invocations of a graph.
//
// Example:
//
//	client := langgraph.GetClient()
//	thread, err := client.threadsClient.Get("thread-id")
//	if err != nil {
//		log.Fatalf("Failed to get thread: %v", err)
//	}
type ThreadsClient struct {
	http *http.HttpClient
}

func NewThreadsClient(httpClient *http.HttpClient) *ThreadsClient {
	return &ThreadsClient{http: httpClient}
}

// Get a thread by ID.
//
// Args:
//
//	threadID: The ID of the thread to get.
//
// Returns:
//
//	Thread: Thread Object
//	error: Any error that occurred while getting the thread.
//
// Example:
//
// ```go
//
//	thread, err := client.threadsClient.Get("thread-id")
//	if err != nil {
//		fmt.Printf("Failed to get thread: %v", err)
//	}
//
// fmt.Printf("Thread: %v", thread)
// ```
// ```json
//
//	{
//	    "thread_id": "my_thread_id",
//	    "created_at": "2024-07-18T18:35:15.540834+00:00",
//	    "updated_at": "2024-07-18T18:35:15.540834+00:00",
//	    "metadata": {"graph_id": "agent"}
//	}
//
// ```
func (c *ThreadsClient) Get(threadID string) (schema.Thread, error) {
	ctx := context.Background()
	resp, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s", threadID), nil)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(resp.Body(), &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

// Create a new thread.
//
// Args:
//
//	metadata: Metadata to associate with the thread.
//	threadID: The ID to assign to the thread. If not provided, a random ID will be generated.
//	ifExists: Behavior to take if a thread with the same ID already exists.
//
// Returns:
//
//	Thread: Thread Object
//	error: Any error that occurred while creating the thread.
//
// Example:
//
// ```go
//
//	thread, err := client.threadsClient.Create(nil, "", "")
//	if err != nil {
//		fmt.Printf("Failed to create thread: %v", err)
//	}
//
// fmt.Printf("Thread: %v", thread)
// ```
func (c *ThreadsClient) Create(metadata schema.Json, threadID string, ifExists schema.OnConflictBehavior) (schema.Thread, error) {
	path := "/threads"

	payload := map[string]any{}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if threadID != "" {
		payload["thread_id"] = threadID
	}
	if ifExists != "" {
		payload["if_exists"] = ifExists
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	ctx := context.Background()
	resp, err := c.http.Post(ctx, path, payload)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(resp.Body(), &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

// Update a thread.
//
// Args:
//
//	threadID: The ID of the thread to update.
//	metadata: Metadata to update the thread with.
//
// Returns:
//
//	Thread: Thread Object
//	error: Any error that occurred while updating the thread.
//
// Example:
//
// ```go
//
//	thread, err := client.threadsClient.Update("thread-id", nil)
//	if err != nil {
//		fmt.Printf("Failed to update thread: %v", err)
//	}
//
// fmt.Printf("Thread: %v", thread)
// ```
func (c *ThreadsClient) Update(threadID string, metadata map[string]any) (schema.Thread, error) {
	path := fmt.Sprintf("/threads/%s", threadID)

	payload := map[string]any{}
	if metadata != nil {
		payload["metadata"] = metadata
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	ctx := context.Background()
	resp, err := c.http.Patch(ctx, path, payload)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(resp.Body(), &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

// Delete a thread.
//
// Args:
//
//	threadID: The ID of the thread to delete.
//
// Returns:
//
//	error: Any error that occurred while deleting the thread.
//
// Example:
//
// ```go
//
//	err := client.threadsClient.Delete("thread-id")
//	if err != nil {
//		fmt.Printf("Failed to delete thread: %v", err)
//	}
//
// fmt.Printf("Thread deleted successfully")
// ```
func (c *ThreadsClient) Delete(threadID string) error {
	path := fmt.Sprintf("/threads/%s", threadID)

	ctx := context.Background()
	err := c.http.Delete(ctx, path, nil)
	if err != nil {
		return err
	}

	return nil
}

// Search for threads.
//
// Args:
//
//	metadata: Metadata to filter threads by.
//	values: Values to filter threads by.
//	status: Status to filter threads by.
//	 		Must be one of "idle", "busy", "interrupted" or "error".
//	limit: The maximum number of threads to return.
//	offset: The number of threads to skip.
//
// Returns:
//
//	[]schema.Thread: List of Thread Objects
//	error: Any error that occurred while searching for threads.
//
// Example:
//
// ```go
//
//	threads, err := client.threadsClient.Search({"number":1}, nil, schema.ThreadStatusInterrupted, 15, 5)
//	if err != nil {
//		fmt.Printf("Failed to search threads: %v", err)
//	}
//
// fmt.Printf("Threads: %v", threads)
// ```
func (c *ThreadsClient) Search(
	metadata schema.Json,
	values schema.Json,
	status schema.ThreadStatus,
	limit int,
	offset int,
) ([]schema.Thread, error) {
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
	if values != nil {
		payload["values"] = values
	}
	if status != "" {
		payload["status"] = status
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	path := "/threads/search"

	ctx := context.Background()
	resp, err := c.http.Post(ctx, path, payload)
	if err != nil {
		return []schema.Thread{}, err
	}

	var threads []schema.Thread

	err = json.Unmarshal(resp.Body(), &threads)
	if err != nil {
		return []schema.Thread{}, err
	}

	return threads, nil
}

// Copy a thread.
//
// Args:
//
//	threadID: The ID of the thread to copy.
//
// Returns:
//
//	error: Any error that occurred while copying the thread.
//
// Example:
//
// ```go
//
//	err := client.threadsClient.Copy("thread-id")
//	if err != nil {
//		fmt.Printf("Failed to copy thread: %v", err)
//	}
//
// fmt.Printf("Thread copied successfully")
// ```
func (c *ThreadsClient) Copy(threadID string) error {
	path := fmt.Sprintf("/threads/%s/copy", threadID)

	ctx := context.Background()
	_, err := c.http.Post(ctx, path, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get the state of a thread.
//
// Args:
//
//		threadID: The ID of the thread to get the state of.
//		checkPoint: The checkpoint to get the state at.
//		checkPointID: The ID of the checkpoint to get the state at.
//	 	subgraphs: Include subgraphs states.
//
// Returns:
//
//	ThreadState: ThreadState Object
//	error: Any error that occurred while getting the thread state.
//
// Example:
//
// ```go
//
//	threadState, err := client.threadsClient.GetState("thread-id", nil, "", false)
//	if err != nil {
//		fmt.Printf("Failed to get thread state: %v", err)
//	}
//
// fmt.Printf("ThreadState: %v", threadState)
// ```
// ```json
//
//	{
//	                "values": {
//	                    "messages": [
//	                        {
//	                            "content": "how are you?",
//	                            "additional_kwargs": {},
//	                            "response_metadata": {},
//	                            "type": "human",
//	                            "name": None,
//	                            "id": "fe0a5778-cfe9-42ee-b807-0adaa1873c10",
//	                            "example": False
//	                        },
//	                        {
//	                            "content": "I"m doing well, thanks for asking! I"m an AI assistant created by Anthropic to be helpful, honest, and harmless.",
//	                            "additional_kwargs": {},
//	                            "response_metadata": {},
//	                            "type": "ai",
//	                            "name": None,
//	                            "id": "run-159b782c-b679-4830-83c6-cef87798fe8b",
//	                            "example": False,
//	                            "tool_calls": [],
//	                            "invalid_tool_calls": [],
//	                            "usage_metadata": None
//	                        }
//	                    ]
//	                },
//	                "next": [],
//	                "checkpoint":
//	                    {
//	                        "thread_id": "e2496803-ecd5-4e0c-a779-3226296181c2",
//	                        "checkpoint_ns": "",
//	                        "checkpoint_id": "1ef4a9b8-e6fb-67b1-8001-abd5184439d1"
//	                    }
//	                "metadata":
//	                    {
//	                        "step": 1,
//	                        "run_id": "1ef4a9b8-d7da-679a-a45a-872054341df2",
//	                        "source": "loop",
//	                        "writes":
//	                            {
//	                                "agent":
//	                                    {
//	                                        "messages": [
//	                                            {
//	                                                "id": "run-159b782c-b679-4830-83c6-cef87798fe8b",
//	                                                "name": None,
//	                                                "type": "ai",
//	                                                "content": "I"m doing well, thanks for asking! I"m an AI assistant created by Anthropic to be helpful, honest, and harmless.",
//	                                                "example": False,
//	                                                "tool_calls": [],
//	                                                "usage_metadata": None,
//	                                                "additional_kwargs": {},
//	                                                "response_metadata": {},
//	                                                "invalid_tool_calls": []
//	                                            }
//	                                        ]
//	                                    }
//	                            },
//	                "user_id": None,
//	                "graph_id": "agent",
//	                "thread_id": "e2496803-ecd5-4e0c-a779-3226296181c2",
//	                "created_by": "system",
//	                "assistant_id": "fe096781-5601-53d2-b2f6-0d3403f7e9ca"},
//	                "created_at": "2024-07-25T15:35:44.184703+00:00",
//	                "parent_config":
//	                    {
//	                        "thread_id": "e2496803-ecd5-4e0c-a779-3226296181c2",
//	                        "checkpoint_ns": "",
//	                        "checkpoint_id": "1ef4a9b8-d80d-6fa7-8000-9300467fad0f"
//	                    }
//	            }
//
// ```
func (c *ThreadsClient) GetState(
	threadID string,
	checkPoint *schema.Checkpoint,
	checkPointID string,
	subgraphs bool,
) (schema.ThreadState, error) {
	if checkPoint != nil {
		path := fmt.Sprintf("/threads/%s/state/checkpoint", threadID)

		payload := map[string]any{
			"checkpoint": *checkPoint,
			"subgraphs":  subgraphs,
		}

		payload, ok := removeEmptyFields(payload).(map[string]any)
		if !ok {
			fmt.Println("Error: cleanedPayload is not a map[string]any")
		}

		ctx := context.Background()
		resp, err := c.http.Post(ctx, path, payload)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState
		err = json.Unmarshal(resp.Body(), &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	} else if checkPointID != "" {
		path := fmt.Sprintf("/threads/%s/state/%s", threadID, checkPointID)

		ctx := context.Background()
		resp, err := c.http.Get(ctx, path, nil)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState
		err = json.Unmarshal(resp.Body(), &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	} else {
		path := fmt.Sprintf("/threads/%s/state", threadID)

		ctx := context.Background()
		resp, err := c.http.Get(ctx, path, nil)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState

		err = json.Unmarshal(resp.Body(), &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	}
}

// Update the state of a thread.
//
// Args:
//
//	threadID: The ID of the thread to update the state of.
//	values: The values to update the thread state with.
//	asNode: The node to update the state as.
//	checkPoint: The checkpoint to update the state at.
//	checkPointID: The ID of the checkpoint to update the state at.
//
// Returns:
//
//	ThreadUpdateStateResponse: ThreadUpdateStateResponse Object
//	error: Any error that occurred while updating the thread state.
//
// Example:
//
// ```go
//
//	threadUpdateStateResponse, err := client.threadsClient.UpdateState("thread-id", nil, "", nil, "")
//	if err != nil {
//		fmt.Printf("Failed to update thread state: %v", err)
//	}
//
// fmt.Printf("ThreadUpdateStateResponse: %v", threadUpdateStateResponse)
// ```
// ```json
//
//	{
//	    "checkpoint": {
//	        "thread_id": "e2496803-ecd5-4e0c-a779-3226296181c2",
//	        "checkpoint_ns": "",
//	        "checkpoint_id": "1ef4a9b8-e6fb-67b1-8001-abd5184439d1",
//	        "checkpoint_map": {}
//	    }
//	}
//
// ```
func (c *ThreadsClient) UpdateState(
	threadID string,
	values any,
	asNode string,
	checkPoint *schema.Checkpoint,
	checkPointID string,
) (schema.ThreadUpdateStateResponse, error) {
	path := fmt.Sprintf("/threads/%s/state", threadID)

	payload := map[string]any{
		"values": values,
	}
	if asNode != "" {
		payload["as_node"] = asNode
	}
	if checkPoint != nil {
		payload["checkpoint"] = *checkPoint
	}
	if checkPointID != "" {
		payload["checkpoint_id"] = checkPointID
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	ctx := context.Background()
	resp, err := c.http.Post(ctx, path, payload)
	if err != nil {
		return schema.ThreadUpdateStateResponse{}, err
	}

	var threadUpdateStateResponse schema.ThreadUpdateStateResponse
	err = json.Unmarshal(resp.Body(), &threadUpdateStateResponse)
	if err != nil {
		return schema.ThreadUpdateStateResponse{}, err
	}

	return threadUpdateStateResponse, nil
}

// Get the state history of a thread.
//
// Args:
//
//	threadID: The ID of the thread to get the state history of.
//	limit: The maximum number of states to return.
//	before: The state to get states before.
//	metadata: Metadata to filter states by.
//	checkPoint: The checkpoint to get the state history at.
//
// Returns:
//
//	[]schema.ThreadState: List of ThreadState Objects
//	error: Any error that occurred while getting the thread state history.
//
// Example:
//
// ```go
//
//	threadStates, err := client.threadsClient.GetHistory("thread-id", 10, nil, nil, nil)
//	if err != nil {
//		fmt.Printf("Failed to get thread state history: %v", err)
//	}
//
// fmt.Printf("ThreadStates: %v", threadStates)
// ```
func (c *ThreadsClient) GetHistory(
	threadID string,
	limit int,
	before any,
	metadata map[string]any,
	checkPoint *schema.Checkpoint,
) ([]schema.ThreadState, error) {
	if limit <= 0 {
		limit = 10
	}

	payload := map[string]any{
		"limit": limit,
	}
	if before != nil {
		payload["before"] = before
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if checkPoint != nil {
		payload["checkpoint"] = *checkPoint
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	path := fmt.Sprintf("/threads/%s/history", threadID)

	ctx := context.Background()
	resp, err := c.http.Post(ctx, path, payload)
	if err != nil {
		return []schema.ThreadState{}, err
	}

	var threadStates []schema.ThreadState

	err = json.Unmarshal(resp.Body(), &threadStates)
	if err != nil {
		return []schema.ThreadState{}, err
	}

	return threadStates, nil
}
