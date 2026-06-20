<?php

declare(strict_types=1);

namespace NexusFlow;

use Predis\Client as PredisClient;

/**
 * Publie des événements métier sur Redis pour le Notification Service.
 *
 * Événements :
 * - order.created      → nouvelle commande créée
 * - order.status_changed → changement de statut
 * - payment.received   → paiement reçu (depuis le Payment Service)
 */
class RedisPubSub
{
    private PredisClient $redis;
    private string $channel;

    public function __construct()
    {
        $redisUrl = $_ENV['REDIS_URL'] ?? 'redis://redis:6379';
        $this->channel = $_ENV['NOTIFY_CHANNEL'] ?? 'nexusflow:notifications';
        $this->redis = new PredisClient($redisUrl);
    }

    /**
     * Publie un événement sur le canal Redis.
     *
     * @param string $type Type d'événement (ex: "order.created")
     * @param array  $data Données métier
     */
    public function publish(string $type, array $data): void
    {
        $payload = json_encode([
            'type'      => $type,
            'data'      => $data,
            'timestamp' => gmdate('Y-m-d\TH:i:s\Z'),
        ], JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);

        $this->redis->publish($this->channel, $payload);
    }

    /**
     * Publie un événement de création de commande.
     */
    public function orderCreated(string $orderId, string $userId, string $customerEmail, float $totalAmount, string $currency, string $customerName = 'Client'): void
    {
        $this->publish('order.created', [
            'order_id'       => $orderId,
            'user_id'        => $userId,
            'customer_email' => $customerEmail,
            'customer_name'  => $customerName,
            'total_amount'   => $totalAmount,
            'currency'       => $currency,
        ]);
    }

    /**
     * Publie un événement de changement de statut.
     */
    public function orderStatusChanged(string $orderId, string $userId, string $customerEmail, string $oldStatus, string $newStatus): void
    {
        $this->publish('order.status_changed', [
            'order_id'       => $orderId,
            'user_id'        => $userId,
            'customer_email' => $customerEmail,
            'old_status'     => $oldStatus,
            'new_status'     => $newStatus,
        ]);
    }

    /**
     * Publie un événement de paiement reçu.
     */
    public function paymentReceived(string $orderId, string $paymentId, float $amount, string $currency): void
    {
        $this->publish('payment.received', [
            'order_id'   => $orderId,
            'payment_id' => $paymentId,
            'amount'     => $amount,
            'currency'   => $currency,
        ]);
    }
}
