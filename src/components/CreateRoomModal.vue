<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { RoomSettings } from '@/types/multiplayer'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', name: string, settings: RoomSettings): void
}>()

const roomName = ref('')
const settings = reactive<RoomSettings>({
  timeLimit: 600,      // 10分钟
  moveTimeLimit: 30,   // 30秒
  firstMove: 'random',
  ratedGame: true,
})

const timeOptions = [
  { label: '5分钟', value: 300 },
  { label: '10分钟', value: 600 },
  { label: '15分钟', value: 900 },
  { label: '20分钟', value: 1200 },
  { label: '30分钟', value: 1800 },
]

const moveTimeOptions = [
  { label: '15秒', value: 15 },
  { label: '30秒', value: 30 },
  { label: '45秒', value: 45 },
  { label: '60秒', value: 60 },
]

const firstMoveOptions: { label: string; value: 'host' | 'guest' | 'random' }[] = [
  { label: '随机', value: 'random' },
  { label: '房主先手', value: 'host' },
  { label: '访客先手', value: 'guest' },
]

function handleCreate() {
  const name = roomName.value.trim() || '未命名房间'
  emit('create', name, { ...settings })
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="modal-overlay" @click.self="handleClose">
    <div class="modal-content glass-panel">
      <div class="modal-header">
        <h3 class="modal-title">创建房间</h3>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>

      <div class="modal-body">
        <!-- 房间名称 -->
        <div class="form-group">
          <label class="form-label">房间名称</label>
          <input
            v-model="roomName"
            type="text"
            class="form-input"
            placeholder="输入房间名称（可选）"
            maxlength="20"
          />
        </div>

        <!-- 总时间 -->
        <div class="form-group">
          <label class="form-label">每人总时间</label>
          <div class="option-buttons">
            <button
              v-for="opt in timeOptions"
              :key="opt.value"
              class="option-btn"
              :class="{ active: settings.timeLimit === opt.value }"
              @click="settings.timeLimit = opt.value"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- 每步限时 -->
        <div class="form-group">
          <label class="form-label">每步限时</label>
          <div class="option-buttons">
            <button
              v-for="opt in moveTimeOptions"
              :key="opt.value"
              class="option-btn"
              :class="{ active: settings.moveTimeLimit === opt.value }"
              @click="settings.moveTimeLimit = opt.value"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- 先手规则 -->
        <div class="form-group">
          <label class="form-label">先手规则</label>
          <div class="option-buttons">
            <button
              v-for="opt in firstMoveOptions"
              :key="opt.value"
              class="option-btn"
              :class="{ active: settings.firstMove === opt.value }"
              @click="settings.firstMove = opt.value"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- 排位模式 -->
        <div class="form-group">
          <label class="checkbox-label">
            <input
              v-model="settings.ratedGame"
              type="checkbox"
              class="checkbox-input"
            />
            <span class="checkbox-custom"></span>
            <span>排位赛模式（计入排行榜）</span>
          </label>
        </div>
      </div>

      <div class="modal-footer">
        <button class="glass-btn" @click="handleClose">取消</button>
        <button class="glass-btn-primary" @click="handleCreate">创建</button>
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
  max-width: 480px;
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
  font-size: 22px;
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
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  font-weight: 500;
}

.form-input {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  padding: 12px 16px;
  color: #fff;
  font-size: 16px;
  transition: all 0.3s;
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.form-input:focus {
  outline: none;
  border-color: #00ffcc;
  background: rgba(255, 255, 255, 0.15);
}

.option-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.option-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 10px 16px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.option-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

.option-btn.active {
  background: linear-gradient(135deg, #00ffcc, #00ccff);
  border-color: transparent;
  color: #000;
  font-weight: 600;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
}

.checkbox-input {
  display: none;
}

.checkbox-custom {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  position: relative;
  transition: all 0.3s;
}

.checkbox-input:checked + .checkbox-custom {
  background: linear-gradient(135deg, #00ffcc, #00ccff);
  border-color: transparent;
}

.checkbox-input:checked + .checkbox-custom::after {
  content: '✓';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #000;
  font-size: 14px;
  font-weight: bold;
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
