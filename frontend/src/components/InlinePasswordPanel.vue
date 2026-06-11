<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useApp } from '../composables/useApp'

const props = defineProps<{
  password: string
  remember: boolean
}>()

const emit = defineEmits<{
  (e: 'update:password', value: string): void
  (e: 'update:remember', value: boolean): void
}>()

const { getPasswordRecords, deletePasswordRecord, getPasswordStats } = useApp()

interface PasswordRecord {
  id: number
  archivePath: string
  password: string
  successCount: number
  lastUsedAt: string
}

const records = ref<PasswordRecord[]>([])
const stats = ref({ totalRecords: 0, totalUsage: 0 })
const showPassword = ref(false)

const localPassword = computed({
  get: () => props.password,
  set: (v: string) => emit('update:password', v)
})

const localRemember = computed({
  get: () => props.remember,
  set: (v: boolean) => emit('update:remember', v)
})

const maskedPassword = (pwd: string) => {
  const len = Math.min(pwd.length, 12)
  return '•'.repeat(len)
}

async function loadRecords() {
  try {
    records.value = await getPasswordRecords()
    stats.value = await getPasswordStats()
  } catch {
    // ignore
  }
}

async function handleDelete(id: number) {
  try {
    await deletePasswordRecord(id)
    await loadRecords()
  } catch {
    // ignore
  }
}

function applyRecord(record: PasswordRecord) {
  localPassword.value = record.password
}

onMounted(loadRecords)

defineExpose({ loadRecords })
</script>

<template>
  <div class="inline-password-panel">
    <!-- 密码输入 -->
    <div class="panel-section">
      <label class="input-label">
        <span class="label-icon">🔒</span>
        密码
      </label>
      <div class="password-input-group">
        <input
          v-model="localPassword"
          :type="showPassword ? 'text' : 'password'"
          placeholder="输入解压密码（可选）"
          class="password-input"
        />
        <button class="toggle-visibility" @click="showPassword = !showPassword">
          {{ showPassword ? '🙈' : '👁️' }}
        </button>
      </div>
    </div>

    <!-- 记住密码开关 -->
    <div class="panel-section remember-row">
      <span class="remember-label">记住此密码</span>
      <button
        class="toggle-switch"
        :class="{ active: localRemember }"
        @click="localRemember = !localRemember"
      >
        <span class="toggle-thumb" />
      </button>
    </div>

    <!-- 历史密码列表 -->
    <div class="panel-section history-section">
      <div class="section-header">
        <span class="section-title">历史密码</span>
        <span class="record-count">{{ records.length }}</span>
      </div>

      <div v-if="records.length === 0" class="empty-history">
        暂无保存的密码
      </div>

      <div v-else class="history-list">
        <div
          v-for="record in records"
          :key="record.id"
          class="history-item"
          :class="{ active: record.password === password }"
          @click="applyRecord(record)"
        >
          <div class="history-info">
            <span class="history-icon">🔑</span>
            <div class="history-text">
              <span class="history-label">{{ record.archivePath.split(/[/\\]/).pop() }}</span>
              <span class="history-password">{{ maskedPassword(record.password) }}</span>
            </div>
          </div>
          <div class="history-meta">
            <span class="history-time">{{ record.lastUsedAt }}</span>
            <button class="delete-btn" @click.stop="handleDelete(record.id)" title="删除">
              ✕
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.inline-password-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.panel-section {
  padding: 0;
}

.input-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8;
  margin-bottom: 8px;
}

.label-icon {
  font-size: 14px;
}

.password-input-group {
  display: flex;
  gap: 8px;
}

.password-input {
  flex: 1;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 10px 14px;
  color: #f1f5f9;
  font-size: 14px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  outline: none;
  transition: all 0.2s ease;
}

.password-input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.password-input::placeholder {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: #64748b;
}

.toggle-visibility {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 10px 12px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.toggle-visibility:hover {
  background: rgba(255, 255, 255, 0.1);
}

.remember-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
}

.remember-label {
  font-size: 13px;
  color: #94a3b8;
}

.toggle-switch {
  position: relative;
  width: 40px;
  height: 22px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 11px;
  cursor: pointer;
  transition: background 0.2s ease;
  padding: 0;
}

.toggle-switch.active {
  background: #6366f1;
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 18px;
  height: 18px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s ease;
}

.toggle-switch.active .toggle-thumb {
  transform: translateX(18px);
}

.history-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8;
}

.record-count {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
}

.empty-history {
  text-align: center;
  padding: 20px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: #64748b;
  font-size: 13px;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.1);
}

.history-item.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.history-info {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.history-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.history-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.history-label {
  font-size: 12px;
  color: #f1f5f9;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-password {
  font-size: 11px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  color: #64748b;
  letter-spacing: 2px;
}

.history-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.history-time {
  font-size: 11px;
  color: #64748b;
}

.delete-btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
  line-height: 1;
}

.delete-btn:hover {
  color: #f43f5e;
  background: rgba(244, 63, 94, 0.1);
}
</style>
