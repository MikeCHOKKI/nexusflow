<?php

declare(strict_types=1);

namespace NexusFlow;

use PDO;

/**
 * Service métier des commandes.
 *
 * Gère le cycle de vie : création, consultation, listing, mise à jour de statut.
 * Calcule les totaux et publie les événements Redis.
 */
class OrderService
{
    private PDO $db;
    private RedisPubSub $pubSub;

    public function __construct()
    {
        $this->db     = Database::getInstance();
        $this->pubSub = new RedisPubSub();
    }

    // ─── CREATE ────────────────────────────────────────────────

    /**
     * Crée une commande avec ses items et calcule les totaux.
     *
     * @param string $userId           UUID du client
     * @param string $customerEmail    Email du client (pour notification)
     * @param string $customerName     Nom du client (pour email)
     * @param array  $items            [{product_id, product_name, quantity, unit_price}, ...]
     * @param string $currency         Code devise (XOF, EUR, USD...)
     * @param string $shippingAddress  Adresse de livraison
     *
     * @return array Données de la commande créée
     */
    public function createOrder(string $userId, string $customerEmail, string $customerName, array $items, string $currency, string $shippingAddress): array
    {
        // Valider les items
        if (empty($items)) {
            json_response(null, 'La commande doit contenir au moins un article', 400);
        }

        // Calculer les totaux
        $totalAmount = 0.0;
        $processedItems = [];

        foreach ($items as $item) {
            $quantity   = (int) ($item['quantity'] ?? 0);
            $unitPrice  = (float) ($item['unit_price'] ?? 0);
            $totalPrice = $quantity * $unitPrice;
            $totalAmount += $totalPrice;

            $processedItems[] = [
                'product_id'   => $item['product_id'],
                'product_name' => $item['product_name'] ?? 'Article',
                'quantity'     => $quantity,
                'unit_price'   => $unitPrice,
                'total_price'  => $totalPrice,
            ];
        }

        if ($totalAmount <= 0) {
            json_response(null, 'Le montant total doit être supérieur à zéro', 400);
        }

        $this->db->beginTransaction();

        try {
            // Insérer la commande
            $stmt = $this->db->prepare(
                "INSERT INTO orders.orders (user_id, total_amount, currency, shipping_address, customer_email, customer_name)
                 VALUES (:user_id, :total_amount, :currency, :shipping_address, :customer_email, :customer_name)
                 RETURNING id, created_at, updated_at"
            );

            $stmt->execute([
                'user_id'          => $userId,
                'total_amount'     => $totalAmount,
                'currency'         => $currency,
                'shipping_address' => $shippingAddress,
                'customer_email'   => $customerEmail,
                'customer_name'    => $customerName,
            ]);

            $order = $stmt->fetch();
            $orderId = $order['id'];

            // Insérer les items
            $itemStmt = $this->db->prepare(
                "INSERT INTO orders.order_items (order_id, product_id, product_name, quantity, unit_price, total_price)
                 VALUES (:order_id, :product_id, :product_name, :quantity, :unit_price, :total_price)"
            );

            $orderItems = [];
            foreach ($processedItems as $pi) {
                $itemStmt->execute([
                    'order_id'    => $orderId,
                    'product_id'  => $pi['product_id'],
                    'product_name' => $pi['product_name'],
                    'quantity'    => $pi['quantity'],
                    'unit_price'  => $pi['unit_price'],
                    'total_price' => $pi['total_price'],
                ]);
                $orderItems[] = $pi;
            }

            $this->db->commit();

            // Publier événement avec les articles
            $this->pubSub->orderCreated($orderId, $userId, $customerEmail, $totalAmount, $currency, $orderItems, $customerName);

            return [
                'id'               => $orderId,
                'user_id'          => $userId,
                'items'            => $orderItems,
                'total_amount'     => (float) $totalAmount,
                'currency'         => $currency,
                'status'           => 'pending',
                'shipping_address' => $shippingAddress,
                'payment_id'       => null,
                'created_at'       => $order['created_at'],
                'updated_at'       => $order['updated_at'],
            ];
        } catch (\Throwable $e) {
            $this->db->rollBack();
            throw $e;
        }
    }

    // ─── GET ───────────────────────────────────────────────────

    /**
     * Récupère une commande avec ses items.
     *
     * @param string      $id     UUID de la commande
     * @param string|null $userId Filtre optionnel par utilisateur
     *
     * @return array|null Données de la commande ou null si introuvable
     */
    public function getOrder(string $id, ?string $userId = null): ?array
    {
        $sql = "SELECT id, user_id, total_amount, currency, status::text, shipping_address,
                       payment_id, created_at, updated_at
                FROM orders.orders
                WHERE id = :id";

        $params = ['id' => $id];

        if ($userId !== null) {
            $sql .= " AND user_id = :user_id";
            $params['user_id'] = $userId;
        }

        $stmt = $this->db->prepare($sql);
        $stmt->execute($params);
        $order = $stmt->fetch();

        if (!$order) {
            return null;
        }

        // Récupérer les items
        $itemStmt = $this->db->prepare(
            "SELECT product_id, product_name, quantity,
                    CAST(unit_price AS DOUBLE PRECISION) AS unit_price,
                    CAST(total_price AS DOUBLE PRECISION) AS total_price
             FROM orders.order_items
             WHERE order_id = :order_id
             ORDER BY id"
        );
        $itemStmt->execute(['order_id' => $id]);
        $items = $itemStmt->fetchAll();

        return [
            'id'               => $order['id'],
            'user_id'          => $order['user_id'],
            'items'            => $items,
            'total_amount'     => (float) $order['total_amount'],
            'currency'         => $order['currency'],
            'status'           => $order['status'],
            'shipping_address' => $order['shipping_address'],
            'payment_id'       => $order['payment_id'],
            'customer_email'   => $order['customer_email'] ?? '',
            'customer_name'    => $order['customer_name'] ?? 'Client',
            'created_at'       => $order['created_at'],
            'updated_at'       => $order['updated_at'],
        ];
    }

    // ─── LIST ──────────────────────────────────────────────────

    /**
     * Liste les commandes avec pagination et filtres optionnels.
     *
     * @param int         $page   Numéro de page (1-based)
     * @param int         $limit  Éléments par page (max 100)
     * @param string|null $userId Filtre par utilisateur
     * @param string|null $status Filtre par statut
     *
     * @return array [$orders, $total, $page, $totalPages]
     */
    public function listOrders(int $page = 1, int $limit = 20, ?string $userId = null, ?string $status = null): array
    {
        $page  = max(1, $page);
        $limit = min(100, max(1, $limit));
        $offset = ($page - 1) * $limit;

        // Construction des conditions
        $conditions = [];
        $params = [];

        if ($userId !== null) {
            $conditions[] = "o.user_id = :user_id";
            $params['user_id'] = $userId;
        }

        if ($status !== null && validate_order_status($status)) {
            $conditions[] = "o.status = :status::orders.order_status";
            $params['status'] = $status;
        }

        $where = $conditions !== [] ? 'WHERE ' . implode(' AND ', $conditions) : '';

        // Total pour pagination
        $countStmt = $this->db->prepare("SELECT COUNT(*) FROM orders.orders o $where");
        $countStmt->execute($params);
        $total = (int) $countStmt->fetchColumn();
        $totalPages = (int) ceil($total / $limit);

        // Requête principale
        $sql = "SELECT o.id, o.user_id,
                       CAST(o.total_amount AS DOUBLE PRECISION) AS total_amount,
                       o.currency, o.status::text, o.shipping_address,
                       o.payment_id, o.customer_email, o.customer_name,
                       o.created_at, o.updated_at
                FROM orders.orders o
                $where
                ORDER BY o.created_at DESC
                LIMIT :limit OFFSET :offset";

        // Les params LIMIT/OFFSET doivent être passés comme entiers
        $stmt = $this->db->prepare($sql);
        foreach ($params as $key => $val) {
            $stmt->bindValue(":$key", $val);
        }
        $stmt->bindValue(':limit', $limit, PDO::PARAM_INT);
        $stmt->bindValue(':offset', $offset, PDO::PARAM_INT);
        $stmt->execute();

        $orders = $stmt->fetchAll();

        // Charger les items pour chaque commande
        if ($orders !== []) {
            $orderIds = array_column($orders, 'id');
            $placeholders = implode(',', array_fill(0, count($orderIds), '?'));

            $itemStmt = $this->db->prepare(
                "SELECT order_id, product_id, product_name, quantity,
                        CAST(unit_price AS DOUBLE PRECISION) AS unit_price,
                        CAST(total_price AS DOUBLE PRECISION) AS total_price
                 FROM orders.order_items
                 WHERE order_id IN ($placeholders)
                 ORDER BY id"
            );
            $itemStmt->execute($orderIds);
            $allItems = $itemStmt->fetchAll();

            // Grouper les items par order_id
            $itemsByOrder = [];
            foreach ($allItems as $item) {
                $itemsByOrder[$item['order_id']][] = $item;
            }

            foreach ($orders as &$order) {
                $order['items'] = $itemsByOrder[$order['id']] ?? [];
                $order['total_amount'] = (float) $order['total_amount'];
            }
            unset($order);
        }

        return [$orders, $total, $page, $totalPages];
    }

    // ─── UPDATE STATUS ─────────────────────────────────────────

    /**
     * Met à jour le statut d'une commande.
     *
     * Publie un événement order.status_changed sur Redis.
     *
     * @param string $id     UUID de la commande
     * @param string $status Nouveau statut
     *
     * @return array|null Données mises à jour ou null si introuvable
     */
    public function updateOrderStatus(string $id, string $status): ?array
    {
        if (!validate_order_status($status)) {
            json_response(null, "Statut invalide : $status", 400);
        }

        // Récupérer le statut actuel + email client
        $stmt = $this->db->prepare(
            "SELECT id, user_id, status::text, customer_email, customer_name FROM orders.orders WHERE id = :id"
        );
        $stmt->execute(['id' => $id]);
        $current = $stmt->fetch();

        if (!$current) {
            return null;
        }

        $oldStatus      = $current['status'];
        $userId         = $current['user_id'];
        $customerEmail  = $current['customer_email'] ?? '';
        $customerName   = $current['customer_name'] ?? 'Client';

        // Mettre à jour
        $updateStmt = $this->db->prepare(
            "UPDATE orders.orders
             SET status = :status::orders.order_status, updated_at = NOW()
             WHERE id = :id
             RETURNING id, user_id,
                       CAST(total_amount AS DOUBLE PRECISION) AS total_amount,
                       currency, status::text, shipping_address,
                       payment_id, customer_email, customer_name,
                       created_at, updated_at"
        );
        $updateStmt->execute([
            'id'     => $id,
            'status' => $status,
        ]);

        $order = $updateStmt->fetch();

        // Récupérer les items
        $itemStmt = $this->db->prepare(
            "SELECT product_id, product_name, quantity,
                    CAST(unit_price AS DOUBLE PRECISION) AS unit_price,
                    CAST(total_price AS DOUBLE PRECISION) AS total_price
             FROM orders.order_items
             WHERE order_id = :order_id
             ORDER BY id"
        );
        $itemStmt->execute(['order_id' => $id]);
        $order['items'] = $itemStmt->fetchAll();

        // Publier événement avec les informations client
        $this->pubSub->orderStatusChanged($id, $userId, $customerEmail, $customerName, $oldStatus, $status);

        return $order;
    }

    /**
     * Retourne l'instance RedisPubSub (pour injection par le Payment Service).
     */
    public function getPubSub(): RedisPubSub
    {
        return $this->pubSub;
    }
}
