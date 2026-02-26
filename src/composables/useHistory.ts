import { ref } from 'vue'
import type { GameHistory, Position, GameMode, Difficulty } from '@/types/game'

const STORAGE_KEY = 'gomoku_history'

// 历史记录管理
export function useHistory() {
  const gameHistories = ref<GameHistory[]>([])

  // 从localStorage加载历史记录
  function loadHistories() {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        gameHistories.value = JSON.parse(stored)
      }
    } catch (e) {
      gameHistories.value = []
    }
  }

  // 保存到localStorage
  function saveHistories() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(gameHistories.value))
    } catch (e) {
      console.error('Failed to save history')
    }
  }

  // 添加游戏记录
  function addHistory(
    mode: GameMode,
    difficulty: Difficulty | undefined,
    winner: 'black' | 'white' | 'draw',
    moves: Position[]
  ) {
    const history: GameHistory = {
      id: Date.now().toString(),
      date: new Date().toLocaleString('zh-CN'),
      mode,
      difficulty,
      winner,
      moves
    }

    gameHistories.value.unshift(history)

    // 只保留最近50条记录
    if (gameHistories.value.length > 50) {
      gameHistories.value = gameHistories.value.slice(0, 50)
    }

    saveHistories()
    return history
  }

  // 多人游戏信息类型
  interface MultiplayerInfo {
    roomName: string
    blackPlayer: string
    whitePlayer: string
    winnerAlias: string
    loserAlias: string
    scoreChanged: number
    reason: 'win' | 'disconnect' | 'timeout' | 'leave'
  }

  // 添加多人游戏记录
  function addMultiplayerHistory(
    winner: 'black' | 'white' | 'draw',
    moves: Position[],
    info: MultiplayerInfo
  ) {
    const history: GameHistory = {
      id: Date.now().toString(),
      date: new Date().toLocaleString('zh-CN'),
      mode: 'pvp',
      winner,
      moves,
      multiplayerInfo: info
    }

    gameHistories.value.unshift(history)

    // 只保留最近50条记录
    if (gameHistories.value.length > 50) {
      gameHistories.value = gameHistories.value.slice(0, 50)
    }

    saveHistories()
    return history
  }

  // 删除单条记录
  function deleteHistory(id: string) {
    const index = gameHistories.value.findIndex(h => h.id === id)
    if (index !== -1) {
      gameHistories.value.splice(index, 1)
      saveHistories()
    }
  }

  // 清空所有记录
  function clearHistories() {
    gameHistories.value = []
    saveHistories()
  }

  // 获取记录详情
  function getHistory(id: string): GameHistory | undefined {
    return gameHistories.value.find(h => h.id === id)
  }

  // 初始化加载
  loadHistories()

  return {
    gameHistories,
    addHistory,
    addMultiplayerHistory,
    deleteHistory,
    clearHistories,
    getHistory
  }
}
