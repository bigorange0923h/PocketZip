<script setup lang="ts">
import { ref, computed, onUnmounted, watch } from 'vue'

const props = defineProps<{
  mode: 'compress' | 'extract'
  logs: string[]
  isRunning: boolean
  result: 'success' | 'error' | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'cancel'): void
  (e: 'openDir'): void
}>()

const elapsed = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

// 从日志中提取当前处理的文件名
const currentFile = computed(() => {
  if (props.logs.length === 0) return ''
  // 7z 输出中通常包含文件路径
  for (let i = props.logs.length - 1; i >= 0; i--) {
    const line = props.logs[i]
    // 跳过纯状态行
    if (line.startsWith('Processing archive') ||
        line.startsWith('Extracting') ||
        line.startsWith('Compressing') ||
        line.startsWith('Everything is Ok') ||
        line.startsWith('Files:') ||
        line.startsWith('Size:') ||
        line.startsWith('---') ||
        line.includes('error') ||
        line.includes('Error')) continue
    // 如果行看起来像文件路径
    if (line.includes('/') || line.includes('\\') || line.includes('.')) {
      return line.trim()
    }
  }
  return ''
})

// 处理的文件数（从日志中解析）
const processedCount = computed(() => {
  let count = 0
  for (const line of props.logs) {
    if (line.startsWith('Extracting') || line.startsWith('Compressing')) {
      count++
    }
  }
  return count
})

const formattedTime = computed(() => {
  const s = elapsed.value
  const m = Math.floor(s / 60)
  const sec = s % 60
  return m > 0 ? `${m}m ${sec.toString().padStart(2, '0')}s` : `${sec}s`
})

const statusText = computed(() => {
  if (props.result === 'success') return props.mode === 'compress' ? '压缩完成' : '解压完成'
  if (props.result === 'error') return props.mode === 'compress' ? '压缩失败' : '解压失败'
  return props.mode === 'compress' ? '正在压缩...' : '正在解压...'
})

const statusIcon = computed(() => {
  if (props.result === 'success') return '✅'
  if (props.result === 'error') return '❌'
  return '⏳'
})

watch(() => props.isRunning, (running) => {
  if (running) {
    elapsed.value = 0
    timer = setInterval(() => elapsed.value++, 1000)
  } else {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }
}, { immediate: true })

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <Teleport to="body">
    <div class="overlay-backdrop" @click.self="result && emit('close')">
      <div class="overlay-card">
        <!-- 头部 -->
        <div class="overlay-header">
          <div class="overlay-title">
            <span class="overlay-icon">{{ statusIcon }}</span>
            <span>{{ statusText }}</span>
          </div>
          <button
            v-if="result"
            class="overlay-close"
            @click="emit('close')"
          >
            ✕
          </button>
        </div>

        <!-- 统计区 -->
        <div class="overlay-stats">
          <div class="stat-item">
            <span class="stat-icon">⏱️</span>
            <div class="stat-content">
              <span class="stat-label">耗时</span>
              <span class="stat-value">{{ formattedTime }}</span>
            </div>
          </div>
          <div class="stat-item">
            <span class="stat-icon">📄</span>
            <div class="stat-content">
              <span class="stat-label">已处理</span>
              <span class="stat-value">{{ processedCount }} 个文件</span>
            </div>
          </div>
        </div>

        <!-- 当前文件 -->
        <div v-if="currentFile && isRunning" class="current-file">
          <span class="current-file-label">当前文件:</span>
          <span class="current-file-name">{{ currentFile }}</span>
        </div>

        <!-- 日志区 -->
        <div class="overlay-log">
          <div
            v-for="(line, index) in logs"
            :key="index"
            class="log-line"
            :class="{ 'log-error': line.toLowerCase().includes('error') }"
          >
            {{ line }}
          </div>
          <div v-if="logs.length === 0" class="log-empty">
            等待开始...
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="overlay-actions">
          <template v-if="isRunning">
            <button class="cancel-btn" @click="emit('cancel')">
              取消
            </button>
          </template>
          <template v-else-if="result">
            <button class="close-btn" @click="emit('close')">
              关闭
            </button>
            <button
              v-if="result === 'success'"
              class="open-dir-btn"
              @click="emit('openDir')"
            >
              📂 打开目录
            </button>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.overlay-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.overlay-card {
  background: linear-gradient(135deg, #0f1320, #1a1040);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 18px;
  width: 100%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.5);
}

.overlay-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.overlay-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: #f1f5f9;
}

.overlay-icon {
  font-size: 20px;
}

.overlay-close {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.overlay-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
}

.overlay-stats {
  display: flex;
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
}

.stat-icon {
  font-size: 18px;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 11px;
  color: #64748b;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #f1f5f9;
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.current-file {
  padding: 10px 20px;
  background: rgba(99, 102, 241, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.current-file-label {
  font-size: 11px;
  color: #64748b;
  margin-right: 8px;
}

.current-file-name {
  font-size: 12px;
  color: #a5b4fc;
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.overlay-log {
  flex: 1;
  min-height: 120px;
  max-height: 300px;
  overflow-y: auto;
  padding: 14px 20px;
  background: rgba(0, 0, 0, 0.2);
}

.log-line {
  font-size: 11px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  color: #94a3b8;
  line-height: 1.6;
  word-break: break-all;
}

.log-line:last-child {
  color: #22d3ee;
}

.log-line.log-error {
  color: #f43f5e;
}

.log-empty {
  text-align: center;
  color: #64748b;
  font-size: 13px;
  padding: 24px 0;
}

.overlay-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.cancel-btn {
  background: rgba(244, 63, 94, 0.15);
  color: #f43f5e;
  border: 1px solid rgba(244, 63, 94, 0.25);
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background: rgba(244, 63, 94, 0.25);
  border-color: #f43f5e;
}

.close-btn {
  background: rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
}

.open-dir-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.open-dir-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4);
}
</style>
