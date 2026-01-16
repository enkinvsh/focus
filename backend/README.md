# Focus Backend

Go + PostgreSQL backend for Focus Telegram Mini App.

## Development

```bash
cp .env.example .env
docker-compose up -d postgres
go run ./cmd/server
```

## Production

```bash
docker-compose up -d
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /api/v1/tasks | Get tasks (query: type, completed) |
| POST | /api/v1/tasks | Create task |
| PATCH | /api/v1/tasks/:id | Update task |
| DELETE | /api/v1/tasks/:id | Delete task |
| GET | /api/v1/user/preferences | Get user preferences |
| PATCH | /api/v1/user/preferences | Update user preferences |

## Authentication

All `/api/v1/*` endpoints require Telegram Mini App `initData` in Authorization header:

```
Authorization: tma <initData>
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| DATABASE_URL | PostgreSQL connection string |
| BOT_TOKEN | Telegram Bot Token |
| GEMINI_KEY | Google Gemini API Key |
| PORT | Server port (default: 8080) |
