# GitHub Actions Deployment — Инструкция

## Шаг 1: Подготовка VPS

```bash
# На VPS создай директорию и .env
mkdir -p /opt/focus
cd /opt/focus

# Создай .env с секретами
cat > .env << 'EOF'
DATABASE_URL=postgres://focus:YOUR_DB_PASSWORD@localhost:5432/focus?sslmode=disable
BOT_TOKEN=YOUR_TELEGRAM_BOT_TOKEN
GEMINI_KEY=YOUR_GEMINI_API_KEY
PORT=8080
EOF

# Установи Docker если нет
curl -fsSL https://get.docker.com | sh
```

## Шаг 2: Создай SSH ключ для деплоя

```bash
# На своей машине (не на VPS)
ssh-keygen -t ed25519 -C "github-deploy" -f ~/.ssh/focus_deploy -N ""

# Скопируй публичный ключ на VPS
ssh-copy-id -i ~/.ssh/focus_deploy.pub user@YOUR_VPS_IP

# Покажи приватный ключ (понадобится для GitHub)
cat ~/.ssh/focus_deploy
```

## Шаг 3: Добавь Secrets в GitHub

Иди в **Repository → Settings → Secrets and variables → Actions → New repository secret**

| Secret Name | Value |
|-------------|-------|
| `VPS_HOST` | IP адрес или домен VPS |
| `VPS_USER` | Пользователь SSH (root или другой) |
| `VPS_SSH_KEY` | Содержимое `~/.ssh/focus_deploy` (приватный ключ) |
| `VPS_PORT` | SSH порт (обычно 22) |

## Шаг 4: Создай workflow файл

Создай файл `.github/workflows/deploy.yml` в репозитории:

```yaml
name: Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}-backend

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    outputs:
      image_tag: ${{ steps.meta.outputs.tags }}

    steps:
      - uses: actions/checkout@v4

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha,prefix=
            type=raw,value=latest

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          push: true
          tags: ${{ steps.meta.outputs.tags }}

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest

    steps:
      - name: Deploy to VPS
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          port: ${{ secrets.VPS_PORT }}
          script: |
            cd /opt/focus
            
            # Pull latest image
            docker pull ghcr.io/${{ github.repository }}-backend:latest
            
            # Stop old container
            docker stop focus-api 2>/dev/null || true
            docker rm focus-api 2>/dev/null || true
            
            # Start new container
            docker run -d \
              --name focus-api \
              --restart unless-stopped \
              --env-file .env \
              -p 8080:8080 \
              ghcr.io/${{ github.repository }}-backend:latest
            
            # Cleanup old images
            docker image prune -f
```

## Шаг 5: Настрой PostgreSQL на VPS

```bash
# На VPS
docker run -d \
  --name focus-db \
  --restart unless-stopped \
  -e POSTGRES_USER=focus \
  -e POSTGRES_PASSWORD=YOUR_DB_PASSWORD \
  -e POSTGRES_DB=focus \
  -v pgdata:/var/lib/postgresql/data \
  -p 127.0.0.1:5432:5432 \
  postgres:15-alpine

# Подожди 10 сек и примени миграции
# (можно вручную или добавить в контейнер)
```

## Шаг 6: Push и проверь

```bash
git add .github/workflows/deploy.yml
git commit -m "Add GitHub Actions deploy workflow"
git push origin main
```

Иди в **Actions** таб в GitHub — увидишь запущенный workflow.

---

## После деплоя: Настрой Telegram Webhook

```bash
curl "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook?url=https://<VPS_IP>:8080/bot/webhook"
```

Или если настроишь Cloudflare Tunnel / nginx с SSL:

```bash
curl "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook?url=https://api.focus.example.com/bot/webhook"
```

---

## Структура на VPS после деплоя

```
/opt/focus/
└── .env          # Секреты (только здесь!)

Docker containers:
- focus-db        # PostgreSQL
- focus-api       # Go backend (pulled from ghcr.io)
```

---

## Troubleshooting

### Проверить логи контейнера
```bash
docker logs focus-api
docker logs focus-db
```

### Проверить что контейнеры запущены
```bash
docker ps
```

### Перезапустить вручную
```bash
docker restart focus-api
```

### Проверить health endpoint
```bash
curl http://localhost:8080/health
```
