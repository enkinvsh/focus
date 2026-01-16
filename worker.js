export default {
    async fetch(request, env) {
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

            if (body.update_id) {
                const userLang = body.message?.from?.language_code || body.callback_query?.from?.language_code || 'en';
                const isRu = userLang.startsWith('ru');
                const t = isRu ? TEXTS.ru : TEXTS.en;

                if (body.message) {
                    const chatId = body.message.chat.id;
                    const text = body.message.text || "";

                    if (text === "/start") {
                        await sendTelegramPhoto(env.BOT_TOKEN, chatId,
                            "https://raw.githubusercontent.com/enkinvsh/focus/main/enter.png",
                            `${t.welcome}\n\n${t.features}\n\n${t.cta}`,
                            {
                                inline_keyboard: [
                                    [{ text: t.btn_launch, web_app: { url: "https://enkinvsh.github.io/focus/" } }],
                                    [{ text: t.btn_breathing, callback_data: "breathing_info" }]
                                ]
                            }
                        );
                    }

                    if (text === "/help") {
                        await sendTelegramMessage(env.BOT_TOKEN, chatId, t.help, {
                            inline_keyboard: [[{ text: t.btn_open, web_app: { url: "https://enkinvsh.github.io/focus/" } }]]
                        });
                    }

                    if (text === "/about") {
                        await sendTelegramMessage(env.BOT_TOKEN, chatId, t.about, {
                            inline_keyboard: [[{ text: "GitHub", url: "https://github.com/enkinvsh/focus" }]]
                        });
                    }

                    return new Response("OK");
                }

                if (body.callback_query) {
                    const callbackId = body.callback_query.id;
                    const chatId = body.callback_query.message.chat.id;
                    const data = body.callback_query.data;

                    await answerCallbackQuery(env.BOT_TOKEN, callbackId);

                    if (data === "breathing_info") {
                        await sendTelegramMessage(env.BOT_TOKEN, chatId, t.breathing_info, {
                            inline_keyboard: [[{ text: t.btn_try, web_app: { url: "https://enkinvsh.github.io/focus/" } }]]
                        });
                    }

                    return new Response("OK");
                }

                return new Response("OK");
            }

            if (body.model && body.contents) {
                const apiKey = env.GEMINI_KEY;
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

const TEXTS = {
    en: {
        welcome: "👋 <b>Welcome to Focus!</b>\n\nMinimalist AI-powered task manager.",
        features: "• Voice input for tasks\n• Smart priority sorting\n• Breathing exercises for focus\n• Cross-device sync",
        cta: "Tap the button below to start:",
        btn_launch: "🎯 Launch Focus",
        btn_breathing: "🧘 Breathing Exercise",
        btn_open: "🚀 Open App",
        btn_try: "🎯 Try Now",
        breathing_info: "🧘 <b>Breathing Exercise</b>\n\n1-minute technique to improve concentration:\n\n• Inhale (4 sec)\n• Hold (4 sec)\n• Exhale (4 sec)\n• 5 cycles\n\nTap the \"Focus\" title in the app to start.",
        help: "📖 <b>Focus Guide</b>\n\n<b>How it works:</b>\n1. Tap «Launch Focus» button\n2. Record tasks by voice or text\n3. AI sorts them by category\n4. Swipe between tabs: Tasks / Long / Routine\n\n<b>Breathing Exercise:</b>\nTap on the «Focus» title in-app\n\n<b>Quick gestures:</b>\n• Swipe left/right — switch tabs\n• Tap a task — action menu",
        about: "ℹ️ <b>About Focus</b>\n\n<b>Version:</b> 0.0.3\n\n<b>Technologies:</b>\n• Telegram CloudStorage for sync\n• Google Gemini AI for task processing\n• Cloudflare Workers for API security\n\n<b>Privacy:</b>\n• Data stored only in Telegram\n• No third-party accounts\n• Secure proxy for AI requests"
    },
    ru: {
        welcome: "👋 <b>Добро пожаловать в Focus!</b>\n\nМинималистичный менеджер задач с ИИ.",
        features: "• Голосовой ввод задач\n• Умное распределение по приоритетам\n• Дыхательные упражнения для фокуса\n• Синхронизация между устройствами",
        cta: "Жми кнопку ниже, чтобы начать:",
        btn_launch: "🎯 Запустить Focus",
        btn_breathing: "🧘 Дыхательное упражнение",
        btn_open: "🚀 Открыть приложение",
        btn_try: "🎯 Попробовать",
        breathing_info: "🧘 <b>Дыхательное упражнение</b>\n\n1-минутная техника для улучшения концентрации:\n\n• Вдох (4 сек)\n• Задержка (4 сек)\n• Выдох (4 сек)\n• 5 циклов\n\nНажми на заголовок «Focus» в приложении, чтобы начать.",
        help: "📖 <b>Руководство по Focus</b>\n\n<b>Как это работает:</b>\n1. Нажми кнопку «Запустить Focus»\n2. Записывай задачи голосом или текстом\n3. ИИ распределит их по категориям\n4. Свайпай между вкладками: Задачи / Долгие / Рутина\n\n<b>Дыхательное упражнение:</b>\nНажми на заголовок «Focus» в приложении\n\n<b>Горячие жесты:</b>\n• Свайп влево/вправо — смена вкладки\n• Нажми на задачу — меню действий",
        about: "ℹ️ <b>О приложении Focus</b>\n\n<b>Версия:</b> 0.0.3\n\n<b>Технологии:</b>\n• Telegram CloudStorage для синхронизации\n• Google Gemini AI для обработки задач\n• Cloudflare Workers для безопасности API\n\n<b>Приватность:</b>\n• Данные хранятся только в Telegram\n• Никаких сторонних аккаунтов\n• Защищённый прокси для AI запросов"
    }
};

async function sendTelegramPhoto(token, chatId, photoUrl, caption, replyMarkup) {
    await fetch(`https://api.telegram.org/bot${token}/sendPhoto`, {
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

async function sendTelegramMessage(token, chatId, text, replyMarkup) {
    await fetch(`https://api.telegram.org/bot${token}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            text: text,
            parse_mode: "HTML",
            reply_markup: replyMarkup
        })
    });
}

async function answerCallbackQuery(token, callbackId) {
    await fetch(`https://api.telegram.org/bot${token}/answerCallbackQuery`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ callback_query_id: callbackId })
    });
}