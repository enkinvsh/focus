# Focus

**Focus** is an AI-powered task manager designed for **Telegram Mini Apps**. It uses Google Gemini to parse voice commands into structured tasks, with a Go backend for cloud persistence and real-time sync.

![License](https://img.shields.io/badge/license-MIT-green) ![Version](https://img.shields.io/badge/version-0.1.0-blue) ![Platform](https://img.shields.io/badge/platform-Telegram-2CA5E0)

**Live Demo:** [https://enkinvsh.github.io/focus/](https://enkinvsh.github.io/focus/)

---

## Key Features

- **AI Voice Input:** Record voice commands with visual progress ring and silence detection. Auto-stops after 2s silence or 8s max.
- **Smart Parsing:** Gemini 2.0 Flash extracts task title, priority, and category from natural language.
- **Cloud Sync:** PostgreSQL backend with Telegram user authentication. Tasks sync across devices.
- **Breathing Exercise:** 1-minute guided breathing (4-4-4 pattern). Tap "Focus" title to start.
- **Premium UI/UX:**
  - Glassmorphism design
  - Swipe gestures for tab navigation
  - Haptic feedback
  - 4 OLED-friendly dark themes
  - Animated start screen
- **Localization:** English and Russian with auto-detection.

---

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Telegram App   │────▶│  GitHub Pages    │────▶│ Cloudflare      │
│  (User)         │     │  (index.html)    │     │ Worker          │
└─────────────────┘     └──────────────────┘     │ (Gemini Proxy)  │
                               │                 └────────┬────────┘
                               │                          │
                               ▼                          ▼
                        ┌──────────────────┐     ┌─────────────────┐
                        │  Go Backend      │     │ Google Gemini   │
                        │  (api.meybz.asia)│     │ API             │
                        │  + PostgreSQL    │     └─────────────────┘
                        │  + Caddy         │
                        └──────────────────┘
```

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Single `index.html`, Vanilla JS, Tailwind CSS |
| Backend | Go 1.23, Gin, PostgreSQL 16 |
| AI Proxy | Cloudflare Worker |
| AI Model | Gemini 2.0 Flash |
| Reverse Proxy | Caddy (auto HTTPS) |
| Deployment | GitHub Actions → Docker → VPS |
| Domain | Cloudflare DNS (api.meybz.asia) |

---

## Installation & Setup

### 1. Frontend (GitHub Pages)

1. Fork this repository
2. Enable **GitHub Pages** (Source: `main` branch)
3. Update `API_URL` in `index.html` to your backend URL

### 2. Backend (VPS)

```bash
# Clone and configure
git clone https://github.com/enkinvsh/focus.git
cd focus/backend

# Create .env file
cat > .env << EOF
DATABASE_URL=postgres://focus:password@focus-db:5432/focus?sslmode=disable
GEMINI_PROXY_URL=https://your-worker.workers.dev
EOF

# Start with Docker Compose
docker compose up -d
```

### 3. Cloudflare Worker (Gemini Proxy)

1. Create Worker at [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. Paste contents of `worker.js`
3. Add environment variable: `GEMINI_KEY`
4. Deploy and copy Worker URL

### 4. GitHub Actions (CI/CD)

Required secrets in repository settings:

| Secret | Description |
|--------|-------------|
| `VPS_HOST` | Server IP address |
| `VPS_USER` | SSH username |
| `VPS_SSH_KEY` | Private SSH key |
| `VPS_PORT` | SSH port |
| `CF_ORIGIN_CERT` | Cloudflare origin certificate |
| `CF_ORIGIN_KEY` | Cloudflare origin key |

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/tasks?limit=50&offset=0` | Get tasks (paginated, max 200) |
| `POST` | `/tasks` | Create task |
| `PUT` | `/tasks/:id` | Update task |
| `DELETE` | `/tasks/:id` | Delete task |
| `POST` | `/transcribe` | Transcribe audio via Gemini |

All endpoints require `X-Telegram-Init-Data` header for authentication.

---

## File Structure

```
focus/
├── index.html                 # Frontend SPA
├── worker.js                  # Cloudflare Worker (Gemini proxy)
├── backend/
│   ├── cmd/server/main.go     # Entry point
│   ├── internal/
│   │   ├── api/handlers.go    # HTTP handlers
│   │   ├── db/postgres.go     # Database layer
│   │   └── services/ai.go     # Gemini integration
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── Caddyfile
├── .github/workflows/
│   └── deploy.yml             # CI/CD pipeline
└── AUDIT_REPORT.md            # Security audit
```

---

## Changelog

### v0.1.0 (Production Ready)
- **Backend:** Go API with PostgreSQL, graceful shutdown, pagination
- **Voice UX:** Progress ring with 8s countdown, silence detection (auto-stop after 2s)
- **AI:** Switched to Gemini 2.0 Flash, direct Worker calls from frontend
- **Security:** CORS whitelist, XSS prevention, input validation, error sanitization
- **DevOps:** GitHub Actions with health polling, rollback on failure, masked secrets
- **Tasks:** Added `createdAt` field, transcript display in modal

### v0.0.3 (Focus Edition)
- Breathing exercise feature (1-min, 4-4-4 pattern)
- Premium animated start screen
- Bot commands /help and /about

### v0.0.2 (Cloud Edition)
- Telegram CloudStorage integration
- Cloudflare Worker proxy

### v0.0.1 (Initial Release)
- Core AI task extraction
- Localization (EN/RU)

---

## Security

- **Authentication:** Telegram `initData` validation with HMAC-SHA256
- **CORS:** Whitelist only (`enkinvsh.github.io`, `web.telegram.org`, `t.me`)
- **XSS:** All user content rendered via `textContent`
- **Input Validation:** 5MB max audio, proper ID parsing
- **Error Handling:** Internal errors logged, generic messages to client
- **Secrets:** API keys in environment variables, never exposed to client

---

## Privacy

- Task data stored in PostgreSQL on private VPS
- Telegram user ID used for authentication (no passwords)
- Gemini API requests proxied (API key hidden)
- No analytics or tracking

---

## License

MIT

---

**Built with AI & Human Collaboration**
