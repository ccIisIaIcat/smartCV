package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"backend/CV"

	"github.com/gin-gonic/gin"
)

type EditCVRequest struct {
	LaTeX           string `json:"latex" binding:"required"`
	EditRequirement string `json:"edit_requirement" binding:"required"`
	APIKey          string `json:"api_key" binding:"required"`
}

type EditCVResponse struct {
	LaTeX    string `json:"latex"`
	Filename string `json:"filename"`
}

// EditCV 处理修改简历的请求
func EditCV(c *gin.Context) {
	var req EditCVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 打印接收到的参数
	fmt.Printf("收到编辑简历请求参数：\n")
	fmt.Printf("LaTeX长度: %d\n", len(req.LaTeX))
	fmt.Printf("修改要求: %s\n", req.EditRequirement)
	fmt.Printf("API Key长度: %d\n", len(req.APIKey))

	// 修改LaTeX代码
	latexCode, err := CV.EditCV(req.APIKey, req.LaTeX, req.EditRequirement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改简历失败: " + err.Error()})
		return
	}

	// 编译PDF
	pdfPath, err := CV.ParseLatexToPDF(latexCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "编译PDF失败: " + err.Error()})
		return
	}

	// 获取文件名
	filename := filepath.Base(pdfPath)

	c.JSON(http.StatusOK, EditCVResponse{
		LaTeX:    latexCode,
		Filename: filename,
	})
}
