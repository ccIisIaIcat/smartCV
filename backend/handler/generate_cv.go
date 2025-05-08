package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"backend/CV"

	"github.com/gin-gonic/gin"
)

type GenerateCVRequest struct {
	Introduction string `json:"introduction" binding:"required"`
	Requirements string `json:"requirements" binding:"required"`
	Level        string `json:"level" binding:"required,oneof=Honest Moderate Aggressive Creative Expert"`
	APIKey       string `json:"api_key" binding:"required"`
}

type GenerateCVResponse struct {
	LaTeX    string `json:"latex"`
	Filename string `json:"filename"`
}

// GenerateCV 处理生成简历的请求
func GenerateCV(c *gin.Context) {
	// 打印原始请求体
	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("原始请求体: %s", string(body))
	// 重新设置请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req GenerateCVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数，level必须是Honest, Moderate, Aggressive, Creative, 或 Expert"})
		return
	}

	log.Printf("收到的请求参数:\n- Level: %v\n- APIKey长度: %d\n- Introduction长度: %d\n- Requirements长度: %d",
		req.Level, len(req.APIKey), len(req.Introduction), len(req.Requirements))

	// 将字符串转换为 PromptLevel
	promptLevel, err := CV.ParsePromptLevel(req.Level)
	if err != nil {
		log.Printf("解析简历级别失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成LaTeX代码
	latexCode, err := CV.GenerateCV(req.APIKey, req.Introduction, req.Requirements, promptLevel)
	if err != nil {
		log.Printf("生成简历失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成简历失败: " + err.Error()})
		return
	}

	// 编译PDF
	pdfPath, err := CV.ParseLatexToPDF(latexCode)
	if err != nil {
		log.Printf("编译PDF失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "编译PDF失败: " + err.Error()})
		return
	}

	// 构建PDF URL
	filename := filepath.Base(pdfPath)

	log.Printf("处理成功:\n- PDF路径: %s\n- 文件名: %s", pdfPath, filename)

	c.JSON(http.StatusOK, GenerateCVResponse{
		LaTeX:    latexCode,
		Filename: filename,
	})
}
