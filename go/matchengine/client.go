package matchengine

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	WebSocketURL string
	HTTPClient   *http.Client
	Dialer       *websocket.Dialer
}

type OrderCommand struct {
	Schema          string          `json:"schema,omitempty"`
	ClientOrderID   string          `json:"client_order_id"`
	IdempotencyKey  string          `json:"idempotency_key"`
	Symbol          string          `json:"symbol"`
	Side            string          `json:"side"`
	Type            string          `json:"type"`
	Quantity        string          `json:"quantity"`
	Price           string          `json:"price,omitempty"`
	TimeInForce     string          `json:"time_in_force,omitempty"`
	AccountID       string          `json:"account_id,omitempty"`
	ClientTimestamp time.Time       `json:"client_timestamp,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
}

type ExecutionReport struct {
	Schema        string          `json:"schema,omitempty"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
	OrderID       string          `json:"order_id,omitempty"`
	Status        string          `json:"status,omitempty"`
	Symbol        string          `json:"symbol,omitempty"`
	FilledQty     string          `json:"filled_qty,omitempty"`
	Raw           json.RawMessage `json:"raw,omitempty"`
}

type MarketDataSubscribeRequest struct {
	Schema    string          `json:"schema,omitempty"`
	Symbols   []string        `json:"symbols,omitempty"`
	Channels  []string        `json:"channels,omitempty"`
	Interval  string          `json:"interval,omitempty"`
	AccountID string          `json:"account_id,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

func New(webSocketURL string, httpClient *http.Client) *Client {
	return &Client{
		WebSocketURL: webSocketURL,
		HTTPClient:   httpClient,
		Dialer:       websocket.DefaultDialer,
	}
}

func NewMarketDataSubscribeRequest(request MarketDataSubscribeRequest) MarketDataSubscribeRequest {
	if request.Schema == "" {
		request.Schema = "matchengine.market_data_subscribe.v1"
	}
	return request
}

func (c *Client) DialOrderEntry(ctx context.Context, headers http.Header) (*websocket.Conn, *http.Response, error) {
	dialer := c.Dialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}
	return dialer.DialContext(ctx, c.WebSocketURL, headers)
}

func (c *Client) DialMarketData(ctx context.Context, headers http.Header) (*websocket.Conn, *http.Response, error) {
	return c.DialOrderEntry(ctx, headers)
}

func (c *Client) SubmitOrder(ctx context.Context, command OrderCommand, headers http.Header) (*ExecutionReport, error) {
	conn, _, err := c.DialOrderEntry(ctx, headers)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if command.Schema == "" {
		command.Schema = "matchengine.order_command.v1"
	}
	if err := conn.WriteJSON(command); err != nil {
		return nil, err
	}
	var report ExecutionReport
	if err := conn.ReadJSON(&report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (c *Client) SubscribeMarketData(ctx context.Context, request MarketDataSubscribeRequest, headers http.Header) (*websocket.Conn, error) {
	conn, _, err := c.DialMarketData(ctx, headers)
	if err != nil {
		return nil, err
	}
	request = NewMarketDataSubscribeRequest(request)
	if err := conn.WriteJSON(request); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
