package main

import (
	"bufio"
	"fmt"
	"strings"
)

const chatSystemPrompt = `你是一个智能助手，可以帮助用户解决各种问题。
请用友好的中文回答用户的问题，回答要清晰、准确、有帮助。`

func handleFreeChat(scanner *bufio.Scanner) {
	fmt.Println("\n── 自由对话 ──")
	fmt.Println("💡 输入 'exit' 或 'quit' 返回主菜单")
	fmt.Println("──────────────────────────────────────────────")

	messages := []ChatMessage{
		{Role: "system", Content: chatSystemPrompt},
	}

	for {
		fmt.Print("\nYou: ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())

		if userInput == "exit" || userInput == "quit" {
			break
		}

		if userInput == "" {
			continue
		}

		messages = append(messages, ChatMessage{Role: "user", Content: userInput})

		reqBody := ChatRequestBody{
			Messages:        messages,
			Model:           deepseekModel,
			ReasoningEffort: "medium",
			Thinking:        &ThinkingOpt{Type: "enabled"},
			MaxTokens:       2048,
			Stream:          false,
		}

		result, err := chatWithDeepSeekRaw(reqBody)
		if err != nil {
			fmt.Printf("❌ 调用 DeepSeek 失败: %v\n", err)
			messages = messages[:len(messages)-1]
			continue
		}

		if result.ReasoningContent != "" {
			fmt.Println("\n──────────────────────────────────────────────")
			fmt.Println("🤔 DeepSeek 思考过程:")
			fmt.Println("──────────────────────────────────────────────")
			fmt.Println(result.ReasoningContent)
			fmt.Println("──────────────────────────────────────────────")
		}

		fmt.Printf("AI: %s\n", result.Content)
		messages = append(messages, ChatMessage{Role: "assistant", Content: result.Content})
	}
}

func chatWithDeepSeekRaw(reqBody ChatRequestBody) (*ChatResult, error) {
	return chatWithDeepSeekMessages(reqBody)
}
