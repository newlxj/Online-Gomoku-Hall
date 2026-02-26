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
    }, 500)
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
    <!-- 外层光晕 -->
    <div class="piece-glow"></div>
    <!-- 主体 -->
    <div class="piece-body">
      <!-- 内部纹理 -->
      <div class="piece-circuit"></div>
      <!-- 高光 -->
      <div class="piece-highlight"></div>
      <!-- 扫描线效果 -->
      <div class="piece-scanline"></div>
    </div>
    <!-- 数据环 -->
    <div class="piece-data-ring" v-if="isLastMove || isWinning"></div>
  </div>
</template>

<style scoped>
.chess-piece {
  width: 85%;
  height: 85%;
  border-radius: 50%;
  position: relative;
  transform: scale(1);
  opacity: 1;
  cursor: default;
}

/* 新棋子放置动画 */
.chess-piece.piece-new {
  animation: placePiece 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes placePiece {
  0% {
    transform: scale(0) rotate(-180deg);
    opacity: 0;
  }
  60% {
    transform: scale(1.2) rotate(20deg);
  }
  100% {
    transform: scale(1) rotate(0deg);
    opacity: 1;
  }
}

/* 外层光晕 */
.piece-glow {
  position: absolute;
  top: -8%;
  left: -8%;
  width: 116%;
  height: 116%;
  border-radius: 50%;
  opacity: 0.5;
  animation: glowPulse 2s ease-in-out infinite;
}

@keyframes glowPulse {
  0%, 100% {
    opacity: 0.3;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.05);
  }
}

/* 黑子样式 */
.piece-black .piece-glow {
  background: radial-gradient(circle, rgba(0, 255, 255, 0.4), transparent 70%);
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
}

/* 白子样式 */
.piece-white .piece-glow {
  background: radial-gradient(circle, rgba(255, 0, 255, 0.4), transparent 70%);
  box-shadow: 0 0 20px rgba(255, 0, 255, 0.3);
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

/* 黑子主体 */
.piece-black .piece-body {
  background:
    radial-gradient(circle at 30% 30%, #333, #000 60%);
  box-shadow:
    0 4px 15px rgba(0, 0, 0, 0.8),
    inset 0 -3px 10px rgba(0, 0, 0, 0.5),
    inset 0 3px 10px rgba(0, 255, 255, 0.1);
  border: 1px solid rgba(0, 255, 255, 0.3);
}

/* 白子主体 */
.piece-white .piece-body {
  background:
    radial-gradient(circle at 30% 30%, #fff, #ddd 60%, #bbb);
  box-shadow:
    0 4px 15px rgba(255, 255, 255, 0.3),
    inset 0 -3px 10px rgba(0, 0, 0, 0.1),
    inset 0 3px 10px rgba(255, 0, 255, 0.2);
  border: 1px solid rgba(255, 0, 255, 0.3);
}

/* 电路纹理 */
.piece-circuit {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0.15;
  background-image:
    linear-gradient(45deg, transparent 40%, currentColor 40%, currentColor 60%, transparent 60%),
    linear-gradient(-45deg, transparent 40%, currentColor 40%, currentColor 60%, transparent 60%);
  background-size: 8px 8px;
}

.piece-black .piece-circuit {
  color: var(--neon-cyan);
}

.piece-white .piece-circuit {
  color: var(--neon-magenta);
}

/* 高光效果 */
.piece-highlight {
  position: absolute;
  width: 40%;
  height: 40%;
  top: 10%;
  left: 15%;
  border-radius: 50%;
}

.piece-black .piece-highlight {
  background: radial-gradient(
    ellipse at 30% 30%,
    rgba(0, 255, 255, 0.4),
    transparent 70%
  );
}

.piece-white .piece-highlight {
  background: radial-gradient(
    ellipse at 30% 30%,
    rgba(255, 255, 255, 0.9),
    rgba(255, 255, 255, 0.3) 50%,
    transparent 70%
  );
}

/* 扫描线效果 */
.piece-scanline {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 0, 0, 0.1) 2px,
    rgba(0, 0, 0, 0.1) 4px
  );
  border-radius: 50%;
  opacity: 0.5;
}

.piece-white .piece-scanline {
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(255, 255, 255, 0.1) 2px,
    rgba(255, 255, 255, 0.1) 4px
  );
}

/* 最后一步标记 */
.last-move .piece-glow {
  animation: lastMoveGlow 1.5s ease-in-out infinite;
}

@keyframes lastMoveGlow {
  0%, 100% {
    opacity: 0.5;
    box-shadow: 0 0 15px var(--neon-cyan);
  }
  50% {
    opacity: 1;
    box-shadow:
      0 0 25px var(--neon-cyan),
      0 0 50px var(--neon-cyan);
  }
}

.piece-white.last-move .piece-glow {
  animation: lastMoveGlowWhite 1.5s ease-in-out infinite;
}

@keyframes lastMoveGlowWhite {
  0%, 100% {
    opacity: 0.5;
    box-shadow: 0 0 15px var(--neon-magenta);
  }
  50% {
    opacity: 1;
    box-shadow:
      0 0 25px var(--neon-magenta),
      0 0 50px var(--neon-magenta);
  }
}

/* 数据环 - 最后一步和获胜棋子 */
.piece-data-ring {
  position: absolute;
  top: -15%;
  left: -15%;
  width: 130%;
  height: 130%;
  border-radius: 50%;
  border: 2px solid transparent;
  animation: dataRingRotate 3s linear infinite;
}

.last-move .piece-data-ring {
  border-color: var(--neon-cyan);
  border-style: dashed;
  opacity: 0.8;
}

.piece-white.last-move .piece-data-ring {
  border-color: var(--neon-magenta);
}

@keyframes dataRingRotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* 获胜棋子特效 */
.winning-piece {
  animation: winningPulse 0.8s ease-in-out infinite;
}

.winning-piece .piece-glow {
  background: radial-gradient(circle, rgba(255, 215, 0, 0.6), transparent 70%) !important;
  box-shadow: 0 0 30px rgba(255, 215, 0, 0.5) !important;
  animation: winningGlow 0.8s ease-in-out infinite !important;
}

@keyframes winningPulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes winningGlow {
  0%, 100% {
    opacity: 0.6;
    box-shadow:
      0 0 20px rgba(255, 215, 0, 0.5),
      0 0 40px rgba(255, 215, 0, 0.3);
  }
  50% {
    opacity: 1;
    box-shadow:
      0 0 30px rgba(255, 215, 0, 0.8),
      0 0 60px rgba(255, 215, 0, 0.5),
      0 0 80px rgba(255, 215, 0, 0.3);
  }
}

.winning-piece .piece-data-ring {
  border-color: var(--cyber-gold) !important;
  border-width: 3px;
  animation: dataRingRotate 1.5s linear infinite, winningRingPulse 0.8s ease-in-out infinite;
}

@keyframes winningRingPulse {
  0%, 100% {
    opacity: 0.6;
    box-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
  }
  50% {
    opacity: 1;
    box-shadow: 0 0 20px rgba(255, 215, 0, 0.8);
  }
}

/* 响应式 - 移动端 */
@media (max-width: 768px) {
  .piece-glow {
    top: -6%;
    left: -6%;
    width: 112%;
    height: 112%;
  }

  .piece-data-ring {
    top: -12%;
    left: -12%;
    width: 124%;
    height: 124%;
  }
}

@media (max-width: 480px) {
  .piece-circuit {
    background-size: 6px 6px;
    opacity: 0.1;
  }

  .piece-scanline {
    opacity: 0.3;
  }
}
</style>
