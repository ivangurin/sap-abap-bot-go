package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// REST API endpoints for models inference
// https://docs.github.com/en/rest/models/inference?apiVersion=2022-11-28#run-an-inference-request
func (c *Client) ChatCompletions(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	if c.config.Debug {
		c.logger.Debugf("ChatCompletions request: %s", string(jsonData))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, host+chatCompletions, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.GitHubToken)
	httpReq.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := httpResp.Body.Close()
		if err != nil {
			c.logger.Errorf("close response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if c.config.Debug {
		c.logger.Debugf("ChatCompletions response: %s", string(body))
	}

	response := &ChatCompletionResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
