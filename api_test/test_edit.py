import requests
import json

# API配置
API_URL = "http://localhost:8080/api/cv/edit"
API_KEY = "sk-iovSklr9Q3aW95E3gwtCoMMxLDCE61gNhWq71JwFSJKyaJ9b"  # 使用实际的API密钥

# 测试数据
data = {
    "latex": r"""\documentclass{article}
\usepackage{xeCJK}
\setCJKmainfont{KaiTi}
\begin{document}
\section{教育背景}
清华大学计算机科学与技术专业博士（2018-2023）
\section{工作经历}
腾讯云 - 高级后端工程师
\end{document}""",
    
    "edit_requirement": "请添加一个技能部分",
    "api_key": API_KEY
}

# 发送请求
response = requests.post(API_URL, json=data)

# 打印结果
print("状态码:", response.status_code)
print("\n响应内容:")
print(json.dumps(response.json(), indent=2, ensure_ascii=False))

# 如果成功，保存文件
if response.status_code == 200:
    result = response.json()
    
    # 保存LaTeX代码
    with open("test_edit_output.tex", "w", encoding="utf-8") as f:
        f.write(result["latex"])
    print("\nLaTeX代码已保存到 test_edit_output.tex")
    
    # 获取PDF文件
    pdf_url = "http://localhost:8080" + result["pdf_url"]
    pdf_response = requests.get(pdf_url)
    if pdf_response.status_code == 200:
        with open("test_edit_output.pdf", "wb") as f:
            f.write(pdf_response.content)
        print("PDF文件已保存到 test_edit_output.pdf") 