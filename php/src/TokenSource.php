<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

interface TokenSource
{
    public function token(): string;
}
