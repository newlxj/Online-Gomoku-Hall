<script setup lang="ts">
import { computed } from 'vue'
import type { Spectator } from '@/types/multiplayer'

defineProps<{
  spectators: Spectator[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()
</script>

<template>
  <div class="spectator-panel glass-panel-dark">
    <div class="panel-header">
      <span class="panel-icon">👁️</span>
      <span class="panel-title">观战者 ({{ spectators.length }})</span>
    </div>
    <div class="spectator-list">
      <div
        v-for="spectator in spectators"
        :key="spectator.id"
        class="spectator-item"
      >
        <span class="spectator-icon">👤</span>
        <span class="spectator-alias">{{ spectator.alias }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.spectator-panel {
  padding: 12px 20px;
  border-radius: 12px;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.panel-icon {
  font-size: 16px;
}

.panel-title {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  font-weight: 500;
}

.spectator-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.spectator-item {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.1);
  padding: 6px 12px;
  border-radius: 16px;
}

.spectator-icon {
  font-size: 12px;
}

.spectator-alias {
  color: #fff;
  font-size: 13px;
}
</style>
