# ─── Order Service ───────────────────────────────────────────
FROM php:8.2-cli

RUN apt-get update && apt-get install -y --no-install-recommends \
        libpq-dev unzip git && \
    docker-php-ext-install pdo_pgsql && \
    rm -rf /var/lib/apt/lists/*

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer
WORKDIR /app

COPY service-order/composer.json service-order/composer.lock ./
RUN composer install --no-autoloader --no-scripts --no-dev

COPY service-order/ ./
RUN composer dump-autoload --optimize

EXPOSE 50053
CMD ["php", "-S", "0.0.0.0:50053", "server.php"]
