<script setup lang="ts">
import { ref } from 'vue'
import FileSelector from './components/FileSelector.vue'
import LogPanel from './components/LogPanel.vue'
import PasswordDialog from './components/PasswordDialog.vue'
import HistoryList from './components/HistoryList.vue'
import { useApp } from './composables/useApp'

const { extract, extractWithPassword, getPasswordCandidates, onExtractLog } = useApp()

const activeTab = ref<'extract' | 'history'>('extract')
const selectedFile = ref('')
const logs = ref<string[]>([])
const isExtracting = ref(false)
const extractResult = ref<'success' | 'error' | null>(null)
const showPasswordDialog = ref(false)
const passwordCandidates = ref<string[]>([])

function handleFileSelect(path: string) {
  selectedFile.value = path
  logs.value = []
  extractResult.value = null
}

async function handleExtract() {
  if (!selectedFile.value || isExtracting.value) return

  isExtracting.value = true
  logs.value = []
  extractResult.value = null

  const unsub = onExtractLog((line) => {
    logs.value.push(line)
  })

  try {
    await extract(selectedFile.value, '')
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

async function handlePasswordSubmit(password: string) {
  showPasswordDialog.value = false
  isExtracting.value = true

  const unsub = onExtractLog((line) => {
    logs.value.push(line)
  })

  try {
    await extractWithPassword(selectedFile.value, '', password)
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
          :class="{ active: activeTab === 'history' }"
          @click="activeTab = 'history'"
        >
          解压历史
        </button>
      </div>
      <div v-if="activeTab === 'extract'" class="tab-content">
        <FileSelector @select="handleFileSelect" />
        <div v-if="selectedFile" class="action-bar">
          <div class="file-info">已选择: {{ selectedFile }}</div>
          <button
            class="extract-btn"
            :disabled="isExtracting"
            @click="handleExtract"
          >
            {{ isExtracting ? '解压中...' : '开始解压' }}
          </button>
        </div>
        <LogPanel v-if="logs.length > 0" :logs="logs" />
        <div v-if="extractResult" :class="['result', extractResult]">
          {{ extractResult === 'success' ? '✅ 解压成功' : '❌ 解压失败' }}
        </div>
        <PasswordDialog
          v-if="showPasswordDialog"
          :archive-path="selectedFile"
          :candidates="passwordCandidates"
          @submit="handlePasswordSubmit"
          @cancel="handlePasswordCancel"
        />
      </div>
      <div v-else-if="activeTab === 'history'" class="tab-content">
        <HistoryList />
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
</style>
