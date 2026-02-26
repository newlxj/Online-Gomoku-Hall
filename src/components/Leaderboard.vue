<script setup lang="ts">
import { computed } from 'vue'
import { useLeaderboard } from '@/composables/useLeaderboard'
import { getRankByScore } from '@/types/multiplayer'

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { entries, isLoading, top10, refresh } = useLeaderboard()

const displayEntries = computed(() => entries.value.slice(0, 50))

function getRankBadge(rank: number): string {
  switch (rank) {
    case 1: return '★'
    case 2: return '◆'
    case 3: return '◇'
    default: return `#${rank}`
  }
}

function getRankStyle(rank: number) {
  switch (rank) {
    case 1: return 'gold'
    case 2: return 'silver'
    case 3: return 'bronze'
    default: return ''
  }
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="modal-overlay" @click.self="handleClose">
    <div class="modal-content glass-panel leaderboard-modal">
      <!-- 顶部装饰线 -->
      <div class="modal-header-line"></div>

      <div class="modal-header">
        <div class="header-left">
          <span class="header-icon">◈</span>
          <h3 class="modal-title">LEADERBOARD</h3>
        </div>
        <div class="header-actions">
          <button class="refresh-btn glass-btn" @click="refresh" :disabled="isLoading">
            <span class="btn-icon" :class="{ spinning: isLoading }">⟳</span>
          </button>
          <button class="close-btn" @click="handleClose">✕</button>
        </div>
      </div>

      <div class="modal-body">
        <!-- Top 3 特效展示 -->
        <div v-if="entries.length >= 3" class="top3-section">
          <!-- 第二名 -->
          <div class="top3-item silver" v-if="entries[1]">
            <div class="top3-rank-num">02</div>
            <div class="top3-avatar">{{ entries[1].alias.charAt(0) }}</div>
            <div class="top3-medal">◆</div>
            <div class="top3-alias">{{ entries[1].alias }}</div>
            <div class="top3-score">{{ entries[1].score }} PTS</div>
            <div class="top3-rank" :style="{ color: getRankByScore(entries[1].score).color }">
              {{ entries[1].rankTitle }}
            </div>
          </div>

          <!-- 第一名 -->
          <div class="top3-item gold" v-if="entries[0]">
            <div class="crown">★</div>
            <div class="top3-rank-num">01</div>
            <div class="top3-avatar">{{ entries[0].alias.charAt(0) }}</div>
            <div class="top3-medal">★</div>
            <div class="top3-alias">{{ entries[0].alias }}</div>
            <div class="top3-score">{{ entries[0].score }} PTS</div>
            <div class="top3-rank" :style="{ color: getRankByScore(entries[0].score).color }">
              {{ entries[0].rankTitle }}
            </div>
          </div>

          <!-- 第三名 -->
          <div class="top3-item bronze" v-if="entries[2]">
            <div class="top3-rank-num">03</div>
            <div class="top3-avatar">{{ entries[2].alias.charAt(0) }}</div>
            <div class="top3-medal">◇</div>
            <div class="top3-alias">{{ entries[2].alias }}</div>
            <div class="top3-score">{{ entries[2].score }} PTS</div>
            <div class="top3-rank" :style="{ color: getRankByScore(entries[2].score).color }">
              {{ entries[2].rankTitle }}
            </div>
          </div>
        </div>

        <!-- 完整排行表 -->
        <div class="leaderboard-table">
          <div class="table-header">
            <span class="col-rank">RANK</span>
            <span class="col-alias">PLAYER</span>
            <span class="col-record">W/L</span>
            <span class="col-score">SCORE</span>
            <span class="col-rank-title">TIER</span>
          </div>

          <div class="table-body">
            <div
              v-for="entry in displayEntries"
              :key="entry.alias"
              class="table-row"
              :class="getRankStyle(entry.rank)"
            >
              <span class="col-rank">{{ getRankBadge(entry.rank) }}</span>
              <span class="col-alias">{{ entry.alias }}</span>
              <span class="col-record">
                <span class="wins">{{ entry.wins }}</span>
                <span class="divider">/</span>
                <span class="losses">{{ entry.losses }}</span>
              </span>
              <span class="col-score">{{ entry.score }}</span>
              <span class="col-rank-title" :style="{ color: getRankByScore(entry.score).color }">
                {{ getRankByScore(entry.score).icon }} {{ entry.rankTitle }}
              </span>
            </div>

            <div v-if="entries.length === 0" class="empty-state">
              <span class="empty-icon">◇</span>
              <span class="empty-text">NO DATA AVAILABLE</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 角落装饰 -->
      <div class="corner-deco tl"></div>
      <div class="corner-deco br"></div>
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
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.leaderboard-modal {
  width: 90%;
  max-width: 650px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}

/* 顶部装饰线 */
.modal-header-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta), var(--neon-cyan));
  animation: lineGlow 3s ease infinite;
}

@keyframes lineGlow {
  0%, 100% {
    box-shadow: 0 0 10px var(--neon-cyan);
  }
  50% {
    box-shadow: 0 0 20px var(--neon-magenta);
  }
}

/* 角落装饰 */
.corner-deco {
  position: absolute;
  width: 20px;
  height: 20px;
  pointer-events: none;
}

.corner-deco::before,
.corner-deco::after {
  content: '';
  position: absolute;
  background: var(--neon-cyan);
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
  top: 10px;
  left: 10px;
}

.corner-deco.tl::before { top: 0; left: 0; }
.corner-deco.tl::after { top: 0; left: 0; }

.corner-deco.br {
  bottom: 10px;
  right: 10px;
}

.corner-deco.br::before { bottom: 0; right: 0; }
.corner-deco.br::after { bottom: 0; right: 0; }

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  color: var(--neon-cyan);
  font-size: 18px;
  text-shadow: 0 0 10px var(--neon-cyan);
}

.modal-title {
  font-family: var(--font-display);
  color: #fff;
  font-size: 20px;
  font-weight: 700;
  margin: 0;
  letter-spacing: 3px;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.refresh-btn {
  width: 36px;
  height: 36px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.refresh-btn .btn-icon {
  font-size: 18px;
}

.refresh-btn .btn-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.close-btn {
  background: transparent;
  border: 1px solid rgba(255, 0, 80, 0.4);
  color: var(--neon-pink);
  width: 36px;
  height: 36px;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(255, 0, 80, 0.15);
  box-shadow: 0 0 15px rgba(255, 0, 80, 0.4);
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 0 24px 24px;
}

/* Top 3 特效 */
.top3-section {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  gap: 16px;
  margin-bottom: 24px;
  padding: 24px 16px;
}

.top3-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  padding: 20px 16px;
  transition: all 0.3s;
}

.top3-item:hover {
  transform: translateY(-4px);
}

.top3-item.gold {
  background: linear-gradient(180deg, rgba(255, 215, 0, 0.15), transparent);
  border: 1px solid rgba(255, 215, 0, 0.4);
  order: 2;
  padding: 28px 20px;
}

.top3-item.silver {
  background: linear-gradient(180deg, rgba(192, 192, 192, 0.1), transparent);
  border: 1px solid rgba(192, 192, 192, 0.3);
  order: 1;
}

.top3-item.bronze {
  background: linear-gradient(180deg, rgba(205, 127, 50, 0.1), transparent);
  border: 1px solid rgba(205, 127, 50, 0.3);
  order: 3;
}

.top3-rank-num {
  font-family: var(--font-display);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.crown {
  font-size: 24px;
  position: absolute;
  top: -5px;
  color: var(--cyber-gold);
  text-shadow: 0 0 20px var(--cyber-gold);
  animation: float 2s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.top3-avatar {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 700;
  color: var(--cyber-black);
  margin-bottom: 8px;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-magenta));
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
}

.top3-item.gold .top3-avatar {
  width: 56px;
  height: 56px;
  font-size: 26px;
  box-shadow: 0 0 30px rgba(255, 215, 0, 0.4);
}

.top3-medal {
  font-size: 20px;
  margin-bottom: 4px;
}

.top3-item.gold .top3-medal {
  color: var(--cyber-gold);
  text-shadow: 0 0 10px var(--cyber-gold);
}

.top3-item.silver .top3-medal {
  color: var(--cyber-silver);
  text-shadow: 0 0 10px var(--cyber-silver);
}

.top3-item.bronze .top3-medal {
  color: var(--cyber-bronze);
  text-shadow: 0 0 10px var(--cyber-bronze);
}

.top3-alias {
  color: #fff;
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 4px;
  letter-spacing: 1px;
}

.top3-score {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  color: var(--neon-cyan);
  text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
  margin-bottom: 4px;
}

.top3-rank {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 500;
}

/* 排行表 */
.leaderboard-table {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.2);
}

.table-header {
  display: grid;
  grid-template-columns: 60px 1fr 80px 80px 100px;
  padding: 12px 16px;
  background: rgba(0, 255, 255, 0.05);
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.table-header span {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--neon-cyan);
  letter-spacing: 1px;
}

.table-body {
  max-height: 280px;
  overflow-y: auto;
}

.table-row {
  display: grid;
  grid-template-columns: 60px 1fr 80px 80px 100px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.05);
  color: #fff;
  font-size: 13px;
  transition: all 0.3s;
}

.table-row:hover {
  background: rgba(0, 255, 255, 0.05);
}

.table-row.gold {
  background: rgba(255, 215, 0, 0.08);
  border-left: 2px solid var(--cyber-gold);
}

.table-row.silver {
  background: rgba(192, 192, 192, 0.05);
  border-left: 2px solid var(--cyber-silver);
}

.table-row.bronze {
  background: rgba(205, 127, 50, 0.05);
  border-left: 2px solid var(--cyber-bronze);
}

.col-rank {
  font-family: var(--font-display);
  font-weight: 600;
}

.table-row.gold .col-rank {
  color: var(--cyber-gold);
}

.table-row.silver .col-rank {
  color: var(--cyber-silver);
}

.table-row.bronze .col-rank {
  color: var(--cyber-bronze);
}

.col-alias {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-display);
  font-weight: 500;
}

.col-record {
  display: flex;
  gap: 4px;
  font-family: var(--font-mono);
}

.wins {
  color: var(--neon-green);
}

.divider {
  color: rgba(255, 255, 255, 0.3);
}

.losses {
  color: var(--neon-pink);
}

.col-score {
  font-family: var(--font-display);
  font-weight: 700;
  color: var(--neon-cyan);
}

.col-rank-title {
  font-family: var(--font-display);
  font-size: 12px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 40px;
  color: rgba(0, 255, 255, 0.2);
}

.empty-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .leaderboard-modal {
    max-height: 90vh;
  }

  .modal-header {
    padding: 16px;
  }

  .modal-body {
    padding: 0 16px 16px;
  }

  .table-header,
  .table-row {
    grid-template-columns: 50px 1fr 70px 70px;
  }

  .col-rank-title {
    display: none;
  }

  .top3-section {
    flex-wrap: wrap;
    gap: 12px;
    padding: 16px 8px;
  }

  .top3-item.gold {
    order: 1;
    width: 100%;
    padding: 20px 16px;
  }

  .top3-item.silver,
  .top3-item.bronze {
    flex: 1;
    min-width: 120px;
  }
}

@media (max-width: 480px) {
  .modal-title {
    font-size: 16px;
    letter-spacing: 2px;
  }

  .top3-section {
    flex-direction: column;
    align-items: center;
  }

  .top3-item {
    width: 100%;
    max-width: 200px;
  }

  .top3-item.gold {
    max-width: 220px;
  }
}
</style>
