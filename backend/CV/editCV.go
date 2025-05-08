package CV

import (
	"backend/LLM"
	"fmt"
	"strings"
)

// EditCV 根据修改要求编辑LaTeX代码
func EditCV(apiKey string, latexCode string, editRequirements string) (string, error) {
	client := LLM.NewOpenAIAPI(apiKey)

	prompt := fmt.Sprintf(`请根据以下修改要求修改LaTeX代码。只修改要求中提到的部分，其他部分保持完全不变。
只返回完整的LaTeX代码，不需要任何解释。

修改要求：
%s

原始LaTeX代码：
%s`, editRequirements, latexCode)

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
