# Go Task Manager (REST API)

Backend-сервис для управления задачами: регистрация/логин, JWT-защита, CRUD задач.
Проект сделан как портфолио: чистая структура (handler/service/repository), PostgreSQL, миграции, Docker Compose.

## Стек
- Go (net/http)
- PostgreSQL
- JWT (golang-jwt/jwt)
- bcrypt (x/crypto/bcrypt)
- Goose migrations
- Docker / docker-compose

## Архитектура
- `internal/handler` — HTTP handlers + middleware (JWT)
- `internal/service` — бизнес-логика и валидация
- `internal/repository` — доступ к PostgreSQL
- `migrations` — миграции БД (goose)
- `cmd/app` — точка входа приложения

---

## Быстрый старт (Docker Compose)

### 1) Запуск
```bash
docker compose up --build -d

## Testing
```bash
go test ./internal/service
