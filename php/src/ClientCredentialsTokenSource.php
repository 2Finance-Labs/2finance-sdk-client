<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

use RuntimeException;

final class ClientCredentialsTokenSource implements TokenSource
{
    private string $accessToken = '';
    private int $expiresAt = 0;

    /**
     * @param list<string> $scopes
     * @param null|callable(HttpRequest):HttpResponse $transport
     */
    public function __construct(
        private readonly string $tokenUrl,
        private readonly string $clientId,
        private readonly string $clientSecret,
        private readonly array $scopes = [],
        private readonly mixed $transport = null,
        private readonly int $expirySkewSeconds = 30,
    ) {
    }

    public function token(): string
    {
        $now = time();
        if ($this->accessToken !== '' && $now < ($this->expiresAt - $this->expirySkewSeconds)) {
            return $this->accessToken;
        }
        if (trim($this->tokenUrl) === '' || trim($this->clientId) === '' || trim($this->clientSecret) === '') {
            throw new RuntimeException('2finance auth: tokenUrl, clientId and clientSecret are required');
        }
        $transport = $this->transport;
        if (!is_callable($transport)) {
            throw new RuntimeException('2finance auth: transport is required');
        }
        $body = http_build_query([
            'grant_type' => 'client_credentials',
            'client_id' => $this->clientId,
            'client_secret' => $this->clientSecret,
            ...($this->scopes === [] ? [] : ['scope' => implode(' ', $this->scopes)]),
        ]);
        $response = $transport(new HttpRequest(
            'POST',
            $this->tokenUrl,
            [
                'Accept' => 'application/json',
                'Content-Type' => 'application/x-www-form-urlencoded',
            ],
            $body,
        ));
        if (!$response instanceof HttpResponse) {
            throw new RuntimeException('2finance auth: transport must return HttpResponse');
        }
        if ($response->statusCode < 200 || $response->statusCode >= 300) {
            throw new RuntimeException(sprintf('2finance auth: token endpoint returned %d', $response->statusCode));
        }
        $payload = json_decode($response->body, true, 512, JSON_THROW_ON_ERROR);
        if (!is_array($payload) || !isset($payload['access_token']) || !is_string($payload['access_token'])) {
            throw new RuntimeException('2finance auth: token response missing access_token');
        }
        $this->accessToken = $payload['access_token'];
        $this->expiresAt = $now + (int) ($payload['expires_in'] ?? 300);
        return $this->accessToken;
    }
}
