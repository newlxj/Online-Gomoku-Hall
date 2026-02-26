<script setup lang="ts">
import { computed } from 'vue'
import type { OnlineUserInfo } from '@/types/websocket'
import { getRankByScore } from '@/types/multiplayer'

const props = defineProps<{
  users: OnlineUserInfo[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

// 按状态分组
const lobbyUsers = computed(() => {
  return props.users.filter(u => u.status === 'lobby')
})

const playingUsers = computed(() => {
  return props.users.filter(u => u.status === 'playing')
})

const spectatingUsers = computed(() => {
  return props.users.filter(u => u.status === 'spectating')
})

// 获取段位颜色
function getRankColor(score: number): string {
  const rank = getRankByScore(score)
  return rank.color
}

function close() {
  emit('close')
}
</script>

<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content glass-panel">
      <div class="modal-header">
        <h3 class="modal-title">
          <span class="title-icon">◉</span>
          <span>ONLINE PLAYERS</span>
          <span class="online-count">{{ users.length }}</span>
        </h3>
        <button class="close-btn" @click="close">
          <span>✕</span>
        </button>
      </div>

      <div class="users-section">
        <!-- 在游戏中 -->
        <div v-if="playingUsers.length > 0" class="status-group">
          <div class="group-header">
            <span class="status-icon playing">⚔</span>
            <span class="status-label">IN GAME</span>
            <span class="status-count">{{ playingUsers.length }}</span>
          </div>
          <div class="user-list">
            <div v-for="user in playingUsers" :key="user.id" class="user-item playing">
              <div class="user-main">
                <span class="user-alias">{{ user.alias }}</span>
                <span class="user-room" v-if="user.roomName">{{ user.roomName }}</span>
              </div>
              <div class="user-stats">
                <span class="user-rank" :style="{ color: getRankColor(user.score) }">{{ user.rank }}</span>
                <span class="user-score">{{ user.score }} PTS</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 观战中 -->
        <div v-if="spectatingUsers.length > 0" class="status-group">
          <div class="group-header">
            <span class="status-icon spectating">👁</span>
            <span class="status-label">SPECTATING</span>
            <span class="status-count">{{ spectatingUsers.length }}</span>
          </div>
          <div class="user-list">
            <div v-for="user in spectatingUsers" :key="user.id" class="user-item spectating">
              <div class="user-main">
                <span class="user-alias">{{ user.alias }}</span>
                <span class="user-room" v-if="user.roomName">{{ user.roomName }}</span>
              </div>
              <div class="user-stats">
                <span class="user-rank" :style="{ color: getRankColor(user.score) }">{{ user.rank }}</span>
                <span class="user-score">{{ user.score }} PTS</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 大厅中 -->
        <div v-if="lobbyUsers.length > 0" class="status-group">
          <div class="group-header">
            <span class="status-icon lobby">◈</span>
            <span class="status-label">IN LOBBY</span>
            <span class="status-count">{{ lobbyUsers.length }}</span>
          </div>
          <div class="user-list">
            <div v-for="user in lobbyUsers" :key="user.id" class="user-item lobby">
              <div class="user-main">
                <span class="user-alias">{{ user.alias }}</span>
                <span class="user-room" v-if="user.roomName">{{ user.roomName }}</span>
              </div>
              <div class="user-stats">
                <span class="user-rank" :style="{ color: getRankColor(user.score) }">{{ user.rank }}</span>
                <span class="user-score">{{ user.score }} PTS</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 无用户 -->
        <div v-if="users.length === 0" class="no-users">
          <span class="no-users-icon">◇</span>
          <span class="no-users-text">No players online</span>
        </div>
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
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  animation: scaleIn 0.3s ease;
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

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  letter-spacing: 2px;
}

.title-icon {
  color: var(--neon-cyan);
  text-shadow: 0 0 10px var(--neon-cyan);
}

.online-count {
  background: var(--neon-cyan);
  color: var(--cyber-black);
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 700;
}

.close-btn {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.6);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
}

.close-btn:hover {
  border-color: var(--neon-magenta);
  color: var(--neon-magenta);
}

.users-section {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.status-group {
  margin-bottom: 20px;
}

.status-group:last-child {
  margin-bottom: 0;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.status-icon {
  font-size: 14px;
}

.status-icon.playing {
  color: var(--neon-orange);
}

.status-icon.spectating {
  color: var(--neon-blue);
}

.status-icon.lobby {
  color: var(--neon-green);
}

.status-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  letter-spacing: 2px;
}

.status-count {
  margin-left: auto;
  font-family: var(--font-display);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.3);
  border-left: 3px solid;
  transition: all 0.2s;
}

.user-item:hover {
  background: rgba(0, 0, 0, 0.5);
}

.user-item.playing {
  border-color: var(--neon-orange);
}

.user-item.spectating {
  border-color: var(--neon-blue);
}

.user-item.lobby {
  border-color: var(--neon-green);
}

.user-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-alias {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.user-room {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
}

.user-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-rank {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 600;
}

.user-score {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
}

.no-users {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
  color: rgba(255, 255, 255, 0.4);
}

.no-users-icon {
  font-size: 32px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.no-users-text {
  font-family: var(--font-mono);
  font-size: 12px;
  letter-spacing: 2px;
}
</style>
