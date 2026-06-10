<script setup lang="ts">
import { ref } from 'vue'
import { useApp } from '../composables/useApp'

const { selectFile } = useApp()

const emit = defineEmits<{
  (e: 'select', path: string): void
}>()

const isDragging = ref(false)

function handleDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
}

function handleDragLeave() {
  isDragging.value = false
}

function handleDrop(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    emit('select', (files[0] as any).path)
  }
}

async function handleClick() {
  const path = await selectFile()
  if (path) {
    emit('select', path)
  }
}
</script>

<template>
  <div
    class="file-selector"
    :class="{ dragging: isDragging }"
    @dragover="handleDragOver"
    @dragleave="handleDragLeave"
    @drop="handleDrop"
    @click="handleClick"
  >
    <div class="icon">📦</div>
    <div class="text">拖拽压缩包到这里，或点击选择文件</div>
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

.icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.text {
  color: #94a3b8;
  font-size: 16px;
}
</style>
