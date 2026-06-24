# Authentication and Authorization

## Papel do repo

`2finance-go-client` deve consumir APIs e protocolos 2Finance a partir de backends, CLIs ou testes. Ele deve suportar bearer tokens OIDC e client credentials para automacao confiavel.

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

## Client credentials para jobs internos

Jobs internos devem obter tokens via OAuth2 client credentials em um provedor
OIDC aprovado. O `client_id`, `client_secret`, issuer/token URL e scopes devem
vir do secret manager ou do ambiente de execucao do job, nunca de constantes no
codigo. O job deve solicitar apenas os scopes necessarios, como
`finance:transaction:submit` ou `finance:mcp:call`, e entregar o access token ao
cliente HTTP por um `TokenSource` com refresh.

O token OIDC autentica o job perante APIs HTTP/MCP. Ele nao substitui a chave da
wallet nem autoriza o cliente a assinar transacoes sem passar pelos fluxos de
wallet/signing existentes.

## Estado atual

- `auth.AuthTransport` injeta bearer token em requests HTTP.
- `auth.RedactAuthorization` mascara `Authorization` antes de logs/asserts.
- O cliente principal ainda usa MQTT para operacoes de rede; HTTP aparece no
  harness MCP/E2E.
