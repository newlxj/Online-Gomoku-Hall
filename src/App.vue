<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import GameBoard from './components/GameBoard.vue'
import GameControls from './components/GameControls.vue'
import GameHistory from './components/GameHistory.vue'
import ReplayModal from './components/ReplayModal.vue'
import Lobby from './components/Lobby.vue'
import GameRoom from './components/GameRoom.vue'
import Leaderboard from './components/Leaderboard.vue'
import { useGame } from './composables/useGame'
import { useAI } from './composables/useAI'
import { useSound } from './composables/useSound'
import { useHistory } from './composables/useHistory'
import { useWebSocket } from './composables/useWebSocket'
import { useRoom } from './composables/useRoom'
import type { GameMode, Difficulty, GameHistory as GameHistoryType, Position } from './types/game'

// 扩展游戏模式类型
type AppMode = GameMode | 'multiplayer' | null

// 游戏状态
const gameMode = ref<AppMode>(null)
const showHistory = ref(false)
const showReplay = ref(false)
const showLeaderboard = ref(false)
const replayHistory = ref<GameHistoryType | null>(null)
const inMultiplayerRoom = ref(false)

// WebSocket连接
const { isConnected } = useWebSocket()

// 多人房间状态
const { isInRoom } = useRoom()

// 监听房间状态变化，自动切换界面
watch(isInRoom, (inRoom) => {
  if (inRoom && gameMode.value === 'multiplayer') {
    inMultiplayerRoom.value = true
    console.log('[App] Entered room, showing GameRoom')
  } else if (!inRoom && inMultiplayerRoom.value) {
    inMultiplayerRoom.value = false
    console.log('[App] Left room, showing Lobby')
  }
})

// 组合式函数
const game = useGame()
const ai = useAI()
const sound = useSound()
const history = useHistory()

// 计算属性
const lastMove = computed<Position | null>(() => {
  const moves = game.moveHistory.value
  return moves.length > 0 ? moves[moves.length - 1].position : null
})

const canUndo = computed(() => {
  return game.moveHistory.value.length > 0
})

// 开始游戏
function startGame(mode: GameMode, difficulty?: Difficulty) {
  gameMode.value = mode
  game.resetGame(mode, difficulty || 'medium')
}

// 开始多人模式
function startMultiplayer() {
  gameMode.value = 'multiplayer'
}

// 加入房间后
function handleJoinRoom() {
  inMultiplayerRoom.value = true
}

// 离开房间
function handleLeaveRoom() {
  inMultiplayerRoom.value = false
}

// 处理落子
async function handlePlacePiece(row: number, col: number) {
  if (!game.placePiece(row, col)) return

  sound.playPlace()

  // 检查游戏是否结束
  if (game.gameStatus.value !== 'playing') {
    handleGameEnd()
    return
  }

  // AI模式下让AI落子
  if (game.mode.value === 'ai' && game.currentPlayer.value === 2) {
    game._isAiThinking.value = true

    // 使用setTimeout让UI有时间更新
    await new Promise(resolve => setTimeout(resolve, 300))

    const aiMove = ai.getAIMove(game._board.value, game.difficulty.value, 2)
    game._isAiThinking.value = false

    if (game.placePiece(aiMove.row, aiMove.col)) {
      sound.playPlace()

      if (game.gameStatus.value !== 'playing') {
        handleGameEnd()
      }
    }
  }
}

// 处理游戏结束
function handleGameEnd() {
  const result = game.getGameResult()
  if (result === 'black') {
    sound.playWin()
  } else if (result === 'white') {
    if (game.mode.value === 'ai') {
      sound.playLose()
    } else {
      sound.playWin()
    }
  }

  // 保存游戏记录
  history.addHistory(
    game.mode.value,
    game.mode.value === 'ai' ? game.difficulty.value : undefined,
    result || 'draw',
    game.getGameRecord() as Position[]
  )
}

// 悔棋
function handleUndo() {
  if (game.undo()) {
    sound.playUndo()
  }
}

// 重新开始
function handleRestart() {
  sound.playClick()
  game.resetGame()
}

// 切换难度
function handleChangeDifficulty(difficulty: Difficulty) {
  sound.playClick()
  game.resetGame(game.mode.value as GameMode, difficulty)
}

// 切换静音
function handleToggleMute() {
  sound.toggleMute()
}

// 打开历史记录
function handleOpenHistory() {
  sound.playClick()
  showHistory.value = true
}

// 关闭历史记录
function handleCloseHistory() {
  showHistory.value = false
}

// 复盘
function handleReplay(h: GameHistoryType) {
  replayHistory.value = h
  showHistory.value = false
  showReplay.value = true
}

// 关闭复盘
function handleCloseReplay() {
  showReplay.value = false
  replayHistory.value = null
}

// 删除记录
function handleDeleteHistory(id: string) {
  history.deleteHistory(id)
}

// 清空记录
function handleClearHistories() {
  history.clearHistories()
}

// 返回菜单
function handleBackToMenu() {
  sound.playClick()
  gameMode.value = null
  inMultiplayerRoom.value = false
  game.resetGame()
}

// 背景动画 - 赛博朋克风格
const bgStyle = computed(() => ({
  background: `
    radial-gradient(ellipse at 20% 80%, rgba(255, 0, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 20%, rgba(0, 255, 255, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 50% 50%, rgba(136, 0, 255, 0.1) 0%, transparent 70%),
    linear-gradient(180deg, #0a0a0f 0%, #12121a 50%, #080810 100%)
  `,
  backgroundSize: '100% 100%, 100% 100%, 100% 100%, 100% 100%'
}))

// 粒子系统
const particles = ref<Array<{ id: number; x: number; y: number; delay: number; duration: number; size: number; color: string }>>([])

onMounted(() => {
  // 生成粒子
  for (let i = 0; i < 30; i++) {
    particles.value.push({
      id: i,
      x: Math.random() * 100,
      y: Math.random() * 100,
      delay: Math.random() * 10,
      duration: 15 + Math.random() * 20,
      size: 1 + Math.random() * 3,
      color: Math.random() > 0.5 ? 'var(--neon-cyan)' : 'var(--neon-magenta)'
    })
  }
})

// 连接线
const connections = ref<Array<{ id: number; x1: number; y1: number; x2: number; y2: number }>>([])

onMounted(() => {
  // 生成网格连接线
  for (let i = 0; i < 10; i++) {
    connections.value.push({
      id: i,
      x1: Math.random() * 100,
      y1: Math.random() * 100,
      x2: Math.random() * 100,
      y2: Math.random() * 100
    })
  }
})
</script>

<template>
  <div class="app" :style="bgStyle">
    <!-- 赛博朋克网格背景 -->
    <div class="cyber-grid"></div>

    <!-- 粒子系统 -->
    <div class="particles-container">
      <div
        v-for="particle in particles"
        :key="particle.id"
        class="cyber-particle"
        :style="{
          left: `${particle.x}%`,
          top: `${particle.y}%`,
          width: `${particle.size}px`,
          height: `${particle.size}px`,
          background: particle.color,
          animationDelay: `${particle.delay}s`,
          animationDuration: `${particle.duration}s`
        }"
      ></div>
    </div>

    <!-- 连接线 -->
    <svg class="connections-svg">
      <line
        v-for="conn in connections"
        :key="conn.id"
        :x1="`${conn.x1}%`"
        :y1="`${conn.y1}%`"
        :x2="`${conn.x2}%`"
        :y2="`${conn.y2}%`"
        stroke="rgba(0, 255, 255, 0.1)"
        stroke-width="1"
      />
    </svg>

    <!-- 扫描线效果 -->
    <div class="scan-line"></div>

    <!-- 主菜单 -->
    <div v-if="!gameMode" class="menu-container">
      <div class="menu-panel glass-panel">
        <!-- 顶部装饰线 -->
        <div class="panel-header-line"></div>

        <!-- 标题区域 -->
        <div class="title-area">
          <div class="title-decoration left">
            <span class="deco-line"></span>
            <span class="deco-diamond"></span>
          </div>
          <h1 class="game-title">
            <span class="title-char" data-text="五">五</span>
            <span class="title-char" data-text="子">子</span>
            <span class="title-char" data-text="棋">棋</span>
          </h1>
          <div class="title-decoration right">
            <span class="deco-diamond"></span>
            <span class="deco-line"></span>
          </div>
        </div>

        <div class="game-subtitle">
          <span class="subtitle-bracket">[</span>
          <span class="subtitle-text">GOMOKU</span>
          <span class="subtitle-version">v2.0</span>
          <span class="subtitle-bracket">]</span>
        </div>

        <!-- 数据流动画 -->
        <div class="data-flow-line"></div>

        <!-- 菜单区域 -->
        <div class="menu-sections">
          <!-- 人机对战 -->
          <div class="menu-section">
            <div class="section-header">
              <span class="section-icon">&lt;/&gt;</span>
              <h3 class="section-title">AI BATTLE</h3>
              <span class="section-line"></span>
            </div>
            <div class="difficulty-buttons">
              <button
                class="glass-btn difficulty-btn"
                data-difficulty="easy"
                @click="startGame('ai', 'easy')"
              >
                <span class="btn-label">EASY</span>
                <span class="btn-sub">简单</span>
              </button>
              <button
                class="glass-btn glass-btn-primary difficulty-btn"
                data-difficulty="medium"
                @click="startGame('ai', 'medium')"
              >
                <span class="btn-label">MEDIUM</span>
                <span class="btn-sub">中等</span>
              </button>
              <button
                class="glass-btn difficulty-btn hard"
                data-difficulty="hard"
                @click="startGame('ai', 'hard')"
              >
                <span class="btn-label">HARD</span>
                <span class="btn-sub">困难</span>
              </button>
            </div>
          </div>

          <!-- 本地对战 -->
          <div class="menu-section">
            <div class="section-header">
              <span class="section-icon">++</span>
              <h3 class="section-title">LOCAL PVP</h3>
              <span class="section-line"></span>
            </div>
            <button class="glass-btn pvp-btn" @click="startGame('pvp')">
              <span class="btn-icon">⚡</span>
              <span class="btn-text">双人对战</span>
              <span class="btn-arrow">→</span>
            </button>
          </div>

          <!-- 在线对战 -->
          <div class="menu-section">
            <div class="section-header">
              <span class="section-icon">◉</span>
              <h3 class="section-title">ONLINE</h3>
              <span class="section-line"></span>
            </div>
            <button class="glass-btn glass-btn-primary multiplayer-btn" @click="startMultiplayer">
              <span class="btn-icon">🌐</span>
              <span class="btn-text">多人在线</span>
              <span class="btn-arrow">→</span>
            </button>
            <div class="ws-status" :class="{ connected: isConnected }">
              <span class="status-dot"></span>
              <span class="status-text">{{ isConnected ? 'ONLINE' : 'OFFLINE' }}</span>
            </div>
          </div>

          <!-- 其他选项 -->
          <div class="menu-section menu-actions">
            <div class="section-header">
              <span class="section-icon">※</span>
              <h3 class="section-title">OPTIONS</h3>
              <span class="section-line"></span>
            </div>
            <div class="action-buttons">
              <button class="glass-btn action-btn" @click="showHistory = true">
                <span class="btn-icon">📜</span>
                <span class="btn-text">历史记录</span>
              </button>
              <button class="glass-btn action-btn" @click="showLeaderboard = true">
                <span class="btn-icon">🏆</span>
                <span class="btn-text">排行榜</span>
              </button>
            </div>
          </div>
        </div>

        <!-- 底部音效控制 -->
        <div class="menu-footer">
          <div class="footer-line"></div>
          <button class="sound-toggle" @click="sound.toggleMute()">
            <span class="sound-icon">{{ sound.isMuted.value ? '🔇' : '🔊' }}</span>
            <span class="sound-label">SOUND</span>
          </button>
          <div class="footer-info">
            <span class="info-text">超级五子棋</span>
            <span class="info-divider">|</span>
            <span class="info-text">2099</span>
          </div>
        </div>

        <!-- 角落装饰 -->
        <div class="corner-deco tl"></div>
        <div class="corner-deco tr"></div>
        <div class="corner-deco bl"></div>
        <div class="corner-deco br"></div>
      </div>
    </div>

    <!-- 多人模式大厅 -->
    <div v-else-if="gameMode === 'multiplayer' && !inMultiplayerRoom" class="lobby-container">
      <button class="back-btn glass-btn" @click="handleBackToMenu">
        <span class="btn-arrow">←</span>
        <span>返回菜单</span>
      </button>
      <Lobby
        @join-room="handleJoinRoom"
        @spectate="handleJoinRoom"
      />
    </div>

    <!-- 多人游戏房间 -->
    <div v-else-if="gameMode === 'multiplayer' && inMultiplayerRoom" class="room-container">
      <GameRoom @leave="handleLeaveRoom" />
    </div>

    <!-- 单机游戏界面 -->
    <div v-else-if="gameMode && gameMode !== 'multiplayer'" class="game-container">
      <div class="game-layout">
        <!-- 游戏信息栏 -->
        <div class="game-info-panel">
          <!-- 当前玩家指示 -->
          <div class="turn-indicator glass-panel">
            <div class="turn-label">CURRENT TURN</div>
            <div class="turn-display">
              <div class="current-piece" :class="{
                'piece-black': game.currentPlayer.value === 1,
                'piece-white': game.currentPlayer.value === 2
              }"></div>
              <span class="turn-text">
                {{ game.currentPlayer.value === 1 ? 'BLACK' : 'WHITE' }}
              </span>
            </div>
            <span v-if="game.mode.value === 'ai'" class="mode-hint">
              {{ game.currentPlayer.value === 1 ? '(YOU)' : '(AI)' }}
            </span>
          </div>

          <!-- 游戏状态 -->
          <div class="game-status-panel glass-panel-dark">
            <div class="status-item">
              <span class="status-label">MODE</span>
              <span class="status-value">{{ game.mode.value === 'ai' ? 'AI BATTLE' : 'LOCAL PVP' }}</span>
            </div>
            <div class="status-item" v-if="game.mode.value === 'ai'">
              <span class="status-label">DIFFICULTY</span>
              <span class="status-value">{{ game.difficulty.value.toUpperCase() }}</span>
            </div>
          </div>
        </div>

        <!-- 棋盘 -->
        <GameBoard
          :board="game.board.value"
          :boardSize="game.BOARD_SIZE"
          :gameStatus="game.gameStatus.value"
          :winningLine="game.winningLine.value"
          :isPlayerTurn="game.isPlayerTurn.value"
          :isAiThinking="game.isAiThinking.value"
          :lastMove="lastMove"
          @placePiece="handlePlacePiece"
        />

        <!-- 控制面板 -->
        <GameControls
          :gameStatus="game.gameStatus.value"
          :mode="game.mode.value"
          :difficulty="game.difficulty.value"
          :isMuted="sound.isMuted.value"
          :canUndo="canUndo"
          @undo="handleUndo"
          @restart="handleRestart"
          @toggleMute="handleToggleMute"
          @changeDifficulty="handleChangeDifficulty"
          @openHistory="handleOpenHistory"
          @backToMenu="handleBackToMenu"
        />
      </div>
    </div>

    <!-- 历史记录弹窗 -->
    <GameHistory
      v-if="showHistory"
      :histories="history.gameHistories.value"
      @close="handleCloseHistory"
      @replay="handleReplay"
      @delete="handleDeleteHistory"
      @clearAll="handleClearHistories"
    />

    <!-- 复盘弹窗 -->
    <ReplayModal
      v-if="showReplay && replayHistory"
      :history="replayHistory"
      @close="handleCloseReplay"
    />

    <!-- 排行榜弹窗 -->
    <Leaderboard
      v-if="showLeaderboard"
      @close="showLeaderboard = false"
    />
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow-x: hidden;
  overflow-y: auto;
  font-family: var(--font-body);
}

/* 赛博朋克网格背景 */
.cyber-grid {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    linear-gradient(rgba(0, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  pointer-events: none;
}

/* 粒子容器 */
.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
}

.cyber-particle {
  position: absolute;
  border-radius: 50%;
  animation: particleFloat linear infinite;
  box-shadow: 0 0 10px currentColor;
}

@keyframes particleFloat {
  0% {
    transform: translateY(100vh) scale(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
    transform: translateY(90vh) scale(1);
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-10vh) scale(0);
    opacity: 0;
  }
}

/* 连接线SVG */
.connections-svg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  opacity: 0.5;
}

/* 扫描线效果 */
.scan-line {
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
    transparent
  );
  animation: scanMove 4s linear infinite;
  opacity: 0.6;
  pointer-events: none;
}

@keyframes scanMove {
  0% {
    top: -5px;
  }
  100% {
    top: 100%;
  }
}

/* ============================================
   主菜单样式
   ============================================ */

.menu-container {
  animation: fadeIn 0.6s ease;
  perspective: 1000px;
}

.menu-panel {
  padding: 40px 50px;
  text-align: center;
  min-width: 420px;
  position: relative;
  overflow: hidden;
}

/* 面板顶部装饰线 */
.panel-header-line {
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
  animation: headerGlow 3s ease infinite;
}

@keyframes headerGlow {
  0%, 100% {
    opacity: 0.7;
    box-shadow: 0 0 10px var(--neon-cyan);
  }
  50% {
    opacity: 1;
    box-shadow: 0 0 20px var(--neon-magenta);
  }
}

/* 标题区域 */
.title-area {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 8px;
}

.title-decoration {
  display: flex;
  align-items: center;
  gap: 8px;
}

.deco-line {
  width: 40px;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--neon-cyan));
}

.title-decoration.right .deco-line {
  background: linear-gradient(90deg, var(--neon-magenta), transparent);
}

.deco-diamond {
  width: 8px;
  height: 8px;
  background: var(--neon-cyan);
  transform: rotate(45deg);
  box-shadow: 0 0 10px var(--neon-cyan);
}

.title-decoration.right .deco-diamond {
  background: var(--neon-magenta);
  box-shadow: 0 0 10px var(--neon-magenta);
}

.game-title {
  font-family: var(--font-display);
  font-size: 56px;
  font-weight: 900;
  margin: 0;
  display: flex;
  gap: 4px;
}

.title-char {
  display: inline-block;
  color: #fff;
  text-shadow:
    0 0 10px var(--neon-cyan),
    0 0 20px var(--neon-cyan),
    0 0 40px var(--neon-cyan);
  animation: titleGlitch 4s ease infinite;
  position: relative;
}

.title-char::before,
.title-char::after {
  content: attr(data-text);
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
}

.title-char::before {
  color: var(--neon-cyan);
  animation: glitchBefore 4s ease infinite;
}

.title-char::after {
  color: var(--neon-magenta);
  animation: glitchAfter 4s ease infinite;
}

@keyframes titleGlitch {
  0%, 90%, 100% {
    transform: translate(0);
    text-shadow:
      0 0 10px var(--neon-cyan),
      0 0 20px var(--neon-cyan);
  }
  92% {
    transform: translate(-2px, 1px);
    text-shadow:
      -2px 0 var(--neon-cyan),
      2px 0 var(--neon-magenta);
  }
  94% {
    transform: translate(2px, -1px);
    text-shadow:
      2px 0 var(--neon-cyan),
      -2px 0 var(--neon-magenta);
  }
  96% {
    transform: translate(-1px, 2px);
  }
  98% {
    transform: translate(1px, -2px);
  }
}

@keyframes glitchBefore {
  0%, 89%, 100% { opacity: 0; }
  90% { opacity: 0.8; transform: translate(-2px, 0); }
  91% { opacity: 0; }
}

@keyframes glitchAfter {
  0%, 91%, 100% { opacity: 0; }
  92% { opacity: 0.8; transform: translate(2px, 0); }
  93% { opacity: 0; }
}

.title-char:nth-child(2) {
  animation-delay: 0.1s;
}

.title-char:nth-child(3) {
  animation-delay: 0.2s;
}

/* 副标题 */
.game-subtitle {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 24px;
  font-family: var(--font-mono);
}

.subtitle-bracket {
  color: var(--neon-cyan);
  font-size: 14px;
  opacity: 0.7;
}

.subtitle-text {
  color: var(--neon-cyan);
  font-size: 16px;
  letter-spacing: 8px;
  text-shadow: 0 0 10px var(--neon-cyan);
}

.subtitle-version {
  color: var(--neon-magenta);
  font-size: 12px;
  padding: 2px 6px;
  border: 1px solid var(--neon-magenta);
  border-radius: 2px;
}

/* 数据流动画 */
.data-flow-line {
  height: 2px;
  margin-bottom: 24px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    var(--neon-cyan) 20%,
    var(--neon-magenta) 50%,
    var(--neon-cyan) 80%,
    transparent 100%
  );
  background-size: 200% 100%;
  animation: dataFlowLine 2s linear infinite;
}

@keyframes dataFlowLine {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 菜单区域 */
.menu-sections {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.menu-section {
  text-align: left;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.section-icon {
  font-family: var(--font-mono);
  color: var(--neon-magenta);
  font-size: 12px;
  opacity: 0.8;
}

.section-title {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 600;
  color: var(--neon-cyan);
  letter-spacing: 3px;
  margin: 0;
  text-shadow: 0 0 5px var(--neon-cyan);
}

.section-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, var(--neon-cyan), transparent);
  opacity: 0.3;
}

/* 难度按钮 */
.difficulty-buttons {
  display: flex;
  gap: 12px;
}

.difficulty-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 14px 12px;
}

.difficulty-btn .btn-label {
  font-size: 14px;
  font-weight: 700;
}

.difficulty-btn .btn-sub {
  font-size: 11px;
  opacity: 0.6;
  letter-spacing: 1px;
}

.difficulty-btn.hard {
  border-color: rgba(255, 0, 80, 0.5);
  color: var(--neon-pink);
}

.difficulty-btn.hard:hover {
  background: rgba(255, 0, 80, 0.15);
  border-color: var(--neon-pink);
  box-shadow: 0 0 15px rgba(255, 0, 80, 0.4);
}

/* PVP和多人按钮 */
.pvp-btn,
.multiplayer-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px 24px;
}

.btn-icon {
  font-size: 18px;
}

.btn-text {
  font-size: 14px;
  font-weight: 600;
  flex: 1;
}

.btn-arrow {
  font-size: 16px;
  opacity: 0.6;
  transition: transform 0.3s ease;
}

.pvp-btn:hover .btn-arrow,
.multiplayer-btn:hover .btn-arrow {
  transform: translateX(5px);
  opacity: 1;
}

/* WebSocket状态 */
.ws-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 10px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 2px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--neon-pink);
  box-shadow: 0 0 10px var(--neon-pink);
  animation: statusPulse 2s ease infinite;
}

.ws-status.connected .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 10px var(--neon-green);
}

@keyframes statusPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--neon-pink);
  letter-spacing: 2px;
}

.ws-status.connected .status-text {
  color: var(--neon-green);
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 12px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

/* 底部区域 */
.menu-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 24px;
  padding-top: 20px;
  position: relative;
}

.footer-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 255, 255, 0.3),
    transparent
  );
}

.sound-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  background: transparent;
  border: 1px solid rgba(0, 255, 255, 0.3);
  color: var(--neon-cyan);
  padding: 8px 16px;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 1px;
  transition: all 0.3s ease;
}

.sound-toggle:hover {
  background: rgba(0, 255, 255, 0.1);
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.3);
}

.footer-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.info-divider {
  opacity: 0.3;
}

/* 角落装饰 */
.corner-deco {
  position: absolute;
  width: 30px;
  height: 30px;
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
  height: 2px;
}

.corner-deco::after {
  width: 2px;
  height: 100%;
}

.corner-deco.tl {
  top: 10px;
  left: 10px;
}

.corner-deco.tl::before { top: 0; left: 0; }
.corner-deco.tl::after { top: 0; left: 0; }

.corner-deco.tr {
  top: 10px;
  right: 10px;
}

.corner-deco.tr::before { top: 0; right: 0; }
.corner-deco.tr::after { top: 0; right: 0; }

.corner-deco.bl {
  bottom: 10px;
  left: 10px;
}

.corner-deco.bl::before { bottom: 0; left: 0; }
.corner-deco.bl::after { bottom: 0; left: 0; }

.corner-deco.br {
  bottom: 10px;
  right: 10px;
}

.corner-deco.br::before { bottom: 0; right: 0; }
.corner-deco.br::after { bottom: 0; right: 0; }

/* ============================================
   大厅容器
   ============================================ */

.lobby-container {
  width: 100%;
  max-width: 1200px;
  padding: 20px;
  animation: fadeIn 0.5s ease;
  min-height: auto;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .lobby-container {
    padding: 12px;
    padding-bottom: 40px;
  }
}

.back-btn {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.back-btn .btn-arrow {
  transition: transform 0.3s ease;
}

.back-btn:hover .btn-arrow {
  transform: translateX(-5px);
}

/* ============================================
   房间容器
   ============================================ */

.room-container {
  width: 100%;
  max-width: 1400px;
  padding: 20px;
  animation: fadeIn 0.5s ease;
  min-height: auto;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .room-container {
    padding: 12px;
    padding-bottom: 40px;
  }
}

/* ============================================
   游戏界面
   ============================================ */

.game-container {
  animation: fadeIn 0.5s ease;
}

.game-layout {
  display: flex;
  align-items: flex-start;
  gap: 32px;
  padding: 20px;
}

/* 游戏信息面板 */
.game-info-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 回合指示 */
.turn-indicator {
  padding: 16px 24px;
  min-width: 180px;
}

.turn-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--neon-cyan);
  letter-spacing: 2px;
  margin-bottom: 8px;
  opacity: 0.8;
}

.turn-display {
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-piece {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  transition: all 0.3s;
  position: relative;
}

.current-piece::after {
  content: '';
  position: absolute;
  top: -3px;
  left: -3px;
  right: -3px;
  bottom: -3px;
  border-radius: 50%;
  border: 2px solid var(--neon-cyan);
  animation: turnPulse 1.5s ease infinite;
}

@keyframes turnPulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.1);
  }
}

.current-piece.piece-black {
  background: radial-gradient(circle at 30% 30%, #444, #000);
  box-shadow:
    0 0 15px rgba(0, 0, 0, 0.8),
    inset 0 0 10px rgba(0, 255, 255, 0.2);
}

.current-piece.piece-white {
  background: radial-gradient(circle at 30% 30%, #fff, #ccc);
  box-shadow:
    0 0 15px rgba(255, 255, 255, 0.5),
    inset 0 0 10px rgba(255, 0, 255, 0.2);
}

.current-piece.piece-white::after {
  border-color: var(--neon-magenta);
}

.turn-text {
  color: #fff;
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
}

.mode-hint {
  color: var(--neon-yellow);
  font-family: var(--font-mono);
  font-size: 12px;
  margin-top: 8px;
  letter-spacing: 1px;
}

/* 游戏状态面板 */
.game-status-panel {
  padding: 12px 16px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}

.status-item + .status-item {
  border-top: 1px solid rgba(0, 255, 255, 0.1);
  padding-top: 8px;
  margin-top: 4px;
}

.status-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.status-value {
  font-family: var(--font-display);
  font-size: 12px;
  color: var(--neon-cyan);
  letter-spacing: 1px;
}

/* ============================================
   响应式设计
   ============================================ */

@media (max-width: 900px) {
  .game-layout {
    flex-direction: column;
    align-items: center;
  }

  .game-info-panel {
    flex-direction: row;
    width: 100%;
    justify-content: center;
  }

  .game-title {
    font-size: 42px;
  }

  .difficulty-buttons {
    flex-direction: column;
  }

  .menu-actions {
    flex-direction: column;
    align-items: center;
  }

  .menu-panel {
    min-width: auto;
    padding: 30px 24px;
  }
}
</style>
