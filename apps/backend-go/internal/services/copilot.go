package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CopilotService struct{ client *http.Client }

func NewCopilotService() *CopilotService {
	return &CopilotService{client: &http.Client{Timeout: 60 * time.Second}}
}

type ChatResponse struct {
	ThreadID string `json:"threadId"`
	Message  string `json:"message"`
	Role     string `json:"role"`
}

func (s *CopilotService) Chat(orgID, userID, message, threadID string) (*ChatResponse, error) {
	if threadID == "" {
		threadID = uuid.New().String()
	}
	ctx := context.Background()

	// Persist message
	msgID := uuid.New().String()
	database.DB.Exec(ctx,
		`INSERT INTO mastra_messages (id, thread_id, content, role, type, "createdAt") VALUES ($1,$2,$3,'user','text',$4) ON CONFLICT DO NOTHING`,
		msgID, threadID, message, time.Now())

	// Call OpenAI
	reply, err := s.callOpenAI(message, threadID)
	if err != nil {
		return nil, err
	}

	// Persist assistant reply
	replyID := uuid.New().String()
	database.DB.Exec(ctx,
		`INSERT INTO mastra_messages (id, thread_id, content, role, type, "createdAt") VALUES ($1,$2,$3,'assistant','text',$4) ON CONFLICT DO NOTHING`,
		replyID, threadID, reply, time.Now())

	return &ChatResponse{ThreadID: threadID, Message: reply, Role: "assistant"}, nil
}

func (s *CopilotService) callOpenAI(message, threadID string) (string, error) {
	ctx := context.Background()

	// Get thread history
	rows, _ := database.DB.Query(ctx,
		`SELECT role, content FROM mastra_messages WHERE thread_id=$1 ORDER BY "createdAt" ASC LIMIT 20`, threadID)
	var messages []map[string]string
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var role, content string
			rows.Scan(&role, &content)
			messages = append(messages, map[string]string{"role": role, "content": content})
		}
	}
	messages = append(messages, map[string]string{"role": "user", "content": message})

	body, _ := json.Marshal(map[string]interface{}{
		"model":    "gpt-4o-mini",
		"messages": append([]map[string]string{{"role": "system", "content": "You are a social media content assistant helping users create engaging posts."}}, messages...),
	})

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+config.App.OpenAIAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct{ Content string }
		}
	}
	if err := json.Unmarshal(b, &result); err != nil || len(result.Choices) == 0 {
		return "", fmt.Errorf("openai error")
	}
	return result.Choices[0].Message.Content, nil
}

func (s *CopilotService) GetCredits(orgID string) (int, error) {
	ctx := context.Background()
	var total int
	database.DB.QueryRow(ctx,
		`SELECT COALESCE(SUM(credits),0) FROM credits WHERE organization_id=$1 AND type='ai_images'`, orgID).Scan(&total)
	return total, nil
}

func (s *CopilotService) ListThreads(userID string) (interface{}, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, title, "createdAt" FROM mastra_threads WHERE "resourceId"=$1 ORDER BY "createdAt" DESC LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var threads []map[string]interface{}
	for rows.Next() {
		var id, title string
		var createdAt time.Time
		rows.Scan(&id, &title, &createdAt)
		threads = append(threads, map[string]interface{}{"id": id, "title": title, "createdAt": createdAt})
	}
	return threads, nil
}

func (s *CopilotService) GetThreadMessages(threadID string) (interface{}, error) {
	ctx := context.Background()
	rows, err := database.DB.Query(ctx,
		`SELECT id, role, content, "createdAt" FROM mastra_messages WHERE thread_id=$1 ORDER BY "createdAt" ASC`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []map[string]interface{}
	for rows.Next() {
		var id, role, content string
		var createdAt time.Time
		rows.Scan(&id, &role, &content, &createdAt)
		messages = append(messages, map[string]interface{}{"id": id, "role": role, "content": content, "createdAt": createdAt})
	}
	return messages, nil
}

func (s *CopilotService) GenerateDraft(orgID, topic, platform, tone string, keywords []string) ([]string, error) {
	prompt := fmt.Sprintf("Generate 3 social media post drafts about: %s. Platform: %s. Tone: %s. Keywords: %s. Return each draft on a new line starting with ---",
		topic, platform, tone, strings.Join(keywords, ", "))
	reply, err := s.callOpenAI(prompt, uuid.New().String())
	if err != nil {
		return nil, err
	}
	parts := strings.Split(reply, "---")
	var drafts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			drafts = append(drafts, p)
		}
	}
	return drafts, nil
}

func (s *CopilotService) GenerateImage(orgID, prompt, style string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Check credits
	var credits int
	database.DB.QueryRow(ctx,
		`SELECT COALESCE(SUM(credits),0) FROM credits WHERE organization_id=$1 AND type='ai_images'`, orgID).Scan(&credits)
	if credits <= 0 {
		return nil, fmt.Errorf("insufficient credits")
	}

	body, _ := json.Marshal(map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": prompt,
		"n":      1,
		"size":   "1024x1024",
		"style":  style,
	})
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+config.App.OpenAIAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var result struct {
		Data []struct{ URL string }
	}
	if err := json.Unmarshal(b, &result); err != nil || len(result.Data) == 0 {
		return nil, fmt.Errorf("image generation failed")
	}

	// Deduct credit
	database.DB.Exec(ctx,
		`INSERT INTO credits (id, organization_id, credits, type, created_at, updated_at) VALUES ($1,$2,-1,'ai_images',$3,$4)`,
		uuid.New().String(), orgID, time.Now(), time.Now())

	return map[string]interface{}{"url": result.Data[0].URL}, nil
}

func (s *CopilotService) SeparatePosts(content string, maxChars int) []string {
	words := strings.Fields(content)
	var parts []string
	current := ""
	for _, w := range words {
		if len(current)+len(w)+1 > maxChars {
			if current != "" {
				parts = append(parts, strings.TrimSpace(current))
			}
			current = w
		} else {
			if current != "" {
				current += " "
			}
			current += w
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}
