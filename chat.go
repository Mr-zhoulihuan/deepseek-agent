package main

import (
	"bufio"
	"fmt"
	"strings"
)

const chatSystemPrompt = `你是一个智能助手，可以帮助用户解决各种问题。
请用友好的中文回答用户的问题，回答要清晰、准确、有帮助。`

func handleFreeChat(scanner *bufio.Scanner) {
	PrintTitle("💬 自由对话")
	PrintInfo("输入 'exit' 或 'quit' 返回主菜单")
	PrintDivider()

	messages := []ChatMessage{
		{Role: "system", Content: chatSystemPrompt},
	}

	for {
		fmt.Println()
		fmt.Print(bold(cyan("You")) + ": ")
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

		PrintWaiting("思考中...")
		result, err := chatWithDeepSeekRaw(reqBody)
		if err != nil {
			PrintError(fmt.Sprintf("调用 DeepSeek 失败: %v", err))
			messages = messages[:len(messages)-1]
			continue
		}

		if result.ReasoningContent != "" {
			fmt.Println()
			PrintDivider()
			fmt.Println(bold(yellow("🤔 DeepSeek 思考过程:")))
			fmt.Println(dim(strings.Repeat("─", 50)))
			fmt.Println(dim(result.ReasoningContent))
			fmt.Println(dim(strings.Repeat("─", 50)))
		}

		fmt.Println()
		fmt.Print(bold(green("AI")) + ": ")
		fmt.Println(result.Content)
		messages = append(messages, ChatMessage{Role: "assistant", Content: result.Content})
	}
}

func chatWithDeepSeekRaw(reqBody ChatRequestBody) (*ChatResult, error) {
	return chatWithDeepSeekMessages(reqBody)
}
