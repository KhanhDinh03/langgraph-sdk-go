package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KhanhDinh03/langgraph-sdk-go/http"
	"github.com/KhanhDinh03/langgraph-sdk-go/schema"
)

// Client for managing recurrent runs (cron jobs) in LangGraph.
//
//	A run is a single invocation of an assistant with optional input and config.
//	This client allows scheduling recurring runs to occur automatically.
type CronsClient struct {
	http *http.HttpClient
}

func NewCronsClient(httpClient *http.HttpClient) *CronsClient {
	return &CronsClient{http: httpClient}
}

// Create a cron job for a thread.
//
// Args:
//
//			threadID: the thread ID to run the cron job on.
//			assistantID: The assistant ID or graph name to use for the cron job.
//	    		If using graph name, will default to first assistant created from that graph.
//			schedule: The cron schedule to execute this job on.
//			input: The input to the graph.
//			metadata: Metadata to assign to the cron job runs.
//			config: The configuration for the assistant.
//			interruptBefore: Nodes to interrupt immediately before they get executed.
//			interruptAfter: Nodes to Nodes to interrupt immediately after they get executed.
//			webhook: Webhook to call after LangGraph API call is done.
//			multitaskStrategy: Multitask strategy to use.
//				Must be one of "reject", "interrupt", "rollback", or "enqueue".
//
// Returns:
//
//	schema.Run: The created cron job run.
//	error: An error if the operation failed.
//
// Example:
//
//	go```
//	cronsClient := client.NewCronsClient(httpClient)
//	run, err := cronsClient.CreatForThread(
//		"threadID",
//		"assistantID",
//		"27 15 * * *",
//		{"messages": [{"role": "user", "content": "hello!"}]},
//		{"name":"my_run"},
//		schema.Config{"configurable": {"model_name": "openai"}},
//		["node_to_stop_before_1","node_to_stop_before_2"],
//		["node_to_stop_after_1","node_to_stop_after_2"],
//		"http://webhook.com",
//		schema.MultitaskStrategyInterrupt,
//	)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(run)
//	}
//	```
func (c *CronsClient) CreatForThread(
	ctx context.Context,
	threadID string,
	assistantID string,
	schedule string,
	input map[string]any,
	metadata map[string]any,
	config schema.Config,
	interruptBefore any,
	interruptAfter any,
	webhook string,
	multitaskStrategy schema.MultitaskStrategy,
) (schema.Run, error) {
	payload := map[string]any{
		"schedule":         schedule,
		"input":            input,
		"config":           config,
		"metadata":         metadata,
		"assistant_id":     assistantID,
		"interrupt_before": interruptBefore,
		"interrupt_after":  interruptAfter,
		"webhook":          webhook,
	}

	if multitaskStrategy != "" {
		payload["multitask_strategy"] = multitaskStrategy
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/crons", threadID), payload)
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

// Create a cron run.
//
// Args:
//
//			assistantID: The assistant ID or graph name to use for the cron job.
//				If using graph name, will default to first assistant created from that graph.
//			schedule: The cron schedule to execute this job on.
//			input: The input to the graph.
//			metadata: Metadata to assign to the cron job runs.
//			config: The configuration for the assistant.
//			interruptBefore: Nodes to interrupt immediately before they get executed.
//			interruptAfter: Nodes to Nodes to interrupt immediately after they get executed.
//			webhook: Webhook to call after LangGraph API call is done.
//			multitaskStrategy: Multitask strategy to use.
//	 			Must be one of "reject", "interrupt", "rollback", or "enqueue".
//
// Returns:
//
//	schema.Run: The created cron job run.
//	error: An error if the operation failed.
//
// Example:
//
//	go```
//	cronsClient := client.NewCronsClient(httpClient)
//	run, err := cronsClient.Creat("assistantID", "27 15 * * *", {"messages": [{"role": "user", "content": "hello!"}]}, {"name":"my_run"}, schema.Config{"configurable": {"model_name": "openai"}}, ["node_to_stop_before_1","node_to_stop_before_2"], ["node_to_stop_after_1","node_to_stop_after_2"], "http://webhook.com", schema.MultitaskStrategyInterrupt)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(run)
//	}
//	```
func (c *CronsClient) Creat(
	ctx context.Context,
	assistantID string,
	schedule string,
	input map[string]any,
	metadata map[string]any,
	config schema.Config,
	interruptBefore schema.All,
	interruptAfter schema.All,
	webhook string,
	multitaskStrategy schema.MultitaskStrategy,
) (schema.Run, error) {
	payload := map[string]any{
		"schedule":         schedule,
		"input":            input,
		"config":           config,
		"metadata":         metadata,
		"assistant_id":     assistantID,
		"interrupt_before": interruptBefore,
		"interrupt_after":  interruptAfter,
		"webhook":          webhook,
	}

	if multitaskStrategy != "" {
		payload["multitask_strategy"] = multitaskStrategy
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "runs/crons", payload)
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

func (c *CronsClient) Delete(cronID string) error {
	path := fmt.Sprintf("/crons/%s", cronID)

	ctx := context.Background()
	err := c.http.Delete(ctx, path, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *CronsClient) Search(
	ctx context.Context,
	assistantID string,
	threadID string,
	limit int,
	offset int,
) ([]schema.Cron, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	payload := map[string]any{
		"aasistant_id": assistantID,
		"thread_id":    threadID,
		"limit":        limit,
		"offset":       offset,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "runs/crons/search", payload)
	if err != nil {
		return []schema.Cron{}, err
	}

	var crons []schema.Cron

	err = json.Unmarshal(resp.Body(), &crons)
	if err != nil {
		return []schema.Cron{}, err
	}

	return crons, nil
}
