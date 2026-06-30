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
use TwoFinance\Sdk\ServiceException;

$client = SdkClient::fromEnv();

try {
    $client->analytics->indicators();
} catch (ServiceException $exception) {
    echo sprintf(
        'request failed with status %d: %s%s',
        $exception->statusCode,
        $exception->body,
        PHP_EOL,
    );
}
