# Authentication and Authorization

## Papel do repo

O pacote Dart blockchain dentro de `2finance-sdk-client` deve ser cliente Dart para APIs 2Finance. Como pode rodar em ambiente publico, nao deve conter segredos e deve receber tokens OIDC ja obtidos por PKCE.

## Avaliacao atual

- Ha ferramentas de smoke para orquestrador, mas a autenticacao deve ser padronizada com bearer token em vez de headers de identidade soltos.
- `HttpFinanceNetworkTransport` aceita `TokenProvider` injetado e envia `Authorization: Bearer` nas chamadas protegidas.
- O mesmo transport continua funcionando sem token provider ou com token vazio, cobrindo modo local/dev com auth desligada.
- Ha teste garantindo envio do header, ausencia de header quando sem token e ausencia do token em erro HTTP.

## Padrao seguro obrigatorio

- Aceitar access token por provider/injecao, nao por constante hardcoded.
- Adicionar `Authorization: Bearer` em chamadas protegidas.
- Nao armazenar refresh token neste pacote, a menos que exista storage seguro especifico da plataforma.
- Mascarar tokens em logs, exceptions e fixtures.

## Proximos passos

1. Atualizar exemplos para usar `2finance-auth-dart-sdk`.
2. Plugar o `TokenProvider` do app/mobile no transport usado em producao.
3. Expandir redaction para outros transports quando eles passarem a carregar tokens.
