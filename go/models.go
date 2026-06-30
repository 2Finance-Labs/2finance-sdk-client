package twofinance

import (
	"errors"
	"net/url"
	"strings"
)

type SDKError struct {
	Error   string         `json:"error"`
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Details map[string]any `json:"details,omitempty"`
}

type PaginationResponse struct {
	Items      []map[string]any `json:"items"`
	Limit      int              `json:"limit"`
	Cursor     string           `json:"cursor,omitempty"`
	NextCursor string           `json:"next_cursor,omitempty"`
}

type IdempotencyRecord struct {
	IdempotencyKey string `json:"idempotency_key"`
	Operation      string `json:"operation"`
	Scope          string `json:"scope"`
	RequestID      string `json:"request_id"`
}

type ServiceCatalogEntry struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type ServiceCatalog struct {
	Services []ServiceCatalogEntry `json:"services"`
}

type ConfiguredServiceEntry struct {
	Name string `json:"name"`
	Env  string `json:"env"`
	URL  string `json:"url"`
}

type DomainOperation struct {
	Name           string   `json:"name"`
	Method         string   `json:"method"`
	Path           string   `json:"path"`
	PathParams     []string `json:"path_params,omitempty"`
	Query          []string `json:"query,omitempty"`
	RequestSchema  string   `json:"request_schema,omitempty"`
	ResponseSchema string   `json:"response_schema,omitempty"`
	Notes          string   `json:"notes,omitempty"`
}

type ResolvedOperation struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (o DomainOperation) Resolve(pathParams map[string]string, query map[string]string) (ResolvedOperation, error) {
	path := o.Path
	for _, name := range o.PathParams {
		value, ok := pathParams[name]
		if !ok {
			return ResolvedOperation{}, errors.New("2finance: missing operation path parameter " + name)
		}
		path = strings.ReplaceAll(path, "{"+name+"}", url.PathEscape(value))
	}

	values := url.Values{}
	for _, name := range o.Query {
		value, ok := query[name]
		if ok {
			values.Set(name, value)
		}
	}
	if encoded := values.Encode(); encoded != "" {
		if strings.Contains(path, "?") {
			path += "&" + encoded
		} else {
			path += "?" + encoded
		}
	}

	return ResolvedOperation{
		Method: strings.ToUpper(strings.TrimSpace(o.Method)),
		Path:   path,
	}, nil
}

type DomainOperationsDomain struct {
	Name        string            `json:"name"`
	Env         string            `json:"env"`
	Transport   string            `json:"transport,omitempty"`
	Description string            `json:"description,omitempty"`
	Operations  []DomainOperation `json:"operations"`
}

type DomainOperationsCatalog struct {
	Schema  string                   `json:"schema"`
	Domains []DomainOperationsDomain `json:"domains"`
}

func (c DomainOperationsCatalog) Operation(domainName string, operationName string) (DomainOperation, bool) {
	for _, domain := range c.Domains {
		if serviceKey(domain.Name) != serviceKey(domainName) {
			continue
		}
		for _, operation := range domain.Operations {
			if operation.Name == operationName {
				return operation, true
			}
		}
		return DomainOperation{}, false
	}
	return DomainOperation{}, false
}

func (c DomainOperationsCatalog) ResolveOperation(domainName string, operationName string, pathParams map[string]string, query map[string]string) (ResolvedOperation, error) {
	operation, ok := c.Operation(domainName, operationName)
	if !ok {
		return ResolvedOperation{}, errors.New("2finance: unknown operation " + domainName + "." + operationName)
	}
	return operation.Resolve(pathParams, query)
}
