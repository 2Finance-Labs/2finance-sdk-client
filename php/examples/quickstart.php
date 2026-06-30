<?php

declare(strict_types=1);

require __DIR__ . '/../src/Auth.php';
require __DIR__ . '/../src/TokenSource.php';
require __DIR__ . '/../src/StaticTokenSource.php';
require __DIR__ . '/../src/ServiceClient.php';
require __DIR__ . '/../src/SdkConfig.php';
require __DIR__ . '/../src/DomainClients.php';
require __DIR__ . '/../src/SdkClient.php';

use TwoFinance\Sdk\SdkClient;

$client = SdkClient::fromEnv();

$indicators = $client->analytics->indicators();
echo 'analytics indicators: ' . json_encode($indicators, JSON_THROW_ON_ERROR) . PHP_EOL;

$plan = $client->planner->tradingPlan([
    'goal' => 'prepare a BTC rebalancing plan',
    'useAnalytics' => true,
    'useTrading' => true,
]);
echo 'planner response: ' . json_encode($plan, JSON_THROW_ON_ERROR) . PHP_EOL;
