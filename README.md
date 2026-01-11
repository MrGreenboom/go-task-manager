# Go Task Manager (REST API)

REST API для управления задачами.

## Стек
- Go (net/http)
- PostgreSQL
- Docker / docker-compose
- Goose migrations

## Запуск БД
```bash
docker compose up -d

## Auth
- POST /auth/register
- POST /auth/login (JWT)

## Protected
Все запросы к /tasks требуют заголовок:
Authorization: Bearer <token>
### Пример: логин
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"gairbek@test.com","password":"qwerty12"}'

## Запуск через Docker Compose
```bash
docker compose up --build -d
