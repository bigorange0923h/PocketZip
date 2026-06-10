<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'

const props = defineProps<{
  logs: string[]
}>()

const logContainer = ref<HTMLElement | null>(null)

watch(() => props.logs.length, async () => {
  await nextTick()
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
})
</script>

<template>
  <div class="log-panel">
    <div class="log-header">实时日志</div>
    <div ref="logContainer" class="log-content">
      <div v-for="(log, index) in logs" :key="index" class="log-line">
        {{ log }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.log-panel {
  margin-top: 16px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.log-header {
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.log-content {
  padding: 16px;
  max-height: 300px;
  overflow-y: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.6;
}

.log-line {
  color: #64748b;
  white-space: pre-wrap;
  word-break: break-all;
}

.log-line:last-child {
  color: #22d3ee;
}
</style>
