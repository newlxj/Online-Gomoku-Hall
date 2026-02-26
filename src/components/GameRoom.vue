<script setup lang="ts">
import { computed, ref, onUnmounted } from 'vue'
import { useRoom } from '@/composables/useRoom'
import { useSound } from '@/composables/useSound'
import GameBoard from './GameBoard.vue'
import GameTimer from './GameTimer.vue'
import PlayerPanel from './PlayerPanel.vue'
import SpectatorPanel from './SpectatorPanel.vue'
import EmojiPanel from './EmojiPanel.vue'
import type { Board, GameStatus } from '@/types/game'
import type { Player } from '@/types/multiplayer'

const emit = defineEmits<{
  (e: 'leave'): void
}>()

const {
  currentRoom,
  isPlaying,
  myPlayer,
  isMyTurn,
  setReady,
  makeMove,
  sendEmoji,
  leaveRoom,
} = useRoom()

// 音效
const { playPlace, playWin, playLose, isMuted } = useSound()

// 本地棋盘状态
const localBoard = computed<Board>(() => {
  if (!currentRoom.value) {
    return Array(15).fill(null).map(() => Array(15).fill(0)) as Board
  }
  return currentRoom.value.board as Board
})

// 赢家信息
const gameWinner = ref<number | null>(null)
const winningLine = ref<{row: number, col: number}[]>([])

// 游戏结束信息
const gameOverInfo = ref<{
  winnerAlias: string
  loserAlias: string
  scoreChanged: number
  reason: string
} | null>(null)

// 俏皮激励语列表
const FUNNY_PHRASES = [
  '对你的敬仰如同滔滔江水连绵不绝！',
  '棋艺精湛，堪称棋圣转世！',
  '这波操作，我愿称之为绝杀！',
  '高手过招，胜负只在一念之间！',
  '今天你是最靓的仔！',
  '这局下得太帅了，建议反复观看学习！',
  '棋盘上的诸葛亮，稳！',
  '这一手，妙啊！给101分，多1分让你骄傲！',
  '天赋异禀，前途无量！',
  '这走位，这意识，绝了！',
  '今日棋神附体，无人能敌！',
  '运筹帷幄之中，决胜千里之外！',
  '这波操作我看傻了，太强了！',
  '棋艺超群，令在下佩服得五体投地！',
  '大师级操作，服了服了！',
  '这局堪称教科书级别的对局！',
  '你是吃了智慧果吗？太厉害了！',
  '棋逢对手，将遇良才，精彩！',
  '这预判能力，开挂了吧？',
  '神之一手，直接封神！',
]

// 获取随机激励语
function getRandomPhrase(): string {
  return FUNNY_PHRASES[Math.floor(Math.random() * FUNNY_PHRASES.length)]
}

const currentPhrase = ref('')

// 处理落子
function handlePlacePiece(row: number, col: number) {
  if (!isMyTurn.value) return
  makeMove({ row, col })
  // 播放落子音效
  playPlace()
}

// 准备
function handleReady() {
  setReady()
}

// 离开房间
function handleLeave() {
  leaveRoom()
  emit('leave')
}

// 再来一局 - 重置准备状态
function handlePlayAgain() {
  // 重置本地游戏结果状态
  gameWinner.value = null
  winningLine.value = []
  gameOverInfo.value = null
  currentPhrase.value = ''
  // 发送准备请求
  setReady()
}

// 发送表情
function handleSendEmoji(emoji: string) {
  sendEmoji(emoji)
}

// 获取玩家信息 - 改进：游戏开始前也能显示所有玩家
const blackPlayer = computed<Player | null>(() => {
  if (!currentRoom.value) return null
  // 游戏开始后按 pieceType 查找
  const playerByPiece = currentRoom.value.players.find(p => p.pieceType === 1)
  if (playerByPiece) return playerByPiece
  // 游戏开始前，第一个玩家作为黑方显示
  if (currentRoom.value.players.length >= 1) {
    return currentRoom.value.players[0]
  }
  return null
})

const whitePlayer = computed<Player | null>(() => {
  if (!currentRoom.value) return null
  // 游戏开始后按 pieceType 查找
  const playerByPiece = currentRoom.value.players.find(p => p.pieceType === 2)
  if (playerByPiece) return playerByPiece
  // 游戏开始前，第二个玩家作为白方显示
  if (currentRoom.value.players.length >= 2) {
    return currentRoom.value.players[1]
  }
  return null
})

const isSpectator = computed(() => {
  return currentRoom.value && !myPlayer.value
})

// 计算最后一步棋的位置 - 用于高亮显示
const lastMovePosition = computed(() => {
  if (!currentRoom.value || currentRoom.value.status !== 'playing') return null
  const history = currentRoom.value.moveHistory
  if (history.length === 0) return null
  const lastMove = history[history.length - 1]
  return lastMove.position
})

// 是否显示玩家面板 - 只要有玩家就显示
const showBlackPlayer = computed(() => {
  return currentRoom.value && currentRoom.value.players.length >= 1
})

const showWhitePlayer = computed(() => {
  return currentRoom.value && currentRoom.value.players.length >= 2
})

// 房间状态文字
const statusText = computed(() => {
  if (!currentRoom.value) return ''
  switch (currentRoom.value.status) {
    case 'waiting': return 'WAITING FOR PLAYERS...'
    case 'ready':
      // 检查双方准备状态
      const allReady = currentRoom.value.players.every(p => p.isReady)
      if (allReady && currentRoom.value.players.length === 2) {
        return 'BOTH READY - STARTING...'
      }
      return 'WAITING FOR READY'
    case 'playing': return isMyTurn.value ? 'YOUR TURN' : 'OPPONENT\'S TURN'
    case 'finished':
      if (gameWinner.value === 1) return 'BLACK WINS!'
      if (gameWinner.value === 2) return 'WHITE WINS!'
      return 'DRAW!'
    default: return ''
  }
})

// 转换 GameStatus 类型
const gameStatusForBoard = computed<GameStatus>(() => {
  if (!currentRoom.value) return 'playing'
  if (currentRoom.value.status === 'playing') return 'playing'
  if (currentRoom.value.status === 'finished') {
    if (gameWinner.value === 1) return 'black-win'
    if (gameWinner.value === 2) return 'white-win'
    return 'draw'
  }
  return 'playing'
})

// 监听 game_over 消息
import { useWebSocket } from '@/composables/useWebSocket'
import { useHistory } from '@/composables/useHistory'
const { subscribe } = useWebSocket()
const { addMultiplayerHistory } = useHistory()

// 表情显示
const receivedEmoji = ref<{ emoji: string; fromAlias: string } | null>(null)
let emojiTimer: ReturnType<typeof setTimeout> | null = null

const unsubscribeGameOver = subscribe<{
  winner: number
  reason: string
  winningLine: {row: number, col: number}[]
  winnerAlias: string
  loserAlias: string
  scoreChanged: number
}>(
  'game_over',
  (payload) => {
    console.log('[GameRoom] Received game_over:', payload)
    gameWinner.value = payload.winner
    winningLine.value = payload.winningLine || []
    gameOverInfo.value = {
      winnerAlias: payload.winnerAlias || '',
      loserAlias: payload.loserAlias || '',
      scoreChanged: payload.scoreChanged || 0,
      reason: payload.reason
    }
    // 设置随机激励语
    currentPhrase.value = getRandomPhrase()
    // 播放胜利/失败音效
    if (myPlayer.value) {
      if (payload.winner === myPlayer.value.pieceType) {
        playWin()
      } else if (payload.winner !== 0) {
        playLose()
      }
    }

    // 保存多人游戏记录到 localStorage
    if (currentRoom.value && blackPlayer.value && whitePlayer.value) {
      const winner: 'black' | 'white' | 'draw' = payload.winner === 1 ? 'black' : payload.winner === 2 ? 'white' : 'draw'
      const moves = currentRoom.value.moveHistory.map(m => m.position)

      addMultiplayerHistory(winner, moves, {
        roomName: currentRoom.value.name,
        blackPlayer: blackPlayer.value.alias,
        whitePlayer: whitePlayer.value.alias,
        winnerAlias: payload.winnerAlias || '',
        loserAlias: payload.loserAlias || '',
        scoreChanged: payload.scoreChanged || 0,
        reason: payload.reason as 'win' | 'disconnect' | 'timeout' | 'leave'
      })
      console.log('[GameRoom] Multiplayer game history saved')
    }
  }
)

// 订阅落子消息 - 对手落子时播放音效
const unsubscribeMove = subscribe<{ roomId: string; position: { row: number; col: number }; playerId: string; pieceType: number }>(
  'move',
  (payload) => {
    console.log('[GameRoom] Received move:', payload)
    // 对手落子时播放音效
    if (payload.playerId !== myPlayer.value?.id) {
      playPlace()
    }
  }
)

const unsubscribeEmoji = subscribe<{ roomId: string; emoji: string; fromPlayer: string; fromAlias: string }>(
  'emoji',
  (payload) => {
    console.log('[GameRoom] Received emoji:', payload)
    receivedEmoji.value = {
      emoji: payload.emoji,
      fromAlias: payload.fromAlias
    }
    // 3秒后清除表情
    if (emojiTimer) {
      clearTimeout(emojiTimer)
    }
    emojiTimer = setTimeout(() => {
      receivedEmoji.value = null
    }, 3000)
  }
)

onUnmounted(() => {
  unsubscribeGameOver()
  unsubscribeMove()
  unsubscribeEmoji()
  if (emojiTimer) {
    clearTimeout(emojiTimer)
  }
})
</script>

<template>
  <div class="game-room" v-if="currentRoom">
    <!-- 表情弹窗 -->
    <Transition name="emoji-popup">
      <div v-if="receivedEmoji" class="emoji-popup">
        <span class="emoji-alias">{{ receivedEmoji.fromAlias }}</span>
        <span class="emoji-content">{{ receivedEmoji.emoji }}</span>
      </div>
    </Transition>

    <!-- 顶部信息栏 -->
    <div class="room-header glass-panel-dark">
      <div class="room-info">
        <div class="room-title">
          <span class="title-icon">◈</span>
          <h3 class="room-name">{{ currentRoom.name }}</h3>
        </div>
        <span class="room-status" :class="currentRoom.status">
          {{ statusText }}
        </span>
      </div>
      <button class="leave-btn glass-btn glass-btn-danger" @click="handleLeave">
        <span class="btn-icon">✕</span>
        <span class="btn-text">离开</span>
      </button>
    </div>

    <!-- 主要游戏区域 -->
    <div class="game-area">
      <!-- 左侧 - 黑方信息 -->
      <div class="player-side left">
        <PlayerPanel
          v-if="showBlackPlayer && blackPlayer"
          :player="blackPlayer"
          :is-current-turn="currentRoom.status === 'playing' && currentRoom.currentTurn === 1"
          :is-me="myPlayer?.id === blackPlayer?.id"
        />
        <div v-else class="player-slot empty">
          <span class="slot-icon">◇</span>
          <span class="slot-text">WAITING FOR BLACK...</span>
        </div>
        <GameTimer
          v-if="isPlaying && blackPlayer"
          :remaining-time="blackPlayer.remainingTime || 0"
          :move-time="blackPlayer.moveTimeLeft || 0"
          :is-active="currentRoom.currentTurn === 1"
        />
      </div>

      <!-- 中间 - 棋盘 -->
      <div class="board-container">
        <GameBoard
          :board="localBoard"
          :board-size="15"
          :game-status="gameStatusForBoard"
          :winning-line="winningLine"
          :is-player-turn="isMyTurn"
          :is-ai-thinking="false"
          :last-move="lastMovePosition"
          @place-piece="handlePlacePiece"
        />

        <!-- 游戏结束提示 -->
        <div v-if="currentRoom.status === 'finished'" class="game-result">
          <!-- 获胜特效 -->
          <div class="victory-effects">
            <div class="cyber-particle" v-for="i in 12" :key="i" :style="{ '--delay': i * 0.1 + 's', '--x': Math.random() * 100 + '%' }"></div>
          </div>

          <div class="result-header">
            <span class="result-label">GAME OVER</span>
          </div>

          <div class="result-winner">
            <span class="winner-icon">{{ gameWinner === 1 ? '◆' : '◇' }}</span>
            <span class="winner-name">{{ gameOverInfo?.winnerAlias || 'UNKNOWN' }}</span>
            <span class="winner-badge">VICTORY</span>
          </div>

          <div class="result-score" v-if="gameOverInfo && gameOverInfo.scoreChanged > 0">
            <span class="score-icon">★</span>
            <span class="score-value">+{{ gameOverInfo.scoreChanged }}</span>
            <span class="score-label">PTS</span>
          </div>

          <div class="result-phrase" v-if="currentPhrase">
            "{{ currentPhrase }}"
          </div>

          <div class="result-reason" v-if="gameOverInfo?.reason && gameOverInfo.reason !== 'win'">
            <template v-if="gameOverInfo.reason === 'disconnect'">
              <span v-if="myPlayer && gameWinner === myPlayer.pieceType">OPPONENT DISCONNECTED</span>
              <span v-else>YOU DISCONNECTED</span>
            </template>
            <template v-else-if="gameOverInfo.reason === 'timeout'">
              <span v-if="myPlayer && gameWinner === myPlayer.pieceType">OPPONENT TIMEOUT</span>
              <span v-else>TIMEOUT</span>
            </template>
            <template v-else-if="gameOverInfo.reason === 'leave'">
              <span v-if="myPlayer && gameWinner === myPlayer.pieceType">OPPONENT SURRENDERED</span>
              <span v-else>YOU LEFT</span>
            </template>
          </div>

          <div class="result-actions">
            <button class="glass-btn glass-btn-primary play-again-btn" @click="handlePlayAgain">
              <span class="btn-icon">⟳</span>
              <span>再来一局</span>
            </button>
            <button class="glass-btn leave-room-btn" @click="handleLeave">
              <span class="btn-icon">←</span>
              <span>离开房间</span>
            </button>
          </div>
        </div>

        <!-- 准备按钮 - 居中显示在棋盘上方 -->
        <Transition name="ready-overlay">
          <div v-if="currentRoom.status === 'waiting' || currentRoom.status === 'ready'" class="ready-overlay">
            <div class="ready-overlay-bg"></div>
            <div class="ready-content">
              <button
                v-if="myPlayer && !myPlayer.isReady"
                class="ready-btn-prominent"
                @click="handleReady"
              >
                <span class="ready-icon">▶</span>
                <span class="ready-main-text">准备</span>
                <span class="ready-sub-text">CLICK TO READY</span>
                <div class="ready-glow"></div>
              </button>
              <div v-else-if="myPlayer?.isReady" class="ready-waiting-prominent">
                <span class="waiting-icon-pulse">◈</span>
                <span class="waiting-main-text">已准备</span>
                <span class="waiting-sub-text">WAITING FOR OPPONENT...</span>
              </div>
            </div>
          </div>
        </Transition>

        <!-- 观战提示 -->
        <div v-if="isSpectator" class="spectator-notice">
          <span class="notice-icon">👁</span>
          <span class="notice-text">SPECTATOR MODE</span>
        </div>

        <!-- 表情面板 -->
        <EmojiPanel
          v-if="isPlaying && myPlayer"
          @send="handleSendEmoji"
        />
      </div>

      <!-- 右侧 - 白方信息 -->
      <div class="player-side right">
        <PlayerPanel
          v-if="showWhitePlayer && whitePlayer"
          :player="whitePlayer"
          :is-current-turn="currentRoom.status === 'playing' && currentRoom.currentTurn === 2"
          :is-me="myPlayer?.id === whitePlayer?.id"
        />
        <div v-else class="player-slot empty">
          <span class="slot-icon">◇</span>
          <span class="slot-text">WAITING FOR WHITE...</span>
        </div>
        <GameTimer
          v-if="isPlaying && whitePlayer"
          :remaining-time="whitePlayer.remainingTime || 0"
          :move-time="whitePlayer.moveTimeLeft || 0"
          :is-active="currentRoom.currentTurn === 2"
        />
      </div>
    </div>

    <!-- 底部 - 观战者 -->
    <SpectatorPanel
      v-if="currentRoom.spectators && currentRoom.spectators.length > 0"
      :spectators="currentRoom.spectators"
    />
  </div>
</template>

<style scoped>
.game-room {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  gap: 16px;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  position: relative;
}

.room-header::before {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta), var(--neon-cyan));
  opacity: 0.5;
}

.room-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.room-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  color: var(--neon-cyan);
  text-shadow: 0 0 10px var(--neon-cyan);
}

.room-name {
  color: #fff;
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  letter-spacing: 1px;
}

.room-status {
  font-family: var(--font-display);
  font-size: 11px;
  padding: 6px 14px;
  font-weight: 600;
  letter-spacing: 2px;
}

.room-status.waiting {
  background: rgba(0, 255, 102, 0.15);
  color: var(--neon-green);
  border: 1px solid rgba(0, 255, 102, 0.3);
}

.room-status.ready {
  background: rgba(255, 204, 0, 0.15);
  color: var(--neon-yellow);
  border: 1px solid rgba(255, 204, 0, 0.3);
  animation: readyPulse 1.5s ease infinite;
}

@keyframes readyPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.room-status.playing {
  background: rgba(0, 136, 255, 0.15);
  color: var(--neon-blue);
  border: 1px solid rgba(0, 136, 255, 0.3);
}

.room-status.finished {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.leave-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
}

.btn-icon {
  font-size: 14px;
}

.btn-text {
  font-size: 12px;
}

.game-area {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 24px;
  flex: 1;
}

.player-side {
  width: 200px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.player-side.left {
  align-items: flex-end;
}

.player-side.right {
  align-items: flex-start;
}

.board-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

/* 准备按钮覆盖层 */
.ready-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  pointer-events: none;
}

.ready-overlay-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at center, transparent 30%, rgba(0, 0, 0, 0.7) 100%);
  pointer-events: none;
}

.ready-content {
  position: relative;
  z-index: 1;
  pointer-events: auto;
}

.ready-btn-prominent {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 80px;
  background: linear-gradient(135deg, rgba(0, 255, 102, 0.2), rgba(0, 136, 255, 0.2));
  border: 2px solid var(--neon-cyan);
  border-radius: 16px;
  cursor: pointer;
  overflow: hidden;
  animation: readyBreathing 2s ease-in-out infinite, readyBorderGlow 2s ease-in-out infinite;
  box-shadow:
    0 0 30px rgba(0, 255, 255, 0.3),
    0 0 60px rgba(0, 255, 102, 0.2),
    inset 0 0 30px rgba(0, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.ready-btn-prominent:hover {
  transform: scale(1.05);
  box-shadow:
    0 0 40px rgba(0, 255, 255, 0.5),
    0 0 80px rgba(0, 255, 102, 0.3),
    inset 0 0 40px rgba(0, 255, 255, 0.2);
}

@keyframes readyBreathing {
  0%, 100% {
    transform: scale(1);
    box-shadow:
      0 0 30px rgba(0, 255, 255, 0.3),
      0 0 60px rgba(0, 255, 102, 0.2),
      inset 0 0 30px rgba(0, 255, 255, 0.1);
  }
  50% {
    transform: scale(1.03);
    box-shadow:
      0 0 50px rgba(0, 255, 255, 0.5),
      0 0 100px rgba(0, 255, 102, 0.4),
      inset 0 0 50px rgba(0, 255, 255, 0.2);
  }
}

@keyframes readyBorderGlow {
  0%, 100% {
    border-color: var(--neon-cyan);
  }
  50% {
    border-color: var(--neon-green);
  }
}

.ready-icon {
  font-size: 32px;
  color: var(--neon-cyan);
  animation: iconPulse 1s ease-in-out infinite;
}

@keyframes iconPulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

.ready-main-text {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  text-shadow: 0 0 20px rgba(255, 255, 255, 0.5);
  letter-spacing: 4px;
}

.ready-sub-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  letter-spacing: 3px;
}

.ready-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200%;
  height: 200%;
  transform: translate(-50%, -50%);
  background: radial-gradient(circle, rgba(0, 255, 255, 0.1) 0%, transparent 50%);
  animation: glowRotate 4s linear infinite;
  pointer-events: none;
}

@keyframes glowRotate {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
}

.ready-waiting-prominent {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px 60px;
  background: rgba(0, 255, 102, 0.15);
  border: 2px solid var(--neon-green);
  border-radius: 16px;
  box-shadow:
    0 0 30px rgba(0, 255, 102, 0.3),
    inset 0 0 30px rgba(0, 255, 102, 0.1);
  animation: waitingBreathing 2s ease-in-out infinite;
}

@keyframes waitingBreathing {
  0%, 100% {
    box-shadow:
      0 0 30px rgba(0, 255, 102, 0.3),
      inset 0 0 30px rgba(0, 255, 102, 0.1);
  }
  50% {
    box-shadow:
      0 0 50px rgba(0, 255, 102, 0.5),
      inset 0 0 50px rgba(0, 255, 102, 0.2);
  }
}

.waiting-icon-pulse {
  font-size: 32px;
  color: var(--neon-green);
  animation: waitingIconPulse 1s ease-in-out infinite;
}

@keyframes waitingIconPulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.2);
  }
}

.waiting-main-text {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 600;
  color: var(--neon-green);
  letter-spacing: 4px;
}

.waiting-sub-text {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 2px;
  animation: textBlink 1.5s ease-in-out infinite;
}

@keyframes textBlink {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

/* 过渡动画 */
.ready-overlay-enter-active {
  animation: readyOverlayIn 0.5s ease;
}

.ready-overlay-leave-active {
  animation: readyOverlayOut 0.3s ease;
}

@keyframes readyOverlayIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes readyOverlayOut {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}

.spectator-notice {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: rgba(0, 136, 255, 0.1);
  border: 1px solid rgba(0, 136, 255, 0.3);
}

.notice-icon {
  font-size: 16px;
}

.notice-text {
  font-family: var(--font-display);
  font-size: 12px;
  color: var(--neon-blue);
  letter-spacing: 2px;
}

.player-slot {
  padding: 20px;
  min-width: 160px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.player-slot.empty {
  background: rgba(0, 0, 0, 0.3);
  border: 1px dashed rgba(0, 255, 255, 0.3);
}

.slot-icon {
  font-size: 24px;
  color: rgba(0, 255, 255, 0.3);
  animation: float 3s ease infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.slot-text {
  font-family: var(--font-mono);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

/* 游戏结果弹窗 */
.game-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 40px 32px;
  background: rgba(10, 10, 20, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--neon-cyan);
  box-shadow:
    0 0 30px rgba(0, 255, 255, 0.3),
    inset 0 0 30px rgba(0, 255, 255, 0.05);
  animation: resultPopIn 0.5s ease;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
  min-width: 340px;
  max-width: 90vw;
  overflow: hidden;
}

/* 赛博朋克粒子效果 */
.victory-effects {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.cyber-particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: var(--neon-cyan);
  left: var(--x);
  animation: particleFall 2s ease-out infinite;
  animation-delay: var(--delay);
  box-shadow: 0 0 10px currentColor;
}

.cyber-particle:nth-child(even) {
  background: var(--neon-magenta);
}

@keyframes particleFall {
  0% {
    top: -10px;
    opacity: 1;
  }
  100% {
    top: 100%;
    opacity: 0;
  }
}

@keyframes resultPopIn {
  0% {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.8);
  }
  50% {
    transform: translate(-50%, -50%) scale(1.02);
  }
  100% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}

.result-header {
  text-align: center;
}

.result-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 3px;
}

.result-winner {
  display: flex;
  align-items: center;
  gap: 12px;
  animation: winnerBounce 0.6s ease;
}

@keyframes winnerBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.winner-icon {
  font-size: 32px;
  color: var(--cyber-gold);
  text-shadow: 0 0 20px var(--cyber-gold);
}

.winner-name {
  color: #fff;
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  text-shadow: 0 0 20px rgba(255, 255, 255, 0.5);
  letter-spacing: 2px;
}

.winner-badge {
  background: linear-gradient(135deg, var(--cyber-gold), var(--neon-orange));
  color: var(--cyber-black);
  padding: 6px 16px;
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 2px;
  animation: badgePulse 1s ease infinite;
}

@keyframes badgePulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.result-score {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: rgba(255, 215, 0, 0.1);
  border: 1px solid rgba(255, 215, 0, 0.3);
  animation: scoreGlow 1.5s ease infinite;
}

@keyframes scoreGlow {
  0%, 100% { box-shadow: 0 0 10px rgba(255, 215, 0, 0.2); }
  50% { box-shadow: 0 0 20px rgba(255, 215, 0, 0.4); }
}

.score-icon {
  font-size: 20px;
  color: var(--cyber-gold);
}

.score-value {
  color: var(--cyber-gold);
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  text-shadow: 0 0 15px rgba(255, 215, 0, 0.5);
}

.score-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.result-phrase {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  font-style: italic;
  text-align: center;
  max-width: 280px;
  line-height: 1.6;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border-left: 2px solid var(--neon-cyan);
  animation: phraseFadeIn 0.8s ease 0.3s both;
}

@keyframes phraseFadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.result-reason {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 1px;
}

.result-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.play-again-btn,
.leave-room-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 24px;
  font-size: 13px;
}

/* 表情弹窗 */
.emoji-popup {
  position: fixed;
  top: 20%;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(10, 10, 20, 0.9);
  padding: 20px 40px;
  border: 1px solid var(--neon-cyan);
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 1000;
  animation: emojiBounce 0.5s ease;
}

.emoji-alias {
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.emoji-content {
  font-size: 56px;
  animation: emojiPop 0.3s ease;
}

@keyframes emojiBounce {
  0% { transform: translateX(-50%) translateY(-20px) scale(0.5); opacity: 0; }
  50% { transform: translateX(-50%) translateY(10px) scale(1.1); }
  100% { transform: translateX(-50%) translateY(0) scale(1); opacity: 1; }
}

@keyframes emojiPop {
  0% { transform: scale(0); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.emoji-popup-enter-active,
.emoji-popup-leave-active {
  transition: all 0.3s ease;
}

.emoji-popup-enter-from,
.emoji-popup-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px) scale(0.5);
}

/* 响应式 */
@media (max-width: 900px) {
  .game-area {
    flex-direction: column;
    align-items: center;
  }

  .player-side {
    width: 100%;
    max-width: 400px;
    flex-direction: row;
    justify-content: center;
    gap: 24px;
  }

  .player-side.left,
  .player-side.right {
    align-items: center;
  }

  .game-result {
    min-width: auto;
    padding: 30px 24px;
  }

  .winner-name {
    font-size: 22px;
  }
}

@media (max-width: 480px) {
  .room-header {
    flex-direction: column;
    gap: 12px;
    padding: 12px 16px;
  }

  .room-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .leave-btn .btn-text {
    display: none;
  }

  .game-result {
    min-width: auto;
    max-width: 90vw;
    left: 5vw;
    right: 5vw;
    transform: translate(0, -50%);
    padding: 24px 16px;
  }

  .result-winner {
    flex-direction: column;
    gap: 8px;
  }

  .winner-badge {
    font-size: 10px;
  }

  .result-actions {
    flex-direction: column;
  }

  .play-again-btn,
  .leave-room-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
