<script setup lang="ts">
import { computed } from 'vue'
import type { Player } from '@/types/multiplayer'
import { getRankByScore } from '@/types/multiplayer'

const props = defineProps<{
  player: Player
  isCurrentTurn: boolean
  isMe: boolean
}>()

const rankInfo = computed(() => getRankByScore(props.player.score))

const pieceIcon = computed(() => {
  return props.player.pieceType === 1 ? '◆' : '◇'
})

const pieceName = computed(() => {
  return props.player.pieceType === 1 ? 'BLACK' : 'WHITE'
})

const statusText = computed(() => {
  if (!props.player.isConnected) return 'OFFLINE'
  if (props.player.isReady) return 'READY'
  return 'NOT READY'
})
</script>

<template>
  <div class="player-panel glass-panel-dark" :class="{ active: isCurrentTurn, me: isMe }">
    <!-- 顶部装饰线 -->
    <div class="panel-top-line" v-if="isCurrentTurn"></div>

    <!-- 棋子图标 -->
    <div class="piece-icon" :class="player.pieceType === 1 ? 'black' : 'white'">
      {{ pieceIcon }}
    </div>

    <!-- 玩家信息 -->
    <div class="player-info">
      <div class="alias-row">
        <span class="alias">{{ player.alias }}</span>
        <span v-if="isMe" class="me-badge">YOU</span>
      </div>
      <div class="piece-name">{{ pieceName }}</div>
      <div class="rank-row">
        <span class="rank-icon">{{ rankInfo.icon }}</span>
        <span class="rank-title" :style="{ color: rankInfo.color }">{{ rankInfo.title }}</span>
      </div>
      <div class="score-row">
        <span class="score-value">{{ player.score }}</span>
        <span class="score-label">PTS</span>
      </div>
    </div>

    <!-- 状态指示 -->
    <div class="status" :class="{ connected: player.isConnected, ready: player.isReady }">
      <span class="status-dot"></span>
      <span class="status-text">{{ statusText }}</span>
    </div>

    <!-- 当前回合指示 -->
    <div v-if="isCurrentTurn" class="turn-indicator">
      <span class="turn-icon">▸</span>
      <span class="turn-text">YOUR TURN</span>
    </div>
  </div>
</template>

<style scoped>
.player-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  min-width: 160px;
  position: relative;
  transition: all 0.3s ease;
}

/* 顶部装饰线 */
.panel-top-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta));
  animation: lineGlow 1.5s ease infinite;
}

@keyframes lineGlow {
  0%, 100% {
    box-shadow: 0 0 10px var(--neon-cyan);
  }
  50% {
    box-shadow: 0 0 20px var(--neon-magenta);
  }
}

.player-panel.active {
  border-color: var(--neon-cyan);
  box-shadow:
    0 0 30px rgba(0, 255, 255, 0.3),
    inset 0 0 20px rgba(0, 255, 255, 0.05);
}

.player-panel.me {
  background: rgba(0, 255, 255, 0.05);
}

/* 棋子图标 */
.piece-icon {
  font-size: 40px;
  margin-bottom: 12px;
  text-shadow: 0 0 20px currentColor;
}

.piece-icon.black {
  color: var(--neon-cyan);
}

.piece-icon.white {
  color: var(--neon-magenta);
}

.player-info {
  text-align: center;
}

.alias-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 4px;
}

.alias {
  color: #fff;
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 1px;
}

.me-badge {
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-magenta));
  color: var(--cyber-black);
  font-family: var(--font-display);
  font-size: 9px;
  padding: 2px 6px;
  font-weight: 700;
  letter-spacing: 1px;
}

.piece-name {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.rank-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-bottom: 4px;
}

.rank-icon {
  font-size: 14px;
}

.rank-title {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 600;
}

.score-row {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 4px;
}

.score-value {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--cyber-gold);
  text-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
}

.score-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

/* 状态指示 */
.status {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.3);
  border-left: 2px solid rgba(255, 255, 255, 0.2);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 0;
  background: rgba(255, 255, 255, 0.4);
}

.status-text {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.status.connected.ready {
  border-color: var(--neon-green);
  background: rgba(0, 255, 102, 0.1);
}

.status.connected.ready .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 10px var(--neon-green);
  animation: dotPulse 1s ease infinite;
}

.status.connected.ready .status-text {
  color: var(--neon-green);
}

@keyframes dotPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status:not(.connected) {
  border-color: var(--neon-pink);
  background: rgba(255, 0, 80, 0.1);
}

.status:not(.connected) .status-dot {
  background: var(--neon-pink);
}

.status:not(.connected) .status-text {
  color: var(--neon-pink);
}

/* 当前回合指示 */
.turn-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  padding: 8px 14px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid var(--neon-cyan);
  animation: turnPulse 1s ease infinite;
}

@keyframes turnPulse {
  0%, 100% {
    box-shadow: 0 0 10px rgba(0, 255, 255, 0.3);
  }
  50% {
    box-shadow: 0 0 20px rgba(0, 255, 255, 0.5);
  }
}

.turn-icon {
  color: var(--neon-cyan);
  font-size: 12px;
  animation: blink 0.8s ease infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.turn-text {
  font-family: var(--font-display);
  font-size: 11px;
  color: var(--neon-cyan);
  font-weight: 600;
  letter-spacing: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .player-panel {
    padding: 16px;
    min-width: auto;
  }

  .piece-icon {
    font-size: 32px;
  }

  .alias {
    font-size: 14px;
  }

  .score-value {
    font-size: 16px;
  }
}
</style>
