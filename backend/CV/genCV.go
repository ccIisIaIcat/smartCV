package CV

import (
	"fmt"
	"strings"

	"backend/LLM"
)

// GenerateCV 生成LaTeX格式的简历
func GenerateCV(apiKey, introduction, requirements string, level PromptLevel) (string, error) {
	client := LLM.NewOpenAIAPI(apiKey)
	generator := NewResumePromptGenerator()

	// 根据等级选择对应的 prompt
	var prompt string
	switch level {
	case Honest:
		prompt = generator.GenerateHonestPrompt(introduction, requirements)
	case Moderate:
		prompt = generator.GenerateModeratePrompt(introduction, requirements)
	case Aggressive:
		prompt = generator.GenerateAggressivePrompt(introduction, requirements)
	case Creative:
		prompt = generator.GenerateCreativePrompt(introduction, requirements)
	case Expert:
		prompt = generator.GenerateExpertPrompt(introduction, requirements)
	default:
		return "", fmt.Errorf("不支持的简历等级：%v", level)
	}

	// 发送请求
	response, err := client.Chat(prompt, "")
	if err != nil {
		return "", fmt.Errorf("调用 API 失败: %w", err)
	}

	// 从响应中提取 LaTeX 代码
	if content, ok := response["choices"].([]interface{}); ok && len(content) > 0 {
		if message, ok := content[0].(map[string]interface{}); ok {
			if messageContent, ok := message["message"].(map[string]interface{}); ok {
				if text, ok := messageContent["content"].(string); ok {
					// 清理代码（移除可能的 markdown 标记）
					text = strings.TrimSpace(text)
					text = strings.TrimPrefix(text, "```latex")
					text = strings.TrimPrefix(text, "```")
					text = strings.TrimSuffix(text, "```")
					return strings.TrimSpace(text), nil
				}
			}
		}
	}

	return "", fmt.Errorf("无法从API响应中提取LaTeX代码")
}
