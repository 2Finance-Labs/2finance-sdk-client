package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/2Finance-Labs/2finance-sdk-client/internal/service"
)

type Client struct {
	service *service.Client
	nextID  atomic.Int64
}

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func New(baseURL string, httpClient *http.Client) *Client {
	return &Client{service: service.New(baseURL, httpClient)}
}

func (c *Client) Call(ctx context.Context, method string, params any) (*Response, error) {
	request := Request{
		JSONRPC: "2.0",
		ID:      c.nextID.Add(1),
		Method:  method,
		Params:  params,
	}
	var response Response
	if err := c.service.DoJSON(ctx, http.MethodPost, "/mcp", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) Request(ctx context.Context, method string, path string, body any, out any) error {
	return c.service.DoJSON(ctx, method, path, body, out)
}

func (c *Client) ListTools(ctx context.Context) (*Response, error) {
	return c.Call(ctx, "tools/list", nil)
}

func (c *Client) CallTool(ctx context.Context, name string, arguments any) (*Response, error) {
	return c.Call(ctx, "tools/call", map[string]any{"name": name, "arguments": arguments})
}

func (c *Client) ConversationPlan(ctx context.Context, arguments any) (*Response, error) {
	return c.CallTool(ctx, "finance_assistant.conversation.plan", arguments)
}
