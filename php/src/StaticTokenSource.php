<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

final class StaticTokenSource implements TokenSource
{
    public function __construct(private readonly string $accessToken)
    {
    }

    public function token(): string
    {
        return $this->accessToken;
    }
}
