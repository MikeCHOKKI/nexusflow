<?php

declare(strict_types=1);

/**
 * NexusFlow — Fonctions utilitaires pour les réponses HTTP JSON.
 */

/**
 * Envoie une réponse JSON uniforme et termine l'exécution.
 *
 * @param mixed       $data   Données à retourner
 * @param string|null $error  Message d'erreur (optionnel)
 * @param int         $status Code HTTP
 */
function json_response(mixed $data = null, ?string $error = null, int $status = 200): void
{
    http_response_code($status);
    header('Content-Type: application/json; charset=utf-8');

    $body = [
        'success' => $error === null,
        'data'    => $data,
    ];

    if ($error !== null) {
        $body['error'] = $error;
    }

    echo json_encode($body, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
    exit;
}

/**
 * Extrait et décode le body JSON d'une requête.
 *
 * @return array Décodé en tableau associatif
 */
function get_json_body(): array
{
    $raw = file_get_contents('php://input');
    if ($raw === false || $raw === '') {
        return [];
    }

    $decoded = json_decode($raw, true);
    if (json_last_error() !== JSON_ERROR_NONE) {
        json_response(null, 'Corps de requête JSON invalide', 400);
    }

    return $decoded ?? [];
}

/**
 * Valide qu'un ensemble de champs est présent dans un tableau.
 *
 * @param array $data    Données à valider
 * @param array $fields  Liste des champs requis
 */
function validate_required(array $data, array $fields): void
{
    $missing = [];
    foreach ($fields as $field) {
        if (!isset($data[$field]) || (is_string($data[$field]) && trim($data[$field]) === '')) {
            $missing[] = $field;
        }
    }

    if ($missing !== []) {
        json_response(null, 'Champs requis manquants : ' . implode(', ', $missing), 400);
    }
}

/**
 * Valide qu'un statut de commande est valide.
 */
function validate_order_status(string $status): bool
{
    $valid = ['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded'];
    return in_array($status, $valid, true);
}
