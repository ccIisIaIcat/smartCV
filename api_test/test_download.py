import requests

# API配置
BASE_URL = "http://localhost:8080/api/cv/download"

def download_pdf(filename):
    url = f"{BASE_URL}/pdf/{filename}"
    response = requests.get(url)
    
    print(f"\n下载PDF文件 {filename}:")
    print("状态码:", response.status_code)
    
    if response.status_code == 200:
        output_file = f"download_test_{filename}"
        with open(output_file, "wb") as f:
            f.write(response.content)
        print(f"文件已保存为: {output_file}")
    else:
        print("错误信息:", response.text)

def download_latex(filename):
    url = f"{BASE_URL}/latex/{filename}"
    response = requests.get(url)
    
    print(f"\n下载LaTeX文件 {filename}:")
    print("状态码:", response.status_code)
    
    if response.status_code == 200:
        output_file = f"download_test_{filename}"
        with open(output_file, "wb") as f:
            f.write(response.content)
        print(f"文件已保存为: {output_file}")
    else:
        print("错误信息:", response.text)

# 测试下载
if __name__ == "__main__":
    # 这里需要替换为实际的文件名
    pdf_filename = "cv_1234567890.pdf"
    latex_filename = "cv_1234567890.tex"
    
    download_pdf(pdf_filename)
    download_latex(latex_filename) 