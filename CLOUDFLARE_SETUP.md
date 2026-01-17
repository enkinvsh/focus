# Cloudflare Security Setup Guide

## Архитектура

```
┌─────────────────────────────────────────────────────────────────┐
│                          Пользователи                           │
└───────────────────────────────┬─────────────────────────────────┘
                                │
        ┌───────────────────────┴───────────────────────┐
        ▼                                               ▼
┌───────────────────┐                   ┌───────────────────────────┐
│   GitHub Pages    │                   │       Cloudflare          │
│   (Frontend)      │                   │   api.oen.meybz.asia      │
│                   │                   │                           │
│ enkinvsh.github   │                   │ • DDoS Protection         │
│ .io/focus         │                   │ • WAF Rules               │
│                   │                   │ • Rate Limiting           │
│ Статика, CDN      │                   │ • Bot Protection          │
│ Fastly            │                   │ • Origin IP скрыт         │
└───────────────────┘                   └─────────────┬─────────────┘
                                                      │
                                                      ▼
                                        ┌─────────────────────────┐
                                        │   VPS (94.232.40.82)    │
                                        │   Caddy + Go API        │
                                        │   PostgreSQL            │
                                        │   Firewall: только CF   │
                                        └─────────────────────────┘
```

---

## Шаг 1: Создать поддомен api.oen.meybz.asia

### В Cloudflare Dashboard:

1. Перейди в **DNS** → **Records**
2. Нажми **Add record**
3. Заполни:

| Поле | Значение |
|------|----------|
| Type | `A` |
| Name | `api` |
| IPv4 address | `94.232.40.82` |
| Proxy status | **Proxied** (оранжевое облако) ☁️ |
| TTL | Auto |

4. Нажми **Save**

> ⚠️ **ВАЖНО**: Облако должно быть **оранжевым** (Proxied) — это скрывает реальный IP сервера.

---

## Шаг 2: Настроить SSL/TLS

### В Cloudflare Dashboard → SSL/TLS:

#### Overview:
```
Encryption mode: Full (strict)
```

#### Edge Certificates:
```
Always Use HTTPS: ON
Minimum TLS Version: TLS 1.2
TLS 1.3: ON
Automatic HTTPS Rewrites: ON
```

---

## Шаг 3: Настроить Rate Limiting

### Cloudflare Dashboard → Security → WAF → Rate limiting rules

#### Правило 1: API Rate Limit

Нажми **Create rule**:

```
Rule name: API Rate Limit

If incoming requests match:
  Field: URI Path
  Operator: starts with
  Value: /api/

Then:
  Action: Block
  Duration: 1 minute

Rate:
  Requests: 100
  Period: 1 minute

With the same:
  IP
```

#### Правило 2: Webhook Protection (только Telegram)

```
Rule name: Webhook Rate Limit

If incoming requests match:
  Field: URI Path
  Operator: equals
  Value: /bot/webhook

Then:
  Action: Block
  Duration: 1 minute

Rate:
  Requests: 60
  Period: 1 minute

With the same:
  IP
```

---

## Шаг 4: Настроить WAF Rules

### Cloudflare Dashboard → Security → WAF → Managed rules

1. Включи **Cloudflare Managed Ruleset**
2. Включи **Cloudflare OWASP Core Ruleset**

### Custom Rules (Security → WAF → Custom rules):

#### Правило: Block Bad Bots

```
Rule name: Block Bad Bots

If:
  (cf.client.bot) and not (cf.verified_bot_category in {"Search Engine Crawler" "Monitoring & Analytics"})

Then:
  Action: Block
```

#### Правило: Block Countries (опционально)

```
Rule name: Block High-Risk Countries

If:
  (ip.geoip.country in {"CN" "RU" "KP"})

Then:
  Action: Managed Challenge
```

> Убери RU если нужны пользователи из России

---

## Шаг 5: Bot Protection

### Cloudflare Dashboard → Security → Bots:

```
Bot Fight Mode: ON
```

### Security → Settings:

```
Security Level: Medium
Challenge Passage: 30 minutes
Browser Integrity Check: ON
```

---

## Шаг 6: Firewall на сервере

Выполни на VPS:

```bash
# Сначала разреши SSH чтобы не потерять доступ!
sudo ufw allow 22/tcp

# Разреши только Cloudflare IPs на порты 80 и 443
# Cloudflare IPv4 ranges:
sudo ufw allow from 173.245.48.0/20 to any port 80,443 proto tcp
sudo ufw allow from 103.21.244.0/22 to any port 80,443 proto tcp
sudo ufw allow from 103.22.200.0/22 to any port 80,443 proto tcp
sudo ufw allow from 103.31.4.0/22 to any port 80,443 proto tcp
sudo ufw allow from 141.101.64.0/18 to any port 80,443 proto tcp
sudo ufw allow from 108.162.192.0/18 to any port 80,443 proto tcp
sudo ufw allow from 190.93.240.0/20 to any port 80,443 proto tcp
sudo ufw allow from 188.114.96.0/20 to any port 80,443 proto tcp
sudo ufw allow from 197.234.240.0/22 to any port 80,443 proto tcp
sudo ufw allow from 198.41.128.0/17 to any port 80,443 proto tcp
sudo ufw allow from 162.158.0.0/15 to any port 80,443 proto tcp
sudo ufw allow from 104.16.0.0/13 to any port 80,443 proto tcp
sudo ufw allow from 104.24.0.0/14 to any port 80,443 proto tcp
sudo ufw allow from 172.64.0.0/13 to any port 80,443 proto tcp
sudo ufw allow from 131.0.72.0/22 to any port 80,443 proto tcp

# Запрети всё остальное
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Включи файрвол
sudo ufw enable

# Проверь правила
sudo ufw status
```

---

## Шаг 7: Обновить .env на сервере

```bash
ssh root@oen.meybz.asia

cat > /opt/focus/.env << 'EOF'
DB_PASSWORD=h7X2m9Qp4Lz8Wv3N
BOT_TOKEN=8243978756:AAGIN5Yax6GDv-YlhlLV25XLKK9b0LqSqdg
GEMINI_KEY=AIzaSyDFAnyvWw5tmV8S4rftS_SP5-oEUcpvzzQ
DOMAIN=api.oen.meybz.asia
EOF
```

---

## Шаг 8: Перезапустить сервисы

```bash
cd /opt/focus
docker compose down
docker compose up -d
```

---

## Шаг 9: Обновить Telegram Webhook

```bash
curl -X POST "https://api.telegram.org/bot8243978756:AAGIN5Yax6GDv-YlhlLV25XLKK9b0LqSqdg/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://api.oen.meybz.asia/bot/webhook"}'
```

---

## Проверка

### 1. Проверить что IP скрыт:
```bash
dig api.oen.meybz.asia +short
# Должен показать Cloudflare IP, не 94.232.40.82
```

### 2. Проверить API:
```bash
curl https://api.oen.meybz.asia/health
# {"status":"ok"}
```

### 3. Проверить Rate Limit:
```bash
# Отправь 150 запросов подряд — после 100 должен получить блок
for i in {1..150}; do curl -s -o /dev/null -w "%{http_code}\n" https://api.oen.meybz.asia/health; done
```

### 4. Проверить Webhook:
```bash
curl "https://api.telegram.org/bot8243978756:AAGIN5Yax6GDv-YlhlLV25XLKK9b0LqSqdg/getWebhookInfo"
```

---

## Мониторинг

### Cloudflare Analytics:

- **Security** → **Events** — атаки и блокировки
- **Analytics** → **Traffic** — трафик и запросы
- **Analytics** → **Security** — угрозы

### Алерты (опционально):

1. Cloudflare Dashboard → **Notifications**
2. Добавь алерты на:
   - DDoS Attack Alerter
   - Security Events Alert
   - Origin Error Rate Alert

---

## Итоговый чеклист

- [ ] A-запись `api.oen.meybz.asia` → Proxied (оранжевое облако)
- [ ] SSL/TLS: Full (strict)
- [ ] Always Use HTTPS: ON
- [ ] Minimum TLS: 1.2
- [ ] Rate Limiting: 100 req/min на /api/*
- [ ] WAF Managed Rules: ON
- [ ] Bot Fight Mode: ON
- [ ] Security Level: Medium
- [ ] UFW на сервере: только Cloudflare IPs
- [ ] .env обновлён: DOMAIN=api.oen.meybz.asia
- [ ] Docker перезапущен
- [ ] Telegram Webhook обновлён
- [ ] dig показывает CF IP, не реальный
