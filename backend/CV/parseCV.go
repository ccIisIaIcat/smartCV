package CV

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ParseLatexToPDF 将LaTeX文本编译为PDF文件
// 返回生成的PDF文件路径和可能的错误
func ParseLatexToPDF(latexContent string) (string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "cv_temp_*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 创建临时tex文件
	texFile := filepath.Join(tempDir, "cv.tex")
	if err := os.WriteFile(texFile, []byte(latexContent), 0644); err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", fmt.Errorf("写入临时tex文件失败: %w", err)
	}

	// 准备xelatex命令
	cmd := exec.Command("xelatex", "-interaction=nonstopmode", "-output-directory="+tempDir, texFile)

	// 运行编译命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", fmt.Errorf("编译PDF失败: %v\n输出: %s", err, string(output))
	}

	// 检查PDF文件是否生成
	pdfFile := filepath.Join(tempDir, "cv.pdf")
	if _, err := os.Stat(pdfFile); err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", fmt.Errorf("PDF文件未生成: %w", err)
	}

	// 确保目标目录存在
	targetDir := filepath.Join("temp")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 使用时间戳生成文件名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	finalPDFPath := filepath.Join(targetDir, fmt.Sprintf("cv_%d.pdf", timestamp))

	// 将PDF文件移动到目标目录
	if err := os.Rename(pdfFile, finalPDFPath); err != nil {
		// 如果跨设备移动失败，尝试复制文件
		pdfContent, readErr := os.ReadFile(pdfFile)
		if readErr != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("读取PDF文件失败: %w", readErr)
		}

		if writeErr := os.WriteFile(finalPDFPath, pdfContent, 0644); writeErr != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("写入PDF文件失败: %w", writeErr)
		}
	}

	// 清理临时目录
	os.RemoveAll(tempDir)

	return finalPDFPath, nil
}
