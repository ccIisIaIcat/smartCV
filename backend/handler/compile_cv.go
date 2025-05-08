package handler

import (
	"net/http"
	"path/filepath"

	"backend/CV"

	"github.com/gin-gonic/gin"
)

type CompileCVRequest struct {
	LaTeX string `json:"latex" binding:"required"`
}

type CompileCVResponse struct {
	Filename string `json:"filename"`
}

// CompileCV 处理编译LaTeX代码的请求
func CompileCV(c *gin.Context) {
	var req CompileCVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 编译PDF
	pdfPath, err := CV.ParseLatexToPDF(req.LaTeX)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "编译PDF失败: " + err.Error()})
		return
	}

	// 获取文件名
	filename := filepath.Base(pdfPath)

	c.JSON(http.StatusOK, CompileCVResponse{
		Filename: filename,
	})
}
