# QuickCV - AI驱动的智能简历生成器

基于大语言模型的智能简历优化工具，能根据职位要求智能调整简历内容，支持多种优化程度（从诚实到“专家”级别），自动生成 LaTeX 格式简历并编译为PDF。

## 环境要求

### LaTeX环境
- 需要安装支持中文的 LaTeX 环境（推荐使用 TeX Live）
- 验证安装：在命令行运行 `xelatex -v`，如有正确版本信息返回则说明安装成功

## 快速开始

### 后端部署
```bash
cd backend
go mod tidy
go run main.go
```
服务将在 http://localhost:4578 启动

### 前端部署
```bash
cd frontend
npm install
npm run serve
```
访问 http://localhost:8080 即可使用

## 主要功能
- 🤖 AI驱动的简历内容优化
- 📝 多级别简历优化策略
- 🎯 基于职位要求的定向优化
- 📄 自动生成LaTeX格式简历
- 🖨️ 实时PDF预览 