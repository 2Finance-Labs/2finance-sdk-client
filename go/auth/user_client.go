package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/2Finance-Labs/2finance-sdk-client/internal/service"
)

// UserClientConfig configures calls to 2finance-auth.
type UserClientConfig struct {
	BaseURL       string
	Realm         string
	ClientID      string
	PhoneClientID string
	HTTPClient    *http.Client
}

// UserClient exposes 2finance-auth user and PKCE flows.
type UserClient struct {
	config  UserClientConfig
	service *service.Client
}

// NewUserClient returns a user auth client.
func NewUserClient(config UserClientConfig) *UserClient {
	return &UserClient{
		config:  config,
		service: service.New(config.BaseURL, config.HTTPClient),
	}
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserInput struct {
	Username   string              `json:"username"`
	Email      string              `json:"email"`
	FirstName  string              `json:"firstName,omitempty"`
	LastName   string              `json:"lastName,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
	Password   string              `json:"password,omitempty"`
}

type SMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	UserID      string `json:"user_id"`
}

type VerifySMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	UserID      string `json:"user_id"`
}

type PhoneLoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token,omitempty"`
}

type DefaultResponse struct {
	Message string          `json:"message,omitempty"`
	Status  string          `json:"status,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type PKCELoginResponse struct {
	AuthURL    string
	State      string
	Cookies    []*http.Cookie
	StatusCode int
}

// Login authenticates a username/password pair.
func (c *UserClient) Login(ctx context.Context, input LoginInput) (*JWT, error) {
	var jwt JWT
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/login"), input, &jwt); err != nil {
		return nil, err
	}
	return &jwt, nil
}

func (c *UserClient) SignUp(ctx context.Context, input CreateUserInput) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/signup"), input, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) RefreshToken(ctx context.Context, refreshToken string) (*JWT, error) {
	var jwt JWT
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/refresh"), map[string]string{"refresh_token": refreshToken}, &jwt); err != nil {
		return nil, err
	}
	return &jwt, nil
}

func (c *UserClient) Logout(ctx context.Context, refreshToken string) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/logout"), map[string]string{"refresh_token": refreshToken}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) PhoneLogin(ctx context.Context, phoneNumber string, code string) (*JWT, error) {
	var jwt JWT
	input := PhoneLoginRequest{PhoneNumber: phoneNumber, Code: code}
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.PhoneClientID, "/phone/sms/login"), input, &jwt); err != nil {
		return nil, err
	}
	return &jwt, nil
}

func (c *UserClient) JWKS(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	if err := c.service.DoJSON(ctx, http.MethodGet, c.oidcPath("/protocol/openid-connect/certs"), nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *UserClient) ValidateToken(ctx context.Context, token string) (json.RawMessage, error) {
	var response json.RawMessage
	if err := c.service.DoJSON(ctx, http.MethodPost, c.oidcPath("/protocol/openid-connect/token/introspect"), map[string]string{"token": token}, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *UserClient) Request(ctx context.Context, method string, path string, body any, out any) error {
	return c.service.DoJSON(ctx, method, path, body, out)
}

func (c *UserClient) RequestAuthenticationCode(ctx context.Context, phoneNumber string) (*DefaultResponse, error) {
	var response DefaultResponse
	input := map[string]string{"phone_number": phoneNumber}
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.PhoneClientID, "/phone/sms/request-code"), input, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) GetUserInfo(ctx context.Context) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodGet, c.path(c.config.ClientID, "/user-info"), nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) CreateUser(ctx context.Context, input CreateUserInput) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/create-user"), input, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) RequestSMSCode(ctx context.Context, input SMSRequest) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/request-sms"), input, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) VerifySMSCode(ctx context.Context, input VerifySMSRequest) (*DefaultResponse, error) {
	var response DefaultResponse
	if err := c.service.DoJSON(ctx, http.MethodPost, c.path(c.config.ClientID, "/verify-sms"), input, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *UserClient) LoginPKCE(ctx context.Context) (*PKCELoginResponse, error) {
	requestURL, err := c.service.URL(c.path(c.config.ClientID, "/pkce/login-redirect"))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	client := *c.service.HTTPClient
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusMultipleChoices || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("auth: pkce redirect returned %d", resp.StatusCode)
	}
	state := ""
	if parsed, err := url.Parse(resp.Header.Get("Location")); err == nil {
		state = parsed.Query().Get("state")
	}
	return &PKCELoginResponse{
		AuthURL:    resp.Header.Get("Location"),
		State:      state,
		Cookies:    resp.Cookies(),
		StatusCode: resp.StatusCode,
	}, nil
}

func (c *UserClient) CallbackPKCE(ctx context.Context, code string, state string) (*JWT, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("auth: pkce code is required")
	}
	if strings.TrimSpace(state) == "" {
		return nil, fmt.Errorf("auth: pkce state is required")
	}
	path := c.path(c.config.ClientID, "/pkce/callback") + "?code=" + url.QueryEscape(code) + "&state=" + url.QueryEscape(state)
	var jwt JWT
	if err := c.service.DoJSON(ctx, http.MethodGet, path, nil, &jwt); err != nil {
		return nil, err
	}
	return &jwt, nil
}

func (c *UserClient) path(clientID string, endpoint string) string {
	realm := strings.Trim(c.config.Realm, "/")
	if realm == "" {
		realm = "2finance"
	}
	if strings.TrimSpace(clientID) == "" {
		clientID = "2finance-network"
	}
	return "/v1/2finance-authenticator/" + realm + "/" + strings.Trim(clientID, "/") + "/" + strings.TrimLeft(endpoint, "/")
}

func (c *UserClient) oidcPath(endpoint string) string {
	realm := strings.Trim(c.config.Realm, "/")
	if realm == "" {
		realm = "2finance"
	}
	return "/realms/" + realm + "/" + strings.TrimLeft(endpoint, "/")
}
