<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useTimer } from '@/composables/useTimer'
import { useSound } from '@/composables/useSound'

const props = defineProps<{
  remainingTime: number
  moveTime: number
  isActive: boolean
}>()

const { formatTime, isLowTime, isCriticalTime } = useTimer()
const { playTimeWarning, playTimeCritical, isMuted } = useSound()

const totalTimeDisplay = computed(() => formatTime(props.remainingTime))
const moveTimeDisplay = computed(() => formatTime(props.moveTime))

const totalTimeClass = computed(() => {
  if (isCriticalTime(props.remainingTime, 30)) return 'critical'
  if (isLowTime(props.remainingTime, 60)) return 'low'
  return ''
})

const moveTimeClass = computed(() => {
  if (isCriticalTime(props.moveTime, 10)) return 'critical'
  if (isLowTime(props.moveTime, 15)) return 'low'
  return ''
})

// 跟踪是否已播放过警告音，避免重复播放
const hasPlayedWarning = ref(false)
const hasPlayedCritical = ref(false)

// 监听步时变化，播放警告音
watch([() => props.moveTime, () => props.isActive], ([moveTime, isActive]) => {
  if (!isActive || isMuted.value) {
    // 非活跃状态或静音时重置标记
    hasPlayedWarning.value = false
    hasPlayedCritical.value = false
    return
  }

  // 步时小于5秒：紧急警告
  if (moveTime <= 5 && moveTime > 0 && !hasPlayedCritical.value) {
    playTimeCritical()
    hasPlayedCritical.value = true
  }
  // 步时小于10秒：普通警告
  else if (moveTime <= 10 && moveTime > 5 && !hasPlayedWarning.value) {
    playTimeWarning()
    hasPlayedWarning.value = true
  }
  // 步时重置时（大于10秒），重置标记
  else if (moveTime > 10) {
    hasPlayedWarning.value = false
    hasPlayedCritical.value = false
  }
})
</script>

<template>
  <div class="game-timer glass-panel-dark" :class="{ active: isActive }">
    <div class="timer-row">
      <span class="timer-label">总时间</span>
      <span class="timer-value" :class="totalTimeClass">
        {{ totalTimeDisplay }}
      </span>
    </div>
    <div class="timer-row">
      <span class="timer-label">步时</span>
      <span class="timer-value move-time" :class="moveTimeClass">
        {{ moveTimeDisplay }}
      </span>
    </div>
    <div v-if="isActive" class="active-indicator">
      <span class="pulse"></span>
      <span class="text">思考中</span>
    </div>
  </div>
</template>

<style scoped>
.game-timer {
  padding: 16px 20px;
  border-radius: 12px;
  min-width: 120px;
  transition: all 0.3s;
}

.game-timer.active {
  border: 2px solid rgba(0, 255, 204, 0.5);
  box-shadow: 0 0 20px rgba(0, 255, 204, 0.2);
}

.timer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.timer-row:last-of-type {
  margin-bottom: 0;
}

.timer-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}

.timer-value {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  font-family: 'Courier New', monospace;
}

.timer-value.low {
  color: #ffc800;
}

.timer-value.critical {
  color: #ff6b6b;
  animation: blink 0.5s ease-in-out infinite;
}

.move-time {
  font-size: 14px;
}

.active-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.pulse {
  width: 8px;
  height: 8px;
  background: #00ffcc;
  border-radius: 50%;
  animation: pulse 1s ease-in-out infinite;
}

.text {
  color: #00ffcc;
  font-size: 12px;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.7;
  }
}
</style>
