import { ref, onUnmounted } from 'vue'

export interface UseTimerReturn {
  formatTime: (seconds: number) => string
  isLowTime: (seconds: number, threshold?: number) => boolean
  isCriticalTime: (seconds: number, threshold?: number) => boolean
}

export function useTimer(): UseTimerReturn {
  /**
   * 格式化时间显示
   * @param seconds 秒数
   * @returns 格式化后的时间字符串 (MM:SS 或 HH:MM:SS)
   */
  function formatTime(seconds: number): string {
    if (seconds < 0) seconds = 0

    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = seconds % 60

    const pad = (n: number) => n.toString().padStart(2, '0')

    if (hours > 0) {
      return `${hours}:${pad(minutes)}:${pad(secs)}`
    }
    return `${pad(minutes)}:${pad(secs)}`
  }

  /**
   * 判断是否是低时间(黄色警告)
   */
  function isLowTime(seconds: number, threshold: number = 60): boolean {
    return seconds > 0 && seconds <= threshold
  }

  /**
   * 判断是否是危险时间(红色警告)
   */
  function isCriticalTime(seconds: number, threshold: number = 30): boolean {
    return seconds > 0 && seconds <= threshold
  }

  return {
    formatTime,
    isLowTime,
    isCriticalTime,
  }
}

/**
 * 游戏计时器Composable
 * 用于管理玩家的总时间和每步限时
 */
export interface UseGameTimerReturn {
  player1Time: ReturnType<typeof ref<number>>
  player2Time: ReturnType<typeof ref<number>>
  player1MoveTime: ReturnType<typeof ref<number>>
  player2MoveTime: ReturnType<typeof ref<number>>
  isRunning: ReturnType<typeof ref<boolean>>
  currentPlayer: ReturnType<typeof ref<1 | 2>>
  start: () => void
  pause: () => void
  reset: (timeLimit: number, moveTimeLimit: number) => void
  switchTurn: () => void
  tick: () => { timeout: boolean; which: 'total' | 'move' | null; player: 1 | 2 | null }
}

export function useGameTimer(): UseGameTimerReturn {
  const player1Time = ref(600) // 默认10分钟
  const player2Time = ref(600)
  const player1MoveTime = ref(30) // 默认30秒
  const player2MoveTime = ref(30)

  const isRunning = ref(false)
  const currentPlayer = ref<1 | 2>(1) // 1=黑棋先行

  let timerInterval: number | null = null

  /**
   * 开始计时
   */
  function start() {
    if (isRunning.value) return
    isRunning.value = true

    timerInterval = window.setInterval(() => {
      tick()
    }, 1000)
  }

  /**
   * 暂停计时
   */
  function pause() {
    isRunning.value = false
    if (timerInterval) {
      clearInterval(timerInterval)
      timerInterval = null
    }
  }

  /**
   * 重置计时器
   */
  function reset(timeLimit: number, moveTimeLimit: number) {
    pause()
    player1Time.value = timeLimit
    player2Time.value = timeLimit
    player1MoveTime.value = moveTimeLimit
    player2MoveTime.value = moveTimeLimit
    currentPlayer.value = 1
  }

  /**
   * 切换回合
   */
  function switchTurn() {
    // 重置当前玩家的步时
    if (currentPlayer.value === 1) {
      // 黑棋切换到白棋时，重置白棋的步时
      player2MoveTime.value = player1MoveTime.value // 使用配置的步时限时
    } else {
      // 白棋切换到黑棋时，重置黑棋的步时
      player1MoveTime.value = player2MoveTime.value
    }
    currentPlayer.value = currentPlayer.value === 1 ? 2 : 1
  }

  /**
   * 计时tick
   * @returns 超时信息
   */
  function tick(): { timeout: boolean; which: 'total' | 'move' | null; player: 1 | 2 | null } {
    if (!isRunning.value) {
      return { timeout: false, which: null, player: null }
    }

    if (currentPlayer.value === 1) {
      // 黑棋回合
      player1Time.value--
      player1MoveTime.value--

      if (player1MoveTime.value <= 0) {
        // 步时超时
        return { timeout: true, which: 'move', player: 1 }
      }
      if (player1Time.value <= 0) {
        // 总时间超时
        return { timeout: true, which: 'total', player: 1 }
      }
    } else {
      // 白棋回合
      player2Time.value--
      player2MoveTime.value--

      if (player2MoveTime.value <= 0) {
        return { timeout: true, which: 'move', player: 2 }
      }
      if (player2Time.value <= 0) {
        return { timeout: true, which: 'total', player: 2 }
      }
    }

    return { timeout: false, which: null, player: null }
  }

  // 清理
  onUnmounted(() => {
    pause()
  })

  return {
    player1Time,
    player2Time,
    player1MoveTime,
    player2MoveTime,
    isRunning,
    currentPlayer,
    start,
    pause,
    reset,
    switchTurn,
    tick,
  }
}
