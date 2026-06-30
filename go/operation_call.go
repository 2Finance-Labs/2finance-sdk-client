package twofinance

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type resolvedRequester interface {
	Request(context.Context, string, string, any, any) error
}

type resolvedGetter interface {
	Get(context.Context, string, any) error
}

type resolvedPoster interface {
	Post(context.Context, string, any, any) error
}

type resolvedPutter interface {
	Put(context.Context, string, any, any) error
}

type resolvedDeleter interface {
	Delete(context.Context, string, any) error
}

func (c *Client) CallResolvedOperation(ctx context.Context, domain string, operation ResolvedOperation, body any) (json.RawMessage, error) {
	var response json.RawMessage
	target, err := c.resolvedOperationTarget(domain, operation)
	if err != nil {
		return nil, err
	}
	if err := callResolvedTarget(ctx, target, operation, body, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) CallOperation(ctx context.Context, domain string, operation DomainOperation, pathParams map[string]string, query map[string]string, body any) (json.RawMessage, error) {
	resolved, err := operation.Resolve(pathParams, query)
	if err != nil {
		return nil, err
	}
	return c.CallResolvedOperation(ctx, domain, resolved, body)
}

func (c *Client) CallCatalogOperation(ctx context.Context, catalog DomainOperationsCatalog, domain string, operationName string, pathParams map[string]string, query map[string]string, body any) (json.RawMessage, error) {
	resolved, err := catalog.ResolveOperation(domain, operationName, pathParams, query)
	if err != nil {
		return nil, err
	}
	return c.CallResolvedOperation(ctx, domain, resolved, body)
}

func (c *Client) resolvedOperationTarget(domain string, operation ResolvedOperation) (any, error) {
	switch serviceKey(domain) {
	case "auth":
		return c.Auth, nil
	case "network":
		return c.Network, nil
	case "analytics":
		return c.Analytics, nil
	case "orchestrator":
		return c.Orchestrator, nil
	case "mcp", "planner":
		return c.MCP, nil
	case "tradingcontrol":
		return c.TradingControl, nil
	case "keystore":
		return c.KeyStore, nil
	case "hummingbot":
		return c.Hummingbot, nil
	case "wise":
		return c.Providers.Wise, nil
	case "airwallex":
		return c.Providers.Airwallex, nil
	case "providers":
		path := strings.ToLower(operation.Path)
		if strings.HasPrefix(path, "/api/v1/") {
			return c.Providers.Airwallex, nil
		}
		return c.Providers.Wise, nil
	case "matchengine":
		return nil, errors.New("2finance: matchengine resolved operations use websocket helpers")
	default:
		return nil, errors.New("2finance: unsupported operation domain " + domain)
	}
}

func callResolvedTarget(ctx context.Context, target any, operation ResolvedOperation, body any, out *json.RawMessage) error {
	if requester, ok := target.(resolvedRequester); ok {
		return requester.Request(ctx, operation.Method, operation.Path, body, out)
	}
	switch operation.Method {
	case http.MethodGet:
		getter, ok := target.(resolvedGetter)
		if !ok {
			return errors.New("2finance: resolved operation target does not support GET")
		}
		return getter.Get(ctx, operation.Path, out)
	case http.MethodPost:
		poster, ok := target.(resolvedPoster)
		if !ok {
			return errors.New("2finance: resolved operation target does not support POST")
		}
		return poster.Post(ctx, operation.Path, body, out)
	case http.MethodPut:
		putter, ok := target.(resolvedPutter)
		if !ok {
			return errors.New("2finance: resolved operation target does not support PUT")
		}
		return putter.Put(ctx, operation.Path, body, out)
	case http.MethodDelete:
		deleter, ok := target.(resolvedDeleter)
		if !ok {
			return errors.New("2finance: resolved operation target does not support DELETE")
		}
		return deleter.Delete(ctx, operation.Path, out)
	default:
		return errors.New("2finance: unsupported resolved operation method " + operation.Method)
	}
}
