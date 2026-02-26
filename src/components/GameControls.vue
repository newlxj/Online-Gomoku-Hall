<script setup lang="ts">
import type { GameMode, Difficulty, GameStatus } from '@/types/game'

defineProps<{
  gameStatus: GameStatus
  mode: GameMode
  difficulty: Difficulty
  isMuted: boolean
  canUndo: boolean
}>()

const emit = defineEmits<{
  (e: 'undo'): void
  (e: 'restart'): void
  (e: 'toggleMute'): void
  (e: 'changeDifficulty', difficulty: Difficulty): void
  (e: 'openHistory'): void
  (e: 'backToMenu'): void
}>()

const difficulties: { value: Difficulty; label: string }[] = [
  { value: 'easy', label: 'EASY' },
  { value: 'medium', label: 'MEDIUM' },
  { value: 'hard', label: 'HARD' }
]

const statusTexts: Record<GameStatus, string> = {
  'playing': 'BATTLE',
  'black-win': 'BLACK WINS',
  'white-win': 'WHITE WINS',
  'draw': 'DRAW'
}
</script>

<template>
  <div class="controls glass-panel">
    <!-- 顶部装饰线 -->
    <div class="controls-header-line"></div>

    <!-- 游戏状态 -->
    <div class="status-section">
      <div class="status-label">STATUS</div>
      <div class="status-indicator" :class="{
        'status-playing': gameStatus === 'playing',
        'status-win': gameStatus === 'black-win' || gameStatus === 'white-win',
        'status-draw': gameStatus === 'draw'
      }">
        <span class="status-icon" v-if="gameStatus !== 'playing'">
          {{ gameStatus.includes('win') ? '◆' : '◇' }}
        </span>
        <span class="status-text">{{ statusTexts[gameStatus] }}</span>
      </div>
    </div>

    <!-- 难度选择 (仅AI模式) -->
    <div v-if="mode === 'ai'" class="difficulty-section">
      <div class="section-label">
        <span class="label-icon">&lt;/&gt;</span>
        <span class="label-text">AI LEVEL</span>
      </div>
      <select
        class="glass-select"
        :value="difficulty"
        @change="emit('changeDifficulty', ($event.target as HTMLSelectElement).value as Difficulty)"
      >
        <option v-for="d in difficulties" :key="d.value" :value="d.value">
          {{ d.label }}
        </option>
      </select>
    </div>

    <!-- 控制按钮 -->
    <div class="buttons-section">
      <button
        class="control-btn glass-btn"
        :disabled="!canUndo || gameStatus !== 'playing'"
        @click="emit('undo')"
      >
        <span class="btn-icon">↩</span>
        <span class="btn-text">悔棋</span>
        <span class="btn-label">UNDO</span>
      </button>

      <button class="control-btn glass-btn glass-btn-primary" @click="emit('restart')">
        <span class="btn-icon">⟳</span>
        <span class="btn-text">重新开始</span>
        <span class="btn-label">RESTART</span>
      </button>

      <button class="control-btn glass-btn" @click="emit('toggleMute')">
        <span class="btn-icon">{{ isMuted ? '🔇' : '🔊' }}</span>
        <span class="btn-text">{{ isMuted ? '开启' : '关闭' }}音效</span>
        <span class="btn-label">SOUND</span>
      </button>

      <button class="control-btn glass-btn" @click="emit('openHistory')">
        <span class="btn-icon">☰</span>
        <span class="btn-text">历史记录</span>
        <span class="btn-label">HISTORY</span>
      </button>

      <button class="control-btn glass-btn glass-btn-danger" @click="emit('backToMenu')">
        <span class="btn-icon">←</span>
        <span class="btn-text">返回菜单</span>
        <span class="btn-label">EXIT</span>
      </button>
    </div>

    <!-- 底部装饰 -->
    <div class="controls-footer">
      <div class="footer-line"></div>
      <span class="footer-text">超级五子棋</span>
    </div>
  </div>
</template>

<style scoped>
.controls {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-width: 200px;
  position: relative;
}

/* 顶部装饰线 */
.controls-header-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(
    90deg,
    transparent,
    var(--neon-cyan),
    var(--neon-magenta),
    transparent
  );
}

/* 状态区域 */
.status-section {
  text-align: center;
  padding: 8px 0;
}

.status-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.status-indicator {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  letter-spacing: 2px;
}

.status-playing {
  color: var(--neon-cyan);
  background: rgba(0, 255, 255, 0.05);
  border: 1px solid rgba(0, 255, 255, 0.3);
  animation: playingPulse 2s ease infinite;
}

@keyframes playingPulse {
  0%, 100% {
    box-shadow: inset 0 0 10px rgba(0, 255, 255, 0.1);
  }
  50% {
    box-shadow: inset 0 0 20px rgba(0, 255, 255, 0.2);
  }
}

.status-win {
  color: var(--cyber-gold);
  background: rgba(255, 215, 0, 0.1);
  border: 1px solid rgba(255, 215, 0, 0.4);
  animation: winGlow 1s ease infinite;
}

@keyframes winGlow {
  0%, 100% {
    box-shadow:
      0 0 10px rgba(255, 215, 0, 0.3),
      inset 0 0 10px rgba(255, 215, 0, 0.1);
  }
  50% {
    box-shadow:
      0 0 20px rgba(255, 215, 0, 0.5),
      inset 0 0 20px rgba(255, 215, 0, 0.2);
  }
}

.status-draw {
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.status-icon {
  font-size: 12px;
}

/* 难度选择区域 */
.difficulty-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label-icon {
  font-family: var(--font-mono);
  color: var(--neon-magenta);
  font-size: 10px;
}

.label-text {
  font-family: var(--font-display);
  font-size: 11px;
  color: var(--neon-cyan);
  letter-spacing: 2px;
}

/* 按钮区域 */
.buttons-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.control-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  position: relative;
}

.control-btn .btn-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}

.control-btn .btn-text {
  flex: 1;
  text-align: left;
  font-size: 13px;
}

.control-btn .btn-label {
  font-family: var(--font-mono);
  font-size: 9px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.control-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.control-btn:disabled:hover {
  background: transparent;
  box-shadow: none;
}

/* 底部装饰 */
.controls-footer {
  margin-top: auto;
  padding-top: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.footer-line {
  width: 60%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 255, 255, 0.3),
    transparent
  );
}

.footer-text {
  font-family: var(--font-mono);
  font-size: 9px;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .controls {
    min-width: auto;
    padding: 16px;
  }

  .status-indicator {
    font-size: 14px;
    padding: 10px 16px;
  }

  .control-btn .btn-label {
    display: none;
  }
}

@media (max-width: 480px) {
  .control-btn {
    padding: 10px 12px;
  }

  .control-btn .btn-text {
    font-size: 12px;
  }
}
</style>
