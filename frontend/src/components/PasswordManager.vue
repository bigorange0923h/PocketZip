<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApp } from '../composables/useApp'

const { getPasswordRecords, deletePasswordRecord, getPasswordStats } = useApp()

interface PasswordRecord {
  id: number
  archivePath: string
  archiveName: string
  successCount: number
  lastUsedAt: string
  createdAt: string
}

interface PasswordStats {
  totalRecords: number
  totalUsed: number
}

const records = ref<PasswordRecord[]>([])
const stats = ref<PasswordStats>({ totalRecords: 0, totalUsed: 0 })
const isLoading = ref(false)

onMounted(async () => {
  await loadData()
})

async function loadData() {
  isLoading.value = true
  try {
    const [recordsData, statsData] = await Promise.all([
      getPasswordRecords(),
      getPasswordStats()
    ])
    records.value = recordsData || []
    stats.value = statsData || { totalRecords: 0, totalUsed: 0 }
  } catch (err) {
    console.error('Failed to load password data:', err)
  } finally {
    isLoading.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定要删除这条密码记录吗？')) return

  try {
    await deletePasswordRecord(id)
    await loadData()
  } catch (err) {
    console.error('Failed to delete record:', err)
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<template>
  <div class="password-manager">
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-value">{{ stats.totalRecords }}</span>
        <span class="stat-label">密码记录</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{{ stats.totalUsed }}</span>
        <span class="stat-label">总使用次数</span>
      </div>
    </div>

    <div v-if="isLoading" class="loading">加载中...</div>

    <div v-else-if="records.length === 0" class="empty">
      暂无密码记录
    </div>

    <div v-else class="records-list">
      <div v-for="record in records" :key="record.id" class="record-item">
        <div class="record-info">
          <div class="record-name">{{ record.archiveName }}</div>
          <div class="record-path">{{ record.archivePath }}</div>
          <div class="record-meta">
            <span class="success-count">使用 {{ record.successCount }} 次</span>
            <span class="last-used">最后使用: {{ formatDate(record.lastUsedAt) }}</span>
          </div>
        </div>
        <button class="delete-btn" @click="handleDelete(record.id)">删除</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.password-manager {
  margin-top: 16px;
}

.stats-bar {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #a5b4fc;
}

.stat-label {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
}

.loading, .empty {
  text-align: center;
  padding: 32px;
  color: #94a3b8;
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.record-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.record-info {
  flex: 1;
  min-width: 0;
}

.record-name {
  font-weight: 600;
  color: #f1f5f9;
  margin-bottom: 4px;
}

.record-path {
  font-size: 12px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 8px;
}

.record-meta {
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: #94a3b8;
}

.success-count {
  color: #34d399;
}

.delete-btn {
  background: rgba(244, 63, 94, 0.2);
  color: #f43f5e;
  border: 1px solid rgba(244, 63, 94, 0.3);
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.delete-btn:hover {
  background: rgba(244, 63, 94, 0.3);
  border-color: #f43f5e;
}
</style>
