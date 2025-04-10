package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func handleError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}
