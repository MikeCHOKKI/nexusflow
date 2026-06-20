# ADR 001 : Architecture Microservices Distribuée

**Statut** : Accepté  
**Date** : 2026-06-20  
**Auteur** : Mike E.S. Chokki

---

## Contexte

Nous souhaitons démontrer une architecture microservices complète, modulaire et scalable pour renforcer un portfolio technique GitHub. Les objectifs sont :

1. Montrer la **maîtrise multi-langages** (Go, Python, PHP, TypeScript)
2. Illustrer les **protocoles de communication légers** (gRPC, Redis Pub/Sub)
3. Démontrer la **gestion des dépendances** (go.mod, requirements.txt, composer.json, package.json)
4. Produire un **projet prêt à l'emploi** (docker compose up)

## Décision

### Architecture

- **API Gateway** (Go) : point d'entrée unique REST, authentification JWT, rate limiting Redis
- **6 microservices indépendants** : User, Catalog, Order, Payment, Notify, Frontend
- **Communication synchrone** : gRPC (HTTP/2, typé, performant)
- **Communication asynchrone** : Redis Pub/Sub (événements, découplage)
- **Stockage** : PostgreSQL (schémas isolés par service) + Redis (cache + messaging)

### Choix technologiques

| Technologie | Justification |
|------------|---------------|
| **Go** (Gateway, User, Payment) | Performance, concurrence native, typage fort, excellent support gRPC |
| **Python/FastAPI** (Catalog, Notify) | Productivité, librairies data, async natif, écosystème ML |
| **PHP** (Order) | Maturité, écosystème commerce, workflow complexes |
| **React + Tailwind** (Frontend) | Écosystème riche, DX, design system flexible |
| **gRPC** | Performant, typé, streaming, compatible multi-langages |
| **Redis Pub/Sub** | Léger, <1ms latence, persistance optionnelle |
| **PostgreSQL** | Schemas isolés, maturité, JSON support, extensions |

### Alternatives considérées

1. **RabbitMQ** → plus lourd à setup, Redis suffit pour le volume de démo
2. **Kubernetes** → overkill pour une démo GitHub, Docker Compose suffit
3. **Full Go** → moins parlant sur un portfolio (moins polyglotte)
4. **REST pur** → moins performant que gRPC, pas de typage fort

## Conséquences

### Positives

- Démonstration complète d'une architecture microservices réaliste
- Chaque service est autonome, dockerisé, testable individuellement
- Le projet sert de référence pour des architectures réelles
- Facile à étendre (ajouter un service = nouveau dossier + proto)

### Négatives

- Complexité de debugging inter-services
- Nécessite de connaître 4 langages pour maintenir
- Overhead Docker : ~2GB d'images
- Pas de déploiement cloud-ready (nécessite adaptation K8s)
