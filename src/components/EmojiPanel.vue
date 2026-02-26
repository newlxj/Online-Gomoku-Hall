<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { EMOJIS, type Emoji } from '@/types/multiplayer'

const emit = defineEmits<{
  (e: 'send', emoji: Emoji): void
}>()

const showPanel = ref(false)
const lastSendTime = ref(0)
const COOLDOWN = 5000 // 5秒冷却
const currentTime = ref(Date.now())

// 定时器更新当前时间
let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(() => {
    currentTime.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})

const canSend = computed(() => {
  return currentTime.value - lastSendTime.value >= COOLDOWN
})

const cooldownRemaining = computed(() => {
  const remaining = Math.ceil((COOLDOWN - (currentTime.value - lastSendTime.value)) / 1000)
  return remaining > 0 ? remaining : 0
})

function handleSend(emoji: Emoji) {
  if (!canSend.value) return

  emit('send', emoji)
  lastSendTime.value = Date.now()
  currentTime.value = Date.now() // 立即更新
  showPanel.value = false
}

function togglePanel() {
  showPanel.value = !showPanel.value
}
</script>

<template>
  <div class="emoji-panel-container">
    <button
      class="emoji-trigger glass-btn"
      :class="{ disabled: !canSend }"
      @click="togglePanel"
    >
      😊 表情
      <span v-if="!canSend" class="cooldown">({{ cooldownRemaining }}s)</span>
    </button>

    <div v-if="showPanel" class="emoji-panel glass-panel-dark">
      <button
        v-for="emoji in EMOJIS"
        :key="emoji"
        class="emoji-btn"
        :disabled="!canSend"
        @click="handleSend(emoji)"
      >
        {{ emoji }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.emoji-panel-container {
  position: relative;
  margin-top: 16px;
}

.emoji-trigger {
  padding: 10px 20px;
  font-size: 14px;
}

.emoji-trigger.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.cooldown {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}

.emoji-panel {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px;
  border-radius: 12px;
  width: 220px;
  margin-bottom: 8px;
  animation: fadeIn 0.2s ease;
}

.emoji-btn {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  font-size: 24px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}

.emoji-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}
</style>
