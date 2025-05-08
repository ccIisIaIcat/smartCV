package main

import (
	"log"
	"os"

	"backend/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 temp 目录（如果不存在）
	if err := os.MkdirAll("temp", 0755); err != nil {
		log.Fatalf("创建temp目录失败: %v", err)
	}

	// 创建gin引擎
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4587", "http://localhost:5173", "http://localhost:8080"} // 允许的前端域名
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-API-Key"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// 注册静态文件目录
	r.Static("/temp", "./temp")

	// API路由组
	api := r.Group("/api")
	{
		// CV相关路由组
		cv := api.Group("/cv")
		{
			// POST请求路由
			cv.POST("/generate", handler.GenerateCV) // 生成简历
			cv.POST("/edit", handler.EditCV)         // 修改简历
			cv.POST("/compile", handler.CompileCV)   // 编译LaTeX

			// 下载路由组
			download := cv.Group("/download")
			{
				download.GET("/pdf/:filename", handler.DownloadPDF)     // 下载PDF
				download.GET("/latex/:filename", handler.DownloadLatex) // 下载LaTeX源码
			}
		}
	}

	// 启动服务器
	if err := r.Run(":4578"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
