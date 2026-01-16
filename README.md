# 🎯 Focus App

**Focus** is a minimalist, AI-powered task manager designed exclusively for **Telegram Web Apps (TWA)**. It uses Google Gemini models to intelligently parse voice or text commands into structured tasks, organizing them by context (Daily Tasks, Long-term goals, Routines).

![License](https://img.shields.io/badge/license-MIT-green) ![Version](https://img.shields.io/badge/version-0.0.2-blue) ![Platform](https://img.shields.io/badge/platform-Telegram-2CA5E0)

🔗 **Live Demo:** [https://enkinvsh.github.io/focus/](https://enkinvsh.github.io/focus/)

---

## ✨ Key Features

*   **🧠 AI-Powered Entry:** Just say "Buy milk" or "Learn Python". The AI extracts the title, assigns priority, and categorizes it based on your current tab context.
*   **☁️ Cloud Sync:** Uses **Telegram CloudStorage** to sync tasks across all your devices (Desktop, Mobile) without requiring a login or external database.
*   **🛡️ Secure Architecture:** API Keys are hidden behind a **Cloudflare Worker** proxy. No sensitive data is exposed to the client.
*   **🗣️ Voice Control:** Native Web Speech API integration for hands-free task addition.
*   **🎨 Pro UI/UX:**
    *   Glassmorphism design with **Tailwind CSS**.
    *   **Swipe Gestures** for tab navigation.
    *   **Haptic Feedback** for a native app feel.
    *   4 OLED-friendly dark themes.
*   **qwerty Localization:** Fully translated into **English** and **Russian**.

---

## 🛠️ Architecture

The application follows a **Serverless** architecture, relying on Telegram's infrastructure for data and Cloudflare for compute.

```mermaid
graph LR
    User[User / Telegram App] -- Voice/Text --> Client[Focus App (SPA)]
    Client -- Sync Tasks --> TG[Telegram CloudStorage]
    Client -- AI Prompt --> Proxy[Cloudflare Worker]
    Proxy -- Secure Request --> Gemini[Google Gemini API]
    Gemini -- JSON Response --> Proxy
    Proxy --> Client
```

---

## 🚀 Installation & Setup

### 1. Frontend (GitHub Pages)
The app is a single `index.html` file. 
1. Fork this repository.
2. Enable **GitHub Pages** in repository settings (Source: `main` branch).
3. The app is ready to be added to Telegram via BotFather.

### 2. Backend (Cloudflare Worker)
To protect your Gemini API Key, we use a lightweight proxy.

1. Log in to [Cloudflare Dashboard](https://dash.cloudflare.com/).
2. Go to **Workers & Pages** -> **Create Worker**.
3. Paste the following code:

```javascript
export default {
  async fetch(request, env) {
    // CORS Headers
    const corsHeaders = {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "POST, OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type",
    };

    if (request.method === "OPTIONS") return new Response(null, { headers: corsHeaders });
    if (request.method !== "POST") return new Response("Method not allowed", { status: 405 });

    try {
      const { model, contents } = await request.json();
      const apiKey = env.GEMINI_KEY; // Secret Variable

      const response = await fetch(`https://generativelanguage.googleapis.com/v1beta/models/${model}:generateContent?key=${apiKey}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ contents }),
      });

      const data = await response.json();
      return new Response(JSON.stringify(data), { headers: { "Content-Type": "application/json", ...corsHeaders } });
    } catch (e) {
      return new Response(JSON.stringify({ error: e.message }), { status: 500, headers: corsHeaders });
    }
  },
};
```

4. **Settings -> Variables**: Add a variable named `GEMINI_KEY` with your Google AI Studio key.
5. Deploy and copy the Worker URL.
6. Update the `PROXY_URL` constant in `index.html` with your new URL.

---

## 📱 Usage

1. **Add to Telegram:** create a new bot via @BotFather, then create a "New Web App" linked to your GitHub Pages URL.
2. **Contexts:**
    *   **TASK:** Immediate actions.
    *   **LONG:** Long-term goals.
    *   **ROUTINE:** Recurring habits.
3. **Gestures:** Swipe Left/Right to switch tabs. Tap empty space to close menus.

---

## 📦 Changelog

### v0.0.2 (Cloud Edition)
*   **New:** Telegram CloudStorage integration for cross-device sync.
*   **Security:** Cloudflare Worker proxy integration.
*   **UX:** Added "Syncing" overlay state.
*   **Fix:** Modal backdrop click issues resolved.

### v0.0.1 (Initial Release)
*   Core AI Task extraction.
*   Localization (EN/RU).
*   Manrope Font & Tailwind styling.

---

**Built with ❤️ by AI & Human Collaboration**
