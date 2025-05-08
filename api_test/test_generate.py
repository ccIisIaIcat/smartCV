import requests
import json
from enum import Enum

# API配置
API_URL = "http://localhost:8080/api/cv/generate"
API_KEY = "sk-iovSklr9Q3aW95E3gwtCoMMxLDCE61gNhWq71JwFSJKyaJ9b"  # 使用实际的API密钥

# 简历级别枚举
class ResumeLevel(str, Enum):
    HONEST = "Honest"
    MODERATE = "Moderate"
    AGGRESSIVE = "Aggressive"
    CREATIVE = "Creative"
    EXPERT = "Expert"

# 测试数据
data = {
    "introduction": """教育背景：
- 清华大学计算机科学与技术专业博士（2018-2023）
- 北京大学计算机科学与技术专业硕士（2015-2018）

工作经历：
1. 腾讯云 - 高级后端工程师（2023至今）
   - 负责云原生微服务平台核心组件开发
   - 优化系统性能，将服务响应时间降低40%""",
    
    "requirements": """职位要求：
1. 精通 Go 语言开发
2. 熟悉微服务架构
3. 有性能优化经验""",
    
    "level": ResumeLevel.HONEST,  # 使用诚实级别
    "api_key": API_KEY
}

# 发送请求
response = requests.post(API_URL, json=data)

# 打印结果
print("状态码:", response.status_code)
print("\n响应内容:")
print(json.dumps(response.json(), indent=2, ensure_ascii=False))

# 如果成功，保存并打印文件
if response.status_code == 200:
    result = response.json()
    
    # 打印LaTeX代码
    print("\nLaTeX代码:")
    print(result["latex"])
    
    # 保存LaTeX代码
    with open("test_output.tex", "w", encoding="utf-8") as f:
        f.write(result["latex"])
    print("\nLaTeX代码已保存到 test_output.tex")
    
    # 获取PDF文件
    pdf_url = "http://localhost:8080" + result["pdf_url"]
    pdf_response = requests.get(pdf_url)
    if pdf_response.status_code == 200:
        with open("test_output.pdf", "wb") as f:
            f.write(pdf_response.content)
        print("PDF文件已保存到 test_output.pdf") 