<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useApp } from '../composables/useApp'
import { OnFileDrop, OnFileDropOff } from '../../wailsjs/runtime/runtime'

const { selectFile } = useApp()

const emit = defineEmits<{
  (e: 'select', path: string): void
}>()

const isDragging = ref(false)
const dragError = ref(false)
const errorMessage = ref('')

function acceptPath(path: string) {
  if (!path) return
  errorMessage.value = ''
  dragError.value = false
  emit('select', path)
}

function showError(message: string) {
  errorMessage.value = message
  dragError.value = true
  setTimeout(() => {
    dragError.value = false
    errorMessage.value = ''
  }, 3000)
}

onMounted(() => {
  OnFileDrop((_x, _y, paths) => {
    isDragging.value = false
    if (paths && paths.length > 0) {
      acceptPath(paths[0])
      return
    }
    showError('未能读取拖拽文件路径，请点击选择文件')
  }, true)
})

onUnmounted(() => {
  OnFileDropOff()
})

function handleDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
  dragError.value = false
  errorMessage.value = ''
}

function handleDragLeave() {
  isDragging.value = false
}

function handleDrop(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false

  const files = e.dataTransfer?.files
  const filePath = files && files.length > 0 ? (files[0] as any).path : ''
  if (filePath) {
    acceptPath(filePath)
  }
}

async function handleClick() {
  dragError.value = false
  errorMessage.value = ''
  try {
    const path = await selectFile()
    if (path) {
      acceptPath(path)
    }
  } catch (err) {
    showError(`打开文件选择器失败: ${err}`)
  }
}
</script>

<template>
  <div
    class="file-selector"
    :class="{ dragging: isDragging, 'drag-error': dragError }"
    style="--wails-drop-target: drop"
    @dragover="handleDragOver"
    @dragleave="handleDragLeave"
    @drop="handleDrop"
    @click="handleClick"
  >
    <div class="icon">📦</div>
    <div v-if="dragError" class="error-text">
      {{ errorMessage || '未能读取文件路径，请点击选择文件' }}
    </div>
    <div v-else class="text">拖拽压缩包到这里，或点击选择文件</div>
  </div>
</template>

<style scoped>
.file-selector {
  border: 2px dashed rgba(99, 102, 241, 0.3);
  border-radius: 16px;
  padding: 48px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.02);
}

.file-selector:hover,
.file-selector.dragging {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
  box-shadow: 0 0 30px rgba(99, 102, 241, 0.2);
}

.file-selector.drag-error {
  border-color: #f43f5e;
  background: rgba(244, 63, 94, 0.05);
}

.icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.text {
  color: #94a3b8;
  font-size: 16px;
}

.error-text {
  color: #f43f5e;
  font-size: 14px;
}
</style>
