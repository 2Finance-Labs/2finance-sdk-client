# Authentication and Authorization

## Papel do repo

`2finance-sdk-client` deve consumir APIs e protocolos 2Finance a partir de backends, CLIs ou testes. Ele deve suportar bearer tokens OIDC e client credentials para automacao confiavel.

## Avaliacao atual

- O repo tem foco em wallet/protocol/lifecycle. Operacoes de assinatura devem continuar separadas de autenticacao OIDC.
- Qualquer identidade externa usada em payloads deve ser conferida contra claims no backend.

## Padrao seguro obrigatorio

- Receber access token por `TokenSource`/provider, nao por global.
- Suportar refresh via OAuth2 quando usado por servico confidencial.
- Nao logar transacoes assinadas junto com tokens ou credenciais.
- Manter separacao entre chave de wallet/signing e token OIDC.
- Exigir scopes especificos para operacoes write/sign/submit no servidor.

## HTTP clients

Use `auth.AuthTransport` for HTTP calls that need bearer auth. The transport
pulls the token from a caller-owned `auth.TokenSource`, injects
`Authorization: Bearer <token>`, and leaves the original request unchanged.

```go
httpClient := &http.Client{
	Transport: auth.AuthTransport{
		Source: tokenSource,
		Base:   http.DefaultTransport,
	},
}
```

For local or E2E MCP calls, set `MCP_ACCESS_TOKEN` when the MCP HTTP endpoint
requires bearer auth. Tests and logs should redact headers with
`auth.RedactAuthorization` before printing request or response metadata.
The MCP E2E harness also accepts client credentials envs when a static token is
not provided: `MCP_OIDC_TOKEN_URL`, `MCP_OIDC_CLIENT_ID`,
`MCP_OIDC_CLIENT_SECRET` and `MCP_OIDC_SCOPES`. The shared infra harness aliases
`AUTH_E2E_TOKEN_URL`, `AUTH_E2E_CLIENT_ID`, `AUTH_E2E_CLIENT_SECRET` and
`AUTH_E2E_SCOPE` are also supported.

## Client credentials para jobs internos

Jobs internos devem obter tokens via OAuth2 client credentials em um provedor
OIDC aprovado. O `client_id`, `client_secret`, issuer/token URL e scopes devem
vir do secret manager ou do ambiente de execucao do job, nunca de constantes no
codigo. O job deve solicitar apenas os scopes necessarios, como
`finance:transaction:submit` ou `finance:mcp:call`, e entregar o access token ao
cliente HTTP por um `TokenSource` com refresh.

```go
tokenSource, err := auth.NewClientCredentialsTokenSource(auth.ClientCredentialsConfig{
	TokenURL:     os.Getenv("TWO_FINANCE_OIDC_TOKEN_URL"),
	ClientID:     os.Getenv("TWO_FINANCE_OIDC_CLIENT_ID"),
	ClientSecret: os.Getenv("TWO_FINANCE_OIDC_CLIENT_SECRET"),
	Scopes:       []string{"network:execute"},
})
```

O token OIDC autentica o job perante APIs HTTP/MCP. Ele nao substitui a chave da
wallet nem autoriza o cliente a assinar transacoes sem passar pelos fluxos de
wallet/signing existentes.

## Estado atual

- `auth.AuthTransport` injeta bearer token em requests HTTP.
- `auth.ClientCredentialsTokenSource` busca e cacheia access tokens via OAuth2 client credentials.
- `auth.BearerAuthorization` normaliza o header sem duplicar o prefixo `Bearer`.
- `auth.RedactAuthorization` mascara `Authorization` antes de logs/asserts.
- `auth.RedactSensitive` mascara `Bearer`, `access_token`, `refresh_token`, `id_token`, `client_secret`, senha e `code` em mensagens de erro/logs.
- Em `TWO_FINANCE_ENV=prod|production|prod_secrets`, token endpoint externo `http://` e rejeitado; `localhost`/`127.0.0.1` segue permitido para testes.
- O E2E MCP usa `MCP_ACCESS_TOKEN` ou client credentials OIDC por envs `MCP_OIDC_*`/`AUTH_E2E_*` e redige erros HTTP/JSON-RPC antes de reportar falhas.
- O workspace Go principal fica em `go/` e agrega auth, network, analytics,
  orchestrator, MCP, planner, trading control, matchengine, hummingbot,
  keystore e providers.
