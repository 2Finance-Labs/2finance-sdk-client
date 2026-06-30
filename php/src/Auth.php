<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

final class Auth
{
    public static function bearerAuthorization(string $accessToken): string
    {
        $trimmed = trim($accessToken);
        if ($trimmed === '') {
            return '';
        }
        if (str_starts_with(strtolower($trimmed), 'bearer ')) {
            return $trimmed;
        }
        return 'Bearer ' . $trimmed;
    }
}
