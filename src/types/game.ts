// 棋子类型: 0-空, 1-黑子, 2-白子
export type PieceType = 0 | 1 | 2

// 坐标
export interface Position {
  row: number
  col: number
}

// 游戏模式
export type GameMode = 'ai' | 'pvp'

// AI难度
export type Difficulty = 'easy' | 'medium' | 'hard'

// 游戏状态
export type GameStatus = 'playing' | 'black-win' | 'white-win' | 'draw'

// 游戏配置
export interface GameConfig {
  mode: GameMode
  difficulty: Difficulty
  boardSize: number
}

// 棋盘类型
export type Board = PieceType[][]

// 历史记录
export interface GameHistory {
  id: string
  date: string
  mode: GameMode
  difficulty?: Difficulty
  winner: 'black' | 'white' | 'draw'
  moves: Position[]
  // 多人游戏额外信息
  multiplayerInfo?: {
    roomName: string
    blackPlayer: string
    whitePlayer: string
    winnerAlias: string
    loserAlias: string
    scoreChanged: number
    reason: 'win' | 'disconnect' | 'timeout' | 'leave'
  }
}

// 一步棋
export interface Move {
  position: Position
  piece: PieceType
  timestamp: number
}

// AI评估分数
export const SCORES = {
  FIVE: 100000,
  LIVE_FOUR: 10000,
  RUSH_FOUR: 1000,
  LIVE_THREE: 1000,
  SLEEP_THREE: 100,
  LIVE_TWO: 100,
  SLEEP_TWO: 10,
  LIVE_ONE: 1
} as const
