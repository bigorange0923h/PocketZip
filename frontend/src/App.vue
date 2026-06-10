<script setup lang="ts">
import { ref } from 'vue'
import FileSelector from './components/FileSelector.vue'
import LogPanel from './components/LogPanel.vue'
import PasswordDialog from './components/PasswordDialog.vue'
import HistoryList from './components/HistoryList.vue'
import PasswordManager from './components/PasswordManager.vue'
import ArchivePreview from './components/ArchivePreview.vue'
import ThemeSwitcher from './components/ThemeSwitcher.vue'
import StrategyManager from './components/StrategyManager.vue'
import { useApp } from './composables/useApp'

const {
  extract,
  extractWithPassword,
  extractWithStrategy,
  batchExtract,
  getPasswordCandidates,
  selectDirectory,
  selectFiles,
  openDirectory,
  testArchive,
  onExtractLog
} = useApp()

const activeTab = ref<'extract' | 'history' | 'passwords' | 'preview' | 'settings'>('extract')
const selectedFile = ref('')
const selectedFiles = ref<string[]>([])
const outputDir = ref('')
const logs = ref<string[]>([])
const isExtracting = ref(false)
const extractResult = ref<'success' | 'error' | null>(null)
const showPasswordDialog = ref(false)
const passwordCandidates = ref<string[]>([])
const isBatchMode = ref(false)
const batchResults = ref<any[]>([])
const selectedStrategy = ref('default')

function handleFileSelect(path: string) {
  selectedFile.value = path
  selectedFiles.value = []
  outputDir.value = ''
  logs.value = []
  extractResult.value = null
  batchResults.value = []
}

async function handleSelectFiles() {
  const files = await selectFiles()
  if (files && files.length > 0) {
    selectedFiles.value = files
    selectedFile.value = ''
    outputDir.value = ''
    logs.value = []
    extractResult.value = null
    batchResults.value = []
    isBatchMode.value = true
  }
}

async function handleSelectOutputDir() {
  const dir = await selectDirectory()
  if (dir) {
    outputDir.value = dir
  }
}

async function handleTestArchive() {
  const fileToTest = selectedFile.value || (selectedFiles.value.length > 0 ? selectedFiles.value[0] : '')
  if (!fileToTest) return

  try {
    const isValid = await testArchive(fileToTest)
    if (isValid) {
      alert('✅ 压缩包完整性验证通过')
    } else {
      alert('❌ 压缩包可能已损坏')
    }
  } catch (err) {
    alert(`测试失败: ${err}`)
  }
}

async function handleExtract() {
  if (isBatchMode.value && selectedFiles.value.length > 0) {
    await handleBatchExtract()
    return
  }

  if (!selectedFile.value || isExtracting.value) return

  isExtracting.value = true
  logs.value = []
  extractResult.value = null

  const unsub = onExtractLog((line) => {
    logs.value.push(line)
  })

  try {
    if (selectedStrategy.value !== 'default') {
      await extractWithStrategy(selectedFile.value, selectedStrategy.value)
    } else {
      await extract(selectedFile.value, outputDir.value)
    }
    extractResult.value = 'success'
  } catch (err: any) {
    if (err?.message?.includes('password required') || String(err).includes('password required')) {
      passwordCandidates.value = await getPasswordCandidates(selectedFile.value)
      showPasswordDialog.value = true
    } else {
      extractResult.value = 'error'
      logs.value.push(`错误: ${err}`)
    }
  } finally {
    isExtracting.value = false
    unsub()
  }
}

async function handleBatchExtract() {
  if (selectedFiles.value.length === 0 || isExtracting.value) return

  isExtracting.value = true
  logs.value = []
  extractResult.value = null
  batchResults.value = []

  logs.value.push(`开始批量解压 ${selectedFiles.value.length} 个文件...`)

  try {
    const results = await batchExtract(selectedFiles.value, outputDir.value)
    batchResults.value = results

    const successCount = results.filter((r: any) => r.success).length
    const failCount = results.filter((r: any) => !r.success).length

    logs.value.push(`批量解压完成: ${successCount} 成功, ${failCount} 失败`)

    if (failCount === 0) {
      extractResult.value = 'success'
    } else {
      extractResult.value = 'error'
    }
  } catch (err) {
    extractResult.value = 'error'
    logs.value.push(`批量解压失败: ${err}`)
  } finally {
    isExtracting.value = false
  }
}

async function handlePasswordSubmit(password: string) {
  showPasswordDialog.value = false
  isExtracting.value = true

  const unsub = onExtractLog((line) => {
    logs.value.push(line)
  })

  try {
    await extractWithPassword(selectedFile.value, outputDir.value, password)
    extractResult.value = 'success'
  } catch (err) {
    extractResult.value = 'error'
    logs.value.push(`错误: ${err}`)
  } finally {
    isExtracting.value = false
    unsub()
  }
}

function handlePasswordCancel() {
  showPasswordDialog.value = false
}

async function handleOpenOutputDir() {
  const dir = outputDir.value || (selectedFile.value ? selectedFile.value.replace(/[/\\][^/\\]+$/, '') : '')
  if (dir) {
    await openDirectory(dir)
  }
}
</script>

<template>
  <div id="app">
    <div class="container">
      <h1 class="title">PocketUnzip</h1>
      <div class="tabs">
        <button
          class="tab"
          :class="{ active: activeTab === 'extract' }"
          @click="activeTab = 'extract'"
        >
          解压文件
        </button>
        <button
          class="tab"
          :class="{ active: activeTab === 'preview' }"
          @click="activeTab = 'preview'"
        >
          预览
        </button>
        <button
          class="tab"
          :class="{ active: activeTab === 'history' }"
          @click="activeTab = 'history'"
        >
          解压历史
        </button>
        <button
          class="tab"
          :class="{ active: activeTab === 'passwords' }"
          @click="activeTab = 'passwords'"
        >
          密码库
        </button>
        <button
          class="tab"
          :class="{ active: activeTab === 'settings' }"
          @click="activeTab = 'settings'"
        >
          设置
        </button>
      </div>
      <div v-if="activeTab === 'extract'" class="tab-content">
        <FileSelector @select="handleFileSelect" />

        <div class="mode-switch">
          <button class="mode-btn" @click="handleSelectFiles">
            📦 批量选择压缩包
          </button>
        </div>

        <!-- 策略选择 -->
        <div v-if="selectedFile || selectedFiles.length > 0" class="strategy-select">
          <label class="strategy-label">解压策略:</label>
          <select v-model="selectedStrategy" class="strategy-dropdown">
            <option value="default">默认</option>
            <option value="retry">自动重试</option>
            <option value="auto-open">完成后打开目录</option>
          </select>
        </div>

        <!-- 单文件模式 -->
        <template v-if="!isBatchMode && selectedFile">
          <div class="action-bar">
            <div class="file-info">已选择: {{ selectedFile }}</div>
            <div class="action-buttons">
              <button class="test-btn" @click="handleTestArchive">测试</button>
              <button
                class="extract-btn"
                :disabled="isExtracting"
                @click="handleExtract"
              >
                {{ isExtracting ? '解压中...' : '开始解压' }}
              </button>
            </div>
          </div>
          <div class="output-dir-bar">
            <div class="output-dir-info">
              输出目录: {{ outputDir || '自动创建同名目录' }}
            </div>
            <button class="select-dir-btn" @click="handleSelectOutputDir">
              选择目录
            </button>
          </div>
        </template>

        <!-- 批量模式 -->
        <template v-if="isBatchMode && selectedFiles.length > 0">
          <div class="batch-info">
            <div class="batch-count">已选择 {{ selectedFiles.length }} 个文件</div>
            <div class="batch-files">
              <div v-for="(file, index) in selectedFiles" :key="index" class="batch-file-item">
                {{ file.split(/[/\\]/).pop() }}
              </div>
            </div>
            <div class="action-bar">
              <button class="test-btn" @click="handleTestArchive">测试全部</button>
              <button
                class="extract-btn"
                :disabled="isExtracting"
                @click="handleExtract"
              >
                {{ isExtracting ? '批量解压中...' : '批量解压' }}
              </button>
            </div>
            <div class="output-dir-bar">
              <div class="output-dir-info">
                输出目录: {{ outputDir || '自动创建同名目录' }}
              </div>
              <button class="select-dir-btn" @click="handleSelectOutputDir">
                选择目录
              </button>
            </div>
          </div>
        </template>

        <LogPanel v-if="logs.length > 0" :logs="logs" />

        <div v-if="extractResult" :class="['result', extractResult]">
          <div>{{ extractResult === 'success' ? '✅ 解压成功' : '❌ 解压失败' }}</div>
          <button
            v-if="extractResult === 'success'"
            class="open-dir-btn"
            @click="handleOpenOutputDir"
          >
            📂 打开目录
          </button>
        </div>

        <!-- 批量结果 -->
        <div v-if="batchResults.length > 0" class="batch-results">
          <div class="batch-results-title">解压结果:</div>
          <div v-for="(result, index) in batchResults" :key="index" class="batch-result-item">
            <span :class="['status', result.success ? 'success' : 'error']">
              {{ result.success ? '✅' : '❌' }}
            </span>
            <span class="file-name">{{ result.archivePath.split(/[/\\]/).pop() }}</span>
            <span v-if="result.error" class="error-msg">{{ result.error }}</span>
          </div>
        </div>

        <PasswordDialog
          v-if="showPasswordDialog"
          :archive-path="selectedFile"
          :candidates="passwordCandidates"
          @submit="handlePasswordSubmit"
          @cancel="handlePasswordCancel"
        />
      </div>
      <div v-else-if="activeTab === 'preview'" class="tab-content">
        <ArchivePreview :archive-path="selectedFile" />
      </div>
      <div v-else-if="activeTab === 'history'" class="tab-content">
        <HistoryList />
      </div>
      <div v-else-if="activeTab === 'passwords'" class="tab-content">
        <PasswordManager />
      </div>
      <div v-else-if="activeTab === 'settings'" class="tab-content">
        <ThemeSwitcher />
        <StrategyManager />
      </div>
    </div>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: linear-gradient(135deg, #0a0e1a 0%, #1a1040 50%, #0d1f2d 100%);
  color: #f1f5f9;
  min-height: 100vh;
}

#app {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.container {
  max-width: 600px;
  width: 100%;
  padding: 32px;
}

.title {
  text-align: center;
  font-size: 32px;
  margin-bottom: 32px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 4px;
}

.tab {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #94a3b8;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab:hover {
  color: #f1f5f9;
  background: rgba(255, 255, 255, 0.05);
}

.tab.active {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

.tab-content {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

.action-bar {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.file-info {
  color: #94a3b8;
  font-size: 14px;
}

.extract-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.extract-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4);
}

.extract-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.result {
  margin-top: 16px;
  padding: 12px;
  border-radius: 8px;
  text-align: center;
  font-weight: 600;
}

.result.success {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
}

.result.error {
  background: rgba(244, 63, 94, 0.1);
  color: #f43f5e;
}

.output-dir-bar {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.output-dir-info {
  color: #94a3b8;
  font-size: 13px;
}

.select-dir-btn {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.select-dir-btn:hover {
  background: rgba(99, 102, 241, 0.3);
  border-color: #6366f1;
}

.mode-switch {
  margin-top: 16px;
  text-align: center;
}

.mode-btn {
  background: rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  color: #a5b4fc;
  border-color: rgba(99, 102, 241, 0.3);
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.test-btn {
  background: rgba(234, 179, 8, 0.2);
  color: #eab308;
  border: 1px solid rgba(234, 179, 8, 0.3);
  border-radius: 8px;
  padding: 10px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.test-btn:hover {
  background: rgba(234, 179, 8, 0.3);
  border-color: #eab308;
}

.open-dir-btn {
  margin-top: 8px;
  background: rgba(52, 211, 153, 0.2);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.3);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.open-dir-btn:hover {
  background: rgba(52, 211, 153, 0.3);
  border-color: #34d399;
}

.batch-info {
  margin-top: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.batch-count {
  font-weight: 600;
  color: #a5b4fc;
  margin-bottom: 12px;
}

.batch-files {
  max-height: 150px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.batch-file-item {
  padding: 6px 12px;
  font-size: 12px;
  color: #94a3b8;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 4px;
  margin-bottom: 4px;
}

.batch-results {
  margin-top: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.batch-results-title {
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 12px;
}

.batch-result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.batch-result-item:last-child {
  border-bottom: none;
}

.batch-result-item .status {
  flex-shrink: 0;
}

.batch-result-item .file-name {
  flex: 1;
  font-size: 13px;
  color: #f1f5f9;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.batch-result-item .error-msg {
  font-size: 11px;
  color: #f43f5e;
  flex-shrink: 0;
}

.strategy-select {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.strategy-label {
  font-size: 13px;
  color: #94a3b8;
  flex-shrink: 0;
}

.strategy-dropdown {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  padding: 8px 12px;
  color: #f1f5f9;
  font-size: 13px;
  cursor: pointer;
}

.strategy-dropdown:focus {
  outline: none;
  border-color: #6366f1;
}

.strategy-dropdown option {
  background: #1a1040;
  color: #f1f5f9;
}
</style>
