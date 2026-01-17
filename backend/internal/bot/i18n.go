package bot

type Texts struct {
	Welcome       string
	Features      string
	CTA           string
	BtnLaunch     string
	BtnBreathing  string
	BtnOpen       string
	BtnTry        string
	BreathingInfo string
	Help          string
	About         string
}

var I18n = map[string]Texts{
	"en": {
		Welcome:       "👋 <b>Welcome to Focus!</b>\n\nMinimalist AI-powered task manager.",
		Features:      "• Voice input for tasks\n• Smart priority sorting\n• Breathing exercises for focus\n• Cross-device sync",
		CTA:           "Tap the button below to start:",
		BtnLaunch:     "🎯 Launch Focus",
		BtnBreathing:  "🧘 Breathing Exercise",
		BtnOpen:       "🚀 Open App",
		BtnTry:        "🎯 Try Now",
		BreathingInfo: "🧘 <b>Breathing Exercise</b>\n\n1-minute technique to improve concentration:\n\n• Inhale (4 sec)\n• Hold (4 sec)\n• Exhale (4 sec)\n• 5 cycles\n\nTap the \"Focus\" title in the app to start.",
		Help:          "📖 <b>Focus Guide</b>\n\n<b>How it works:</b>\n1. Tap «Launch Focus» button\n2. Record tasks by voice or text\n3. AI sorts them by category\n4. Swipe between tabs: Tasks / Long / Routine\n\n<b>Breathing Exercise:</b>\nTap on the «Focus» title in-app\n\n<b>Quick gestures:</b>\n• Swipe left/right — switch tabs\n• Tap a task — action menu",
		About:         "ℹ️ <b>About Focus</b>\n\n<b>Version:</b> 0.0.4\n\n<b>Technologies:</b>\n• PostgreSQL for data storage\n• Google Gemini AI for task processing\n• Go backend for API\n\n<b>Privacy:</b>\n• Data stored securely on our servers\n• No third-party accounts required\n• Secure API for all requests",
	},
	"ru": {
		Welcome:       "👋 <b>Добро пожаловать в Focus!</b>\n\nМинималистичный менеджер задач с ИИ.",
		Features:      "• Голосовой ввод задач\n• Умное распределение по приоритетам\n• Дыхательные упражнения для фокуса\n• Синхронизация между устройствами",
		CTA:           "Жми кнопку ниже, чтобы начать:",
		BtnLaunch:     "🎯 Запустить Focus",
		BtnBreathing:  "🧘 Дыхательное упражнение",
		BtnOpen:       "🚀 Открыть приложение",
		BtnTry:        "🎯 Попробовать",
		BreathingInfo: "🧘 <b>Дыхательное упражнение</b>\n\n1-минутная техника для улучшения концентрации:\n\n• Вдох (4 сек)\n• Задержка (4 сек)\n• Выдох (4 сек)\n• 5 циклов\n\nНажми на заголовок «Focus» в приложении, чтобы начать.",
		Help:          "📖 <b>Руководство по Focus</b>\n\n<b>Как это работает:</b>\n1. Нажми кнопку «Запустить Focus»\n2. Записывай задачи голосом или текстом\n3. ИИ распределит их по категориям\n4. Свайпай между вкладками: Задачи / Долгие / Рутина\n\n<b>Дыхательное упражнение:</b>\nНажми на заголовок «Focus» в приложении\n\n<b>Горячие жесты:</b>\n• Свайп влево/вправо — смена вкладки\n• Нажми на задачу — меню действий",
		About:         "ℹ️ <b>О приложении Focus</b>\n\n<b>Версия:</b> 0.0.4\n\n<b>Технологии:</b>\n• PostgreSQL для хранения данных\n• Google Gemini AI для обработки задач\n• Go бэкенд для API\n\n<b>Приватность:</b>\n• Данные хранятся безопасно на наших серверах\n• Никаких сторонних аккаунтов\n• Защищённый API для всех запросов",
	},
}

func GetTexts(langCode string) Texts {
	if len(langCode) >= 2 && langCode[:2] == "ru" {
		return I18n["ru"]
	}
	return I18n["en"]
}
