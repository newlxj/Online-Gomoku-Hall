import type { Board, PieceType, Position, Difficulty } from '@/types/game'
import { SCORES } from '@/types/game'

const BOARD_SIZE = 15

// AI玩家
export function useAI() {
  // 获取AI落子位置
  function getAIMove(board: Board, difficulty: Difficulty, aiPiece: PieceType): Position {
    const validMoves = getValidMoves(board)

    if (validMoves.length === 0) {
      return { row: Math.floor(BOARD_SIZE / 2), col: Math.floor(BOARD_SIZE / 2) }
    }

    switch (difficulty) {
      case 'easy':
        return getEasyMove(board, validMoves, aiPiece)
      case 'medium':
        return getMediumMove(board, validMoves, aiPiece)
      case 'hard':
        return getHardMove(board, validMoves, aiPiece)
      default:
        return getMediumMove(board, validMoves, aiPiece)
    }
  }

  // 简单难度: 随机 + 基础防守
  function getEasyMove(board: Board, validMoves: Position[], aiPiece: PieceType): Position {
    const playerPiece = aiPiece === 1 ? 2 : 1

    // 检查是否能赢
    for (const move of validMoves) {
      if (wouldWin(board, move, aiPiece)) {
        return move
      }
    }

    // 检查是否需要防守
    for (const move of validMoves) {
      if (wouldWin(board, move, playerPiece)) {
        return move
      }
    }

    // 随机选择
    return validMoves[Math.floor(Math.random() * validMoves.length)]
  }

  // 中等难度: 启发式评估 + 2层搜索
  function getMediumMove(board: Board, validMoves: Position[], aiPiece: PieceType): Position {
    const playerPiece = aiPiece === 1 ? 2 : 1
    let bestMove = validMoves[0]
    let bestScore = -Infinity

    // 对每个可能的位置进行评估
    for (const move of validMoves) {
      const score = evaluatePosition(board, move, aiPiece, playerPiece)
      if (score > bestScore) {
        bestScore = score
        bestMove = move
      }
    }

    return bestMove
  }

  // 困难难度: Alpha-Beta剪枝 + 4层搜索
  function getHardMove(board: Board, validMoves: Position[], aiPiece: PieceType): Position {
    const playerPiece = aiPiece === 1 ? 2 : 1
    const depth = 4

    // 先检查必赢和必防
    for (const move of validMoves) {
      if (wouldWin(board, move, aiPiece)) {
        return move
      }
    }

    for (const move of validMoves) {
      if (wouldWin(board, move, playerPiece)) {
        return move
      }
    }

    // 对候选位置进行评分和排序
    const scoredMoves = validMoves.map(move => ({
      move,
      score: evaluatePosition(board, move, aiPiece, playerPiece)
    })).sort((a, b) => b.score - a.score)

    // 只考虑前20个最佳位置
    const candidateMoves = scoredMoves.slice(0, 20).map(s => s.move)

    let bestMove = candidateMoves[0]
    let bestScore = -Infinity
    const alpha = -Infinity
    const beta = Infinity

    for (const move of candidateMoves) {
      const newBoard = makeMove(board, move, aiPiece)
      const score = minimax(newBoard, depth - 1, alpha, beta, false, aiPiece, playerPiece)

      if (score > bestScore) {
        bestScore = score
        bestMove = move
      }
    }

    return bestMove
  }

  // Minimax + Alpha-Beta剪枝
  function minimax(
    board: Board,
    depth: number,
    alpha: number,
    beta: number,
    isMaximizing: boolean,
    aiPiece: PieceType,
    playerPiece: PieceType
  ): number {
    // 终止条件
    if (depth === 0) {
      return evaluateBoard(board, aiPiece, playerPiece)
    }

    const validMoves = getValidMoves(board)
    if (validMoves.length === 0) {
      return evaluateBoard(board, aiPiece, playerPiece)
    }

    if (isMaximizing) {
      let maxScore = -Infinity
      for (const move of validMoves.slice(0, 15)) {
        const newBoard = makeMove(board, move, aiPiece)
        if (wouldWinOnBoard(newBoard, move, aiPiece)) {
          return SCORES.FIVE
        }
        const score = minimax(newBoard, depth - 1, alpha, beta, false, aiPiece, playerPiece)
        maxScore = Math.max(maxScore, score)
        alpha = Math.max(alpha, score)
        if (beta <= alpha) break
      }
      return maxScore
    } else {
      let minScore = Infinity
      for (const move of validMoves.slice(0, 15)) {
        const newBoard = makeMove(board, move, playerPiece)
        if (wouldWinOnBoard(newBoard, move, playerPiece)) {
          return -SCORES.FIVE
        }
        const score = minimax(newBoard, depth - 1, alpha, beta, true, aiPiece, playerPiece)
        minScore = Math.min(minScore, score)
        beta = Math.min(beta, score)
        if (beta <= alpha) break
      }
      return minScore
    }
  }

  // 获取有效的落子位置
  function getValidMoves(board: Board): Position[] {
    const moves: Position[] = []
    const range = 2
    const checked = new Set<string>()

    for (let row = 0; row < BOARD_SIZE; row++) {
      for (let col = 0; col < BOARD_SIZE; col++) {
        if (board[row][col] !== 0) {
          for (let dr = -range; dr <= range; dr++) {
            for (let dc = -range; dc <= range; dc++) {
              const r = row + dr
              const c = col + dc
              const key = `${r},${c}`
              if (
                r >= 0 && r < BOARD_SIZE &&
                c >= 0 && c < BOARD_SIZE &&
                board[r][c] === 0 &&
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

    if (moves.length === 0) {
      moves.push({ row: Math.floor(BOARD_SIZE / 2), col: Math.floor(BOARD_SIZE / 2) })
    }

    return moves
  }

  // 检查某个位置落子后是否能赢
  function wouldWin(board: Board, pos: Position, piece: PieceType): boolean {
    const newBoard = board.map(row => [...row])
    newBoard[pos.row][pos.col] = piece
    return wouldWinOnBoard(newBoard, pos, piece)
  }

  function wouldWinOnBoard(board: Board, pos: Position, piece: PieceType): boolean {
    const directions = [
      [[0, 1], [0, -1]],
      [[1, 0], [-1, 0]],
      [[1, 1], [-1, -1]],
      [[1, -1], [-1, 1]]
    ]

    for (const [dir1, dir2] of directions) {
      let count = 1
      for (const [dr, dc] of [dir1, dir2]) {
        let r = pos.row + dr
        let c = pos.col + dc
        while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE && board[r][c] === piece) {
          count++
          r += dr
          c += dc
        }
      }
      if (count >= 5) return true
    }
    return false
  }

  // 在棋盘上落子
  function makeMove(board: Board, pos: Position, piece: PieceType): Board {
    const newBoard = board.map(row => [...row]) as Board
    newBoard[pos.row][pos.col] = piece
    return newBoard
  }

  // 评估单个位置的分数
  function evaluatePosition(board: Board, pos: Position, aiPiece: PieceType, playerPiece: PieceType): number {
    let score = 0

    // AI进攻分数
    score += evaluatePoint(board, pos.row, pos.col, aiPiece) * 1.1
    // 防守分数
    score += evaluatePoint(board, pos.row, pos.col, playerPiece)

    // 中心位置加分
    const centerRow = Math.floor(BOARD_SIZE / 2)
    const centerCol = Math.floor(BOARD_SIZE / 2)
    const distToCenter = Math.abs(pos.row - centerRow) + Math.abs(pos.col - centerCol)
    score += Math.max(0, 10 - distToCenter)

    return score
  }

  // 评估某个点对某方的价值
  function evaluatePoint(board: Board, row: number, col: number, piece: PieceType): number {
    const directions = [
      [0, 1], [1, 0], [1, 1], [1, -1]
    ]

    let totalScore = 0

    for (const [dr, dc] of directions) {
      const line = getLine(board, row, col, dr, dc, piece)
      totalScore += evaluateLine(line, piece)
    }

    return totalScore
  }

  // 获取一条线上的棋子情况
  function getLine(board: Board, row: number, col: number, dr: number, dc: number, piece: PieceType): (PieceType | -1)[] {
    const line: (PieceType | -1)[] = []

    // 向负方向延伸4格
    for (let i = 4; i >= 1; i--) {
      const r = row - dr * i
      const c = col - dc * i
      if (r < 0 || r >= BOARD_SIZE || c < 0 || c >= BOARD_SIZE) {
        line.push(-1) // 边界
      } else {
        line.push(board[r][c])
      }
    }

    // 当前位置（假设落子）
    line.push(piece)

    // 向正方向延伸4格
    for (let i = 1; i <= 4; i++) {
      const r = row + dr * i
      const c = col + dc * i
      if (r < 0 || r >= BOARD_SIZE || c < 0 || c >= BOARD_SIZE) {
        line.push(-1)
      } else {
        line.push(board[r][c])
      }
    }

    return line
  }

  // 评估一条线的棋型
  function evaluateLine(line: (PieceType | -1)[], piece: PieceType): number {
    let score = 0

    // 检查所有可能的五连位置
    for (let i = 0; i <= 4; i++) {
      const window = line.slice(i, i + 5)
      let pieceCount = 0
      let emptyCount = 0
      let blocked = false

      for (const cell of window) {
        if (cell === piece) pieceCount++
        else if (cell === 0) emptyCount++
        else blocked = true
      }

      if (blocked) continue

      if (pieceCount === 5) score += SCORES.FIVE
      else if (pieceCount === 4 && emptyCount === 1) score += SCORES.LIVE_FOUR
      else if (pieceCount === 3 && emptyCount === 2) score += SCORES.LIVE_THREE
      else if (pieceCount === 2 && emptyCount === 3) score += SCORES.LIVE_TWO
      else if (pieceCount === 1 && emptyCount === 4) score += SCORES.LIVE_ONE
    }

    return score
  }

  // 评估整个棋盘
  function evaluateBoard(board: Board, aiPiece: PieceType, playerPiece: PieceType): number {
    let aiScore = 0
    let playerScore = 0

    for (let row = 0; row < BOARD_SIZE; row++) {
      for (let col = 0; col < BOARD_SIZE; col++) {
        if (board[row][col] === aiPiece) {
          aiScore += evaluatePoint(board, row, col, aiPiece)
        } else if (board[row][col] === playerPiece) {
          playerScore += evaluatePoint(board, row, col, playerPiece)
        }
      }
    }

    return aiScore - playerScore * 1.1
  }

  return {
    getAIMove
  }
}
