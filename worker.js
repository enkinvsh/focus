export default {
    async fetch(request, env) {
        // === 1. Обработка CORS (для Web App) ===
        if (request.method === "OPTIONS") {
            return new Response(null, {
                headers: {
                    "Access-Control-Allow-Origin": "*",
                    "Access-Control-Allow-Methods": "POST, OPTIONS",
                    "Access-Control-Allow-Headers": "Content-Type",
                },
            });
        }

        if (request.method !== "POST") return new Response("Only POST allowed", { status: 405 });

        try {
            const body = await request.json();

            // === 2. ЛОГИКА TELEGRAM BOTA (если пришло сообщение) ===
            // Проверяем наличие update_id, чтобы отличить запрос от Telegram
            if (body.update_id && body.message) {
                const chatId = body.message.chat.id;
                const text = body.message.text || "";

                if (text === "/start") {
                    await sendTelegramPhoto(env.BOT_TOKEN, chatId,
                        // Ссылка на raw-изображение в репозитории
                        "https://raw.githubusercontent.com/enkinvsh/focus/main/promo_banner.png",
                        "👋 <b>Добро пожаловать в Focus!</b>\n\nМинималистичный менеджер задач с ИИ.\nЖми кнопку ниже, чтобы начать:",
                        {
                            inline_keyboard: [[
                                // Ссылка на GitHub Pages
                                { text: "🎯 Запустить Focus", web_app: { url: "https://enkinvsh.github.io/focus/" } }
                            ]]
                        }
                    );
                }
                return new Response("OK");
            }

            // === 3. ЛОГИКА GEMINI PROXY (если пришел запрос от приложения) ===
            if (body.model && body.contents) {
                const apiKey = env.GEMINI_KEY;
                // Убедитесь, что GEMINI_KEY установлен в Cloudflare Dashboard -> Settings -> Variables

                const url = `https://generativelanguage.googleapis.com/v1beta/models/${body.model}:generateContent?key=${apiKey}`;

                const response = await fetch(url, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ contents: body.contents }),
                });

                const data = await response.json();
                return new Response(JSON.stringify(data), {
                    headers: {
                        "Content-Type": "application/json",
                        "Access-Control-Allow-Origin": "*",
                    },
                });
            }

            return new Response("Unknown request", { status: 400 });

        } catch (e) {
            return new Response(JSON.stringify({ error: e.message }), {
                status: 500,
                headers: { "Content-Type": "application/json" }
            });
        }
    },
};

// Вспомогательная функция для отправки фото
async function sendTelegramPhoto(token, chatId, photoUrl, caption, replyMarkup) {
    const url = `https://api.telegram.org/bot${token}/sendPhoto`;
    await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            photo: photoUrl,
            caption: caption,
            parse_mode: "HTML",
            reply_markup: replyMarkup
        })
    });
}