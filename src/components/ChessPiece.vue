<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { PieceType } from '@/types/game'

const props = defineProps<{
  piece: PieceType
  isLastMove?: boolean
  isWinning?: boolean
}>()

// 跟踪是否是新棋子（刚放置的）
const isNewPiece = ref(false)
const previousPiece = ref<PieceType>(0)

// 监听棋子变化，只有从0变为1或2时才播放动画
watch(() => props.piece, (newPiece, oldPiece) => {
  if (oldPiece === 0 && newPiece !== 0) {
    // 棋子刚放置，需要播放动画
    isNewPiece.value = true
    // 动画结束后移除新棋子标记
    setTimeout(() => {
      isNewPiece.value = false
    }, 300)
  }
  previousPiece.value = newPiece
}, { immediate: true })

const pieceClass = computed(() => ({
  'piece-black': props.piece === 1,
  'piece-white': props.piece === 2,
  'last-move': props.isLastMove,
  'winning-piece': props.isWinning,
  'piece-new': isNewPiece.value
}))
</script>

<template>
  <div v-if="piece !== 0" class="chess-piece" :class="pieceClass">
    <!-- 主体 -->
    <div class="piece-body">
      <!-- 高光 -->
      <div class="piece-highlight"></div>
    </div>
    <!-- 最后一步标记 -->
    <div class="last-move-marker" v-if="isLastMove && !isWinning"></div>
  </div>
</template>

<style scoped>
.chess-piece {
  width: 90%;
  height: 90%;
  border-radius: 50%;
  position: relative;
  transform: scale(1);
  opacity: 1;
  cursor: default;
}

/* 新棋子放置动画 */
.chess-piece.piece-new {
  animation: placePiece 0.3s ease-out forwards;
}

@keyframes placePiece {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  60% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 棋子主体 */
.piece-body {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  overflow: hidden;
}

/* 黑子主体 - 传统风格 */
.piece-black .piece-body {
  background: radial-gradient(circle at 35% 35%, #4a4a4a 0%, #1a1a1a 50%, #000 100%);
  box-shadow:
    2px 2px 4px rgba(0, 0, 0, 0.5),
    inset -2px -2px 4px rgba(0, 0, 0, 0.3),
    inset 2px 2px 4px rgba(255, 255, 255, 0.1);
  border: 1px solid #333;
}

/* 白子主体 - 传统风格 */
.piece-white .piece-body {
  background: radial-gradient(circle at 35% 35%, #ffffff 0%, #f0f0f0 50%, #d0d0d0 100%);
  box-shadow:
    2px 2px 4px rgba(0, 0, 0, 0.3),
    inset -2px -2px 4px rgba(0, 0, 0, 0.05),
    inset 2px 2px 4px rgba(255, 255, 255, 0.8);
  border: 1px solid #aaa;
}

/* 高光效果 */
.piece-highlight {
  position: absolute;
  width: 35%;
  height: 35%;
  top: 12%;
  left: 18%;
  border-radius: 50%;
}

.piece-black .piece-highlight {
  background: radial-gradient(
    ellipse at 40% 40%,
    rgba(255, 255, 255, 0.25),
    transparent 70%
  );
}

.piece-white .piece-highlight {
  background: radial-gradient(
    ellipse at 40% 40%,
    rgba(255, 255, 255, 0.9),
    rgba(255, 255, 255, 0.3) 50%,
    transparent 70%
  );
}

/* 最后一步标记 - 红色小圆点 */
.last-move-marker {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 25%;
  height: 25%;
  border-radius: 50%;
  background: #e74c3c;
  box-shadow: 0 0 4px rgba(231, 76, 60, 0.5);
  animation: lastMovePulse 1.5s ease-in-out infinite;
}

@keyframes lastMovePulse {
  0%, 100% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: 0.7;
    transform: translate(-50%, -50%) scale(1.1);
  }
}

/* 获胜棋子特效 */
.winning-piece {
  animation: winningPulse 0.8s ease-in-out infinite;
}

.winning-piece .piece-body {
  box-shadow:
    0 0 15px rgba(255, 215, 0, 0.6),
    2px 2px 4px rgba(0, 0, 0, 0.5),
    inset -2px -2px 4px rgba(0, 0, 0, 0.3),
    inset 2px 2px 4px rgba(255, 255, 255, 0.1);
}

@keyframes winningPulse {
  0%, 100% {
    transform: scale(1);
    filter: brightness(1);
  }
  50% {
    transform: scale(1.05);
    filter: brightness(1.1);
  }
}

.winning-piece .last-move-marker {
  display: none;
}

/* 响应式 - 移动端 */
@media (max-width: 768px) {
  .piece-highlight {
    width: 30%;
    height: 30%;
    top: 15%;
    left: 20%;
  }

  .last-move-marker {
    width: 20%;
    height: 20%;
  }
}

@media (max-width: 480px) {
  .piece-highlight {
    width: 25%;
    height: 25%;
  }

  .last-move-marker {
    width: 18%;
    height: 18%;
  }
}
</style>
