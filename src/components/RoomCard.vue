<script setup lang="ts">
import { ref, computed } from 'vue'
import type { RoomInfo } from '@/types/multiplayer'

const props = defineProps<{
  room: RoomInfo
}>()

const emit = defineEmits<{
  (e: 'join'): void
  (e: 'spectate'): void
}>()

const isHovered = ref(false)

const statusText = computed(() => {
  switch (props.room.status) {
    case 'waiting': return 'WAITING'
    case 'ready': return 'READY'
    case 'playing': return 'PLAYING'
    case 'finished': return 'FINISHED'
    default: return String(props.room.status).toUpperCase()
  }
})

const statusClass = computed(() => {
  return `status-${props.room.status}`
})

const canJoin = computed(() => {
  // waiting 或 ready 状态且人数小于2时可以加入
  return (props.room.status === 'waiting' || props.room.status === 'ready') && props.room.playerCount < 2
})

const timeText = computed(() => {
  const mins = Math.floor(props.room.settings.timeLimit / 60)
  const moveTime = props.room.settings.moveTimeLimit
  return `${mins}M / ${moveTime}S`
})

const createdTime = computed(() => {
  const date = new Date(props.room.createdAt)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

function handleJoin() {
  if (canJoin.value) {
    emit('join')
  }
}

function handleSpectate() {
  emit('spectate')
}
</script>

<template>
  <div
    class="room-card glass-panel-dark"
    :class="{ hovered: isHovered }"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
  >
    <!-- 顶部装饰线 -->
    <div class="card-top-line"></div>

    <!-- 角落装饰 -->
    <div class="corner-deco tl"></div>
    <div class="corner-deco br"></div>

    <!-- 房间名称和状态 -->
    <div class="room-header">
      <div class="room-id">
        <span class="id-label">ID</span>
        <span class="id-value">{{ room.id.slice(0, 6).toUpperCase() }}</span>
      </div>
      <span class="room-status" :class="statusClass">{{ statusText }}</span>
    </div>

    <h4 class="room-name">{{ room.name }}</h4>

    <!-- 房间信息 -->
    <div class="room-info">
      <div class="info-item">
        <span class="info-icon">👤</span>
        <div class="info-content">
          <span class="info-label">HOST</span>
          <span class="info-value">{{ room.hostName || 'UNKNOWN' }}</span>
        </div>
      </div>
      <div class="info-item">
        <span class="info-icon">👥</span>
        <div class="info-content">
          <span class="info-label">PLAYERS</span>
          <span class="info-value">{{ room.playerCount }}/2</span>
        </div>
      </div>
      <div class="info-item">
        <span class="info-icon">👁</span>
        <div class="info-content">
          <span class="info-label">SPECTATORS</span>
          <span class="info-value">{{ room.spectatorCount }}</span>
        </div>
      </div>
      <div class="info-item">
        <span class="info-icon">⏱</span>
        <div class="info-content">
          <span class="info-label">TIME</span>
          <span class="info-value">{{ timeText }}</span>
        </div>
      </div>
    </div>

    <!-- 时间显示 -->
    <div class="room-time">
      <span class="time-label">CREATED</span>
      <span class="time-value">{{ createdTime }}</span>
    </div>

    <!-- 操作按钮 -->
    <div class="room-actions" :class="{ visible: isHovered || !canJoin }">
      <button
        v-if="canJoin"
        class="glass-btn glass-btn-primary join-btn"
        @click="handleJoin"
      >
        <span class="btn-icon">→</span>
        <span class="btn-text">加入房间</span>
      </button>
      <button
        class="glass-btn spectate-btn"
        @click="handleSpectate"
      >
        <span class="btn-icon">👁</span>
        <span class="btn-text">观战</span>
      </button>
    </div>

    <!-- 计分模式标识 -->
    <div v-if="room.settings.ratedGame" class="rated-badge">
      <span class="badge-icon">★</span>
      <span class="badge-text">RATED</span>
    </div>
  </div>
</template>

<style scoped>
.room-card {
  position: relative;
  padding: 16px;
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
}

/* 顶部装饰线 */
.card-top-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta));
  opacity: 0.5;
  transition: opacity 0.3s ease;
}

.room-card:hover .card-top-line,
.room-card.hovered .card-top-line {
  opacity: 1;
  animation: lineGlow 2s ease infinite;
}

@keyframes lineGlow {
  0%, 100% {
    box-shadow: 0 0 5px var(--neon-cyan);
  }
  50% {
    box-shadow: 0 0 15px var(--neon-magenta);
  }
}

/* 角落装饰 */
.corner-deco {
  position: absolute;
  width: 12px;
  height: 12px;
  pointer-events: none;
}

.corner-deco::before,
.corner-deco::after {
  content: '';
  position: absolute;
  background: var(--neon-cyan);
  opacity: 0.5;
}

.corner-deco::before {
  width: 100%;
  height: 1px;
}

.corner-deco::after {
  width: 1px;
  height: 100%;
}

.corner-deco.tl {
  top: 8px;
  left: 8px;
}

.corner-deco.tl::before { top: 0; left: 0; }
.corner-deco.tl::after { top: 0; left: 0; }

.corner-deco.br {
  bottom: 8px;
  right: 8px;
}

.corner-deco.br::before { bottom: 0; right: 0; }
.corner-deco.br::after { bottom: 0; right: 0; }

.room-card:hover,
.room-card.hovered {
  transform: translateY(-4px);
  box-shadow:
    0 0 20px rgba(0, 255, 255, 0.2),
    0 0 40px rgba(255, 0, 255, 0.1);
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.room-id {
  display: flex;
  align-items: center;
  gap: 6px;
}

.id-label {
  font-family: var(--font-mono);
  font-size: 9px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.id-value {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--neon-cyan);
  letter-spacing: 1px;
}

.room-status {
  font-family: var(--font-display);
  font-size: 10px;
  padding: 4px 10px;
  font-weight: 600;
  letter-spacing: 1px;
}

.status-waiting {
  background: rgba(0, 255, 102, 0.15);
  color: var(--neon-green);
  border: 1px solid rgba(0, 255, 102, 0.3);
}

.status-ready {
  background: rgba(255, 204, 0, 0.15);
  color: var(--neon-yellow);
  border: 1px solid rgba(255, 204, 0, 0.3);
}

.status-playing {
  background: rgba(0, 136, 255, 0.15);
  color: var(--neon-blue);
  border: 1px solid rgba(0, 136, 255, 0.3);
  animation: statusPulse 1.5s ease infinite;
}

@keyframes statusPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.status-finished {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.room-name {
  color: #fff;
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 16px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
}

/* 房间信息 */
.room-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.info-icon {
  font-size: 14px;
  opacity: 0.8;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-label {
  font-family: var(--font-mono);
  font-size: 9px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.info-value {
  font-family: var(--font-body);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
}

/* 时间显示 */
.room-time {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.3);
  border-left: 2px solid var(--neon-magenta);
}

.time-label {
  font-family: var(--font-mono);
  font-size: 9px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.time-value {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--neon-cyan);
}

/* 操作按钮 */
.room-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  opacity: 0;
  transform: translateY(10px);
  transition: all 0.3s ease;
}

.room-actions.visible {
  opacity: 1;
  transform: translateY(0);
}

.join-btn,
.spectate-btn {
  flex: 1;
  padding: 10px 12px;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.spectate-btn {
  border-color: rgba(0, 136, 255, 0.4);
  color: var(--neon-blue);
}

.spectate-btn:hover {
  background: rgba(0, 136, 255, 0.15);
  border-color: var(--neon-blue);
  box-shadow: 0 0 15px rgba(0, 136, 255, 0.3);
}

/* 排位标识 */
.rated-badge {
  position: absolute;
  top: 8px;
  right: -24px;
  background: linear-gradient(135deg, var(--cyber-gold), var(--neon-orange));
  color: var(--cyber-black);
  font-size: 9px;
  padding: 4px 28px;
  transform: rotate(45deg);
  font-family: var(--font-display);
  font-weight: 700;
  letter-spacing: 1px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.badge-icon {
  font-size: 10px;
}

/* 响应式 */
@media (max-width: 768px) {
  .room-info {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .room-actions {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 480px) {
  .room-card {
    padding: 12px;
  }

  .room-name {
    font-size: 14px;
  }

  .join-btn,
  .spectate-btn {
    padding: 8px 10px;
    font-size: 11px;
  }

  .btn-text {
    display: none;
  }
}
</style>
