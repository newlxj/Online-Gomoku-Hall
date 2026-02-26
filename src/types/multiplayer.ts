import type { PieceType, Position, Board } from './game'

// 玩家信息
export interface Player {
  id: string
  alias: string
  pieceType: PieceType
  isReady: boolean
  remainingTime: number
  moveTimeLeft: number
  isConnected: boolean
  score: number
  rank: string
}

// 观战者信息
export interface Spectator {
  id: string
  alias: string
  joinedAt: number
}

// 房间设置
export interface RoomSettings {
  timeLimit: number      // 每人总时间(秒)
  moveTimeLimit: number  // 每步限时(秒)
  firstMove: 'host' | 'guest' | 'random'
  ratedGame: boolean
}

// 房间状态
export type RoomStatus = 'waiting' | 'ready' | 'playing' | 'finished'

// 落子记录
export interface MoveRecord {
  position: Position
  pieceType: PieceType
  timestamp: number
  playerId: string
}

// 房间完整信息
export interface Room {
  id: string
  name: string
  hostId: string
  players: Player[]
  spectators: Spectator[]
  status: RoomStatus
  currentTurn: PieceType
  board: Board
  settings: RoomSettings
  moveHistory: MoveRecord[]
}

// 房间简要信息(用于列表)
export interface RoomInfo {
  id: string
  name: string
  hostId: string
  hostName: string
  playerCount: number
  spectatorCount: number
  status: RoomStatus
  createdAt: number
  settings: RoomSettings
}

// 排行榜条目
export interface LeaderboardEntry {
  rank: number
  alias: string
  wins: number
  losses: number
  draws: number
  winRate: number
  score: number
  rankTitle: string
}

// 段位定义
export interface RankDefinition {
  title: string
  min: number
  max: number
  icon: string
  color: string
}

// 段位列表
export const RANKS: RankDefinition[] = [
  { title: '大师', min: 1000, max: 999999, icon: '👑', color: '#FFD700' },
  { title: '钻石', min: 600, max: 999, icon: '💎', color: '#00CED1' },
  { title: '黄金', min: 300, max: 599, icon: '🟡', color: '#FFD700' },
  { title: '白银', min: 100, max: 299, icon: '⚪', color: '#C0C0C0' },
  { title: '青铜', min: 0, max: 99, icon: '🥉', color: '#CD7F32' },
]

// 获取段位信息
export function getRankByScore(score: number): RankDefinition {
  for (const rank of RANKS) {
    if (score >= rank.min && score <= rank.max) {
      return rank
    }
  }
  return RANKS[RANKS.length - 1]
}

// 成就定义
export interface AchievementDefinition {
  id: string
  name: string
  description: string
  reward: number
  icon: string
}

// 成就列表
export const ACHIEVEMENTS: AchievementDefinition[] = [
  { id: 'first_win', name: '初出茅庐', description: '首次获胜', reward: 20, icon: '🏆' },
  { id: 'ten_wins', name: '初窥门径', description: '累计10胜', reward: 50, icon: '🎖️' },
  { id: 'hundred_wins', name: '登峰造极', description: '累计100胜', reward: 200, icon: '🏅' },
  { id: 'perfect_win', name: '完美胜利', description: '对手未下一子获胜', reward: 30, icon: '✨' },
  { id: 'speed_win', name: '闪电战', description: '60秒内获胜', reward: 25, icon: '⚡' },
  { id: 'streak_3', name: '三连胜', description: '连续3胜', reward: 40, icon: '🔥' },
  { id: 'streak_5', name: '五连胜', description: '连续5胜', reward: 100, icon: '💥' },
  { id: 'spectator', name: '围观群众', description: '观战10局', reward: 10, icon: '👀' },
]

// 每日任务
export interface DailyTask {
  id: string
  name: string
  description: string
  reward: number
  target: number
  progress: number
  completed: boolean
  claimed: boolean
}

// 默认每日任务
export const DEFAULT_DAILY_TASKS: Omit<DailyTask, 'progress' | 'completed' | 'claimed'>[] = [
  { id: 'play_3', name: '棋局爱好者', description: '完成3局对局', reward: 10, target: 3 },
  { id: 'win_1', name: '首胜达成', description: '获得1场胜利', reward: 15, target: 1 },
  { id: 'login_3', name: '连续登录', description: '连续登录3天', reward: 20, target: 3 },
  { id: 'spectate_2', name: '观战达人', description: '观战2局', reward: 5, target: 2 },
]

// 表情列表
export const EMOJIS = ['👍', '😂', '😮', '😢', '🤔', '💪', '🎉', '🙏'] as const
export type Emoji = typeof EMOJIS[number]

// 随机别名生成器
const ADJECTIVES = [
  '疾风', '暗影', '烈焰', '冰霜', '雷霆',
  '星空', '月影', '日光', '云霄', '海浪',
  '山岳', '森林', '沙漠', '极光', '彩虹',
  '神秘', '勇敢', '智慧', '沉稳', '迅猛',
  '精准', '冷静', '热情', '坚韧', '灵动',
]

const NOUNS = [
  '棋手', '大师', '王者', '侠客', '剑士',
  '法师', '猎手', '骑士', '战士', '智者',
  '圣者', '仙人', '英雄', '传奇', '宗师',
  '高手', '达人', '专家', '新星', '霸主',
]

export function generateRandomAlias(): string {
  const adj = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)]
  const noun = NOUNS[Math.floor(Math.random() * NOUNS.length)]
  return adj + noun
}
