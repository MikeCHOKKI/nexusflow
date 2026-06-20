<?php

declare(strict_types=1);

/**
 * NexusFlow — Order Service Configuration
 *
 * Charge et expose la configuration depuis les variables d'environnement.
 * Compatible .env (via phpdotenv) et variables système.
 */

use Dotenv\Dotenv;

// Charger .env depuis la racine du projet si présent
$rootDir = dirname(__DIR__);
if (file_exists($rootDir . '/.env')) {
    $dotenv = Dotenv::createImmutable($rootDir);
    $dotenv->safeLoad();
}

// Surcharger avec .env local du service
$serviceEnv = __DIR__ . '/.env';
if (file_exists($serviceEnv)) {
    $dotenv = Dotenv::createImmutable(__DIR__);
    $dotenv->safeLoad();
}

return [
    'app' => [
        'port'     => (int) ($_ENV['APP_PORT'] ?? 50053),
        'env'      => $_ENV['APP_ENV'] ?? 'production',
        'debug'    => filter_var($_ENV['APP_DEBUG'] ?? false, FILTER_VALIDATE_BOOLEAN),
    ],

    'database' => [
        'url' => $_ENV['DATABASE_URL'] ?? 'postgresql://nexusflow:nexusflow_dev@postgres:5432/nexusflow',
    ],

    'redis' => [
        'url'     => $_ENV['REDIS_URL'] ?? 'redis://redis:6379',
        'channel' => $_ENV['NOTIFY_CHANNEL'] ?? 'nexusflow:notifications',
    ],

    'services' => [
        'payment' => [
            'addr' => $_ENV['PAYMENT_SERVICE_ADDR'] ?? 'service-payment:50054',
        ],
    ],
];
