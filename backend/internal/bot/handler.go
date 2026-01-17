package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	WebAppURL = "https://enkinvsh.github.io/focus/"
	PhotoURL  = "https://raw.githubusercontent.com/enkinvsh/focus/main/enter.png"
	GitHubURL = "https://github.com/enkinvsh/focus"
)

type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	From      *User  `json:"from,omitempty"`
	Text      string `json:"text,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data,omitempty"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineButton `json:"inline_keyboard"`
}

type InlineButton struct {
	Text         string  `json:"text"`
	URL          string  `json:"url,omitempty"`
	CallbackData string  `json:"callback_data,omitempty"`
	WebApp       *WebApp `json:"web_app,omitempty"`
}

type WebApp struct {
	URL string `json:"url"`
}

func HandleWebhook(update *Update) error {
	var langCode string
	if update.Message != nil && update.Message.From != nil {
		langCode = update.Message.From.LanguageCode
	} else if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
		langCode = update.CallbackQuery.From.LanguageCode
	}
	t := GetTexts(langCode)

	if update.Message != nil {
		return handleMessage(update.Message, t)
	}

	if update.CallbackQuery != nil {
		return handleCallback(update.CallbackQuery, t)
	}

	return nil
}

func handleMessage(msg *Message, t Texts) error {
	chatID := msg.Chat.ID
	text := msg.Text

	switch text {
	case "/start":
		caption := fmt.Sprintf("%s\n\n%s\n\n%s", t.Welcome, t.Features, t.CTA)
		keyboard := InlineKeyboard{
			InlineKeyboard: [][]InlineButton{
				{{Text: t.BtnLaunch, WebApp: &WebApp{URL: WebAppURL}}},
				{{Text: t.BtnBreathing, CallbackData: "breathing_info"}},
			},
		}
		return sendPhoto(chatID, PhotoURL, caption, keyboard)

	case "/help":
		keyboard := InlineKeyboard{
			InlineKeyboard: [][]InlineButton{
				{{Text: t.BtnOpen, WebApp: &WebApp{URL: WebAppURL}}},
			},
		}
		return sendMessage(chatID, t.Help, keyboard)

	case "/about":
		keyboard := InlineKeyboard{
			InlineKeyboard: [][]InlineButton{
				{{Text: "GitHub", URL: GitHubURL}},
			},
		}
		return sendMessage(chatID, t.About, keyboard)
	}

	return nil
}

func handleCallback(cb *CallbackQuery, t Texts) error {
	// Answer callback to remove loading state
	if err := answerCallback(cb.ID); err != nil {
		return err
	}

	if cb.Message == nil {
		return nil
	}

	chatID := cb.Message.Chat.ID

	switch cb.Data {
	case "breathing_info":
		keyboard := InlineKeyboard{
			InlineKeyboard: [][]InlineButton{
				{{Text: t.BtnTry, WebApp: &WebApp{URL: WebAppURL}}},
			},
		}
		return sendMessage(chatID, t.BreathingInfo, keyboard)
	}

	return nil
}

func sendPhoto(chatID int64, photoURL, caption string, keyboard InlineKeyboard) error {
	token := os.Getenv("BOT_TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", token)

	body := map[string]interface{}{
		"chat_id":      chatID,
		"photo":        photoURL,
		"caption":      caption,
		"parse_mode":   "HTML",
		"reply_markup": keyboard,
	}

	return postJSON(url, body)
}

func sendMessage(chatID int64, text string, keyboard InlineKeyboard) error {
	token := os.Getenv("BOT_TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"parse_mode":   "HTML",
		"reply_markup": keyboard,
	}

	return postJSON(url, body)
}

func answerCallback(callbackID string) error {
	token := os.Getenv("BOT_TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)

	body := map[string]interface{}{
		"callback_query_id": callbackID,
	}

	return postJSON(url, body)
}

func postJSON(url string, body map[string]interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram API error: %d", resp.StatusCode)
	}

	return nil
}
