import { ref, computed, readonly } from 'vue'
import type { Board, PieceType, Position, GameStatus, GameMode, Difficulty, Move } from '@/types/game'

const BOARD_SIZE = 15

export function useGame() {
  // 游戏配置
  const mode = ref<GameMode>('ai')
  const difficulty = ref<Difficulty>('medium')

  // 棋盘状态
  const board = ref<Board>(createEmptyBoard())
  const currentPlayer = ref<PieceType>(1) // 1-黑子先手
  const gameStatus = ref<GameStatus>('playing')
  const moveHistory = ref<Move[]>([])
  const winningLine = ref<Position[]>([])
  const isAiThinking = ref(false)

  // 创建空棋盘
  function createEmptyBoard(): Board {
    return Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0) as PieceType[])
  }

  // 当前是否是玩家回合 (AI模式下玩家执黑)
  const isPlayerTurn = computed(() => {
    if (mode.value === 'pvp') return true
    return currentPlayer.value === 1
  })

  // 落子
  function placePiece(row: number, col: number): boolean {
    if (gameStatus.value !== 'playing') return false
    if (board.value[row][col] !== 0) return false

    board.value[row][col] = currentPlayer.value
    moveHistory.value.push({
      position: { row, col },
      piece: currentPlayer.value,
      timestamp: Date.now()
    })

    // 检查胜负
    const winner = checkWinner(row, col)
    if (winner) {
      gameStatus.value = winner === 1 ? 'black-win' : 'white-win'
    } else if (isBoardFull()) {
      gameStatus.value = 'draw'
    } else {
      currentPlayer.value = currentPlayer.value === 1 ? 2 : 1
    }

    return true
  }

  // 检查胜负
  function checkWinner(row: number, col: number): PieceType | null {
    const piece = board.value[row][col]
    if (piece === 0) return null

    const directions = [
      [[0, 1], [0, -1]],   // 横向
      [[1, 0], [-1, 0]],   // 纵向
      [[1, 1], [-1, -1]],  // 主对角线
      [[1, -1], [-1, 1]]   // 副对角线
    ]

    for (const [dir1, dir2] of directions) {
      const line: Position[] = [{ row, col }]

      // 向两个方向延伸
      for (const [dr, dc] of [dir1, dir2]) {
        let r = row + dr
        let c = col + dc
        while (
          r >= 0 && r < BOARD_SIZE &&
          c >= 0 && c < BOARD_SIZE &&
          board.value[r][c] === piece
        ) {
          line.push({ row: r, col: c })
          r += dr
          c += dc
        }
      }

      if (line.length >= 5) {
        winningLine.value = line
        return piece
      }
    }

    return null
  }

  // 检查棋盘是否已满
  function isBoardFull(): boolean {
    return board.value.every(row => row.every(cell => cell !== 0))
  }

  // 悔棋
  function undo(): boolean {
    if (moveHistory.value.length === 0) return false
    if (gameStatus.value !== 'playing' && gameStatus.value !== 'draw') return false

    // AI模式下需要悔两步
    const stepsToUndo = mode.value === 'ai' && moveHistory.value.length >= 2 ? 2 : 1

    for (let i = 0; i < stepsToUndo && moveHistory.value.length > 0; i++) {
      const lastMove = moveHistory.value.pop()!
      board.value[lastMove.position.row][lastMove.position.col] = 0
    }

    currentPlayer.value = moveHistory.value.length % 2 === 0 ? 1 : 2
    gameStatus.value = 'playing'
    winningLine.value = []
    return true
  }

  // 重新开始
  function resetGame(gameMode?: GameMode, gameDifficulty?: Difficulty) {
    board.value = createEmptyBoard()
    currentPlayer.value = 1
    gameStatus.value = 'playing'
    moveHistory.value = []
    winningLine.value = []
    isAiThinking.value = false

    if (gameMode) mode.value = gameMode
    if (gameDifficulty) difficulty.value = gameDifficulty
  }

  // 获取所有空位
  function getEmptyPositions(): Position[] {
    const positions: Position[] = []
    for (let row = 0; row < BOARD_SIZE; row++) {
      for (let col = 0; col < BOARD_SIZE; col++) {
        if (board.value[row][col] === 0) {
          positions.push({ row, col })
        }
      }
    }
    return positions
  }

  // 获取有效的落子位置 (周围有棋子的空位)
  function getValidMoves(): Position[] {
    const moves: Position[] = []
    const range = 2
    const checked = new Set<string>()

    for (let row = 0; row < BOARD_SIZE; row++) {
      for (let col = 0; col < BOARD_SIZE; col++) {
        if (board.value[row][col] !== 0) {
          // 检查周围的位置
          for (let dr = -range; dr <= range; dr++) {
            for (let dc = -range; dc <= range; dc++) {
              const r = row + dr
              const c = col + dc
              const key = `${r},${c}`
              if (
                r >= 0 && r < BOARD_SIZE &&
                c >= 0 && c < BOARD_SIZE &&
                board.value[r][c] === 0 &&
                !checked.has(key)
              ) {
                checked.add(key)
                moves.push({ row: r, col: c })
              }
            }
          }
        }
      }
    }

    // 如果棋盘为空，返回中心点
    if (moves.length === 0) {
      moves.push({ row: Math.floor(BOARD_SIZE / 2), col: Math.floor(BOARD_SIZE / 2) })
    }

    return moves
  }

  // 获取游戏记录用于保存
  function getGameRecord(): Position[] {
    return moveHistory.value.map(m => m.position)
  }

  // 获取游戏结果
  function getGameResult(): 'black' | 'white' | 'draw' | null {
    if (gameStatus.value === 'black-win') return 'black'
    if (gameStatus.value === 'white-win') return 'white'
    if (gameStatus.value === 'draw') return 'draw'
    return null
  }

  return {
    // 状态
    board: readonly(board),
    currentPlayer: readonly(currentPlayer),
    gameStatus: readonly(gameStatus),
    moveHistory: readonly(moveHistory),
    winningLine: readonly(winningLine),
    isAiThinking: readonly(isAiThinking),
    mode: readonly(mode),
    difficulty: readonly(difficulty),
    isPlayerTurn,
    BOARD_SIZE,

    // 方法
    placePiece,
    undo,
    resetGame,
    getEmptyPositions,
    getValidMoves,
    getGameRecord,
    getGameResult,

    // 内部使用
    _board: board,
    _currentPlayer: currentPlayer,
    _gameStatus: gameStatus,
    _isAiThinking: isAiThinking
  }
}
