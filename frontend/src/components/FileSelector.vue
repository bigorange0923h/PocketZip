<script setup lang="ts">
import { ref } from 'vue'
import { useApp } from '../composables/useApp'

const { selectFile } = useApp()

const emit = defineEmits<{
  (e: 'select', path: string): void
}>()

const isDragging = ref(false)
const dragError = ref(false)

function handleDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
  dragError.value = false
}

function handleDragLeave() {
  isDragging.value = false
}

function handleDrop(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false

  // WebView2 中 File 对象没有 path 属性，需要提示用户点击选择
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    // 尝试获取路径（某些环境下可用）
    const filePath = (files[0] as any).path
    if (filePath) {
      emit('select', filePath)
    } else {
      // 无法获取路径，提示用户点击选择
      dragError.value = true
      setTimeout(() => {
        dragError.value = false
      }, 3000)
    }
  }
}

async function handleClick() {
  dragError.value = false
  const path = await selectFile()
  if (path) {
    emit('select', path)
  }
}
</script>

<template>
  <div
    class="file-selector"
    :class="{ dragging: isDragging, 'drag-error': dragError }"
    @dragover="handleDragOver"
    @dragleave="handleDragLeave"
    @drop="handleDrop"
    @click="handleClick"
  >
    <div class="icon">📦</div>
    <div v-if="dragError" class="error-text">
      ⚠️ 拖拽无法获取文件路径，请点击选择文件
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
