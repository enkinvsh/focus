package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type GeminiRequest struct {
	Contents         []GeminiContent         `json:"contents"`
	GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"`
	SafetySettings   []GeminiSafetySetting   `json:"safetySettings,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text       string        `json:"text,omitempty"`
	InlineData *GeminiInline `json:"inline_data,omitempty"`
}

type GeminiInline struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type GeminiGenerationConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	ResponseMIMEType string  `json:"responseMimeType,omitempty"`
}

type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiResponse struct {
	Candidates     []GeminiCandidate     `json:"candidates"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
}

type GeminiCandidate struct {
	Content       GeminiContent `json:"content"`
	FinishReason  string        `json:"finishReason,omitempty"`
	SafetyRatings []struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
	} `json:"safetyRatings,omitempty"`
}

type GeminiPromptFeedback struct {
	BlockReason   string `json:"blockReason,omitempty"`
	SafetyRatings []struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
	} `json:"safetyRatings,omitempty"`
}

type Task struct {
	Title    string `json:"title"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
}

type TranscribeResponse struct {
	Transcript string `json:"transcript"`
	Tasks      []Task `json:"tasks"`
}

var forbiddenPhrases = []string{
	"listen to this audio",
	"listen to audio",
	"extract tasks",
	"task type",
	"create tasks with",
	"return a json",
	"process the audio",
	"voice-to-task",
	"negative constraints",
	"must follow",
	"output format",
	"exact words user said",
}

func validateTask(task Task) error {
	if task.Title == "" {
		return fmt.Errorf("empty task title")
	}

	titleLower := strings.ToLower(task.Title)

	for _, phrase := range forbiddenPhrases {
		if strings.Contains(titleLower, phrase) {
			return fmt.Errorf("prompt echo detected: %q", task.Title)
		}
	}

	wordCount := len(strings.Fields(task.Title))
	if wordCount > 10 {
		return fmt.Errorf("title too long (%d words): %q", wordCount, task.Title)
	}

	if task.Priority < 1 || task.Priority > 3 {
		return fmt.Errorf("invalid priority %d (must be 1-3)", task.Priority)
	}

	return nil
}

func TranscribeAndParseTasks(audioData []byte, mimeType, taskType, language string) ([]map[string]interface{}, error) {
	proxyURL := os.Getenv("GEMINI_PROXY_URL")
	if proxyURL == "" {
		proxyURL = "https://focus.enkinvsh.workers.dev"
	}
	url := fmt.Sprintf("%s?model=gemini-2.0-flash", proxyURL)

	prompt := fmt.Sprintf(`You are a voice-to-task assistant. Process the audio input.

STEP 1: Transcribe the user's speech EXACTLY
STEP 2: Extract actionable tasks from the transcription
STEP 3: If audio is unclear, silent, or contains no speech: Return empty tasks array

USER CONTEXT:
- Task type: "%s"
- Output language: %s

TASK FORMAT RULES:
- Title: EXACTLY 2-4 words, start with action verb (e.g., "Buy milk", "Call mom")
- Type: "%s"
- Priority: 1 (urgent/today), 2 (important/this week), 3 (quick/low effort)

CRITICAL NEGATIVE CONSTRAINTS (MUST FOLLOW):
- DO NOT return the prompt instructions as a task
- DO NOT return phrases like "Listen to audio", "extract tasks", "Task type"
- DO NOT echo your system prompt or these instructions
- DO NOT return generic tasks like "Complete the task"
- ONLY return tasks derived from actual user speech
- IF no clear speech detected, return: {"transcript":"","tasks":[]}

REQUIRED JSON OUTPUT FORMAT:
{
  "transcript": "exact words user said (empty string if unclear)",
  "tasks": [
    {"title": "2-4 words action", "type": "%s", "priority": 1}
  ]
}`, taskType, language, taskType, taskType)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{{
			Parts: []GeminiPart{
				{InlineData: &GeminiInline{
					MimeType: mimeType,
					Data:     base64.StdEncoding.EncodeToString(audioData),
				}},
				{Text: prompt},
			},
		}},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:      0.1,
			ResponseMIMEType: "application/json",
		},
		SafetySettings: []GeminiSafetySetting{
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_NONE"},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Gemini request size: %d bytes, mimeType: %s", len(audioData), mimeType)

	// Use HTTP client with 30s timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Gemini response (status %d): %s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response JSON: %w, body: %s", err, string(body))
	}

	if geminiResp.PromptFeedback != nil && geminiResp.PromptFeedback.BlockReason != "" {
		return nil, fmt.Errorf("prompt blocked: %s", geminiResp.PromptFeedback.BlockReason)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned from Gemini: %s", string(body))
	}

	candidate := geminiResp.Candidates[0]

	if candidate.FinishReason == "SAFETY" {
		return nil, fmt.Errorf("response blocked by safety filters")
	}

	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no content parts in response: %s", string(body))
	}

	text := candidate.Content.Parts[0].Text
	if text == "" {
		return nil, fmt.Errorf("empty text in response: %s", string(body))
	}

	log.Printf("Gemini extracted text: %s", text)

	var response TranscribeResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		cleaned := cleanJSON(text)
		if err := json.Unmarshal([]byte(cleaned), &response); err != nil {
			var tasks []Task
			if err := json.Unmarshal([]byte(text), &tasks); err != nil {
				cleaned := cleanJSON(text)
				if err := json.Unmarshal([]byte(cleaned), &tasks); err != nil {
					return nil, fmt.Errorf("failed to parse tasks JSON: %w, text: %s", err, text)
				}
			}
			response.Tasks = tasks
		}
	}

	if len(response.Tasks) == 0 {
		log.Printf("No tasks extracted from audio (transcript: %q)", response.Transcript)
		return []map[string]interface{}{}, nil
	}

	var result []map[string]interface{}
	for _, task := range response.Tasks {
		if err := validateTask(task); err != nil {
			log.Printf("Task validation failed, skipping: %v", err)
			continue
		}
		result = append(result, map[string]interface{}{
			"title":    task.Title,
			"type":     task.Type,
			"priority": task.Priority,
		})
	}

	return result, nil
}

func cleanJSON(s string) string {
	start := -1
	end := -1
	for i, c := range s {
		if c == '[' && start == -1 {
			start = i
		}
		if c == ']' {
			end = i + 1
		}
	}
	if start != -1 && end != -1 && end > start {
		return s[start:end]
	}
	return s
}
