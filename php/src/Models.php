<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

final class SdkErrorPayload
{
    public function __construct(
        public readonly string $error,
        public readonly string $message,
        public readonly string $code,
        /** @var array<string,mixed> */
        public readonly array $details = [],
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            (string) $payload['error'],
            (string) $payload['message'],
            (string) $payload['code'],
            $payload['details'] ?? [],
        );
    }
}

final class PaginationResponse
{
    /**
     * @param array<int,array<string,mixed>> $items
     */
    public function __construct(
        public readonly array $items,
        public readonly int $limit,
        public readonly ?string $cursor = null,
        public readonly ?string $nextCursor = null,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            $payload['items'] ?? [],
            (int) $payload['limit'],
            $payload['cursor'] ?? null,
            $payload['next_cursor'] ?? null,
        );
    }
}

final class IdempotencyRecord
{
    public function __construct(
        public readonly string $idempotencyKey,
        public readonly string $operation,
        public readonly string $scope,
        public readonly string $requestId,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            (string) $payload['idempotency_key'],
            (string) $payload['operation'],
            (string) $payload['scope'],
            (string) $payload['request_id'],
        );
    }
}

final class ServiceCatalogEntry
{
    public function __construct(
        public readonly string $name,
        public readonly string $env,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self((string) $payload['name'], (string) $payload['env']);
    }
}

final class ServiceCatalog
{
    /** @param array<int,ServiceCatalogEntry> $services */
    public function __construct(public readonly array $services)
    {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(array_map(
            static fn(array $entry): ServiceCatalogEntry => ServiceCatalogEntry::fromArray($entry),
            $payload['services'] ?? [],
        ));
    }
}

final class ConfiguredServiceEntry
{
    public function __construct(
        public readonly string $name,
        public readonly string $env,
        public readonly string $url,
    ) {
    }
}

final class DomainOperation
{
    /**
     * @param array<int,string> $pathParams
     * @param array<int,string> $query
     */
    public function __construct(
        public readonly string $name,
        public readonly string $method,
        public readonly string $path,
        public readonly array $pathParams = [],
        public readonly array $query = [],
        public readonly ?string $requestSchema = null,
        public readonly ?string $responseSchema = null,
        public readonly ?string $notes = null,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            (string) $payload['name'],
            (string) $payload['method'],
            (string) $payload['path'],
            $payload['path_params'] ?? [],
            $payload['query'] ?? [],
            $payload['request_schema'] ?? null,
            $payload['response_schema'] ?? null,
            $payload['notes'] ?? null,
        );
    }

    /**
     * @param array<string,scalar> $pathParams
     * @param array<string,scalar|null> $query
     */
    public function resolve(array $pathParams = [], array $query = []): ResolvedOperation
    {
        $path = $this->path;
        foreach ($this->pathParams as $name) {
            if (!array_key_exists($name, $pathParams)) {
                throw new \InvalidArgumentException("2finance: missing operation path parameter {$name}");
            }
            $path = str_replace('{' . $name . '}', rawurlencode((string) $pathParams[$name]), $path);
        }

        $queryParams = [];
        foreach ($this->query as $name) {
            if (array_key_exists($name, $query) && $query[$name] !== null) {
                $queryParams[$name] = $query[$name];
            }
        }
        if ($queryParams !== []) {
            $path .= str_contains($path, '?') ? '&' : '?';
            $path .= http_build_query($queryParams, '', '&', PHP_QUERY_RFC3986);
        }

        return new ResolvedOperation(strtoupper(trim($this->method)), $path);
    }
}

final class ResolvedOperation
{
    public function __construct(
        public readonly string $method,
        public readonly string $path,
    ) {
    }
}

final class DomainOperationsDomain
{
    /** @param array<int,DomainOperation> $operations */
    public function __construct(
        public readonly string $name,
        public readonly string $env,
        public readonly array $operations,
        public readonly ?string $transport = null,
        public readonly ?string $description = null,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            (string) $payload['name'],
            (string) $payload['env'],
            array_map(
                static fn(array $entry): DomainOperation => DomainOperation::fromArray($entry),
                $payload['operations'] ?? [],
            ),
            $payload['transport'] ?? null,
            $payload['description'] ?? null,
        );
    }
}

final class DomainOperationsCatalog
{
    /** @param array<int,DomainOperationsDomain> $domains */
    public function __construct(
        public readonly string $schema,
        public readonly array $domains,
    ) {
    }

    /** @param array<string,mixed> $payload */
    public static function fromArray(array $payload): self
    {
        return new self(
            (string) $payload['schema'],
            array_map(
                static fn(array $entry): DomainOperationsDomain => DomainOperationsDomain::fromArray($entry),
                $payload['domains'] ?? [],
            ),
        );
    }

    public function operation(string $domainName, string $operationName): ?DomainOperation
    {
        foreach ($this->domains as $domain) {
            if (self::domainKey($domain->name) !== self::domainKey($domainName)) {
                continue;
            }
            foreach ($domain->operations as $operation) {
                if ($operation->name === $operationName) {
                    return $operation;
                }
            }
            return null;
        }
        return null;
    }

    /**
     * @param array<string,scalar> $pathParams
     * @param array<string,scalar|null> $query
     */
    public function resolveOperation(string $domainName, string $operationName, array $pathParams = [], array $query = []): ResolvedOperation
    {
        $operation = $this->operation($domainName, $operationName);
        if ($operation === null) {
            throw new \InvalidArgumentException("2finance: unknown operation {$domainName}.{$operationName}");
        }
        return $operation->resolve($pathParams, $query);
    }

    private static function domainKey(string $domain): string
    {
        return str_replace(['-', '_', ' '], '', strtolower(trim($domain)));
    }
}
