# SkillTracker Backend (Go + Echo + PostgreSQL)

Готовый бэкенд под твой фронт и требования JWT + creator_id.

## Запуск
```bash
docker-compose up --build
```

- БД: `localhost:5432`, db=skillstracker, user=postgres, pass=12345678
- Backend API: `http://localhost:8080/api/v1`

Создаётся пользователь:
- username: `admin`
- password: `admin123`
- role: `manager`

## Эндпоинты
- `POST /api/v1/login` -> `{ token }`
- Tasks:
  - `GET /api/v1/tasks/my`
  - `GET /api/v1/tasks/:id`
  - `POST /api/v1/tasks` (только manager)
  - `PUT /api/v1/tasks/:id`
  - `DELETE /api/v1/tasks/:id`
- Users (только manager):
  - CRUD
- Comments:
  - `POST /api/v1/comments`
  - `GET /api/v1/tasks/:task_id/comments`
  - `PUT /api/v1/comments/:id`
  - `DELETE /api/v1/comments/:id`

## Конфиг
`config/config.yaml` — DSN, порт, секрет JWT.
