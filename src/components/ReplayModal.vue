<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import ChessPiece from './ChessPiece.vue'
import type { GameHistory, Board, PieceType } from '@/types/game'

const props = defineProps<{
  history: GameHistory
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const BOARD_SIZE = 15
const currentStep = ref(0)
const isPlaying = ref(false)
const playInterval = ref<number | null>(null)

const board = ref<Board>(createEmptyBoard())
const isComplete = computed(() => currentStep.value >= props.history.moves.length)

function createEmptyBoard(): Board {
  return Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0) as PieceType[])
}

// 更新棋盘显示
watch(currentStep, (step) => {
  const newBoard = createEmptyBoard()
  for (let i = 0; i < step && i < props.history.moves.length; i++) {
    const move = props.history.moves[i]
    const piece: PieceType = ((i % 2) + 1) as PieceType
    newBoard[move.row][move.col] = piece
  }
  board.value = newBoard
}, { immediate: true })

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function nextStep() {
  if (currentStep.value < props.history.moves.length) {
    currentStep.value++
  }
}

function goToEnd() {
  currentStep.value = props.history.moves.length
}

function goToStart() {
  currentStep.value = 0
}

function togglePlay() {
  if (isPlaying.value) {
    stopPlay()
  } else {
    startPlay()
  }
}

function startPlay() {
  if (isComplete.value) {
    currentStep.value = 0
  }
  isPlaying.value = true
  playInterval.value = window.setInterval(() => {
    if (currentStep.value < props.history.moves.length) {
      currentStep.value++
    } else {
      stopPlay()
    }
  }, 800)
}

function stopPlay() {
  isPlaying.value = false
  if (playInterval.value) {
    clearInterval(playInterval.value)
    playInterval.value = null
  }
}

function isCurrentMove(row: number, col: number): boolean {
  if (currentStep.value === 0) return false
  const move = props.history.moves[currentStep.value - 1]
  return move?.row === row && move?.col === col
}

onUnmounted(() => {
  stopPlay()
})
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="replay-modal glass-panel">
      <div class="modal-header">
        <h2>🔄 复盘</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-content">
        <!-- 棋盘 -->
        <div class="replay-board glass-panel-dark">
          <div class="board-grid">
            <svg class="grid-lines" viewBox="0 0 540 540">
              <line v-for="i in 15" :key="`h-${i}`"
                x1="18" :y1="(i - 0.5) * 36"
                x2="522" :y2="(i - 0.5) * 36"
                stroke="rgba(50, 40, 30, 0.6)" stroke-width="1" />
              <line v-for="i in 15" :key="`v-${i}`"
                :x1="(i - 0.5) * 36" y1="18"
                :x2="(i - 0.5) * 36" y2="522"
                stroke="rgba(50, 40, 30, 0.6)" stroke-width="1" />
              <circle v-for="star in [
                { x: 3, y: 3 }, { x: 3, y: 7 }, { x: 3, y: 11 },
                { x: 7, y: 3 }, { x: 7, y: 7 }, { x: 7, y: 11 },
                { x: 11, y: 3 }, { x: 11, y: 7 }, { x: 11, y: 11 }
              ]" :key="`star-${star.x}-${star.y}`"
                :cx="(star.x + 0.5) * 36" :cy="(star.y + 0.5) * 36" r="4"
                fill="rgba(50, 40, 30, 0.7)" />
            </svg>
          </div>
          <div class="pieces-layer">
            <div v-for="(row, rowIndex) in board" :key="rowIndex" class="board-row">
              <div v-for="(cell, colIndex) in row" :key="colIndex" class="board-cell">
                <ChessPiece
                  :piece="cell"
                  :isLastMove="isCurrentMove(rowIndex, colIndex)"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 步数信息 -->
        <div class="step-info">
          <span class="step-current">{{ currentStep }}</span>
          <span class="step-divider">/</span>
          <span class="step-total">{{ history.moves.length }}</span>
          <span class="step-label">步</span>
        </div>

        <!-- 控制按钮 -->
        <div class="replay-controls">
          <button class="control-btn" @click="goToStart" :disabled="currentStep === 0">
            ⏮
          </button>
          <button class="control-btn" @click="prevStep" :disabled="currentStep === 0">
            ◀
          </button>
          <button class="control-btn play-btn" @click="togglePlay">
            {{ isPlaying ? '⏸' : '▶' }}
          </button>
          <button class="control-btn" @click="nextStep" :disabled="isComplete">
            ▶
          </button>
          <button class="control-btn" @click="goToEnd" :disabled="isComplete">
            ⏭
          </button>
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
  animation: fadeIn 0.3s ease;
}

.replay-modal {
  width: 90%;
  max-width: 620px;
  animation: scaleIn 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h2 {
  color: white;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 28px;
  cursor: pointer;
}

.close-btn:hover {
  color: white;
}

.modal-content {
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.replay-board {
  position: relative;
  padding: 10px;
  background:
    linear-gradient(135deg, rgba(210, 180, 140, 0.9), rgba(180, 140, 90, 0.9)),
    repeating-linear-gradient(45deg, rgba(0, 0, 0, 0.02), rgba(0, 0, 0, 0.02) 2px, transparent 2px, transparent 4px);
}

.board-grid {
  position: relative;
}

.grid-lines {
  display: block;
  width: 540px;
  height: 540px;
}

.pieces-layer {
  position: absolute;
  top: 10px;
  left: 10px;
}

.board-row {
  display: flex;
}

.board-cell {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-info {
  display: flex;
  align-items: baseline;
  gap: 4px;
  color: white;
}

.step-current {
  font-size: 32px;
  font-weight: bold;
  color: #00d9ff;
}

.step-divider {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.5);
}

.step-total {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.7);
}

.step-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  margin-left: 4px;
}

.replay-controls {
  display: flex;
  gap: 12px;
}

.control-btn {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.2s;
}

.control-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

.control-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.play-btn {
  width: 56px;
  height: 56px;
  background: rgba(0, 217, 255, 0.3);
  border-color: rgba(0, 217, 255, 0.5);
}

.play-btn:hover:not(:disabled) {
  background: rgba(0, 217, 255, 0.5);
  box-shadow: 0 0 20px rgba(0, 217, 255, 0.3);
}
</style>
