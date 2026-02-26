import { ref } from 'vue'
import { generateRandomAlias } from '@/types/multiplayer'
import { useWebSocket } from './useWebSocket'

const ALIAS_STORAGE_KEY = 'gomoku_alias_history'
const CURRENT_ALIAS_KEY = 'gomoku_current_alias'
const MAX_ALIAS_HISTORY = 5

export interface UseAliasReturn {
  currentAlias: ReturnType<typeof ref<string>>
  aliasHistory: ReturnType<typeof ref<string[]>>
  isAvailable: ReturnType<typeof ref<boolean | null>>
  checking: ReturnType<typeof ref<boolean>>
  setCurrentAlias: (alias: string) => void
  generateNewAlias: () => string
  checkAlias: (alias: string) => Promise<boolean>
  saveAliasToHistory: (alias: string) => void
  removeFromHistory: (alias: string) => void
  clearHistory: () => void
}

// ========== 单例状态 ==========
// 当前别名 - 模块级别单例
const currentAlias = ref<string>('')

// 历史别名列表
const aliasHistory = ref<string[]>([])

// 别名可用性检查状态
const isAvailable = ref<boolean | null>(null)
const checking = ref(false)

// 是否已初始化
let aliasInitialized = false

// 从localStorage加载历史
function loadHistory() {
  try {
    const saved = localStorage.getItem(ALIAS_STORAGE_KEY)
    if (saved) {
      aliasHistory.value = JSON.parse(saved)
    }
  } catch (e) {
    console.error('Failed to load alias history', e)
  }
}

// 保存历史到localStorage
function saveHistory() {
  try {
    localStorage.setItem(ALIAS_STORAGE_KEY, JSON.stringify(aliasHistory.value))
  } catch (e) {
    console.error('Failed to save alias history', e)
  }
}

// 从localStorage加载当前别名 - 只执行一次
function initializeAlias() {
  if (aliasInitialized) return

  loadHistory()

  try {
    const saved = localStorage.getItem(CURRENT_ALIAS_KEY)
    if (saved) {
      currentAlias.value = saved
      console.log('[useAlias] Loaded current alias from localStorage:', saved)
    }
  } catch (e) {
    console.error('[useAlias] Failed to load current alias', e)
  }

  aliasInitialized = true
}

// 立即初始化
initializeAlias()

export function useAlias(): UseAliasReturn {
  const { send, subscribe, isConnected } = useWebSocket()

  // 设置当前别名
  function setCurrentAlias(alias: string) {
    console.log('[useAlias] setCurrentAlias called with:', alias)
    currentAlias.value = alias
    isAvailable.value = null
    // 同时保存到 localStorage 以便下次使用
    try {
      localStorage.setItem(CURRENT_ALIAS_KEY, alias)
      console.log('[useAlias] Saved alias to localStorage')
    } catch (e) {
      console.error('[useAlias] Failed to save alias', e)
    }
  }

  // 生成新别名
  function generateNewAlias(): string {
    const newAlias = generateRandomAlias()
    currentAlias.value = newAlias
    isAvailable.value = null
    return newAlias
  }

  // 检查别名是否可用
  function checkAlias(alias: string): Promise<boolean> {
    return new Promise((resolve) => {
      if (!isConnected.value) {
        // 未连接时假设可用
        console.log('[useAlias] WebSocket not connected, assuming alias is available')
        isAvailable.value = true
        checking.value = false
        resolve(true)
        return
      }

      console.log('[useAlias] Checking alias availability:', alias)
      checking.value = true
      isAvailable.value = null
      send('check_alias', { alias })

      // 设置超时
      const timeout = setTimeout(() => {
        console.log('[useAlias] Alias check timeout, assuming available')
        checking.value = false
        isAvailable.value = true // 超时时假设可用
        resolve(true)
      }, 5000)

      // 监听结果
      const unsubscribe = subscribe<{ alias: string; available: boolean }>(
        'alias_check_result',
        (payload) => {
          console.log('[useAlias] Received alias_check_result:', payload)
          if (payload.alias === alias) {
            clearTimeout(timeout)
            checking.value = false
            isAvailable.value = payload.available
            unsubscribe()
            resolve(payload.available)
          }
        }
      )
    })
  }

  // 保存别名到历史
  function saveAliasToHistory(alias: string) {
    if (!alias.trim()) return

    // 移除重复
    const index = aliasHistory.value.indexOf(alias)
    if (index > -1) {
      aliasHistory.value.splice(index, 1)
    }

    // 添加到开头
    aliasHistory.value.unshift(alias)

    // 限制数量
    if (aliasHistory.value.length > MAX_ALIAS_HISTORY) {
      aliasHistory.value = aliasHistory.value.slice(0, MAX_ALIAS_HISTORY)
    }

    saveHistory()
  }

  // 从历史中移除
  function removeFromHistory(alias: string) {
    const index = aliasHistory.value.indexOf(alias)
    if (index > -1) {
      aliasHistory.value.splice(index, 1)
      saveHistory()
    }
  }

  // 清空历史
  function clearHistory() {
    aliasHistory.value = []
    saveHistory()
  }

  return {
    currentAlias,
    aliasHistory,
    isAvailable,
    checking,
    setCurrentAlias,
    generateNewAlias,
    checkAlias,
    saveAliasToHistory,
    removeFromHistory,
    clearHistory,
  }
}
