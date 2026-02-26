<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import ChessPiece from './ChessPiece.vue'
import type { Position, GameStatus, PieceType } from '@/types/game'

const props = defineProps<{
  board: readonly (readonly PieceType[])[]
  boardSize: number
  gameStatus: GameStatus
  winningLine: readonly Position[]
  isPlayerTurn: boolean
  isAiThinking: boolean
  lastMove: Position | null
}>()

const emit = defineEmits<{
  (e: 'placePiece', row: number, col: number): void
}>()

// 响应式单元格大小
const cellSize = ref(36)
const boardRef = ref<HTMLElement | null>(null)

// 根据屏幕大小计算单元格尺寸
function calculateCellSize() {
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  // 移动端适配
  if (viewportWidth < 768) {
    // 手机端 - 根据屏幕宽度计算
    const padding = 40 // 左右padding
    const availableWidth = viewportWidth - padding
    const calculatedSize = Math.floor(availableWidth / props.boardSize)
    cellSize.value = Math.min(calculatedSize, 28) // 最大28px
  } else if (viewportWidth < 1024) {
    // 平板端
    const maxBoardSize = Math.min(viewportWidth, viewportHeight) * 0.7
    const calculatedSize = Math.floor(maxBoardSize / props.boardSize)
    cellSize.value = Math.min(calculatedSize, 32)
  } else {
    // 桌面端
    cellSize.value = 36
  }
}

onMounted(() => {
  calculateCellSize()
  window.addEventListener('resize', calculateCellSize)
})

onUnmounted(() => {
  window.removeEventListener('resize', calculateCellSize)
})

function getCellStyle(_row: number, _col: number) {
  return {
    width: `${cellSize.value}px`,
    height: `${cellSize.value}px`
  }
}

function isWinningCell(row: number, col: number): boolean {
  return props.winningLine.some(p => p.row === row && p.col === col)
}

function handleCellClick(row: number, col: number) {
  if (props.gameStatus !== 'playing') return
  if (props.board[row][col] !== 0) return
  if (props.isAiThinking) return
  if (!props.isPlayerTurn) return

  emit('placePiece', row, col)
}

function isLastMove(row: number, col: number): boolean {
  return props.lastMove?.row === row && props.lastMove?.col === col
}

// 计算棋盘总尺寸
const boardPixelSize = computed(() => props.boardSize * cellSize.value)
</script>

<template>
  <div class="board-container">
 

    <!-- 棋盘 -->
    <div class="board glass-panel-dark" ref="boardRef">
      <!-- 霓虹边框效果 -->
      <div class="board-border-glow"></div>

      <!-- 棋盘背景 - 赛博网格 -->
      <div class="board-background">
        <!-- 底层电路纹理 -->
        <div class="circuit-texture"></div>

        <!-- 网格线 - 霓虹风格 -->
        <svg
          class="grid-lines"
          :viewBox="`0 0 ${boardPixelSize} ${boardPixelSize}`"
          :width="boardPixelSize"
          :height="boardPixelSize"
        >
          <!-- 渐变定义 -->
          <defs>
            <linearGradient id="lineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" style="stop-color: rgba(0, 255, 255, 0);"/>
              <stop offset="50%" style="stop-color: rgba(0, 255, 255, 0.6);"/>
              <stop offset="100%" style="stop-color: rgba(0, 255, 255, 0);"/>
            </linearGradient>
            <linearGradient id="vLineGradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" style="stop-color: rgba(0, 255, 255, 0);"/>
              <stop offset="50%" style="stop-color: rgba(0, 255, 255, 0.6);"/>
              <stop offset="100%" style="stop-color: rgba(0, 255, 255, 0);"/>
            </linearGradient>
            <filter id="glow">
              <feGaussianBlur stdDeviation="1" result="coloredBlur"/>
              <feMerge>
                <feMergeNode in="coloredBlur"/>
                <feMergeNode in="SourceGraphic"/>
              </feMerge>
            </filter>
          </defs>

          <!-- 横线 -->
          <line
            v-for="i in boardSize"
            :key="`h-${i}`"
            :x1="cellSize / 2"
            :y1="(i - 0.5) * cellSize"
            :x2="(boardSize - 0.5) * cellSize"
            :y2="(i - 0.5) * cellSize"
            stroke="url(#lineGradient)"
            stroke-width="1"
            filter="url(#glow)"
          />
          <!-- 竖线 -->
          <line
            v-for="i in boardSize"
            :key="`v-${i}`"
            :x1="(i - 0.5) * cellSize"
            :y1="cellSize / 2"
            :x2="(i - 0.5) * cellSize"
            :y2="(boardSize - 0.5) * cellSize"
            stroke="url(#vLineGradient)"
            stroke-width="1"
            filter="url(#glow)"
          />

          <!-- 星位点 - 霓虹菱形 -->
          <g v-for="star in [
            { x: 3, y: 3 }, { x: 3, y: 7 }, { x: 3, y: 11 },
            { x: 7, y: 3 }, { x: 7, y: 7 }, { x: 7, y: 11 },
            { x: 11, y: 3 }, { x: 11, y: 7 }, { x: 11, y: 11 }
          ]" :key="`star-${star.x}-${star.y}`">
            <rect
              :x="(star.x + 0.5) * cellSize - 4"
              :y="(star.y + 0.5) * cellSize - 4"
              width="8"
              height="8"
              fill="var(--neon-magenta)"
              transform="rotate(45, (star.x + 0.5) * cellSize, (star.y + 0.5) * cellSize)"
              opacity="0.8"
              filter="url(#glow)"
            />
          </g>
        </svg>
      </div>

      <!-- 棋子层 -->
      <div class="pieces-layer">
        <div
          v-for="(row, rowIndex) in board"
          :key="rowIndex"
          class="board-row"
        >
          <div
            v-for="(cell, colIndex) in row"
            :key="colIndex"
            class="board-cell"
            :style="getCellStyle(rowIndex, colIndex)"
            :class="{
              'clickable': cell === 0 && gameStatus === 'playing' && isPlayerTurn && !isAiThinking,
              'hover-black': cell === 0 && isPlayerTurn
            }"
            @click="handleCellClick(rowIndex, colIndex)"
          >
            <ChessPiece
              :piece="cell"
              :isLastMove="isLastMove(rowIndex, colIndex)"
              :isWinning="isWinningCell(rowIndex, colIndex)"
            />
          </div>
        </div>
      </div>

      <!-- 坐标标注 -->
      <div class="board-coordinates">
        <span v-for="i in boardSize" :key="`top-${i}`" class="coord coord-top">
          {{ String.fromCharCode(64 + i) }}
        </span>
        <span v-for="i in boardSize" :key="`left-${i}`" class="coord coord-left">
          {{ i }}
        </span>
      </div>
    </div>
       <!-- AI思考提示 - 赛博朋克风格 -->
    <div v-if="isAiThinking" class="ai-thinking">
      <div class="thinking-animation">
        <span class="thinking-bracket">[</span>
        <span class="thinking-dots">
          <span class="dot"></span>
          <span class="dot"></span>
          <span class="dot"></span>
        </span>
        <span class="thinking-bracket">]</span>
      </div>
      <span class="thinking-text">AI PROCESSING</span>
      <div class="thinking-progress">
        <div class="progress-bar"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.board-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

/* AI思考提示 - 赛博朋克风格 */
.ai-thinking {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 32px;
  background: rgba(0, 255, 255, 0.05);
  border: 1px solid var(--neon-cyan);
  border-radius: 0;
  clip-path: polygon(
    0 0,
    calc(100% - 10px) 0,
    100% 10px,
    100% 100%,
    10px 100%,
    0 calc(100% - 10px)
  );
  animation: thinkingPulse 1.5s ease infinite;
}

@keyframes thinkingPulse {
  0%, 100% {
    box-shadow:
      0 0 10px rgba(0, 255, 255, 0.3),
      inset 0 0 20px rgba(0, 255, 255, 0.05);
  }
  50% {
    box-shadow:
      0 0 20px rgba(0, 255, 255, 0.5),
      inset 0 0 30px rgba(0, 255, 255, 0.1);
  }
}

.thinking-animation {
  display: flex;
  align-items: center;
  gap: 8px;
}

.thinking-bracket {
  font-family: var(--font-mono);
  font-size: 18px;
  color: var(--neon-cyan);
  text-shadow: 0 0 10px var(--neon-cyan);
}

.thinking-dots {
  display: flex;
  gap: 6px;
}

.dot {
  width: 8px;
  height: 8px;
  background: var(--neon-cyan);
  border-radius: 0;
  animation: dotPulse 1.4s ease-in-out infinite;
  box-shadow: 0 0 10px var(--neon-cyan);
}

.dot:nth-child(2) {
  animation-delay: 0.2s;
}

.dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes dotPulse {
  0%, 80%, 100% {
    transform: scale(0.6);
    opacity: 0.3;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.thinking-text {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 600;
  color: var(--neon-cyan);
  letter-spacing: 3px;
  text-shadow: 0 0 5px var(--neon-cyan);
}

.thinking-progress {
  width: 100%;
  height: 2px;
  background: rgba(0, 255, 255, 0.2);
  overflow: hidden;
}

.progress-bar {
  width: 30%;
  height: 100%;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta));
  animation: progressMove 1.5s ease-in-out infinite;
}

@keyframes progressMove {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(400%);
  }
}

/* 棋盘 */
.board {
  position: relative;
  padding: 20px;
  background: var(--cyber-darker);
  border: 1px solid rgba(0, 255, 255, 0.3);
}

/* 霓虹边框效果 */
.board-border-glow {
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border: 2px solid transparent;
  background: linear-gradient(45deg, var(--neon-cyan), var(--neon-magenta), var(--neon-cyan)) border-box;
  -webkit-mask:
    linear-gradient(#fff 0 0) padding-box,
    linear-gradient(#fff 0 0);
  -webkit-mask-composite: destination-out;
  mask-composite: exclude;
  pointer-events: none;
  animation: borderRotate 4s linear infinite;
  opacity: 0.6;
}

@keyframes borderRotate {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}

.board-background {
  position: relative;
}

/* 电路纹理 */
.circuit-texture {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    linear-gradient(90deg, rgba(0, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(rgba(0, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 10px 10px;
  pointer-events: none;
}

.grid-lines {
  display: block;
}

/* 棋子层 */
.pieces-layer {
  position: absolute;
  top: 20px;
  left: 20px;
  pointer-events: none;
}

.board-row {
  display: flex;
}

.board-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.board-cell.clickable {
  cursor: pointer;
  pointer-events: auto;
}

/* 悬停效果 - 赛博朋克风格 */
.board-cell.clickable:hover::before {
  content: '';
  position: absolute;
  width: 70%;
  height: 70%;
  border: 1px solid var(--neon-cyan);
  background: rgba(0, 255, 255, 0.1);
  box-shadow:
    0 0 10px rgba(0, 255, 255, 0.3),
    inset 0 0 10px rgba(0, 255, 255, 0.1);
  clip-path: polygon(
    0 0,
    calc(100% - 4px) 0,
    100% 4px,
    100% 100%,
    4px 100%,
    0 calc(100% - 4px)
  );
  animation: hoverPulse 1s ease infinite;
  pointer-events: none;
}

@keyframes hoverPulse {
  0%, 100% {
    opacity: 0.6;
    transform: scale(0.95);
  }
  50% {
    opacity: 1;
    transform: scale(1);
  }
}

.board-cell.hover-black.clickable:hover::before {
  border-color: var(--neon-magenta);
  background: rgba(255, 0, 255, 0.1);
  box-shadow:
    0 0 10px rgba(255, 0, 255, 0.3),
    inset 0 0 10px rgba(255, 0, 255, 0.1);
}

/* 坐标标注 */
.board-coordinates {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.coord {
  position: absolute;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--neon-cyan);
  opacity: 0.5;
}

.coord-top {
  top: 4px;
  transform: translateX(-50%);
}

.coord-left {
  left: 4px;
  transform: translateY(-50%);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .board {
    padding: 12px;
  }

  .pieces-layer {
    top: 12px;
    left: 12px;
  }

  .ai-thinking {
    padding: 12px 20px;
  }

  .thinking-text {
    font-size: 10px;
    letter-spacing: 2px;
  }

  .coord {
    font-size: 8px;
  }

  .coord-top {
    top: 2px;
  }

  .coord-left {
    left: 2px;
  }
}

@media (max-width: 480px) {
  .board {
    padding: 8px;
  }

  .pieces-layer {
    top: 8px;
    left: 8px;
  }

  /* 隐藏坐标在很小的屏幕上 */
  .board-coordinates {
    display: none;
  }
}
</style>
