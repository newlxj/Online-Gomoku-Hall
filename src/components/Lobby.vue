<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoom } from '@/composables/useRoom'
import { useLeaderboard } from '@/composables/useLeaderboard'
import { useAlias } from '@/composables/useAlias'
import { useWebSocket } from '@/composables/useWebSocket'
import RoomCard from './RoomCard.vue'
import CreateRoomModal from './CreateRoomModal.vue'
import Leaderboard from './Leaderboard.vue'
import AliasModal from './AliasModal.vue'
import OnlineUsersModal from './OnlineUsersModal.vue'
import type { RoomSettings, RoomInfo } from '@/types/multiplayer'
import type { OnlineUserInfo } from '@/types/websocket'

const emit = defineEmits<{
  (e: 'joinRoom', roomId: string): void
  (e: 'spectate', roomId: string): void
}>()

const {
  roomList,
  createRoom,
  joinRoom,
  spectate,
  refreshRoomList,
} = useRoom()

const { refresh: refreshLeaderboard, entries, getRankByScore } = useLeaderboard()
const { currentAlias } = useAlias()

// 当前用户积分
const myScore = computed(() => {
  if (!currentAlias.value) return 0
  const entry = entries.value.find(e => e.alias === currentAlias.value)
  return entry?.score || 0
})

// 当前用户段位
const myRank = computed(() => getRankByScore(myScore.value))

// UI状态
const showCreateModal = ref(false)
const showLeaderboard = ref(false)
const showAliasModal = ref(false)
const showOnlineUsersModal = ref(false)
const isRefreshing = ref(false)

// 在线用户状态
const onlineUsers = ref<OnlineUserInfo[]>([])
const onlineCount = computed(() => onlineUsers.value.length)

// WebSocket订阅
const { subscribe, send } = useWebSocket()
const { setCurrentAlias } = useAlias()

// 订阅在线用户更新
const unsubscribeOnlineUsers = subscribe<{ count: number; users: OnlineUserInfo[] }>(
  'online_users',
  (payload) => {
    console.log('[Lobby] Received online_users:', payload)
    onlineUsers.value = payload.users
  }
)

// 订阅进入大厅确认消息
const unsubscribeLobbyEntered = subscribe<{ alias: string; playerId: string }>(
  'lobby_entered',
  (payload) => {
    console.log('[Lobby] Received lobby_entered:', payload)
    // 如果服务器分配了别名，更新本地别名
    if (payload.alias && payload.alias !== currentAlias.value) {
      setCurrentAlias(payload.alias)
    }
  }
)

// 待执行的操作（别名确认后执行）
const pendingAction = ref<(() => void) | null>(null)

// 创建房间
function handleCreateRoom(name: string, settings: RoomSettings) {
  createRoom(name, settings)
  showCreateModal.value = false
}

// 加入房间
function handleJoinRoom(room: RoomInfo) {
  // 先检查别名
  checkAliasAndProceed(() => {
    console.log('[Lobby] Joining room:', room.id)
    joinRoom(room.id)
    emit('joinRoom', room.id)
  })
}

// 观战
function handleSpectate(roomId: string) {
  // 先检查别名
  checkAliasAndProceed(() => {
    console.log('[Lobby] Spectating room:', roomId)
    spectate(roomId)
    emit('spectate', roomId)
  })
}

// 快速匹配
function handleQuickMatch() {
  // 找出所有可加入的房间
  const availableRooms = roomList.value.filter(r => r.status === 'waiting' && r.playerCount < 2)
  if (availableRooms.length > 0) {
    // 随机选择一个房间
    const randomIndex = Math.floor(Math.random() * availableRooms.length)
    handleJoinRoom(availableRooms[randomIndex])
  } else {
    alert('暂无可用房间，请创建新房间')
  }
}

// 刷新列表
async function handleRefresh() {
  isRefreshing.value = true
  refreshRoomList()
  setTimeout(() => {
    isRefreshing.value = false
  }, 500)
}

// 检查别名 - 确保别名可用后才执行操作
function checkAliasAndProceed(action: () => void) {
  console.log('[Lobby] checkAliasAndProceed called, currentAlias:', currentAlias.value)
  if (!currentAlias.value || currentAlias.value.trim() === '') {
    console.log('[Lobby] No alias, setting pendingAction and showing modal')
    pendingAction.value = action
    showAliasModal.value = true
    return
  }
  console.log('[Lobby] Alias exists, executing action')
  // 直接执行，不需要延迟
  action()
}

// 别名确认后
function handleAliasConfirm(alias: string) {
  console.log('[Lobby] handleAliasConfirm called with alias:', alias)
  console.log('[Lobby] pendingAction:', pendingAction.value)
  showAliasModal.value = false
  if (pendingAction.value) {
    const action = pendingAction.value
    pendingAction.value = null
    // 延迟执行，确保别名已保存
    setTimeout(() => {
      console.log('[Lobby] Executing pending action')
      action()
    }, 100)
  }
}

// 初始化
onMounted(() => {
  // 进入大厅时注册在线用户
  if (currentAlias.value) {
    send('enter_lobby', { alias: currentAlias.value })
  } else {
    // 如果没有别名，发送空别名让服务器生成
    send('enter_lobby', { alias: '' })
  }

  refreshRoomList()
  refreshLeaderboard()
  // 请求在线用户列表
  send('get_online_users', {})
})

// 清理
onUnmounted(() => {
  // 离开大厅时注销
  send('leave_lobby', {})
  unsubscribeOnlineUsers()
  unsubscribeLobbyEntered()
})
</script>

<template>
  <div class="lobby glass-panel">
    <!-- 顶部装饰线 -->
    <div class="lobby-header-line"></div>

    <!-- 顶部标题和操作 -->
    <div class="lobby-header">
      <div class="header-left">
        <h2 class="lobby-title">
          <span class="title-icon">◈</span>
          <span class="title-text">ONLINE LOBBY</span>
        </h2>
        <span class="title-subtitle">多人在线大厅</span>
      </div>
      <div class="header-right">
        <!-- 在线人数 -->
        <button
          class="online-count-btn"
          @click="showOnlineUsersModal = true"
          title="Click to view online players"
        >
          <span class="online-dot"></span>
          <span class="online-label">ONLINE</span>
          <span class="online-number">{{ onlineCount }}</span>
        </button>
        <div class="lobby-actions">
          <button
            class="glass-btn action-btn"
            @click="showLeaderboard = true"
          >
            <span class="btn-icon">🏆</span>
            <span class="btn-text">排行榜</span>
          </button>
          <button
            class="glass-btn action-btn refresh-btn"
            :class="{ 'is-refreshing': isRefreshing }"
            :disabled="isRefreshing"
            @click="handleRefresh"
          >
            <span class="btn-icon" :class="{ 'spinning': isRefreshing }">🔄</span>
            <span class="btn-text">{{ isRefreshing ? '刷新中' : '刷新' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 当前别名显示 -->
    <div class="alias-display glass-panel-dark" v-if="currentAlias">
      <div class="alias-main">
        <span class="alias-label">PLAYER</span>
        <span class="alias-name">{{ currentAlias }}</span>
      </div>
      <div class="alias-stats">
        <div class="stat-item rank-item">
          <span class="stat-icon">{{ myRank.icon }}</span>
          <span class="stat-value" :style="{ color: myRank.color }">{{ myRank.title }}</span>
        </div>
        <div class="stat-item score-item">
          <span class="stat-value">{{ myScore }}</span>
          <span class="stat-label">PTS</span>
        </div>
      </div>
      <button class="change-alias-btn" @click="showAliasModal = true">
        <span class="btn-icon">✎</span>
      </button>
    </div>

    <!-- 创建房间和快速匹配 -->
    <div class="lobby-main-actions">
      <button
        class="glass-btn glass-btn-primary create-btn"
        @click="checkAliasAndProceed(() => showCreateModal = true)"
      >
        <span class="btn-icon">+</span>
        <span class="btn-text">创建房间</span>
        <span class="btn-arrow">→</span>
      </button>
      <button
        class="glass-btn quick-btn"
        @click="checkAliasAndProceed(handleQuickMatch)"
      >
        <span class="btn-icon">⚡</span>
        <span class="btn-text">快速匹配</span>
        <span class="btn-arrow">→</span>
      </button>
    </div>

    <!-- 房间列表 -->
    <div class="room-list-section">
      <div class="section-header">
        <span class="section-icon">▣</span>
        <h3 class="section-title">ROOM LIST</h3>
        <span class="section-count">[{{ roomList.length }}]</span>
        <span class="section-line"></span>
      </div>

      <div v-if="roomList.length === 0" class="no-rooms">
        <div class="no-rooms-icon">◇</div>
        <p class="no-rooms-text">暂无房间</p>
        <p class="no-rooms-hint">创建一个房间开始游戏吧</p>
      </div>

      <div v-else class="room-grid">
        <RoomCard
          v-for="room in roomList"
          :key="room.id"
          :room="room"
          @join="handleJoinRoom(room)"
          @spectate="() => handleSpectate(room.id)"
        />
      </div>
    </div>

    <!-- 创建房间弹窗 -->
    <CreateRoomModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @create="handleCreateRoom"
    />

    <!-- 排行榜弹窗 -->
    <Leaderboard
      v-if="showLeaderboard"
      @close="showLeaderboard = false"
    />

    <!-- 别名设置弹窗 -->
    <AliasModal
      v-if="showAliasModal"
      @close="showAliasModal = false; pendingAction = null"
      @confirm="handleAliasConfirm"
    />

    <!-- 在线用户弹窗 -->
    <OnlineUsersModal
      v-if="showOnlineUsersModal"
      :users="onlineUsers"
      @close="showOnlineUsersModal = false"
    />
  </div>
</template>

<style scoped>
.lobby {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 80vh;
  position: relative;
}

/* 顶部装饰线 */
.lobby-header-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(
    90deg,
    transparent,
    var(--neon-cyan),
    var(--neon-magenta),
    var(--neon-cyan),
    transparent
  );
}

.lobby-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 20px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.lobby-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  color: var(--neon-cyan);
  text-shadow: 0 0 10px var(--neon-cyan);
}

.title-text {
  background: linear-gradient(90deg, #fff, var(--neon-cyan));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 3px;
}

.title-subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-left: 24px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.online-count-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(0, 255, 102, 0.1);
  border: 1px solid rgba(0, 255, 102, 0.3);
  cursor: pointer;
  transition: all 0.3s;
}

.online-count-btn:hover {
  background: rgba(0, 255, 102, 0.2);
  border-color: var(--neon-green);
  box-shadow: 0 0 15px rgba(0, 255, 102, 0.3);
}

.online-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--neon-green);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    box-shadow: 0 0 5px var(--neon-green);
  }
  50% {
    opacity: 0.6;
    box-shadow: 0 0 10px var(--neon-green);
  }
}

.online-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.online-number {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--neon-green);
  text-shadow: 0 0 10px rgba(0, 255, 102, 0.5);
}

.lobby-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
}

.action-btn .btn-icon {
  font-size: 14px;
}

.action-btn .btn-text {
  font-size: 12px;
}

.refresh-btn .btn-icon.spinning {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 别名显示区 */
.alias-display {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 12px 20px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.alias-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.alias-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--neon-cyan);
  letter-spacing: 2px;
  opacity: 0.8;
}

.alias-name {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
}

.alias-stats {
  display: flex;
  gap: 16px;
  margin-left: auto;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.3);
  border-left: 2px solid var(--neon-cyan);
}

.stat-icon {
  font-size: 16px;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 600;
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
}

.score-item {
  border-color: var(--cyber-gold);
}

.score-item .stat-value {
  color: var(--cyber-gold);
  font-size: 18px;
  text-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
}

.change-alias-btn {
  background: transparent;
  border: 1px solid rgba(0, 255, 255, 0.3);
  color: var(--neon-cyan);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 16px;
}

.change-alias-btn:hover {
  background: rgba(0, 255, 255, 0.1);
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.3);
}

/* 主要操作按钮 */
.lobby-main-actions {
  display: flex;
  gap: 16px;
  margin-bottom: 28px;
}

.create-btn,
.quick-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 18px 24px;
  font-size: 14px;
}

.create-btn .btn-icon {
  font-size: 20px;
  font-weight: bold;
}

.quick-btn {
  border-color: var(--neon-yellow);
  color: var(--neon-yellow);
}

.quick-btn:hover {
  background: rgba(255, 255, 0, 0.1);
  border-color: var(--neon-yellow);
  box-shadow: 0 0 15px rgba(255, 255, 0, 0.4);
}

.btn-arrow {
  opacity: 0.5;
  transition: all 0.3s ease;
}

.create-btn:hover .btn-arrow,
.quick-btn:hover .btn-arrow {
  opacity: 1;
  transform: translateX(5px);
}

/* 房间列表区域 */
.room-list-section {
  margin-top: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.section-icon {
  color: var(--neon-magenta);
  font-size: 14px;
  text-shadow: 0 0 10px var(--neon-magenta);
}

.section-title {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 600;
  color: var(--neon-cyan);
  letter-spacing: 2px;
  margin: 0;
}

.section-count {
  font-family: var(--font-mono);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.section-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, var(--neon-cyan), transparent);
  opacity: 0.3;
}

/* 无房间提示 */
.no-rooms {
  text-align: center;
  padding: 60px 20px;
}

.no-rooms-icon {
  font-size: 48px;
  color: rgba(0, 255, 255, 0.3);
  margin-bottom: 16px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.no-rooms-text {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 8px;
}

.no-rooms-hint {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
}

/* 房间网格 */
.room-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .lobby {
    padding: 16px;
    min-height: auto;
  }

  .lobby-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-left {
    width: 100%;
  }

  .lobby-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .alias-display {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .alias-stats {
    margin-left: 0;
    width: 100%;
    justify-content: flex-start;
  }

  .lobby-main-actions {
    flex-direction: column;
  }

  .room-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .lobby-title {
    font-size: 20px;
  }

  .action-btn {
    padding: 8px 12px;
  }

  .action-btn .btn-text {
    display: none;
  }

  .create-btn,
  .quick-btn {
    padding: 14px 20px;
  }
}
</style>
