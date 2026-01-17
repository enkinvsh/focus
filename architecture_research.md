# Architecture Research: Focus TMA

## Current State
- Frontend: Single HTML (likely `index.html`)
- Storage: Telegram CloudStorage
- Proxy: Cloudflare Worker

## Key Questions Analysis
1. DB Approach: Telegram CloudStorage vs Postgres
2. API Split: CF Worker vs Go VPS
3. Auth: Worker to VPS
4. Scheduling: Push notifications
5. Audio Pipeline: Transcription flow
