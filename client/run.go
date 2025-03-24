package client

import (
	"context"
	"encoding/json"
	"fmt"
	"langgraph-sdk/http"
	"langgraph-sdk/schema"
	"net/url"
)

// Client for managing runs in LangGraph.
//
// A run is a single assistant invocation with optional input, config, and metadata.
// This client manages runs, which can be stateful (on threads) or stateless.
//
// Example:
//
//	client := langgraph.GetClient()
//	run := client.Runs.Create("my-assistant", map[string]interface{}{"text": "Hello, world!"}, schema.Command{Type: "message"})
type RunsClient struct {
	http *http.HttpClient
}

func NewRunsClient(httpClient *http.HttpClient) *RunsClient {
	return &RunsClient{http: httpClient}
}

// Create a run and stream the results.
//
// Args:
//
//	threadID: the thread ID to assign to the thread.
//		If None will create a stateless run.
//	assistantID: The assistant ID or graph name to stream from.
//		If using graph name, will default to first assistant created from that graph.
//	input: The input to the graph.
//	command: A command to execute. Cannot be combined with input.
//	streamMode: The stream mode(s) to use.
//	streamSubgraphs: Whether to stream output from subgraphs.
//	metadata: Metadata to assign to the run.
//	config: The configuration for the assistant.
//	checkpoint: The checkpoint to resume from.
//	interruptBefore: Nodes to interrupt immediately before they get executed.
//	interruptAfter: Nodes to Nodes to interrupt immediately after they get executed.
//	feedbackKeys: Feedback keys to assign to run.
//	onDisconnect: The disconnect mode to use.
//		Must be one of "cancel" or "continue".
//	onCompletion: Whether to delete or keep the thread created for a stateless run.
//		Must be one of "delete" or "keep".
//	webhook: Webhook to call after LangGraph API call is done.
//	multitaskStrategy: Multitask strategy to use.
//		Must be one of "reject", "interrupt", "rollback", or "enqueue".
//	ifNotExists: How to handle missing thread. Defaults to "reject".
//		Must be either "reject" (raise error if missing), or "create" (create new thread).
//	afterSeconds: The number of seconds to wait before starting the run.
//		Use to schedule future runs.
//
// Returns:
//
//	A channel of stream parts and a cancel function.
//
// Example:
//
//	stream, cancel := client.Runs.Stream("my-thread", "my-assistant", map[string]interface{}{"text": "Hello, world!"}, schema.Command{Type: "message"}, schema.StreamModeAll, true, nil, schema.Config{}, nil, "", "", "", nil, schema.DisconnectModeCancel, schema.OnCompletionBehaviorDelete, "", schema.MultitaskStrategyReject, schema.IfNotExistsReject, 0)
//	for part := range stream {
//		fmt.Println(part)
//	}
//	cancel()
//
// StreamPart(event="metadata", data={"run_id": "1ef4a9b8-d7da-679a-a45a-872054341df2"})
//
//	StreamPart(event="values", data={"messages": [{"content": "how are you?", "additional_kwargs": {}, "response_metadata": {}, "type": "human", "name": None, "id": "fe0a5778-cfe9-42ee-b807-0adaa1873c10", "example": False}]})
//	StreamPart(event="values", data={"messages": [{"content": "how are you?", "additional_kwargs": {}, "response_metadata": {}, "type": "human", "name": None, "id": "fe0a5778-cfe9-42ee-b807-0adaa1873c10", "example": False}, {"content": "I"m doing well, thanks for asking! I"m an AI assistant created by Anthropic to be helpful, honest, and harmless.", "additional_kwargs": {}, "response_metadata": {}, "type": "ai", "name": None, "id": "run-159b782c-b679-4830-83c6-cef87798fe8b", "example": False, "tool_calls": [], "invalid_tool_calls": [], "usage_metadata": None}]})
//	StreamPart(event="end", data=None)
func (c *RunsClient) Stream(
	ctx context.Context,
	threadID string,
	assistantID string,
	input map[string]any,
	command schema.Command,
	streamMode any,
	streamSubgraphs bool,
	metadata map[string]any,
	config schema.Config,
	checkpoint schema.Checkpoint,
	checkpointID string,
	interruptBefore any,
	interruptAfter any,
	feedbackKeys []string,
	onDisconnect schema.DisconnectMode,
	onCompletion schema.OnCompletionBehavior,
	webhook string,
	multitaskStrategy schema.MultitaskStrategy,
	ifNotExists schema.IfNotExists,
	afterSeconds int,
) (chan schema.StreamPart, context.CancelFunc) {
	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"stream_mode":        streamMode,
		"stream_subgraphs":   streamSubgraphs,
		"assistant_id":       assistantID,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"feedback_keys":      feedbackKeys,
		"webhook":            webhook,
		"checkpoint":         checkpoint,
		"checkpoint_id":      checkpointID,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"on_disconnect":      onDisconnect,
		"on_completion":      onCompletion,
		"after_seconds":      afterSeconds,
	}
	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs/stream", threadID)
	} else {
		endPoint = "/runs/stream"
	}

	ctx, cancel := context.WithCancel(ctx)

	streamCh, errCh, err := c.http.Stream(ctx, endPoint, "POST", payload, nil)
	if err != nil {
		cancel()
		errCh <- err
		close(errCh)
		return nil, cancel
	}

	return streamCh, cancel
}

// Create a background run.
//
// Args:
//
//	threadID: the thread ID to assign to the thread.
//		If None will create a stateless run.
//	assistantID: The assistant ID or graph name to stream from.
//		If using graph name, will default to first assistant created from that graph.
//	input: The input to the graph.
//	command: A command to execute. Cannot be combined with input.
//	streamMode: The stream mode(s) to use.
//	streamSubgraphs: Whether to stream output from subgraphs.
//	metadata: Metadata to assign to the run.
//	config: The configuration for the assistant.
//	checkpoint: The checkpoint to resume from.
//	interruptBefore: Nodes to interrupt immediately before they get executed.
//	interruptAfter: Nodes to Nodes to interrupt immediately after they get executed.
//	webhook: Webhook to call after LangGraph API call is done.
//	multitaskStrategy: Multitask strategy to use.
//		Must be one of "reject", "interrupt", "rollback", or "enqueue".
//	onCompletion: Whether to delete or keep the thread created for a stateless run.
//		Must be one of "delete" or "keep".
//	ifNotExists: How to handle missing thread. Defaults to "reject".
//		Must be either "reject" (raise error if missing), or "create" (create new thread).
//	afterSeconds: The number of seconds to wait before starting the run.
//		Use to schedule future runs.
//
// Returns:
//
//	schema.Run: The created background run.
//
// Example:
//
// ```go
//
//	run, err := client.Runs.Create("my-thread", "my-assistant", map[string]interface{}{"text": "Hello, world!"}, schema.Command{Type: "message"}, schema.StreamModeAll, true, nil, schema.Config{}, nil, "", "", "", nil, schema.MultitaskStrategyReject, schema.IfNotExistsReject, 0, "", schema.OnCompletionBehaviorDelete, 0)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(run)
//
// ```
// ```json
//
//	{
//	               "run_id": "my_run_id",
//	               "thread_id": "my_thread_id",
//	               "assistant_id": "my_assistant_id",
//	               "created_at": "2024-07-25T15:35:42.598503+00:00",
//	               "updated_at": "2024-07-25T15:35:42.598503+00:00",
//	               "metadata": {},
//	               "status": "pending",
//	               "kwargs":
//	                   {
//	                       "input":
//	                           {
//	                               "messages": [
//	                                   {
//	                                       "role": "user",
//	                                       "content": "how are you?"
//	                                   }
//	                               ]
//	                           },
//	                       "config":
//	                           {
//	                               "metadata":
//	                                   {
//	                                       "created_by": "system"
//	                                   },
//	                               "configurable":
//	                                   {
//	                                       "run_id": "my_run_id",
//	                                       "user_id": None,
//	                                       "graph_id": "agent",
//	                                       "thread_id": "my_thread_id",
//	                                       "checkpoint_id": None,
//	                                       "model_name": "openai",
//	                                       "assistant_id": "my_assistant_id"
//	                                   }
//	                           },
//	                       "webhook": "https://my.fake.webhook.com",
//	                       "temporary": False,
//	                       "stream_mode": ["values"],
//	                       "feedback_keys": None,
//	                       "interrupt_after": ["node_to_stop_after_1","node_to_stop_after_2"],
//	                       "interrupt_before": ["node_to_stop_before_1","node_to_stop_before_2"]
//	                   },
//	               "multitask_strategy": "interrupt"
//	           }
//
// ```
func (c *RunsClient) Create(
	ctx context.Context,
	threadID string,
	assistantID string,
	input map[string]any,
	command schema.Command,
	streamMode schema.StreamMode,
	streamSubgraphs bool,
	metadata map[string]any,
	config schema.Config,
	checkpoint schema.Checkpoint,
	checkpointID string,
	interruptBefore schema.All,
	interruptAfter schema.All,
	webhook string,
	multitaskStrategy schema.MultitaskStrategy,
	ifNotExists schema.IfNotExists,
	onCompletion schema.OnCompletionBehavior,
	afterSeconds int,
) (schema.Run, error) {
	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"stream_mode":        streamMode,
		"stream_subgraphs":   streamSubgraphs,
		"assistant_id":       assistantID,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"webhook":            webhook,
		"checkpoint":         checkpoint,
		"checkpoint_id":      checkpointID,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"after_seconds":      afterSeconds,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs", threadID)
	} else {
		endPoint = "/runs"
	}

	resp, err := c.http.Post(ctx, endPoint, payload)
	if err != nil {
		return schema.Run{}, err
	}

	var run schema.Run
	err = json.Unmarshal(resp.Body(), &run)
	if err != nil {
		return schema.Run{}, err
	}

	return run, nil
}

func filterPayload(payload map[string]any) map[string]any {
	filtered := make(map[string]any)
	for k, v := range payload {
		if v != nil {
			filtered[k] = v
		}
	}
	return filtered
}

// Create a batch of stateless background runs.
func (c *RunsClient) CreateBatch(ctx context.Context, payloads []map[string]any) ([]schema.Run, error) {
	filteredPayloads := make([]map[string]any, 0, len(payloads))
	for _, payload := range payloads {
		filteredPayloads = append(filteredPayloads, filterPayload(payload))
	}

	jsonData := map[string]any{"batch": filteredPayloads}

	resp, err := c.http.Post(ctx, "/runs/batch", jsonData)
	if err != nil {
		return nil, err
	}

	var runs []schema.Run
	err = json.Unmarshal(resp.Body(), &runs)
	if err != nil {
		return nil, err
	}

	return runs, nil
}

// Create a run, wait until it finishes and return the final state.
//
// Args:
//
//		threadID: the thread ID to create the run on.
//			If None will create a stateless run.
//		assistantID: The assistant ID or graph name to run.
//			If using graph name, will default to first assistant created from that graph.
//		input: The input to the graph.
//		command: A command to execute. Cannot be combined with input.
//		metadata: Metadata to assign to the run.
//		config: The configuration for the assistant.
//		checkpoint: The checkpoint to resume from.
//		interruptBefore: Nodes to interrupt immediately before they get executed.
//		interruptAfter: Nodes to Nodes to interrupt immediately after they get executed.
//		webhook: Webhook to call after LangGraph API call is done.
//		onDisconnect: The disconnect mode to use.
//			Must be one of "cancel" or "continue".
//		onCompletion: Whether to delete or keep the thread created for a stateless run.
//			Must be one of "delete" or "keep".
//		multitaskStrategy: Multitask strategy to use.
//			Must be one of "reject", "interrupt", "rollback", or "enqueue".
//		ifNotExists: How to handle missing thread. Defaults to "reject".
//			Must be either "reject" (raise error if missing), or "create" (create new thread).
//		afterSeconds: The number of seconds to wait before starting the run.
//			Use to schedule future runs.
//		raiseError: Whether to raise an error if the run fails.
//	Returns:
//
//		map[string]any || map[string]any: The final state of the run.
//
// Example:
//
//	```go
//
//	result, err := client.Runs.Wait("my-thread", "my-assistant", map[string]interface{}{"text": "Hello, world!"}, schema.Command{Type: "message"}, nil, schema.Config{}, nil, "", "", "", nil, schema.DisconnectModeCancel, schema.OnCompletionBehaviorDelete, schema.MultitaskStrategyReject, schema.IfNotExistsReject, 0, true)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(result)
//
//	```
//	```json
//
//	{
//	                "messages": [
//	                    {
//	                        "content": "how are you?",
//	                        "additional_kwargs": {},
//	                        "response_metadata": {},
//	                        "type": "human",
//	                        "name": None,
//	                        "id": "f51a862c-62fe-4866-863b-b0863e8ad78a",
//	                        "example": False
//	                    },
//	                    {
//	                        "content": "I"m doing well, thanks for asking! I'm an AI assistant created by Anthropic to be helpful, honest, and harmless.",
//	                        "additional_kwargs": {},
//	                        "response_metadata": {},
//	                        "type": "ai",
//	                        "name": None,
//	                        "id": "run-bf1cd3c6-768f-4c16-b62d-ba6f17ad8b36",
//	                        "example": False,
//	                        "tool_calls": [],
//	                        "invalid_tool_calls": [],
//	                        "usage_metadata": None
//	                    }
//	                ]
//	            }
//
//	```
func (c *RunsClient) Wait(
	ctx context.Context,
	threadID string,
	assistantID string,
	input map[string]any,
	command schema.Command,
	metadata map[string]any,
	config schema.Config,
	checkPoint schema.Checkpoint,
	checkPointID string,
	interruptBefore any,
	interruptAfter any,
	webhook string,
	onDisconnect schema.DisconnectMode,
	onCompletion schema.OnCompletionBehavior,
	multitaskStrategy schema.MultitaskStrategy,
	ifNotExists schema.IfNotExists,
	afterSeconds int,
	raiseError bool,
) (any, error) {
	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"assistant_id":       assistantID,
		"checkpoint":         checkPoint,
		"checkpoint_id":      checkPointID,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"webhook":            webhook,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"on_disconnect":      onDisconnect,
		"on_completion":      onCompletion,
		"after_seconds":      afterSeconds,
		"raise_error":        raiseError,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs/wait", threadID)
	} else {
		endPoint = "/runs/wait"
	}

	resp, err := c.http.Post(ctx, endPoint, payload)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if raiseError {
		if errData, exists := result["__error__"].(map[string]any); exists {
			return nil, fmt.Errorf("%s", errData["message"])
		}
	}

	return result, nil
}

// List runs in a thread.
//
// Args:
//
//		threadID: The thread ID to list runs for.
//		limit: The maximum number of results to return.
//		offset: The number of results to skip.
//		status: The status of the run to filter by.
//
//	Returns:
//
//		[]schema.Run: The list of runs.
//
// Example:
//
//	```go
//
//	runs, err := client.Runs.List("my-thread", 10, 0, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(runs)
//
//	```
func (c *RunsClient) List(ctx context.Context, threadID string, limit int, offset int, status *schema.RunStatus) ([]schema.Run, error) {
	if limit <= 0 {
		limit = 10
	}

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))

	if status != nil {
		params.Add("status", string(*status))
	}

	resp, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs", threadID), params)
	if err != nil {
		return []schema.Run{}, err
	}

	var runs []schema.Run
	err = json.Unmarshal(resp.Body(), &runs)
	if err != nil {
		return []schema.Run{}, err
	}

	return runs, nil
}

// Get a run.
//
// Args:
//
//	threadID: The thread ID the run is on.
//	runID: The run ID to get.
//
// Returns:
//
//	schema.Run: The run.
//
// Example:
//
//	```go
//
//	run, err := client.Runs.Get("my-thread", "my-run")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(run)
//
//	```
func (c *RunsClient) Get(ctx context.Context, threadID string, runID string) (schema.Run, error) {
	resp, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs/%s", threadID, runID), nil)
	if err != nil {
		return schema.Run{}, err
	}

	var run schema.Run
	err = json.Unmarshal(resp.Body(), &run)
	if err != nil {
		return schema.Run{}, err
	}

	return run, nil
}

// Get a run.
//
// Args:
//
//		threadID: The thread ID to cancel.
//		run_id: The run ID to cancek.
//		wait: Whether to wait until run has completed.
//		action: Action to take when cancelling the run. Possible values are `interrupt` or `rollback`. Default is `interrupt`.
//
//	Returns:
//
//		error: An error if the operation failed.
//
// Example:
//
//	```go
//
//	err := client.Runs.Cancel("my-thread", "my-run", false, schema.CancelActionInterrupt)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	```
func (c *RunsClient) Cancel(ctx context.Context, threadID string, runID string, wait bool, action schema.CancelAction) error {
	if action == "" {
		action = schema.CancelActionInterrupt
	}

	payload := map[string]any{
		"wait":   wait,
		"action": action,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	_, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/runs/%s/cancel", threadID, runID), payload)
	if err != nil {
		return err
	}

	return nil
}

// Block until a run is done. Returns the final state of the thread.
//
// Args:
//
//	threadID: The thread ID to wait for.
//	runID: The run ID to wait for.
//
// Returns:
//
//	map[string]any: The final state of the run.
//	error: An error if the operation failed.
//
// Example:
//
// ```go
//
//	result, err := client.Runs.Wait("my-thread", "my-run")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(result)
//
// ```
func (c *RunsClient) Join(ctx context.Context, threadID string, runID string) (map[string]any, error) {
	resp, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs/%s/join", threadID, runID), nil)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Stream output from a run in real-time, until the run is done.
// Output is not buffered, so any output produced before this call will not be received here.
//
// Args:
//
//	threadID: The thread ID to stream from.
//	runID: The run ID to stream from.
//	cancelOnDisconnect: Whether to cancel the stream when disconnected.
//
// Returns:
//
//	chan http.StreamPart: A channel of stream parts.
//	context.CancelFunc: A cancel function to stop the stream.
//
// Example:
//
// ```go
//
//	stream, cancel := client.Runs.JoinStream("my-thread", "my-run", false)
//	for part := range stream {
//		fmt.Println(part)
//	}
//	cancel()
//
// ```
func (c *RunsClient) JoinStream(ctx context.Context, threadID string, runID string, cancelOnDisconnect bool) (chan schema.StreamPart, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	streamCh, errCh, err := c.http.Stream(ctx, fmt.Sprintf("/threads/%s/runs/%s/join/stream", threadID, runID), "GET", nil, nil)
	if err != nil {
		cancel()
		errCh <- err
		close(errCh)
		return nil, cancel
	}

	convertedCh := make(chan schema.StreamPart)

	go func() {
		defer close(convertedCh)
		for streamPart := range streamCh {
			convertedCh <- schema.StreamPart{
				Event: streamPart.Event,
				Data:  streamPart.Data,
			}
		}
	}()

	return convertedCh, cancel
}

// Delete a run.
//
// Args:
//
//	threadID: The thread ID to delete the run from.
//	runID: The run ID to delete.
//
// Returns:
//
//	error: An error if the operation failed.
//
// Example:
//
//	```go
//
//	err := client.Runs.Delete("my-thread", "my-run")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	```
func (c *RunsClient) Delete(ctx context.Context, threadID string, runID string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/threads/%s/runs/%s", threadID, runID), nil)
	if err != nil {
		return err
	}

	return nil
}
