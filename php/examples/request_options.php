<?php

declare(strict_types=1);

require __DIR__ . '/../src/Auth.php';
require __DIR__ . '/../src/TokenSource.php';
require __DIR__ . '/../src/StaticTokenSource.php';
require __DIR__ . '/../src/ServiceClient.php';
require __DIR__ . '/../src/SdkConfig.php';
require __DIR__ . '/../src/DomainClients.php';
require __DIR__ . '/../src/SdkClient.php';

use TwoFinance\Sdk\RequestOptions;
use TwoFinance\Sdk\SdkClient;

$client = SdkClient::fromEnv();
$response = $client->analytics->post(
    '/analytics/candles:upsert',
    ['symbol' => 'BTC-USDT'],
    new RequestOptions(
        headers: ['X-Trace-ID' => 'trace-1'],
        idempotencyKey: 'candles-upsert-001',
        query: ['source' => 'sdk-example'],
        timeoutSeconds: 5.0,
        maxRetries: 1,
        page: 1,
        limit: 25,
    ),
);

echo 'response: ' . json_encode($response, JSON_THROW_ON_ERROR) . PHP_EOL;
