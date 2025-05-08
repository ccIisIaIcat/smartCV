import requests
import json

# API配置
API_URL = "http://localhost:8080/api/cv/compile"

# 测试数据
data = {
    "latex": r"""\documentclass{article}
\usepackage{xeCJK}
\setCJKmainfont{KaiTi}
\begin{document}
\title{测试编译功能}
\author{测试用户}
\maketitle
\section{测试章节}
这是一个用于测试编译功能的简单文档。
\end{document}"""
}

# 发送请求
response = requests.post(API_URL, json=data)

# 打印结果
print("状态码:", response.status_code)
print("\n响应内容:")
print(json.dumps(response.json(), indent=2, ensure_ascii=False))

# 如果成功，下载PDF文件
if response.status_code == 200:
    result = response.json()
    pdf_url = "http://localhost:8080" + result["pdf_url"]
    
    # 获取PDF文件
    pdf_response = requests.get(pdf_url)
    if pdf_response.status_code == 200:
        with open("test_compile_output.pdf", "wb") as f:
            f.write(pdf_response.content)
        print("\nPDF文件已保存到 test_compile_output.pdf") 