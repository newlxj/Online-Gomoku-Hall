<script setup lang="ts">
import type { GameHistory } from '@/types/game'

defineProps<{
  histories: GameHistory[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'replay', history: GameHistory): void
  (e: 'delete', id: string): void
  (e: 'clearAll'): void
}>()

function getModeText(mode: string): string {
  return mode === 'ai' ? '人机对战' : '在线对战'
}

function getDifficultyText(difficulty?: string): string {
  if (!difficulty) return '-'
  const map: Record<string, string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return map[difficulty] || difficulty
}

function getWinnerText(winner: string): string {
  const map: Record<string, string> = {
    black: '黑子胜',
    white: '白子胜',
    draw: '平局'
  }
  return map[winner] || winner
}

function getWinnerColor(winner: string): string {
  const map: Record<string, string> = {
    black: '#333',
    white: '#fff',
    draw: '#888'
  }
  return map[winner] || '#888'
}

function getReasonText(reason?: string): string {
  const map: Record<string, string> = {
    win: '胜利',
    disconnect: '对手断线',
    timeout: '超时',
    leave: '对手离开'
  }
  return reason ? map[reason] || reason : ''
}
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="history-modal glass-panel">
      <div class="modal-header">
        <h2>📜 游戏记录</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-content">
        <div v-if="histories.length === 0" class="empty-state">
          <p>暂无游戏记录</p>
        </div>

        <div v-else class="history-list">
          <div
            v-for="history in histories"
            :key="history.id"
            class="history-item glass-panel-dark"
          >
            <div class="history-info">
              <div class="history-date">{{ history.date }}</div>
              <div class="history-details">
                <span class="tag" :class="{ 'multiplayer': history.mode === 'pvp' }">
                  {{ getModeText(history.mode) }}
                </span>
                <span v-if="history.difficulty" class="tag">
                  {{ getDifficultyText(history.difficulty) }}
                </span>
                <span class="tag winner-tag" :style="{ background: getWinnerColor(history.winner), color: history.winner === 'black' ? '#fff' : '#333' }">
                  {{ getWinnerText(history.winner) }}
                </span>
              </div>
              <!-- 多人游戏额外信息 -->
              <div v-if="history.multiplayerInfo" class="multiplayer-info">
                <div class="players-row">
                  <span class="player black">◆ {{ history.multiplayerInfo.blackPlayer }}</span>
                  <span class="vs">VS</span>
                  <span class="player white">◇ {{ history.multiplayerInfo.whitePlayer }}</span>
                </div>
                <div class="result-row" v-if="history.multiplayerInfo.winnerAlias">
                  <span class="winner-name">{{ history.multiplayerInfo.winnerAlias }}</span>
                  <span class="score-change" v-if="history.multiplayerInfo.scoreChanged > 0">
                    +{{ history.multiplayerInfo.scoreChanged }} PTS
                  </span>
                  <span class="reason" v-if="history.multiplayerInfo.reason !== 'win'">
                    ({{ getReasonText(history.multiplayerInfo.reason) }})
                  </span>
                </div>
              </div>
              <div class="history-moves">共 {{ history.moves.length }} 步</div>
            </div>
            <div class="history-actions">
              <button class="action-btn replay-btn" @click="emit('replay', history)">
                ▶ 复盘
              </button>
              <button class="action-btn delete-btn" @click="emit('delete', history.id)">
                🗑
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="glass-btn glass-btn-danger" @click="emit('clearAll')" :disabled="histories.length === 0">
          清空所有记录
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.history-modal {
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  animation: scaleIn 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h2 {
  color: white;
  font-size: 20px;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 28px;
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: white;
}

.modal-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.5);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.history-info {
  flex: 1;
}

.history-date {
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
  margin-bottom: 8px;
}

.history-details {
  display: flex;
  gap: 8px;
  margin-bottom: 6px;
}

.tag {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  background: rgba(0, 217, 255, 0.2);
  color: #00d9ff;
}

.tag.multiplayer {
  background: rgba(0, 255, 102, 0.2);
  color: #00ff66;
}

.winner-tag {
  font-weight: 600;
}

.multiplayer-info {
  margin-top: 8px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  margin-bottom: 6px;
}

.players-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}

.player {
  font-size: 13px;
  font-weight: 500;
}

.player.black {
  color: #888;
}

.player.white {
  color: #fff;
}

.vs {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
}

.result-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.winner-name {
  font-size: 12px;
  color: var(--neon-cyan);
  font-weight: 600;
}

.score-change {
  font-size: 11px;
  color: var(--cyber-gold);
  font-weight: 600;
}

.reason {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
}

.history-moves {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.history-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.replay-btn {
  background: rgba(0, 217, 255, 0.2);
  color: #00d9ff;
}

.replay-btn:hover {
  background: rgba(0, 217, 255, 0.4);
}

.delete-btn {
  background: rgba(255, 100, 100, 0.2);
  color: #ff6464;
}

.delete-btn:hover {
  background: rgba(255, 100, 100, 0.4);
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
