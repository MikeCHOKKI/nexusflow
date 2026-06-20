<?php

declare(strict_types=1);

namespace NexusFlow;

/**
 * Routeur HTTP simple.
 *
 * Associe les verbes et patterns d'URL aux méthodes du service.
 * Supporte les paramètres dynamiques dans l'URL : /orders/{id}
 */
class Router
{
    private OrderService $service;

    public function __construct()
    {
        $this->service = new OrderService();
    }

    /**
     * Traite une requête entrante.
     */
    public function dispatch(string $method, string $uri): void
    {
        // Nettoyer l'URI
        $uri = parse_url($uri, PHP_URL_PATH);
        $uri = '/' . trim($uri, '/');

        // POST /orders → createOrder
        if ($method === 'POST' && $uri === '/orders') {
            $this->handleCreateOrder();
            return;
        }

        // GET /orders → listOrders
        if ($method === 'GET' && $uri === '/orders') {
            $this->handleListOrders();
            return;
        }

        // GET /orders/{id} → getOrder
        if ($method === 'GET' && preg_match('#^/orders/([a-f0-9-]+)$#i', $uri, $m)) {
            $this->handleGetOrder($m[1]);
            return;
        }

        // PUT /orders/{id}/status → updateOrderStatus
        if ($method === 'PUT' && preg_match('#^/orders/([a-f0-9-]+)/status$#i', $uri, $m)) {
            $this->handleUpdateOrderStatus($m[1]);
            return;
        }

        json_response(null, 'Route non trouvée', 404);
    }

    // ─── Handlers ──────────────────────────────────────────────

    private function handleCreateOrder(): void
    {
        $body = get_json_body();

        validate_required($body, ['user_id', 'items', 'currency', 'shipping_address']);

        $result = $this->service->createOrder(
            userId:          $body['user_id'],
            items:           $body['items'],
            currency:        $body['currency'],
            shippingAddress: $body['shipping_address'],
        );

        json_response($result, null, 201);
    }

    private function handleGetOrder(string $id): void
    {
        $userId = $_GET['user_id'] ?? null;

        $order = $this->service->getOrder($id, $userId);

        if ($order === null) {
            json_response(null, 'Commande introuvable', 404);
        }

        json_response($order);
    }

    private function handleListOrders(): void
    {
        $page   = (int) ($_GET['page'] ?? 1);
        $limit  = (int) ($_GET['limit'] ?? 20);
        $userId = $_GET['user_id'] ?? null;
        $status = $_GET['status'] ?? null;

        [$orders, $total, $currentPage, $totalPages] = $this->service->listOrders(
            page:   $page,
            limit:  $limit,
            userId: $userId,
            status: $status,
        );

        json_response([
            'orders'      => $orders,
            'pagination'  => [
                'page'        => $currentPage,
                'limit'       => $limit,
                'total'       => $total,
                'total_pages' => $totalPages,
            ],
        ]);
    }

    private function handleUpdateOrderStatus(string $id): void
    {
        $body = get_json_body();
        validate_required($body, ['status']);

        $order = $this->service->updateOrderStatus($id, $body['status']);

        if ($order === null) {
            json_response(null, 'Commande introuvable', 404);
        }

        json_response($order);
    }
}
