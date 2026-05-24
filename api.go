package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var deepseekApiKey string
var deepseekModel string

func init() {
	deepseekApiKey = os.Getenv("DEEPSEEK_API_KEY")
	if deepseekApiKey == "" {
		deepseekApiKey = "**************************"
	}
	deepseekModel = os.Getenv("DEEPSEEK_MODEL")
	if deepseekModel == "" {
		deepseekModel = "deepseek-v4-pro"
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ThinkingOpt struct {
	Type string `json:"type"`
}

type ChatRequestBody struct {
	Messages        []ChatMessage `json:"messages"`
	Model           string        `json:"model"`
	ReasoningEffort string        `json:"reasoning_effort,omitempty"`
	Thinking        *ThinkingOpt  `json:"thinking,omitempty"`
	MaxTokens       int           `json:"max_tokens"`
	Stream          bool          `json:"stream"`
}

type ChatResponseChoiceMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
}

type ChatResponseChoice struct {
	Message ChatResponseChoiceMessage `json:"message"`
}

type ChatResponseBody struct {
	Choices []ChatResponseChoice `json:"choices"`
}

type ChatResult struct {
	Content          string
	ReasoningContent string
}

func chatWithDeepSeek(systemPrompt, userPrompt string) (*ChatResult, error) {
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	reqBody := ChatRequestBody{
		Messages:        messages,
		Model:           deepseekModel,
		ReasoningEffort: "medium",
		Thinking:        &ThinkingOpt{Type: "enabled"},
		MaxTokens:       2048,
		Stream:          false,
	}

	return callDeepSeek(reqBody)
}

func chatWithDeepSeekMessages(reqBody ChatRequestBody) (*ChatResult, error) {
	return callDeepSeek(reqBody)
}

func callDeepSeek(reqBody ChatRequestBody) (*ChatResult, error) {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+deepseekApiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 请求失败: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("DeepSeek API 错误 (status %d): %s", res.StatusCode, string(body))
	}

	var chatResp ChatResponseBody
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("API 返回为空")
	}

	reasoning := chatResp.Choices[0].Message.ReasoningContent

	return &ChatResult{
		Content:          chatResp.Choices[0].Message.Content,
		ReasoningContent: reasoning,
	}, nil
}

func printWithReasoning(result *ChatResult) {
	if result.ReasoningContent != "" {
		PrintDivider()
		fmt.Println(bold(yellow("🤔 DeepSeek 思考过程:")))
		fmt.Println(dim(strings.Repeat("─", 50)))
		fmt.Println(dim(result.ReasoningContent))
		fmt.Println(dim(strings.Repeat("─", 50)))
	}
}
