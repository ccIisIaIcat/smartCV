package LLM

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIAPI 结构体定义了OpenAI API客户端
type OpenAIAPI struct {
	apiKey  string
	baseURL string
	headers map[string]string
}

// Message 定义了聊天消息的结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatPayload 定义了聊天请求的负载
type ChatPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse 定义了聊天响应的结构
type ChatResponse struct {
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	// 可以根据需要添加更多响应字段
}

// NewOpenAIAPI 创建一个新的OpenAI API客户端
func NewOpenAIAPI(apiKey string) *OpenAIAPI {
	return &OpenAIAPI{
		apiKey:  apiKey,
		baseURL: "https://api.openai-proxy.org/v1",
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", apiKey),
			"Content-Type":  "application/json",
		},
	}
}

// Chat 发送聊天请求
func (api *OpenAIAPI) Chat(userMessage, systemMessage string) (map[string]interface{}, error) {
	if systemMessage == "" {
		systemMessage = "You are a helpful assistant."
	}

	payload := ChatPayload{
		Model: "gpt-4-0125-preview",
		Messages: []Message{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", api.baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range api.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ChatResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("状态码: %d, 解析错误响应失败: %w", resp.StatusCode, err)
		}
		if errorResp.Error != nil {
			return nil, fmt.Errorf("API错误: %s", errorResp.Error.Message)
		}
		return nil, fmt.Errorf("状态码: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON解码失败: %w", err)
	}

	return result, nil
}
