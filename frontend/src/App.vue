<script setup lang="ts">
import { ref, computed } from 'vue'
import FileSelector from './components/FileSelector.vue'
import LogPanel from './components/LogPanel.vue'
import PasswordDialog from './components/PasswordDialog.vue'
import HistoryList from './components/HistoryList.vue'
import PasswordManager from './components/PasswordManager.vue'
import ThemeSwitcher from './components/ThemeSwitcher.vue'
import StrategyManager from './components/StrategyManager.vue'
import InlinePasswordPanel from './components/InlinePasswordPanel.vue'
import ProgressOverlay from './components/ProgressOverlay.vue'
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
  onExtractLog,
  compress,
  selectFilesForCompress,
  selectFolderForCompress,
  selectSavePath,
  onCompressLog
} = useApp()

// --- 状态 ---
const activeTab = ref<'work' | 'history' | 'passwords' | 'settings'>('work')
const mode = ref<'compress' | 'extract'>('extract')

// 解压相关
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

// 压缩相关
const compressFiles = ref<string[]>([])
const archiveName = ref('archive')
const compressFormat = ref('zip')
const isCompressing = ref(false)
const compressResult = ref<'success' | 'error' | null>(null)

// 共享
const password = ref('')
const rememberPassword = ref(false)

// --- 计算属性 ---
const formats = ['zip', '7z', 'tar', 'gz']

const hasSelectedFiles = computed(() => {
  if (mode.value === 'extract') {
    return selectedFile.value || selectedFiles.value.length > 0
  }
  return compressFiles.value.length > 0
})

const actionLabel = computed(() => {
  if (mode.value === 'compress') {
    const count = compressFiles.value.length
    return isCompressing.value ? '压缩中...' : `开始压缩${count ? ` (${count})` : ''}`
  }
  const count = isBatchMode.value ? selectedFiles.value.length : (selectedFile.value ? 1 : 0)
  return isExtracting.value ? '解压中...' : `开始解压${count ? ` (${count})` : ''}`
})

const isBusy = computed(() => isExtracting.value || isCompressing.value)

// --- 解压逻辑 ---
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
    } else if (password.value) {
      await extractWithPassword(selectedFile.value, outputDir.value, password.value)
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

async function handlePasswordSubmit(pwd: string) {
  showPasswordDialog.value = false
  isExtracting.value = true

  const unsub = onExtractLog((line) => {
    logs.value.push(line)
  })

  try {
    await extractWithPassword(selectedFile.value, outputDir.value, pwd)
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

// --- 压缩逻辑 ---
async function handleSelectCompressFiles() {
  const files = await selectFilesForCompress()
  if (files && files.length > 0) {
    compressFiles.value = [...compressFiles.value, ...files]
    compressResult.value = null
    logs.value = []
  }
}

async function handleSelectCompressFolder() {
  const folder = await selectFolderForCompress()
  if (folder) {
    compressFiles.value = [...compressFiles.value, folder]
    compressResult.value = null
    logs.value = []
  }
}

function handleRemoveCompressFile(index: number) {
  compressFiles.value.splice(index, 1)
}

function handleClearCompressFiles() {
  compressFiles.value = []
}

function getFileName(path: string): string {
  return path.split(/[/\\]/).pop() || path
}

async function handleCompress() {
  if (compressFiles.value.length === 0 || isCompressing.value) return

  // 确定保存路径
  let savePath = ''
  const ext = compressFormat.value === 'gz' ? '.tar.gz' : `.${compressFormat.value}`
  const defaultName = archiveName.value ? `${archiveName.value}${ext}` : `archive${ext}`

  savePath = await selectSavePath(defaultName)
  if (!savePath) return

  isCompressing.value = true
  logs.value = []
  compressResult.value = null

  const unsub = onCompressLog((line) => {
    logs.value.push(line)
  })

  try {
    await compress(compressFiles.value, savePath, compressFormat.value, password.value)
    compressResult.value = 'success'
    logs.value.push(`✅ 压缩完成: ${savePath}`)
  } catch (err) {
    compressResult.value = 'error'
    logs.value.push(`❌ 压缩失败: ${err}`)
  } finally {
    isCompressing.value = false
    unsub()
  }
}

function handleAction() {
  if (mode.value === 'compress') {
    handleCompress()
  } else {
    handleExtract()
  }
}

function handleCloseOverlay() {
  extractResult.value = null
  compressResult.value = null
  logs.value = []
}

function handleCancelOperation() {
  // 取消操作（目前无法真正取消 7z 进程，仅重置状态）
  isExtracting.value = false
  isCompressing.value = false
  logs.value.push('⚠️ 操作已取消')
}
</script>

<template>
  <div id="app">
    <!-- 顶部 Header -->
    <header class="app-header">
      <div class="header-inner">
        <div class="header-left">
          <span class="app-icon">
            <span class="icon-emoji">📦</span>
          </span>
          <div class="app-info">
            <span class="app-name">PocketZip</span>
            <span class="app-subtitle">文件压缩 / 解压工具</span>
          </div>
        </div>

        <nav class="header-nav">
          <button
            class="nav-btn"
            :class="{ active: activeTab === 'work' }"
            @click="activeTab = 'work'"
          >
            {{ mode === 'compress' ? '压缩文件' : '解压文件' }}
          </button>
          <button
            class="nav-btn"
            :class="{ active: activeTab === 'history' }"
            @click="activeTab = 'history'"
          >
            历史记录
          </button>
          <button
            class="nav-btn"
            :class="{ active: activeTab === 'passwords' }"
            @click="activeTab = 'passwords'"
          >
            密码库
          </button>
          <button
            class="nav-btn"
            :class="{ active: activeTab === 'settings' }"
            @click="activeTab = 'settings'"
          >
            设置
          </button>
        </nav>

        <button
          class="settings-btn"
          :class="{ active: activeTab === 'settings' }"
          @click="activeTab = 'settings'"
          title="设置"
        >
          ⚙️
        </button>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 工作区 Tab（压缩/解压） -->
      <div v-if="activeTab === 'work'" class="tab-content">
        <!-- 标题区 + 模式切换 -->
        <div class="section-header">
          <div class="section-header-row">
            <div>
              <h1 class="page-title">{{ mode === 'compress' ? '压缩文件' : '解压文件' }}</h1>
              <p class="page-subtitle">
                {{ mode === 'compress' ? '选择文件或文件夹，打包为压缩包' : '拖入压缩包或选择文件，设置密码后一键提取' }}
              </p>
            </div>
            <div class="mode-tabs">
              <button
                class="mode-tab"
                :class="{ active: mode === 'compress' }"
                @click="mode = 'compress'"
              >
                📁 压缩
              </button>
              <button
                class="mode-tab"
                :class="{ active: mode === 'extract' }"
                @click="mode = 'extract'"
              >
                📦 解压
              </button>
            </div>
          </div>
        </div>

        <!-- 双栏布局 -->
        <div class="content-grid">
          <!-- 左侧主区 -->
          <div class="left-column">

            <!-- ===== 压缩模式 ===== -->
            <template v-if="mode === 'compress'">
              <!-- 文件选择区 -->
              <div class="compress-picker">
                <button class="picker-btn" @click="handleSelectCompressFiles">
                  📄 选择文件
                </button>
                <button class="picker-btn" @click="handleSelectCompressFolder">
                  📂 选择文件夹
                </button>
              </div>

              <!-- 已选文件列表 -->
              <div v-if="compressFiles.length > 0" class="file-list-card">
                <div class="file-list-header">
                  <span class="file-list-count">已选 {{ compressFiles.length }} 个</span>
                  <button class="clear-btn" @click="handleClearCompressFiles">清空</button>
                </div>
                <div class="file-list-body">
                  <div v-for="(file, index) in compressFiles" :key="index" class="file-list-item">
                    <span class="file-list-icon">{{ file.includes('.') ? '📄' : '📂' }}</span>
                    <span class="file-list-name">{{ getFileName(file) }}</span>
                    <button class="remove-btn" @click="handleRemoveCompressFile(index)">✕</button>
                  </div>
                </div>
              </div>

              <!-- 压缩选项 -->
              <div v-if="compressFiles.length > 0" class="compress-options">
                <div class="option-row">
                  <label class="option-label">归档名称</label>
                  <input
                    v-model="archiveName"
                    class="option-input"
                    placeholder="archive"
                  />
                </div>
                <div class="option-row">
                  <label class="option-label">压缩格式</label>
                  <div class="format-buttons">
                    <button
                      v-for="f in formats"
                      :key="f"
                      class="format-btn"
                      :class="{ active: compressFormat === f }"
                      @click="compressFormat = f"
                    >
                      {{ f }}
                    </button>
                  </div>
                </div>
              </div>
            </template>

            <!-- ===== 解压模式 ===== -->
            <template v-else>
              <FileSelector @select="handleFileSelect" />

              <div v-if="isBatchMode" class="mode-switch">
                <button class="mode-btn" @click="handleSelectFiles">
                  📦 批量选择压缩包
                </button>
              </div>

              <!-- 策略选择 -->
              <div v-if="hasSelectedFiles" class="strategy-select">
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
                  <div class="file-info">已选择: {{ selectedFile.split(/[/\\]/).pop() }}</div>
                  <div class="action-buttons">
                    <button class="test-btn" @click="handleTestArchive">测试</button>
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
            </template>

            <!-- 共用：日志面板 -->
            <LogPanel v-if="logs.length > 0" :logs="logs" />

            <!-- 解压结果 -->
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

            <!-- 压缩结果 -->
            <div v-if="compressResult" :class="['result', compressResult]">
              {{ compressResult === 'success' ? '✅ 压缩成功' : '❌ 压缩失败' }}
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
          </div>

          <!-- 右侧面板 -->
          <div class="right-column">
            <!-- 密码面板 -->
            <div class="panel-card">
              <InlinePasswordPanel
                v-model:password="password"
                v-model:remember="rememberPassword"
              />
            </div>

            <!-- 操作按钮 -->
            <button
              class="action-btn"
              :disabled="!hasSelectedFiles || isBusy"
              @click="handleAction"
            >
              <span class="btn-icon">✨</span>
              {{ actionLabel }}
            </button>
          </div>
        </div>
      </div>

      <!-- 历史记录 Tab -->
      <div v-else-if="activeTab === 'history'" class="tab-content">
        <div class="section-header">
          <div>
            <h1 class="page-title">历史记录</h1>
            <p class="page-subtitle">查看过往操作记录</p>
          </div>
        </div>
        <div class="single-column">
          <HistoryList />
        </div>
      </div>

      <!-- 密码库 Tab -->
      <div v-else-if="activeTab === 'passwords'" class="tab-content">
        <div class="section-header">
          <div>
            <h1 class="page-title">密码库</h1>
            <p class="page-subtitle">管理已保存的解压密码</p>
          </div>
        </div>
        <div class="single-column">
          <PasswordManager />
        </div>
      </div>

      <!-- 设置 Tab -->
      <div v-else-if="activeTab === 'settings'" class="tab-content">
        <div class="section-header">
          <div>
            <h1 class="page-title">设置</h1>
            <p class="page-subtitle">自定义行为和外观</p>
          </div>
        </div>
        <div class="single-column">
          <ThemeSwitcher />
          <StrategyManager />
        </div>
      </div>
    </main>

    <!-- 密码弹窗（仅在需要时显示） -->
    <PasswordDialog
      v-if="showPasswordDialog"
      :archive-path="selectedFile"
      :candidates="passwordCandidates"
      @submit="handlePasswordSubmit"
      @cancel="handlePasswordCancel"
    />

    <!-- 进度弹窗 -->
    <ProgressOverlay
      v-if="isBusy || extractResult || compressResult"
      :mode="mode"
      :logs="logs"
      :is-running="isBusy"
      :result="mode === 'compress' ? compressResult : extractResult"
      @close="handleCloseOverlay"
      @cancel="handleCancelOperation"
      @open-dir="handleOpenOutputDir"
    />
  </div>
</template>

<style>
/* 全局样式 */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: linear-gradient(135deg, #0a0e1a 0%, #1a1040 50%, #0d1f2d 100%);
  color: #f1f5f9;
  min-height: 100vh;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header */
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 26, 0.8);
  backdrop-filter: blur(12px);
}

.header-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  gap: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.app-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

.icon-emoji {
  font-size: 18px;
}

.app-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.app-name {
  font-size: 14px;
  font-weight: 600;
  color: #f1f5f9;
}

.app-subtitle {
  font-size: 11px;
  color: #64748b;
}

/* Header Navigation */
.header-nav {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  padding: 3px;
}

.nav-btn {
  padding: 7px 14px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.nav-btn:hover {
  color: #f1f5f9;
  background: rgba(255, 255, 255, 0.05);
}

.nav-btn.active {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

.settings-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  padding: 7px 10px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.settings-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.settings-btn.active {
  background: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.3);
}

/* Main Content */
.main-content {
  flex: 1;
  max-width: 1100px;
  margin: 0 auto;
  width: 100%;
  padding: 24px;
}

.tab-content {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Section Header */
.section-header {
  margin-bottom: 24px;
}

.section-header-row {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (min-width: 640px) {
  .section-header-row {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #f1f5f9;
  margin: 0 0 4px 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 0;
}

/* Mode Tabs */
.mode-tabs {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  padding: 3px;
}

.mode-tab {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.mode-tab:hover {
  color: #f1f5f9;
  background: rgba(255, 255, 255, 0.05);
}

.mode-tab.active {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

/* Content Grid */
.content-grid {
  display: grid;
  gap: 24px;
  grid-template-columns: 1fr;
}

@media (min-width: 900px) {
  .content-grid {
    grid-template-columns: 1fr 340px;
  }
}

/* 左侧主区 */
.left-column {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* 右侧面板 */
.right-column {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.panel-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  padding: 18px;
}

/* 单栏布局 */
.single-column {
  max-width: 700px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 压缩文件选择器 */
.compress-picker {
  display: flex;
  gap: 12px;
}

.picker-btn {
  flex: 1;
  padding: 24px 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 2px dashed rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  color: #94a3b8;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.picker-btn:hover {
  background: rgba(99, 102, 241, 0.05);
  border-color: rgba(99, 102, 241, 0.3);
  color: #a5b4fc;
}

/* 文件列表卡片 */
.file-list-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  overflow: hidden;
}

.file-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.file-list-count {
  font-size: 13px;
  font-weight: 600;
  color: #a5b4fc;
}

.clear-btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  color: #f43f5e;
  background: rgba(244, 63, 94, 0.1);
}

.file-list-body {
  max-height: 200px;
  overflow-y: auto;
}

.file-list-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  transition: background 0.15s ease;
}

.file-list-item:last-child {
  border-bottom: none;
}

.file-list-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.file-list-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.file-list-name {
  flex: 1;
  font-size: 13px;
  color: #f1f5f9;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.remove-btn:hover {
  color: #f43f5e;
  background: rgba(244, 63, 94, 0.1);
}

/* 压缩选项 */
.compress-options {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
}

.option-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.option-label {
  font-size: 13px;
  color: #94a3b8;
  flex-shrink: 0;
  min-width: 70px;
}

.option-input {
  flex: 1;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  color: #f1f5f9;
  font-size: 13px;
  outline: none;
  transition: border-color 0.2s ease;
}

.option-input:focus {
  border-color: #6366f1;
}

.format-buttons {
  display: flex;
  gap: 6px;
  flex: 1;
}

.format-btn {
  flex: 1;
  padding: 7px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 0.2s ease;
}

.format-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.3);
  color: #a5b4fc;
}

.format-btn.active {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-color: transparent;
  color: white;
}

/* 模式切换 */
.mode-switch {
  text-align: center;
}

.mode-btn {
  background: rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 10px;
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

/* 策略选择 */
.strategy-select {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
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
  border-radius: 8px;
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

/* Action Bar */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
}

.file-info {
  color: #94a3b8;
  font-size: 13px;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.test-btn {
  background: rgba(234, 179, 8, 0.15);
  color: #eab308;
  border: 1px solid rgba(234, 179, 8, 0.25);
  border-radius: 8px;
  padding: 8px 14px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.test-btn:hover {
  background: rgba(234, 179, 8, 0.25);
  border-color: #eab308;
}

/* Output Dir Bar */
.output-dir-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 10px;
}

.output-dir-info {
  color: #64748b;
  font-size: 12px;
}

.select-dir-btn {
  background: rgba(99, 102, 241, 0.15);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.25);
  border-radius: 6px;
  padding: 5px 10px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.select-dir-btn:hover {
  background: rgba(99, 102, 241, 0.25);
  border-color: #6366f1;
}

/* Action Button */
.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px 24px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 16px;
}

/* Result */
.result {
  padding: 14px;
  border-radius: 10px;
  text-align: center;
  font-weight: 600;
  font-size: 14px;
}

.result.success {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.2);
}

.result.error {
  background: rgba(244, 63, 94, 0.1);
  color: #f43f5e;
  border: 1px solid rgba(244, 63, 94, 0.2);
}

.open-dir-btn {
  margin-top: 10px;
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.25);
  border-radius: 8px;
  padding: 7px 14px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.open-dir-btn:hover {
  background: rgba(52, 211, 153, 0.25);
  border-color: #34d399;
}

/* Batch Info */
.batch-info {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
}

.batch-count {
  font-weight: 600;
  color: #a5b4fc;
  margin-bottom: 12px;
  font-size: 14px;
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
  border-radius: 6px;
  margin-bottom: 4px;
}

/* Batch Results */
.batch-results {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
}

.batch-results-title {
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 12px;
  font-size: 14px;
}

.batch-result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
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

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
