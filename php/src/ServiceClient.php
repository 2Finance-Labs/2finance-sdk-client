<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

use RuntimeException;

final class HttpResponse
{
    public function __construct(
        public readonly int $statusCode,
        public readonly string $body,
    ) {
    }
}

final class HttpRequest
{
    /**
     * @param array<string,string> $headers
     */
    public function __construct(
        public readonly string $method,
        public readonly string $url,
        public readonly array $headers,
        public readonly ?string $body,
        public readonly ?float $timeoutSeconds = null,
    ) {
    }
}

final class RequestOptions
{
    /**
     * @param array<string,string> $headers
     * @param array<string,scalar|null> $query
     */
    public function __construct(
        public readonly array $headers = [],
        public readonly ?string $idempotencyKey = null,
        public readonly array $query = [],
        public readonly ?float $timeoutSeconds = null,
        public readonly int $maxRetries = 0,
        public readonly ?int $page = null,
        public readonly ?int $limit = null,
    ) {
    }
}

final class ServiceException extends RuntimeException
{
    public function __construct(
        public readonly string $method,
        public readonly string $url,
        public readonly int $statusCode,
        public readonly string $body,
    ) {
        parent::__construct(sprintf('2finance: %s %s returned %d: %s', $method, $url, $statusCode, $body));
    }
}

final class ServiceClient
{
    /** @var callable(HttpRequest):HttpResponse */
    private $transport;

    public function __construct(
        private readonly string $baseUrl,
        ?callable $transport = null,
        private readonly ?TokenSource $tokenSource = null,
    ) {
        $this->transport = $transport ?? self::defaultTransport();
    }

    public function url(string $path): string
    {
        if (str_starts_with($path, 'http://') || str_starts_with($path, 'https://')) {
            return $path;
        }
        $base = rtrim(trim($this->baseUrl), '/');
        if ($base === '') {
            throw new RuntimeException('2finance: baseUrl is required');
        }
        return $base . '/' . ltrim($path, '/');
    }

    /**
     * @return mixed
     */
    public function get(string $path, ?RequestOptions $options = null): mixed
    {
        return $this->request('GET', $path, null, $options);
    }

    /**
     * @param mixed $body
     * @return mixed
     */
    public function post(string $path, mixed $body = null, ?RequestOptions $options = null): mixed
    {
        return $this->request('POST', $path, $body, $options);
    }

    /**
     * @param mixed $body
     * @return mixed
     */
    public function put(string $path, mixed $body = null, ?RequestOptions $options = null): mixed
    {
        return $this->request('PUT', $path, $body, $options);
    }

    /**
     * @return mixed
     */
    public function delete(string $path, ?RequestOptions $options = null): mixed
    {
        return $this->request('DELETE', $path, null, $options);
    }

    /**
     * @param mixed $body
     * @return mixed
     */
    public function request(string $method, string $path, mixed $body = null, ?RequestOptions $options = null): mixed
    {
        $headers = ['Accept' => 'application/json'];
        $payload = null;
        if ($body !== null) {
            $headers['Content-Type'] = 'application/json';
            $payload = is_string($body) ? $body : json_encode($body, JSON_THROW_ON_ERROR);
        }
        if ($this->tokenSource !== null) {
            $authorization = Auth::bearerAuthorization($this->tokenSource->token());
            if ($authorization !== '') {
                $headers['Authorization'] = $authorization;
            }
        }
        if ($options !== null) {
            foreach ($options->headers as $key => $value) {
                $headers[$key] = $value;
            }
            $idempotencyKey = trim($options->idempotencyKey ?? '');
            if ($idempotencyKey !== '') {
                $headers['Idempotency-Key'] = $idempotencyKey;
            }
        }
        $requestUrl = $this->url($path);
        if ($options !== null && ($options->query !== [] || $options->page !== null || $options->limit !== null)) {
            $requestUrl = self::urlWithQuery($requestUrl, $options);
        }
        $maxRetries = max(0, $options?->maxRetries ?? 0);
        for ($attempt = 0; $attempt <= $maxRetries; $attempt++) {
            $response = ($this->transport)(new HttpRequest($method, $requestUrl, $headers, $payload, $options?->timeoutSeconds));
            if ($response->statusCode >= 200 && $response->statusCode < 300) {
                return $response->body === '' ? null : json_decode($response->body, true, 512, JSON_THROW_ON_ERROR);
            }
            if ($attempt >= $maxRetries || !self::isRetryableStatus($response->statusCode)) {
                throw new ServiceException($method, $requestUrl, $response->statusCode, $response->body);
            }
        }
        throw new ServiceException($method, $requestUrl, 0, '');
    }

    /**
     * @param mixed $body
     * @return mixed
     */
    public function requestOperation(ResolvedOperation $operation, mixed $body = null, ?RequestOptions $options = null): mixed
    {
        return $this->request($operation->method, $operation->path, $body, $options);
    }

    /**
     * @param array<string,scalar> $pathParams
     * @param array<string,scalar|null> $query
     * @param mixed $body
     * @return mixed
     */
    public function requestCatalogOperation(
        DomainOperationsCatalog $catalog,
        string $domainName,
        string $operationName,
        array $pathParams = [],
        array $query = [],
        mixed $body = null,
        ?RequestOptions $options = null,
    ): mixed {
        return $this->requestOperation(
            $catalog->resolveOperation($domainName, $operationName, $pathParams, $query),
            $body,
            $options,
        );
    }

    /**
     * @return callable(HttpRequest):HttpResponse
     */
    private static function defaultTransport(): callable
    {
        return static function (HttpRequest $request): HttpResponse {
            $headerLines = [];
            foreach ($request->headers as $key => $value) {
                $headerLines[] = $key . ': ' . $value;
            }
            $httpOptions = [
                'method' => $request->method,
                'header' => implode("\r\n", $headerLines),
                'content' => $request->body ?? '',
                'ignore_errors' => true,
            ];
            if ($request->timeoutSeconds !== null) {
                $httpOptions['timeout'] = $request->timeoutSeconds;
            }
            $context = stream_context_create(['http' => $httpOptions]);
            $body = file_get_contents($request->url, false, $context);
            $statusCode = 0;
            foreach ($http_response_header ?? [] as $header) {
                if (preg_match('/^HTTP\/\S+\s+(\d+)/', $header, $matches)) {
                    $statusCode = (int) $matches[1];
                    break;
                }
            }
            return new HttpResponse($statusCode, $body === false ? '' : $body);
        };
    }

    /**
     */
    private static function urlWithQuery(string $url, RequestOptions $options): string
    {
        $parts = parse_url($url);
        $existing = [];
        if (isset($parts['query'])) {
            parse_str($parts['query'], $existing);
        }
        foreach ($options->query as $key => $value) {
            if ($value !== null) {
                $existing[$key] = $value;
            }
        }
        if ($options->page !== null) {
            $existing['page'] = $options->page;
        }
        if ($options->limit !== null) {
            $existing['limit'] = $options->limit;
        }
        $scheme = isset($parts['scheme']) ? $parts['scheme'] . '://' : '';
        $host = $parts['host'] ?? '';
        $port = isset($parts['port']) ? ':' . $parts['port'] : '';
        $path = $parts['path'] ?? '';
        $fragment = isset($parts['fragment']) ? '#' . $parts['fragment'] : '';
        $queryString = $existing === [] ? '' : '?' . http_build_query($existing);
        return $scheme . $host . $port . $path . $queryString . $fragment;
    }

    private static function isRetryableStatus(int $statusCode): bool
    {
        return $statusCode === 429 || $statusCode >= 500;
    }
}
