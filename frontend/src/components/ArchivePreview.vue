<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApp } from '../composables/useApp'

const { previewArchive } = useApp()

const props = defineProps<{
  archivePath: string
}>()

interface ArchiveEntry {
  path: string
  size: number
  isDir: boolean
  modified: string
}

const entries = ref<ArchiveEntry[]>([])
const isLoading = ref(false)
const error = ref('')

watch(() => props.archivePath, async (newPath) => {
  if (newPath) {
    await loadPreview(newPath)
  } else {
    entries.value = []
  }
}, { immediate: true })

async function loadPreview(path: string) {
  isLoading.value = true
  error.value = ''
  entries.value = []

  try {
    const result = await previewArchive(path)
    entries.value = result || []
  } catch (err: any) {
    error.value = err?.message || String(err)
  } finally {
    isLoading.value = false
  }
}

function formatSize(size: number): string {
  if (size === 0) return '-'
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(1) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

function getFileIcon(entry: ArchiveEntry): string {
  if (entry.isDir) return '📁'
  const ext = entry.path.split('.').pop()?.toLowerCase() || ''
  const iconMap: Record<string, string> = {
    'txt': '📄', 'doc': '📝', 'docx': '📝', 'pdf': '📕',
    'jpg': '🖼️', 'jpeg': '🖼️', 'png': '🖼️', 'gif': '🖼️',
    'mp3': '🎵', 'mp4': '🎬', 'avi': '🎬', 'mkv': '🎬',
    'zip': '📦', '7z': '📦', 'rar': '📦', 'tar': '📦',
    'exe': '⚙️', 'msi': '⚙️', 'dll': '⚙️',
    'js': '📜', 'ts': '📜', 'py': '📜', 'go': '📜',
  }
  return iconMap[ext] || '📄'
}
</script>

<template>
  <div class="archive-preview">
    <div v-if="isLoading" class="loading">加载预览中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="entries.length === 0" class="empty">无内容或不支持预览</div>
    <div v-else class="entries-list">
      <div class="entries-header">
        <span class="col-icon"></span>
        <span class="col-name">文件名</span>
        <span class="col-size">大小</span>
        <span class="col-modified">修改时间</span>
      </div>
      <div v-for="(entry, index) in entries" :key="index" class="entry-item">
        <span class="col-icon">{{ getFileIcon(entry) }}</span>
        <span class="col-name" :title="entry.path">{{ entry.path }}</span>
        <span class="col-size">{{ formatSize(entry.size) }}</span>
        <span class="col-modified">{{ entry.modified }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.archive-preview {
  margin-top: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.loading, .error, .empty {
  text-align: center;
  padding: 20px;
  color: #94a3b8;
}

.error {
  color: #f43f5e;
}

.entries-header {
  display: flex;
  padding: 8px 12px;
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.entry-item {
  display: flex;
  padding: 8px 12px;
  font-size: 13px;
  color: #f1f5f9;
  border-bottom: 1px solid rgba(255, 255, 255, 0.02);
  transition: background 0.2s ease;
}

.entry-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.col-icon {
  width: 24px;
  flex-shrink: 0;
}

.col-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.col-size {
  width: 80px;
  text-align: right;
  color: #94a3b8;
  flex-shrink: 0;
}

.col-modified {
  width: 140px;
  text-align: right;
  color: #64748b;
  flex-shrink: 0;
}
</style>
