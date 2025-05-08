<template>
  <div class="app-container">
    <el-card class="header-card">
      <h1>Smart CV <span class="chinese">聪明简历</span></h1>
    </el-card>

    <div class="main-content">
      <el-row :gutter="20">
        <!-- 左侧输入区域 -->
        <el-col :span="12">
          <el-card class="input-card">
            <!-- API Key输入 -->
            <div class="input-group">
              <div class="input-header">
                <span>API Key</span>
              </div>
              <el-input
                v-model="apiKey"
                type="password"
                placeholder="请输入CloseAI（大模型国内镜像） API Key"
                clearable
              />
            </div>

            <!-- 简历级别选择 -->
            <div class="input-group">
              <div class="input-header">
                <span>简历级别</span>
              </div>
              <el-select v-model="promptLevel" placeholder="请选择简历优化程度">
                <el-option label="诚实" value="Honest" />
                <el-option label="适中" value="Moderate" />
                <el-option label="积极" value="Aggressive" />
                <el-option label="创意" value="Creative" />
                <el-option label="专家" value="Expert" />
              </el-select>
            </div>

            <!-- 个人介绍输入 -->
            <div class="input-group mb-4">
              <div class="input-header">
                <span>个人介绍</span>
                <el-button type="primary" size="small" @click="showDialog('introduction')">
                  <span>📝 放大编辑</span>
                </el-button>
              </div>
              <el-input
                v-model="introduction"
                type="textarea"
                :rows="6"
                placeholder="请输入个人介绍...（什么样的都行，从就简历把文本复制过来也行，不需要排版）"
              />
            </div>

            <!-- 工作要求输入 -->
            <div class="input-group mb-4">
              <div class="input-header">
                <span>工作要求</span>
                <el-button type="primary" size="small" @click="showDialog('requirements')">
                  <span>📝 放大编辑</span>
                </el-button>
              </div>
              <el-input
                v-model="jobRequirements"
                type="textarea"
                :rows="4"
                placeholder="请输入目标职位要求...（把jd贴过来就行）"
              />
            </div>

            <!-- LaTeX代码编辑器 -->
            <div class="input-group mb-4">
              <div class="input-header">
                <span>LaTeX代码</span>
                <div class="button-group-inline">
                  <el-button type="primary" size="small" @click="showDialog('latex')">
                    <span>📝 放大编辑</span>
                  </el-button>
                  <el-button type="primary" size="small" @click="compileLatex">
                    <span>📝 编译预览</span>
                  </el-button>
                </div>
              </div>
              <el-input
                v-model="latexCode"
                type="textarea"
                :rows="8"
                placeholder="LaTeX代码将在这里生成..."
              />
            </div>

            <!-- 修改建议输入 -->
            <div class="input-group mb-4">
              <div class="input-header">
                <span>修改建议</span>
                <el-button type="primary" size="small" @click="showDialog('modifications')">
                  <span>📝 放大编辑</span>
                </el-button>
              </div>
              <el-input
                v-model="modifications"
                type="textarea"
                :rows="4"
                placeholder="请输入需要修改的内容..."
              />
            </div>

            <!-- 操作按钮 -->
            <div class="button-group">
              <el-button type="primary" @click="generateResume">
                生成简历
              </el-button>
              <el-button type="success" @click="updateResume">
                更新简历
              </el-button>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧预览区域 -->
        <el-col :span="12">
          <el-card class="preview-card">
            <template #header>
              <div class="preview-header">
                <span>预览</span>
              </div>
            </template>
            <div class="preview-container" v-loading="loading">
              <el-empty v-if="!previewUrl" description="预览将显示在这里..." />
              <iframe
                v-else
                :src="previewUrl"
                type="application/pdf"
                width="100%"
                height="800"
                class="preview-frame"
                @load="() => console.log('iframe加载完成')"
                @error="(e) => console.error('iframe加载失败:', e)"
              ></iframe>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 弹出式对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="80%"
      :close-on-click-modal="false"
    >
      <el-input
        v-model="dialogContent"
        type="textarea"
        :rows="20"
        :placeholder="dialogPlaceholder"
      />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveDialogContent">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// API 基础URL
const API_BASE_URL = 'http://localhost:4578'
console.log('初始化 API_BASE_URL:', API_BASE_URL)

// 数据
const apiKey = ref('')
const introduction = ref('')
const jobRequirements = ref('')
const latexCode = ref('')
const modifications = ref('')
const previewUrl = ref(`${API_BASE_URL}/temp/check.pdf`) // 设置默认预览URL
const loading = ref(false)
const promptLevel = ref('Moderate') // 默认使用适中级别

// 弹出框相关数据
const dialogVisible = ref(false)
const dialogTitle = ref('')
const dialogContent = ref('')
const dialogPlaceholder = ref('')
const currentEditType = ref('')

// 监听预览URL变化
watch(previewUrl, (newUrl) => {
  console.log('预览URL发生变化:', newUrl)
})

// 页面加载完成时
onMounted(() => {
  console.log('页面加载完成')
  console.log('当前预览URL:', previewUrl.value)
  // 测试PDF文件是否可访问
  fetch(previewUrl.value)
    .then(response => {
      console.log('测试PDF访问状态:', response.status, response.statusText)
      if (!response.ok) {
        console.error('PDF文件不可访问')
        // 如果PDF不可访问，清空预览URL
        previewUrl.value = ''
      }
    })
    .catch(error => {
      console.error('测试PDF访问失败:', error)
      // 如果发生错误，清空预览URL
      previewUrl.value = ''
    })
})

// API调用工具函数
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

api.interceptors.request.use(config => {
  if (apiKey.value) {
    config.headers['X-API-Key'] = apiKey.value
  }
  return config
})

// 生成简历
const generateResume = async () => {
  if (!introduction.value || !jobRequirements.value || !apiKey.value) {
    ElMessage.warning('请填写个人介绍、职位要求和API Key')
    return
  }

  loading.value = true
  try {
    console.log('发送请求到:', `${API_BASE_URL}/api/cv/generate`)
    const response = await api.post('/api/cv/generate', {
      introduction: introduction.value,
      requirements: jobRequirements.value,
      level: promptLevel.value,
      api_key: apiKey.value
    })

    console.log('后端响应数据:', response.data)
    latexCode.value = response.data.latex
    
    // 使用后端返回的文件名构建预览URL
    if (response.data.filename) {
      const pdfUrl = `${API_BASE_URL}/temp/${response.data.filename}`
      console.log('设置预览URL:', pdfUrl)
      previewUrl.value = pdfUrl
    } else {
      console.error('后端未返回文件名')
      ElMessage.warning('未获取到PDF文件名')
    }

    ElMessage.success('简历生成成功')
  } catch (error) {
    console.error('生成简历失败:', error.response || error)
    ElMessage.error(error.response?.data?.error || '生成简历失败')
  } finally {
    loading.value = false
  }
}

// 更新简历
const updateResume = async () => {
  if (!latexCode.value || !modifications.value || !apiKey.value) {
    ElMessage.warning('请确保有LaTeX代码、修改建议和API Key')
    return
  }

  loading.value = true
  try {
    const response = await api.post('/api/cv/edit', {
      latex: latexCode.value,
      edit_requirement: modifications.value,
      api_key: apiKey.value
    })

    latexCode.value = response.data.latex
    if (response.data.filename) {
      const pdfUrl = `${API_BASE_URL}/temp/${response.data.filename}`
      console.log('设置预览URL:', pdfUrl)
      previewUrl.value = pdfUrl
    } else {
      console.error('后端未返回文件名')
      ElMessage.warning('未获取到PDF文件名')
    }
    ElMessage.success('简历更新成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '更新简历失败')
  } finally {
    loading.value = false
  }
}

// 编译LaTeX并预览
const compileLatex = async () => {
  if (!latexCode.value) {
    ElMessage.warning('请先生成或输入LaTeX代码')
    return
  }
  await updatePreview()
}

// 更新预览
const updatePreview = async () => {
  loading.value = true
  try {
    const response = await api.post('/api/cv/compile', {
      latex: latexCode.value
    })

    if (response.data.filename) {
      const pdfUrl = `${API_BASE_URL}/temp/${response.data.filename}`
      console.log('设置预览URL:', pdfUrl)
      previewUrl.value = pdfUrl
    } else {
      console.error('后端未返回文件名')
      ElMessage.warning('未获取到PDF文件名')
    }
  } catch (error) {
    ElMessage.error('预览生成失败')
  } finally {
    loading.value = false
  }
}

// 下载PDF
const downloadPDF = async () => {
  if (!latexCode.value) {
    ElMessage.warning('请先生成简历')
    return
  }

  loading.value = true
  try {
    const response = await api.post('/api/cv/compile', {
      latex: latexCode.value
    }, {
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'resume.pdf')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    ElMessage.error('PDF下载失败')
  } finally {
    loading.value = false
  }
}

// 下载LaTeX源码
const downloadLatex = () => {
  if (!latexCode.value) {
    ElMessage.warning('请先生成简历')
    return
  }

  const blob = new Blob([latexCode.value], { type: 'text/plain' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', 'resume.tex')
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 显示弹出框
const showDialog = (type) => {
  currentEditType.value = type
  switch (type) {
    case 'introduction':
      dialogTitle.value = '编辑个人介绍'
      dialogContent.value = introduction.value
      dialogPlaceholder.value = '请输入个人介绍...'
      break
    case 'requirements':
      dialogTitle.value = '编辑工作要求'
      dialogContent.value = jobRequirements.value
      dialogPlaceholder.value = '请输入目标职位要求...'
      break
    case 'latex':
      dialogTitle.value = '编辑LaTeX代码'
      dialogContent.value = latexCode.value
      dialogPlaceholder.value = 'LaTeX代码将在这里生成...'
      break
    case 'modifications':
      dialogTitle.value = '编辑修改建议'
      dialogContent.value = modifications.value
      dialogPlaceholder.value = '请输入需要修改的内容...'
      break
  }
  dialogVisible.value = true
}

// 保存弹出框内容
const saveDialogContent = () => {
  switch (currentEditType.value) {
    case 'introduction':
      introduction.value = dialogContent.value
      break
    case 'requirements':
      jobRequirements.value = dialogContent.value
      break
    case 'latex':
      latexCode.value = dialogContent.value
      break
    case 'modifications':
      modifications.value = dialogContent.value
      break
  }
  dialogVisible.value = false
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  min-height: 100vh;
  background-color: #1e1e2d;
  color: #e0e0e0;
}

.header-card {
  margin-bottom: 20px;
  background-color: #2b2b3c;
  border: none;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  padding: 24px 0;
  position: relative;
  overflow: hidden;
}

.header-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #4f46e5, #60a5fa, #4f46e5);
  background-size: 200% 100%;
  animation: shimmer 2s infinite linear;
}

.header-card h1 {
  margin: 0;
  text-align: center;
  color: transparent;
  background: linear-gradient(135deg, #60a5fa 0%, #4f46e5 100%);
  -webkit-background-clip: text;
  background-clip: text;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  font-size: 2.5em;
  font-weight: 700;
  letter-spacing: 2px;
  text-transform: uppercase;
  position: relative;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
}

.header-card h1 span.chinese {
  font-family: 'Microsoft YaHei', '微软雅黑', sans-serif;
  font-weight: 500;
  margin-left: 10px;
  font-size: 0.9em;
}

@keyframes shimmer {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}

.main-content {
  margin-top: 20px;
}

.input-card,
.preview-card {
  height: 100%;
  background-color: #2b2b3c;
  border: 1px solid #3d3d56;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.mb-4 {
  margin-bottom: 16px;
}

.input-group {
  background: #2f2f45;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #3d3d56;
}

.input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding: 8px 12px;
  background-color: #323248;
  border-radius: 8px;
  border-bottom: 2px solid #4f46e5;
}

.input-header span {
  font-size: 15px;
  font-weight: 500;
  color: #e0e0e0;
}

:deep(.el-select) {
  width: 200px;
  margin-top: 8px;
}

:deep(.el-select .el-input) {
  margin: 0;
}

:deep(.el-select .el-input__wrapper) {
  background-color: #323248 !important;
  border: 1px solid #3d3d56 !important;
  box-shadow: none !important;
  --el-input-bg-color: #323248 !important;
  border-radius: 8px;
  padding: 0 12px !important;
  height: 40px !important;
  margin: 0 !important;
  width: 200px !important;
}

:deep(.el-select .el-input__inner) {
  height: 40px !important;
  line-height: 40px !important;
  color: #e0e0e0 !important;
  font-size: 14px !important;
}

:deep(.el-select-dropdown) {
  background-color: #2b2b3c !important;
  border: 1px solid #3d3d56 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2) !important;
  margin-top: 4px !important;
}

:deep(.el-select-dropdown__item) {
  color: #e0e0e0 !important;
  padding: 8px 16px !important;
  height: 36px !important;
  line-height: 20px !important;
}

:deep(.el-select-dropdown__item.hover) {
  background-color: #323248 !important;
}

:deep(.el-select-dropdown__item.selected) {
  background-color: #4f46e5 !important;
  color: #ffffff !important;
  font-weight: 500 !important;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__inner) {
  background-color: #323248 !important;
  border: 1px solid #3d3d56 !important;
  box-shadow: none !important;
  --el-input-bg-color: #323248 !important;
  color: #e0e0e0 !important;
  border-radius: 8px;
  margin-top: 8px;
}

.button-group {
  background-color: #2f2f45;
  border-radius: 12px;
  padding: 20px;
  margin-top: 20px;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #3d3d56;
  display: flex;
  gap: 12px;
}

.button-group-inline {
  display: flex;
  gap: 12px;
  align-items: center;
}

.preview-card {
  height: 100%;
  background-color: #2b2b3c;
  border: 1px solid #3d3d56;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.preview-container {
  flex: 1;
  min-height: 600px;
  border: 1px solid #3d3d56;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: stretch;
  background: #323248;
  overflow: hidden;
}

.preview-frame {
  border: none;
  width: 100%;
  height: 100%;
  min-height: 800px;
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

:deep(.el-card__body) {
  height: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  margin-bottom: 16px;
}

:deep(.el-empty) {
  width: 100%;
  padding: 40px;
  background: #323248;
  border-radius: 8px;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__inner) {
  background-color: #323248 !important;
  border: 1px solid #3d3d56 !important;
  box-shadow: none !important;
  --el-input-bg-color: #323248 !important;
  color: #e0e0e0 !important;
}

:deep(.el-input__inner) {
  height: 40px !important;
  line-height: 40px !important;
  color: #e0e0e0 !important;
  font-size: 14px !important;
}

:deep(.el-select .el-select__tags) {
  background-color: transparent !important;
}

:deep(.el-select-dropdown) {
  background-color: #2b2b3c !important;
  border: 1px solid #3d3d56 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2) !important;
  margin-top: 4px !important;
}

:deep(.el-select-dropdown__item) {
  color: #e0e0e0 !important;
  padding: 8px 16px !important;
  height: 36px !important;
  line-height: 20px !important;
}

:deep(.el-select-dropdown__item.hover) {
  background-color: #323248 !important;
}

:deep(.el-select-dropdown__item.selected) {
  background-color: #4f46e5 !important;
  color: #ffffff !important;
  font-weight: 500 !important;
}

:deep(.el-input.el-input--large) {
  margin: 0 !important;
  padding: 0 !important;
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
}

.prompt-level-container {
  background: #2f2f45;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #3d3d56;
}

.prompt-level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 12px;
  background-color: #323248;
  border-radius: 8px;
  border-bottom: 2px solid #4f46e5;
}

.prompt-level-header span {
  font-size: 15px;
  font-weight: 500;
  color: #e0e0e0;
}

.input-card {
  padding: 20px;
  background-color: #2b2b3c;
}

:deep(.el-select) {
  margin-bottom: 20px;
  background: #2f2f45;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #3d3d56;
}

:deep(.el-input.el-input--large) {
  margin-bottom: 20px;
  background: #2f2f45;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #3d3d56;
}

:deep(.el-input__wrapper) {
  box-shadow: none !important;
}

:deep(.el-dialog) {
  background-color: #2b2b3c;
  border: 1px solid #3d3d56;
  border-radius: 8px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #3d3d56;
  padding: 20px;
  margin-right: 0;
}

:deep(.el-dialog__title) {
  color: #e0e0e0;
  font-size: 18px;
  font-weight: 500;
}

:deep(.el-dialog__body) {
  padding: 20px;
  background-color: #2b2b3c;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #3d3d56;
  padding: 20px;
}
</style> 