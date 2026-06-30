package planner

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/2Finance-Labs/2finance-sdk-client/analytics"
	"github.com/2Finance-Labs/2finance-sdk-client/mcp"
	"github.com/2Finance-Labs/2finance-sdk-client/orchestrator"
	"github.com/2Finance-Labs/2finance-sdk-client/tradingcontrol"
)

type Config struct {
	MCP            *mcp.Client
	Orchestrator   *orchestrator.Client
	TradingControl *tradingcontrol.Client
	Analytics      *analytics.Client
}

type Client struct {
	mcp            *mcp.Client
	orchestrator   *orchestrator.Client
	tradingControl *tradingcontrol.Client
	analytics      *analytics.Client
}

type Request struct {
	Goal         string         `json:"goal"`
	Context      map[string]any `json:"context,omitempty"`
	UseTrading   bool           `json:"use_trading,omitempty"`
	UseAnalytics bool           `json:"use_analytics,omitempty"`
}

type Plan struct {
	Source string          `json:"source"`
	Raw    json.RawMessage `json:"raw"`
}

func New(config Config) *Client {
	return &Client{
		mcp:            config.MCP,
		orchestrator:   config.Orchestrator,
		tradingControl: config.TradingControl,
		analytics:      config.Analytics,
	}
}

func (c *Client) ConversationPlan(ctx context.Context, request Request) (*Plan, error) {
	if c.mcp == nil {
		return nil, errors.New("planner: mcp client is required")
	}
	response, err := c.mcp.ConversationPlan(ctx, request)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return &Plan{Source: "mcp.finance_assistant.conversation.plan", Raw: raw}, nil
}

func (c *Client) OrchestratedPlan(ctx context.Context, request any) (*Plan, error) {
	if c.orchestrator == nil {
		return nil, errors.New("planner: orchestrator client is required")
	}
	response, err := c.orchestrator.SendMessage(ctx, request)
	if err != nil {
		return nil, err
	}
	return &Plan{Source: "orchestrator.messages", Raw: response}, nil
}

func (c *Client) OperationalPlan(ctx context.Context, request any) (*Plan, error) {
	return c.OrchestratedPlan(ctx, request)
}

func (c *Client) TradingPlan(ctx context.Context, request Request) (*Plan, error) {
	if request.Context == nil {
		request.Context = map[string]any{}
	}
	if request.UseTrading && c.tradingControl != nil {
		if robots, err := c.tradingControl.Robots(ctx); err == nil {
			request.Context["trading_robots"] = json.RawMessage(robots)
		}
	}
	if request.UseAnalytics && c.analytics != nil {
		if indicators, err := c.analytics.ListIndicators(ctx); err == nil {
			request.Context["analytics_indicators"] = json.RawMessage(indicators)
		}
	}
	return c.ConversationPlan(ctx, request)
}
