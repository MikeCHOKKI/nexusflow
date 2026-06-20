<?php

declare(strict_types=1);

namespace NexusFlow;

use PDO;
use PDOException;

/**
 * Connexion PostgreSQL via PDO.
 *
 * Singleton léger — une seule connexion par processus.
 */
class Database
{
    private static ?PDO $instance = null;

    /**
     * Retourne l'instance PDO unique.
     */
    public static function getInstance(): PDO
    {
        if (self::$instance === null) {
            $dsn = self::buildDsn();

            try {
                self::$instance = new PDO($dsn, options: [
                    PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
                    PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
                    PDO::ATTR_EMULATE_PREPARES   => false,
                    PDO::ATTR_STRINGIFY_FETCHES  => false,
                ]);
            } catch (PDOException $e) {
                http_response_code(500);
                echo json_encode([
                    'success' => false,
                    'error'   => 'Erreur de connexion à la base de données',
                ]);
                exit(1);
            }
        }

        return self::$instance;
    }

    /**
     * Construit le DSN à partir de DATABASE_URL.
     */
    private static function buildDsn(): string
    {
        $url = $_ENV['DATABASE_URL'] ?? 'postgresql://nexusflow:nexusflow_dev@postgres:5432/nexusflow';
        $parts = parse_url($url);

        $host     = $parts['host'] ?? 'postgres';
        $port     = $parts['port'] ?? 5432;
        $dbname   = trim($parts['path'] ?? '/nexusflow', '/');
        $user     = $parts['user'] ?? 'nexusflow';
        $password = $parts['pass'] ?? 'nexusflow_dev';

        return sprintf(
            'pgsql:host=%s;port=%d;dbname=%s;user=%s;password=%s',
            $host,
            $port,
            $dbname,
            $user,
            $password
        );
    }

    /**
     * Réinitialise la connexion (utile pour les tests).
     */
    public static function reset(): void
    {
        self::$instance = null;
    }
}
