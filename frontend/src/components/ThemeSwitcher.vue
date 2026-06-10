<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useApp } from '../composables/useApp'

const { getTheme, setTheme } = useApp()

const currentTheme = ref('dark')
const themes = [
  { id: 'dark', name: '深色', icon: '🌙' },
  { id: 'light', name: '浅色', icon: '☀️' },
  { id: 'auto', name: '跟随系统', icon: '💻' }
]

onMounted(async () => {
  try {
    currentTheme.value = await getTheme()
    applyTheme(currentTheme.value)
  } catch (err) {
    console.error('Failed to load theme:', err)
  }
})

watch(currentTheme, async (newTheme) => {
  applyTheme(newTheme)
  try {
    await setTheme(newTheme)
  } catch (err) {
    console.error('Failed to save theme:', err)
  }
})

function applyTheme(theme: string) {
  document.documentElement.setAttribute('data-theme', theme)
}

function handleThemeChange(themeId: string) {
  currentTheme.value = themeId
}
</script>

<template>
  <div class="theme-switcher">
    <div class="theme-label">主题设置</div>
    <div class="theme-options">
      <button
        v-for="theme in themes"
        :key="theme.id"
        class="theme-btn"
        :class="{ active: currentTheme === theme.id }"
        @click="handleThemeChange(theme.id)"
      >
        <span class="theme-icon">{{ theme.icon }}</span>
        <span class="theme-name">{{ theme.name }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.theme-switcher {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.theme-label {
  font-size: 14px;
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 12px;
}

.theme-options {
  display: flex;
  gap: 8px;
}

.theme-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.theme-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.3);
}

.theme-btn.active {
  background: rgba(99, 102, 241, 0.2);
  border-color: #6366f1;
}

.theme-icon {
  font-size: 24px;
}

.theme-name {
  font-size: 12px;
  color: #f1f5f9;
}
</style>
