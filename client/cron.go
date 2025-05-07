package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type CronsClient struct {
	http *http.HttpClient
}

func NewCronsClient(httpClient *http.HttpClient) *CronsClient {
	return &CronsClient{http: httpClient}
}

func (c *CronsClient) CreatForThread(ctx context.Context, threadID string, assistantID string, schedule string, input *map[string]any, metadata *map[string]any, config *schema.Config, interruptBefore *any, interruptAfter *any, webhook *string, multitaskStrategy *schema.MultitaskStrategy, headers *map[string]string) (schema.Run, error) {
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

	if multitaskStrategy != nil {
		payload["multitask_strategy"] = *multitaskStrategy
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/crons", threadID), payload, headers)
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

func (c *CronsClient) Creat(ctx context.Context, assistantID string, schedule string, input *map[string]any, metadata *map[string]any, config *schema.Config, interruptBefore *schema.All, interruptAfter *schema.All, webhook *string, multitaskStrategy *schema.MultitaskStrategy, headers *map[string]string) (schema.Run, error) {
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

	if multitaskStrategy != nil {
		payload["multitask_strategy"] = *multitaskStrategy
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	resp, err := c.http.Post(ctx, "runs/crons", payload, headers)
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

func (c *CronsClient) Delete(ctx context.Context, cronID string, headers *map[string]string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/crons/%s", cronID), nil, headers)
	if err != nil {
		return err
	}

	return nil
}

func (c *CronsClient) Search(ctx context.Context, assistantID *string, threadID *string, limit *int, offset *int, headers *map[string]string) ([]schema.Cron, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	if offset != nil && *offset < 0 {
		*offset = 0
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

	resp, err := c.http.Post(ctx, "runs/crons/search", payload, headers)
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
