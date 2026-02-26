<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAlias } from '@/composables/useAlias'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm', alias: string): void
}>()

const {
  currentAlias,
  aliasHistory,
  isAvailable,
  checking,
  setCurrentAlias,
  generateNewAlias,
  checkAlias,
  saveAliasToHistory,
} = useAlias()

const inputAlias = ref(currentAlias.value || '')
const errorMessage = ref('')

const canConfirm = computed(() => {
  return inputAlias.value.trim().length >= 2 && isAvailable.value === true
})

async function handleCheck() {
  const alias = inputAlias.value.trim()
  if (alias.length < 2) {
    errorMessage.value = '别名至少需要2个字符'
    return
  }
  if (alias.length > 12) {
    errorMessage.value = '别名最多12个字符'
    return
  }

  errorMessage.value = ''
  setCurrentAlias(alias)
  await checkAlias(alias)
}

function handleGenerate() {
  const newAlias = generateNewAlias()
  inputAlias.value = newAlias
  handleCheck()
}

function handleSelectHistory(alias: string) {
  inputAlias.value = alias
  handleCheck()
}

async function handleConfirm() {
  const alias = inputAlias.value.trim()

  if (alias.length < 2) {
    errorMessage.value = '别名至少需要2个字符'
    return
  }
  if (alias.length > 12) {
    errorMessage.value = '别名最多12个字符'
    return
  }

  // 如果还没检查或者检查结果不可用，先检查
  if (isAvailable.value !== true) {
    await checkAlias(alias)
    // 检查完成后再次验证 - 使用 === false 来避免类型收窄问题
    if (isAvailable.value === false) {
      errorMessage.value = '该昵称已被使用，请换一个'
      return
    }
    // 如果检查结果不是 true（可能是 null），也返回
    if (isAvailable.value !== true) {
      errorMessage.value = '无法验证昵称，请重试'
      return
    }
  }

  console.log('[AliasModal] Confirming alias:', alias)
  setCurrentAlias(alias)
  saveAliasToHistory(alias)
  emit('confirm', alias)
  emit('close')
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="modal-overlay" @click.self="handleClose">
    <div class="modal-content glass-panel">
      <div class="modal-header">
        <h3 class="modal-title">设置你的游戏昵称</h3>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>

      <div class="modal-body">
        <!-- 输入框 -->
        <div class="input-group">
          <input
            v-model="inputAlias"
            type="text"
            class="alias-input"
            placeholder="输入你的昵称"
            maxlength="12"
            @blur="handleCheck"
            @keyup.enter="handleConfirm"
          />
          <button class="generate-btn" @click="handleGenerate">
            🎲 随机
          </button>
        </div>

        <!-- 状态提示 -->
        <div class="status-message">
          <template v-if="checking">
            <span class="checking">检查中...</span>
          </template>
          <template v-else-if="errorMessage">
            <span class="error">{{ errorMessage }}</span>
          </template>
          <template v-else-if="isAvailable === true">
            <span class="available">✓ 该昵称可用</span>
          </template>
          <template v-else-if="isAvailable === false">
            <span class="unavailable">✗ 该昵称已被使用</span>
          </template>
        </div>

        <!-- 历史记录 -->
        <div v-if="aliasHistory.length > 0" class="history-section">
          <h4 class="history-title">最近使用</h4>
          <div class="history-list">
            <button
              v-for="alias in aliasHistory"
              :key="alias"
              class="history-item"
              @click="handleSelectHistory(alias)"
            >
              {{ alias }}
            </button>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="glass-btn" @click="handleClose">取消</button>
        <button
          class="glass-btn-primary"
          :disabled="!canConfirm"
          @click="handleConfirm"
        >
          确认
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal-content {
  width: 90%;
  max-width: 400px;
  padding: 24px;
  border-radius: 20px;
  animation: scaleIn 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-title {
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  font-size: 24px;
  cursor: pointer;
  padding: 4px 8px;
  transition: color 0.3s;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  display: flex;
  gap: 12px;
}

.alias-input {
  flex: 1;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  padding: 14px 16px;
  color: #fff;
  font-size: 18px;
  transition: all 0.3s;
}

.alias-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.alias-input:focus {
  outline: none;
  border-color: #00ffcc;
  background: rgba(255, 255, 255, 0.15);
}

.generate-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  padding: 14px 20px;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.3s;
}

.generate-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.status-message {
  min-height: 24px;
  font-size: 14px;
}

.checking {
  color: rgba(255, 255, 255, 0.6);
}

.error {
  color: #ff6b6b;
}

.available {
  color: #00ff88;
}

.unavailable {
  color: #ff6b6b;
}

.history-section {
  margin-top: 8px;
}

.history-title {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  font-weight: normal;
  margin: 0 0 8px 0;
}

.history-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.history-item {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  padding: 6px 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-footer .glass-btn,
.modal-footer .glass-btn-primary {
  padding: 12px 24px;
}

.modal-footer .glass-btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
