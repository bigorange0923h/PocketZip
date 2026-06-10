<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApp } from '../composables/useApp'

const { getExtractStrategies, saveExtractStrategy } = useApp()

interface ExtractStrategy {
  name: string
  outputDir: string
  autoRetry: boolean
  maxRetries: number
  autoOpen: boolean
}

const strategies = ref<ExtractStrategy[]>([])
const editingStrategy = ref<ExtractStrategy | null>(null)
const isLoading = ref(false)

onMounted(async () => {
  await loadStrategies()
})

async function loadStrategies() {
  isLoading.value = true
  try {
    strategies.value = await getExtractStrategies()
  } catch (err) {
    console.error('Failed to load strategies:', err)
  } finally {
    isLoading.value = false
  }
}

function editStrategy(strategy: ExtractStrategy) {
  editingStrategy.value = { ...strategy }
}

function cancelEdit() {
  editingStrategy.value = null
}

async function saveStrategy() {
  if (!editingStrategy.value) return

  try {
    await saveExtractStrategy(editingStrategy.value)
    await loadStrategies()
    editingStrategy.value = null
  } catch (err) {
    console.error('Failed to save strategy:', err)
  }
}

function getStrategyDescription(strategy: ExtractStrategy): string {
  const parts = []
  if (strategy.autoRetry) parts.push(`自动重试 ${strategy.maxRetries} 次`)
  if (strategy.autoOpen) parts.push('完成后打开目录')
  if (strategy.outputDir) parts.push(`输出到 ${strategy.outputDir}`)
  return parts.join('，') || '默认设置'
}

function getStrategyIcon(strategy: ExtractStrategy): string {
  if (strategy.name === 'default') return '📦'
  if (strategy.name === 'retry') return '🔄'
  if (strategy.name === 'auto-open') return '📂'
  return '⚙️'
}
</script>

<template>
  <div class="strategy-manager">
    <div class="strategy-header">
      <div class="strategy-title">解压策略模板</div>
      <div class="strategy-desc">预设不同的解压配置，快速应用</div>
    </div>

    <div v-if="isLoading" class="loading">加载中...</div>

    <div v-else class="strategies-list">
      <div
        v-for="strategy in strategies"
        :key="strategy.name"
        class="strategy-item"
        :class="{ editing: editingStrategy?.name === strategy.name }"
      >
        <div class="strategy-info">
          <div class="strategy-icon">{{ getStrategyIcon(strategy) }}</div>
          <div class="strategy-details">
            <div class="strategy-name">{{ strategy.name }}</div>
            <div class="strategy-desc">{{ getStrategyDescription(strategy) }}</div>
          </div>
        </div>
        <button class="edit-btn" @click="editStrategy(strategy)">编辑</button>
      </div>
    </div>

    <!-- 编辑面板 -->
    <div v-if="editingStrategy" class="edit-panel">
      <div class="edit-header">编辑策略: {{ editingStrategy.name }}</div>
      <div class="edit-form">
        <div class="form-group">
          <label>输出目录</label>
          <input
            v-model="editingStrategy.outputDir"
            type="text"
            placeholder="留空则使用默认目录"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label>
            <input
              v-model="editingStrategy.autoRetry"
              type="checkbox"
              class="form-checkbox"
            />
            自动重试失败的解压
          </label>
        </div>
        <div v-if="editingStrategy.autoRetry" class="form-group">
          <label>最大重试次数</label>
          <input
            v-model.number="editingStrategy.maxRetries"
            type="number"
            min="1"
            max="10"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label>
            <input
              v-model="editingStrategy.autoOpen"
              type="checkbox"
              class="form-checkbox"
            />
            解压完成后自动打开目录
          </label>
        </div>
      </div>
      <div class="edit-actions">
        <button class="cancel-btn" @click="cancelEdit">取消</button>
        <button class="save-btn" @click="saveStrategy">保存</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.strategy-manager {
  margin-top: 16px;
}

.strategy-header {
  margin-bottom: 16px;
}

.strategy-title {
  font-size: 16px;
  font-weight: 600;
  color: #f1f5f9;
  margin-bottom: 4px;
}

.strategy-desc {
  font-size: 12px;
  color: #64748b;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #94a3b8;
}

.strategies-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.strategy-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.2s ease;
}

.strategy-item:hover {
  background: rgba(99, 102, 241, 0.05);
  border-color: rgba(99, 102, 241, 0.2);
}

.strategy-item.editing {
  border-color: #6366f1;
}

.strategy-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.strategy-icon {
  font-size: 24px;
}

.strategy-name {
  font-weight: 600;
  color: #f1f5f9;
  text-transform: capitalize;
}

.strategy-desc {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 2px;
}

.edit-btn {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.edit-btn:hover {
  background: rgba(99, 102, 241, 0.3);
  border-color: #6366f1;
}

.edit-panel {
  margin-top: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.edit-header {
  font-size: 14px;
  font-weight: 600;
  color: #a5b4fc;
  margin-bottom: 16px;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 13px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-input {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  padding: 8px 12px;
  color: #f1f5f9;
  font-size: 13px;
}

.form-input:focus {
  outline: none;
  border-color: #6366f1;
}

.form-checkbox {
  width: 16px;
  height: 16px;
  accent-color: #6366f1;
}

.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

.cancel-btn {
  background: rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.save-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.save-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}
</style>
