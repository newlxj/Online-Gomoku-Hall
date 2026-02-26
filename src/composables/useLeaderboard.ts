import { ref, computed, type Ref, type ComputedRef } from 'vue'
import type { LeaderboardEntry, RankDefinition } from '@/types/multiplayer'
import { getRankByScore } from '@/types/multiplayer'
import { useWebSocket } from './useWebSocket'

export interface UseLeaderboardReturn {
  entries: Ref<LeaderboardEntry[]>
  isLoading: Ref<boolean>
  myRank: ComputedRef<number>
  top3: ComputedRef<LeaderboardEntry[]>
  top10: ComputedRef<LeaderboardEntry[]>
  refresh: () => void
  getRankByScore: (score: number) => RankDefinition
}

const LEADERBOARD_CACHE_KEY = 'gomoku_leaderboard_cache'
const CACHE_DURATION = 60000 // 1分钟缓存

// 单例状态
const entries = ref<LeaderboardEntry[]>([])
const isLoading = ref(false)
const lastUpdate = ref(0)
let subscriptionsSetup = false

export function useLeaderboard(): UseLeaderboardReturn {
  const { send, subscribe, isConnected } = useWebSocket()

  // 从缓存加载
  function loadFromCache() {
    try {
      const cached = localStorage.getItem(LEADERBOARD_CACHE_KEY)
      if (cached) {
        const data = JSON.parse(cached)
        if (Date.now() - data.timestamp < CACHE_DURATION) {
          entries.value = data.entries
          lastUpdate.value = data.timestamp
        }
      }
    } catch (e) {
      console.error('Failed to load leaderboard cache', e)
    }
  }

  // 保存到缓存
  function saveToCache() {
    try {
      localStorage.setItem(LEADERBOARD_CACHE_KEY, JSON.stringify({
        entries: entries.value,
        timestamp: lastUpdate.value,
      }))
    } catch (e) {
      console.error('Failed to save leaderboard cache', e)
    }
  }

  // 初始化加载缓存
  loadFromCache()

  // 只设置一次订阅
  if (!subscriptionsSetup) {
    subscribe<{ entries: LeaderboardEntry[] }>('leaderboard_update', (payload) => {
      entries.value = payload.entries
      lastUpdate.value = Date.now()
      isLoading.value = false
      saveToCache()
    })
    subscriptionsSetup = true
  }

  // 计算属性
  const myRank = computed(() => {
    // TODO: 根据当前玩家别名获取排名
    return -1
  })

  const top3 = computed(() => entries.value.slice(0, 3))
  const top10 = computed(() => entries.value.slice(0, 10))

  // 刷新排行榜
  function refresh() {
    if (!isConnected.value) {
      isLoading.value = false
      return
    }

    isLoading.value = true
    send('get_leaderboard', {})
  }

  // 连接后自动刷新
  if (isConnected.value) {
    refresh()
  }

  return {
    entries,
    isLoading,
    myRank,
    top3,
    top10,
    refresh,
    getRankByScore
  }
}
