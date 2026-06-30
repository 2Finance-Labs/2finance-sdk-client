# Specs

Este diretório guarda especificações executáveis e revisáveis do SDK.

As specs usam JSON para evitar dependências novas no projeto. O formato é propositalmente simples:

```json
{
  "id": "wallet_v1",
  "title": "Wallet V1 client contract",
  "owner": "sdk",
  "status": "active",
  "layers": ["harness-small", "contract-unit"],
  "cases": [
    {
      "id": "wallet.add",
      "title": "Add wallet sends Go-compatible payload",
      "tags": ["wallet", "mqtt"],
      "test_refs": ["test/blockchain/contract/walletV1/wallet_unit_test.dart"],
      "given": {},
      "when": {},
      "then": {}
    }
  ]
}
```

Campos obrigatórios:

- `id`: identificador estável em snake/dot/kebab case.
- `title`: descrição curta.
- `owner`: dono da spec.
- `status`: `draft`, `active` ou `deprecated`.
- `layers`: onde a spec deve ser validada.
- `cases`: lista de exemplos concretos.

Campos obrigatórios em cada case:

- `id`
- `title`
- `tags`
- `test_refs` quando `layers` inclui `contract-unit`
- `given`
- `when`
- `then`

Rode:

```bash
HOME=/tmp dart test test/spec_harness/spec_files_test.dart
```
