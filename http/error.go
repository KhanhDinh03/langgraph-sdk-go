package http

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func handleError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.IsError() {
		log.Printf("Error from langgraph-api: %s", string(resp.Body()))
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}
