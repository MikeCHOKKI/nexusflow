<div align="center">
  <img src="frontend/public/favicon.svg" width="80" alt="NexusFlow" />
  <h1>NexusFlow</h1>
  <p><strong>Architecture Microservices Distribuée</strong> — Polyglotte · gRPC · Redis · Docker</p>

  <a href="#-architecture"><img src="https://img.shields.io/badge/Architecture-gRPC%20%2F%20REST%20%2F%20Redis-06B6D4?style=for-the-badge&labelColor=0a0a1a" /></a>
  <a href="#-stack"><img src="https://img.shields.io/badge/Stack-Go%20%2F%20Python%20%2F%20PHP%20%2F%20React-4f46e5?style=for-the-badge&labelColor=0a0a1a" /></a>
  <a href="#-quick-start"><img src="https://img.shields.io/badge/Docker%20Compose-1%20cmd-22c55e?style=for-the-badge&labelColor=0a0a1a" /></a>
  <a href="https://github.com/MikeCHOKKI"><img src="https://img.shields.io/badge/By%20MikeCHOKKI-0a0a1a?style=for-the-badge&logo=github&logoColor=white" /></a>

  <br/><br/>
  <img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=12,20,24&height=160&section=header&text=NexusFlow&fontSize=48&fontColor=ffffff&fontAlignY=35&desc=Distributed%20Microservice%20Architecture%20%E2%80%94%20Polyglot%20%E2%80%A2%20gRPC%20%E2%80%A2%20Event-Driven&descAlignY=55&descSize=14&animation=fadeIn" />
</div>

---

## 🎯 Vision

**NexusFlow** est une architecture microservices distribuée complète, conçue pour démontrer les principes fondamentaux de la **modularité**, de la **communication inter-services** et de la **résilience** dans un environnement polyglotte.

> 6 services · 4 langages · 2 protocoles de communication · 1 commande pour lancer

Chaque service est autonome, dockerisé, et communique via des protocoles légers :
- **gRPC** pour les appels synchrones (performants, typés, bidirectionnels)
- **Redis Pub/Sub** pour les événements asynchrones (découplage, scalabilité)
- **REST** pour l'API externe (point d'entrée unique)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        🌐 Frontend                              │
│                  React + Tailwind + Recharts                    │
│                        :5173 / :3000                            │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP REST
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                     🚪 API Gateway (Go)                         │
│              :8080 — JWT · Rate Limiting · Router               │
│                                                                 │
│  ┌─────────┬──────────┬──────────┬──────────┐                  │
│  │  Auth   │ Catalog  │  Order   │ Payment  │                  │
│  └────┬────┴────┬─────┴────┬─────┴────┬─────┘                  │
└───────┼─────────┼──────────┼──────────┼─────────────────────────┘
        │ gRPC    │ gRPC     │ gRPC     │ gRPC
        ▼         ▼          ▼          ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│   👤     │ │  📦     │ │  📋     │ │  💳     │
│  User    │ │ Catalog  │ │  Order   │ │ Payment  │
│ Service  │ │ Service  │ │  Service │ │ Service  │
│  (Go)    │ │ (Python) │ │  (PHP)   │ │  (Go)    │
│ :50051   │ │ :50052   │ │ :50053   │ │ :50054   │
└────┬─────┘ └──────────┘ └────┬─────┘ └──────────┘
     │                         │
     │              ┌──────────┤
     │              │ Redis    │ (Pub/Sub)
     │              ▼          ▼
     │         ┌──────────────────┐
     │         │   🔔 Notify      │
     │         │   Service        │
     │         │   (Python)       │
     │         │   :50055         │
     │         └────────┬─────────┘
     │                  │ SMTP
     │                  ▼
     │            ┌──────────┐
     │            │  MailHog │
     │            │  :8025   │
     │            └──────────┘
     │
     ├──────────────────────────────────┐
     │            PostgreSQL :5432      │
     │  ┌────────┐ ┌────────┐ ┌──────┐ │
     │  │ users  │ │catalog │ │orders│ │
     │  │ schema │ │ schema │ │schema│ │
     │  └────────┘ └────────┘ └──────┘ │
     │  ┌────────┐                     │
     │  │payments│                     │
     │  │ schema │                     │
     │  └────────┘                     │
     └──────────────────────────────────┘
```

### Flux de communication

| Type | Technologie | Cas d'usage | Latence |
|------|------------|-------------|---------|
| **Sync** | gRPC (HTTP/2) | User→Auth, Order→Payment | <5ms |
| **Async** | Redis Pub/Sub | Order→Notify | <1ms |
| **REST** | HTTP/1.1 | Frontend→Gateway | <10ms |
| **Cache** | Redis | Produits, Sessions | <1ms |

---

## 🧩 Services

| Service | Langage | Port | Protocole | Rôle |
|---------|---------|------|-----------|------|
| **API Gateway** | Go | `:8080` | REST + gRPC | Point d'entrée, auth JWT, rate limiting, routage |
| **User Service** | Go | `:50051` | gRPC | Auth, CRUD utilisateurs, JWT |
| **Catalog Service** | Python | `:50052` | gRPC REST | Gestion produits, catégories, stock |
| **Order Service** | PHP | `:50053` | HTTP | Commandes, workflow, événements Redis |
| **Payment Service** | Go | `:50054` | gRPC | Stub paiement (Wave, OM, MTN, Visa) |
| **Notification** | Python | `:50055` | Redis Sub | Emails (MailHog), notifications |
| **Frontend** | React/TS | `:3000` | HTTP | Dashboard admin, visualisation temps réel |

---

## 🛠️ Stack

<table>
<tr>
<td width="33%" valign="top"><strong>Backend</strong><br/>
Go 1.22 · Python 3.12 · PHP 8.2<br/><br/>
<strong>Frameworks</strong><br/>
FastAPI · Gorilla Mux · pgx/v5
</td>
<td width="33%" valign="top"><strong>Frontend</strong><br/>
React 19 · TypeScript 5.6 · Vite 6<br/>
Tailwind CSS v4 · Motion · Recharts<br/>
Phosphor Icons
</td>
<td width="33%" valign="top"><strong>Infra</strong><br/>
Docker · Docker Compose<br/>
PostgreSQL 16 · Redis 7<br/>
MailHog · gRPC · Protobuf
</td>
</tr>
</table>

### Dépendances clés

| Service | Dépendances | Gestion |
|---------|------------|---------|
| Go (gateway, user, payment) | `go-redis`, `pgx`, `jwt-go`, `grpc` | `go.mod` / `go.sum` |
| Python (catalog, notify) | `fastapi`, `asyncpg`, `grpcio`, `aiosmtplib` | `requirements.txt` |
| PHP (order) | `grpc/grpc`, `predis/predis` | `composer.json` |
| Frontend | `react`, `motion`, `recharts`, `phosphor-icons` | `package.json` |

---

## 🚀 Quick Start

### Prérequis
- Docker & Docker Compose v2
- Git
- Make (optionnel)

### Lancer tout le projet

```bash
# Cloner
git clone https://github.com/MikeCHOKKI/nexusflow.git
cd nexusflow

# Démarrer tous les services
make up
# ou : docker compose up --build -d

# Vérifier l'état
docker compose ps
make ps
```

### Accès

| Service | URL |
|---------|-----|
| **Dashboard Frontend** | [http://localhost:3000](http://localhost:3000) |
| **API Gateway** | [http://localhost:8080](http://localhost:8080) |
| **Health Check** | [http://localhost:8080/health](http://localhost:8080/health) |
| **MailHog (emails)** | [http://localhost:8025](http://localhost:8025) |
| **PostgreSQL** | `localhost:5432` (nexusflow/nexusflow_dev) |
| **Redis** | `localhost:6379` |

### Arrêter

```bash
make down
# ou : docker compose down -v  # avec suppression volumes
```

### API Endpoints

```bash
# Auth
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@nexusflow.local", "password": "password123"}'

# Produits (avec token)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/products

# Commandes
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/orders

# Santé des services
curl http://localhost:8080/health
```

---

## 📊 Dashboard

Le dashboard NexusFlow offre une interface professionnelle pour visualiser et interagir avec l'ensemble des microservices :

| Page | Fonctionnalités |
|------|----------------|
| **Dashboard** | Stats clés (produits, commandes, revenus, users), graphique d'évolution, statut des services en temps réel |
| **Produits** | CRUD complet, recherche, filtres, pagination |
| **Commandes** | Workflow de statuts (pending → delivered), détails, mise à jour |
| **Paiements** | Transactions, méthodes multiples (Wave, Orange Money, MTN, Visa), remboursement |
| **Utilisateurs** | Gestion des rôles, activation/désactivation |
| **Paramètres** | Configuration, API URL, notifications |

---

## 📁 Structure du projet

```
nexusflow/
├── api-gateway/           # Go — Point d'entrée REST → gRPC
├── service-user/          # Go — Auth & utilisateurs
├── service-catalog/       # Python/FastAPI — Produits & catalogue
├── service-order/         # PHP — Commandes & workflow
├── service-payment/       # Go — Paiements (stub)
├── service-notify/        # Python — Notifications emails
├── frontend/              # React + Tailwind — Dashboard
├── protos/                # Définitions Protobuf partagées
├── db/                    # Migrations SQL + seeds
│   ├── migrations/
│   └── seeds/
├── docs/                  # Documentation & ADRs
│   ├── architecture.md
│   └── adr/
├── docker-compose.yml     # Orchestration complète
├── Makefile               # Commandes automatisées
├── .env                   # Variables d'environnement
└── README.md              # Ce fichier
```

---

## 🔄 Patterns implémentés

| Pattern | Implémentation |
|---------|---------------|
| **API Gateway** | Point d'entrée unique, routage, auth déléguée |
| **Service Registry** | Adresses statiques via docker-compose + env vars |
| **Saga (future)** | Workflow distribué via Redis Pub/Sub |
| **CQRS** | Lectures/écritures séparées (Redis cache + PG writes) |
| **Event-Driven** | Notifications asynchrones via Redis Pub/Sub |
| **Polyglot Persistence** | PostgreSQL (schémas isolés) + Redis (cache/messaging) |
| **Health Check** | Endpoint agrégé + Docker healthchecks |
| **Graceful Shutdown** | Capture SIGINT/SIGTERM dans tous les services |

---

## 🧪 Tests

```bash
# Lancer tous les tests
make test

# Vérifier le lint
make lint

# Construire les protos
make proto
```

---

## 🛣️ Roadmap

- [ ] **Kubernetes manifests** — Déploiement sur cluster K8s
- [ ] **Service Mesh** — Istio/Linkerd pour observabilité
- [ ] **Prometheus + Grafana** — Métriques et alerting
- [ ] **Saga Pattern** — Transactions distribuées complètes
- [ ] **CI/CD** — GitHub Actions → build → test → deploy
- [ ] **Rate Limiting avancé** — Par user, par endpoint
- [ ] **Circuit Breaker** — Résilience inter-services

---

## 📄 Licence

MIT — Fait avec OpenCode par [Mike E.S. Chokki](https://github.com/MikeCHOKKI)

<p align="center">
  <sub>Construisons des solutions qui comptent pour l'Afrique 🌍</sub>
</p>
