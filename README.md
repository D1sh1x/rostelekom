## SkillsTracker — техническая документация

### Описание
SkillsTracker — REST API сервис для управления пользователями, задачами и комментариями. Проект написан на Go c использованием Echo, PostgreSQL, JWT-аутентификации, разделён на слои handler/service/repository и поддерживает автоматическую документацию Swagger.

### Технологии
- Go (Echo v4, sqlx, zerolog, viper)
- PostgreSQL (lib/pq)
- JWT (github.com/golang-jwt/jwt/v5)
- Swagger (swaggo/swag, echo-swagger)
- Docker, Docker Compose

### Структура проекта (основное)
- `cmd/main.go` — точка входа
- `internal/config` — конфигурация (viper)
- `internal/transport` — HTTP-сервер и роутер (Echo)
- `internal/handler` — HTTP-обработчики
- `internal/service` — бизнес-логика
- `internal/repository` — интерфейсы репозиториев
- `internal/storage/postgres` — реализация репозитория на PostgreSQL (sqlx)
- `internal/middleware` — JWT-аутентификация и проверка ролей
- `internal/utils/jwt` — генерация/валидация JWT
- `migrations` — SQL-миграции (используются init-скриптами postgres)
- `config/config.yaml` — конфиг приложения
- `docs` — Swagger-спека (генерируется swag’ом)

### Конфигурация
Конфигурация читается из `./config/config.yaml` (можно переопределить переменной окружения `CONFIG_PATH`).

Пример `config/config.yaml`:
```yaml
env: dev
jwt_secret: andrey
http_server:
  port: :8080
  write_timeout: 5s
  read_timeout: 5s
  idle_timeout: 60s
database:
  dsn: "postgres://postgres:12345678@db:5432/skillstracker?sslmode=disable"
  host: db
  port: 5432
  user: postgres
  password: 12345678
  name: skillstracker
  sslmode: disable
```

### Запуск в Docker
1) Собрать и запустить:
```bash
docker-compose up -d
```
2) База данных: первый запуск применит SQL из `migrations/` через `/docker-entrypoint-initdb.d` (исполняется только при создании нового volume). Убедитесь, что в файлах миграций отсутствуют секции Down.
3) Приложение доступно на `http://localhost:8080`.

Порты:
- API — `:8080`
- Postgres — `:5432`

### Локальный запуск (без Docker)
```bash
# из корня проекта
export CONFIG_PATH=./config/config.yaml
go run ./cmd/main.go
```
Требуется локальный PostgreSQL и корректный DSN в конфиге.

### Swagger (OpenAPI)
- UI: `http://localhost:8080/swagger/index.html`
- JSON: `http://localhost:8080/swagger/doc.json`

Генерация спеков (локально):
```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.6
# сгенерировать файлы в папке docs
swag init -g cmd/main.go -o docs
```
Важно:
- В коде должен быть импорт `_ "SkillsTracker/docs"` и роут: `router.GET("/swagger/*", echoSwagger.WrapHandler)`.
- В `router` перед запуском сервера можно указать базовый путь: `docs.SwaggerInfo.BasePath = "/api/v1"`.
- Если собираете в Docker, лучше закоммитить папку `docs/` или запускать `swag init` в стадии `builder`.
- Если видите ошибку про `LeftDelim/RightDelim` — версии swag CLI и библиотеки не совпадают. Выровняйте их и пересгенерируйте `docs/`.

### Аутентификация и роли
- Формат: `Authorization: Bearer <JWT>`
- Токен содержит: `user_id` (int), `username`, `role`; срок жизни 24 часа.
- Middleware:
  - `AuthRequired` — проверяет токен, кладёт в контекст `user_id`, `username`, `role`.
  - `RequireRole("manager")` — ограничение по ролям.

### Основные эндпоинты (сокращённо)
Базовый путь API — `/api/v1`.

Публичные:
- `POST /api/v1/register` — регистрация пользователя
- `POST /api/v1/login` — логин, возвращает JWT

Требуют авторизации (Bearer JWT):
- Users (только `manager`):
  - `GET /api/v1/users`
  - `GET /api/v1/users/{id}`
  - `PUT /api/v1/users/{id}`
  - `DELETE /api/v1/users/{id}`
- Tasks:
  - `POST /api/v1/tasks`
  - `GET /api/v1/tasks/{id}`
  - `GET /api/v1/tasks?employee_id=<id>`
  - `PUT /api/v1/tasks/{id}`
  - `DELETE /api/v1/tasks/{id}`
- Comments:
  - `POST /api/v1/comments`
  - `GET /api/v1/comments/{id}`
  - `GET /api/v1/comments?task_id=<id>`
  - `PUT /api/v1/comments/{id}`
  - `DELETE /api/v1/comments/{id}`

### Быстрый старт (curl)
```bash
# регистрация
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"andrey","password":"123456","role":"manager","name":"Андрей"}'

# логин
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"andrey","password":"123456"}' | jq -r .token)

# список пользователей (manager)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/users
```

### Миграции
- В Docker Compose применяются автоматически при первом создании БД (тома).
- Для последующих изменений схемы рекомендуется использовать инструмент миграций (goose/migrate) как отдельный сервис или вручную.

### Тестирование (рекомендации)
Типы тестов:
- Unit: тестировать бизнес-логику (service) с моками репозитория; JWT-утилиты; middleware (через `httptest`).
- Интеграционные: поднятая БД (docker) + реальные запросы к серверу (Echo + `httptest.NewServer`).
- E2E: коллекция Postman + Newman.

Примеры сценариев:
- JWT: генерация/валидация, истёкший токен, неверный секрет.
- Middleware: корректное извлечение `user_id/username/role`, ответы 401/403.
- UserService: логин (успех/ошибка), создание (дубликат пользователя), получение списка пользователей.
- Task/Comment: happy-path и ошибки (невалидные параметры, не найдено, нет прав).

Запуск тестов:
```bash
go test ./...
```

### Производственная эксплуатация
- Храните `jwt_secret` и параметры БД в переменных окружения/секретах.
- Логирование: zerolog с уровнями по окружению.
- Настройте ротацию логов и мониторинг (Prometheus/Grafana) при необходимости.

### Лицензия
MIT (или добавьте свою).
