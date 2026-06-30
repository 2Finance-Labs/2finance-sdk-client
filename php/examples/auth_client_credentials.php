<?php

declare(strict_types=1);

require __DIR__ . '/../src/Auth.php';
require __DIR__ . '/../src/TokenSource.php';
require __DIR__ . '/../src/StaticTokenSource.php';
require __DIR__ . '/../src/ServiceClient.php';
require __DIR__ . '/../src/ClientCredentialsTokenSource.php';
require __DIR__ . '/../src/SdkConfig.php';
require __DIR__ . '/../src/DomainClients.php';
require __DIR__ . '/../src/SdkClient.php';

use TwoFinance\Sdk\ClientCredentialsTokenSource;
use TwoFinance\Sdk\HttpRequest;
use TwoFinance\Sdk\HttpResponse;
use TwoFinance\Sdk\SdkClient;
use TwoFinance\Sdk\SdkConfig;

$tokenTransport = static function (HttpRequest $request): HttpResponse {
    return new HttpResponse(200, '{"access_token":"example-token","expires_in":300}');
};

$config = SdkConfig::fromEnv();
$config->tokenSource = new ClientCredentialsTokenSource(
    getenv('TWO_FINANCE_AUTH_TOKEN_URL') ?: '',
    getenv('TWO_FINANCE_AUTH_CLIENT_ID') ?: '',
    getenv('TWO_FINANCE_AUTH_CLIENT_SECRET') ?: '',
    ['2finance.sdk'],
    $tokenTransport,
);

$client = new SdkClient($config);
echo json_encode($client->analytics->indicators(), JSON_THROW_ON_ERROR) . PHP_EOL;
