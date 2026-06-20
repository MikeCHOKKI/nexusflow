<?php

declare(strict_types=1);

/**
 * NexusFlow — Order Service Entry Point
 *
 * Point d'entrée du serveur PHP built-in (php -S 0.0.0.0:50053 server.php).
 * Route les requêtes HTTP vers le routeur REST.
 *
 * Usage :
 *   php -S 0.0.0.0:50053 server.php
 */

// ─── Autoload ─────────────────────────────────────────────────
$autoload = __DIR__ . '/vendor/autoload.php';
if (!file_exists($autoload)) {
    http_response_code(500);
    echo json_encode([
        'success' => false,
        'error'   => 'Vendor autoload introuvable. Lancez "composer install"',
    ]);
    exit(1);
}

require_once $autoload;

// ─── Configuration ────────────────────────────────────────────
$config = require __DIR__ . '/config.php';

// ─── Health Check (avant routing) ─────────────────────────────
$uri = parse_url($_SERVER['REQUEST_URI'] ?? '/', PHP_URL_PATH);

if ($uri === '/health') {
    header('Content-Type: application/json');
    echo json_encode([
        'success' => true,
        'data'    => [
            'status'  => 'ok',
            'service' => 'order-service',
            'port'    => $config['app']['port'],
        ],
    ]);
    exit;
}

// ─── CORS Headers (dev) ───────────────────────────────────────
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

// ─── Routing ──────────────────────────────────────────────────
try {
    $router = new NexusFlow\Router();
    $router->dispatch($_SERVER['REQUEST_METHOD'], $_SERVER['REQUEST_URI']);
} catch (\Throwable $e) {
    http_response_code(500);
    header('Content-Type: application/json; charset=utf-8');

    $error = $config['app']['debug']
        ? $e->getMessage() . ' in ' . $e->getFile() . ':' . $e->getLine()
        : 'Erreur interne du serveur';

    echo json_encode([
        'success' => false,
        'error'   => $error,
    ]);
}
