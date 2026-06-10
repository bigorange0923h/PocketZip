<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApp } from '../composables/useApp'

const { getHistory } = useApp()

interface History {
  id: number
  archivePath: string
  outputDir: string
  success: boolean
  usedPassword: boolean
  errorMessage: string
  createdAt: string
}

const histories = ref<History[]>([])
const loading = ref(false)

async function loadHistory() {
  loading.value = true
  try {
    histories.value = await getHistory(50)
  } catch (err) {
    console.error('Failed to load history:', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadHistory)

function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

function truncatePath(path: string, maxLen: number = 40): string {
  if (path.length <= maxLen) return path
  return '...' + path.slice(path.length - maxLen + 3)
}
</script>

<template>
  <div class="history-list">
    <div class="history-header">
      <h3>解压历史</h3>
      <button class="refresh-btn" @click="loadHistory" :disabled="loading">
        {{ loading ? '加载中...' : '刷新' }}
      </button>
    </div>
    <div v-if="histories.length === 0" class="empty">
      暂无解压记录
    </div>
    <div v-else class="history-items">
      <div
        v-for="item in histories"
        :key="item.id"
        class="history-item"
        :class="{ success: item.success, error: !item.success }"
      >
        <div class="item-header">
          <span class="status">{{ item.success ? '✅' : '❌' }}</span>
          <span class="path" :title="item.archivePath">
            {{ truncatePath(item.archivePath) }}
          </span>
          <span class="time">{{ formatTime(item.createdAt) }}</span>
        </div>
        <div class="item-details">
          <span v-if="item.usedPassword" class="tag">使用密码</span>
          <span v-if="item.errorMessage" class="error-msg">
            {{ item.errorMessage }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-list {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.history-header h3 {
  margin: 0;
  font-size: 16px;
  color: #f1f5f9;
}

.refresh-btn {
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 8px;
  padding: 6px 12px;
  color: #6366f1;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(99, 102, 241, 0.2);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.empty {
  padding: 32px;
  text-align: center;
  color: #64748b;
}

.history-items {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  padding: 12px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: background 0.2s ease;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

.item-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status {
  font-size: 14px;
}

.path {
  flex: 1;
  font-size: 13px;
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.time {
  font-size: 11px;
  color: #64748b;
  white-space: nowrap;
}

.item-details {
  margin-top: 4px;
  margin-left: 24px;
  display: flex;
  gap: 8px;
}

.tag {
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
}

.error-msg {
  color: #f43f5e;
  font-size: 11px;
}
</style>
