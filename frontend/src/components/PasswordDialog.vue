<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  archivePath: string
  candidates: string[]
}>()

const emit = defineEmits<{
  (e: 'submit', password: string): void
  (e: 'cancel'): void
}>()

const password = ref('')
const showPassword = ref(false)
const selectedCandidate = ref('')

function handleSubmit() {
  const pwd = selectedCandidate.value || password.value
  if (pwd) {
    emit('submit', pwd)
  }
}

function handleCancel() {
  emit('cancel')
}

function selectCandidate(candidate: string) {
  selectedCandidate.value = candidate
  password.value = candidate
}
</script>

<template>
  <div class="dialog-overlay" @click.self="handleCancel">
    <div class="dialog">
      <div class="dialog-header">
        <h3>需要密码</h3>
        <button class="close-btn" @click="handleCancel">&times;</button>
      </div>
      <div class="dialog-body">
        <p class="archive-path">{{ archivePath }}</p>
        <div v-if="candidates.length > 0" class="candidates">
          <p class="candidates-title">历史密码：</p>
          <div class="candidate-list">
            <button
              v-for="candidate in candidates"
              :key="candidate"
              class="candidate-btn"
              :class="{ active: selectedCandidate === candidate }"
              @click="selectCandidate(candidate)"
            >
              {{ candidate }}
            </button>
          </div>
        </div>
        <div class="input-group">
          <input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="输入密码"
            class="password-input"
            @keyup.enter="handleSubmit"
          />
          <button class="toggle-btn" @click="showPassword = !showPassword">
            {{ showPassword ? '🙈' : '👁️' }}
          </button>
        </div>
      </div>
      <div class="dialog-footer">
        <button class="cancel-btn" @click="handleCancel">取消</button>
        <button class="submit-btn" @click="handleSubmit" :disabled="!password">
          解压
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.dialog {
  background: linear-gradient(135deg, #1a1040, #0d1f2d);
  border-radius: 16px;
  width: 400px;
  max-width: 90vw;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.dialog-header h3 {
  margin: 0;
  font-size: 18px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.close-btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #f1f5f9;
}

.dialog-body {
  padding: 20px;
}

.archive-path {
  color: #94a3b8;
  font-size: 12px;
  margin: 0 0 16px 0;
  word-break: break-all;
}

.candidates {
  margin-bottom: 16px;
}

.candidates-title {
  color: #94a3b8;
  font-size: 12px;
  margin: 0 0 8px 0;
}

.candidate-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.candidate-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 6px 12px;
  color: #94a3b8;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.candidate-btn:hover,
.candidate-btn.active {
  background: rgba(99, 102, 241, 0.2);
  border-color: #6366f1;
  color: #f1f5f9;
}

.input-group {
  display: flex;
  gap: 8px;
}

.password-input {
  flex: 1;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 12px 16px;
  color: #f1f5f9;
  font-size: 14px;
  outline: none;
  transition: all 0.3s ease;
}

.password-input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2);
}

.toggle-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 12px;
  cursor: pointer;
  font-size: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.cancel-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 10px 20px;
  color: #94a3b8;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
}

.submit-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border: none;
  border-radius: 10px;
  padding: 10px 20px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
