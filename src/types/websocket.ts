import type {
  Room,
  RoomInfo,
  RoomSettings,
  LeaderboardEntry,
  Emoji,
} from './multiplayer'
import type { Position } from './game'

// 消息类型
export type MessageType =
  // 客户端 -> 服务器
  | 'create_room'
  | 'join_room'
  | 'leave_room'
  | 'spectate'
  | 'ready'
  | 'move'
  | 'emoji'
  | 'check_alias'
  | 'heartbeat'
  | 'get_room_list'
  | 'get_leaderboard'
  | 'get_online_users'
  | 'enter_lobby'
  | 'leave_lobby'
  | 'surrender_request'
  | 'surrender_response'
  | 'undo_request'
  | 'undo_response'
  // 服务器 -> 客户端
  | 'room_list'
  | 'room_update'
  | 'game_start'
  | 'move'
  | 'game_over'
  | 'time_update'
  | 'player_disconnected'
  | 'alias_check_result'
  | 'emoji'
  | 'achievement_unlock'
  | 'leaderboard_update'
  | 'online_users'
  | 'lobby_entered'
  | 'surrender_request'
  | 'surrender_response'
  | 'undo_request'
  | 'undo_response'
  | 'error'

// 基础消息结构
export interface WSMessage<T = unknown> {
  type: MessageType
  payload: T
  timestamp: number
}

// 创建消息
export function createMessage<T>(type: MessageType, payload: T): WSMessage<T> {
  return {
    type,
    payload,
    timestamp: Date.now(),
  }
}

// ===== 客户端 -> 服务器 Payloads =====

export interface CreateRoomPayload {
  name: string
  settings: RoomSettings
  alias: string
}

export interface JoinRoomPayload {
  roomId: string
  alias: string
}

export interface SpectatePayload {
  roomId: string
  alias: string
}

export interface ReadyPayload {
  roomId: string
}

export interface MovePayload {
  roomId: string
  position: Position
}

export interface EmojiPayload {
  roomId: string
  emoji: Emoji
}

export interface CheckAliasPayload {
  alias: string
}

// ===== 服务器 -> 客户端 Payloads =====

export interface RoomListPayload {
  rooms: RoomInfo[]
}

export interface RoomUpdatePayload extends Room {
  // 继承 Room 的所有字段
}

export interface MoveSyncPayload {
  roomId: string
  position: Position
  playerId: string
  pieceType: number
}

export interface GameOverPayload {
  roomId: string
  winner: number // 0=平局, 1=黑, 2=白
  reason: 'win' | 'disconnect' | 'timeout'
  winningLine?: Position[]
}

export interface TimeUpdatePayload {
  roomId: string
  players: Array<{
    playerId: string
    remainingTime: number
    moveTimeLeft: number
  }>
}

export interface PlayerDisconnectedPayload {
  roomId: string
  playerId: string
}

export interface AliasCheckResultPayload {
  alias: string
  available: boolean
}

export interface EmojiBroadcastPayload {
  roomId: string
  emoji: Emoji
  fromPlayer: string
  fromAlias: string
}

export interface AchievementUnlockPayload {
  achievementId: string
  achievementName: string
  description: string
  reward: number
}

export interface LeaderboardUpdatePayload {
  entries: LeaderboardEntry[]
}

export interface ErrorPayload {
  code: number
  message: string
}

// 认输/悔棋请求 - 服务器转发
export interface GameActionRequestPayload {
  roomId: string
  action: 'surrender' | 'undo'
  fromPlayer: string
  fromAlias: string
  movePosition?: Position // 悔棋时使用
}

// 认输/悔棋响应 - 服务器转发
export interface GameActionResponsePayload {
  roomId: string
  action: 'surrender' | 'undo'
  fromPlayer: string
  fromAlias: string
  accept: boolean
}

// 认输请求 payload
export interface SurrenderRequestPayload {
  roomId: string
}

// 认输响应 payload
export interface SurrenderResponsePayload {
  roomId: string
  accept: boolean
}

// 悔棋请求 payload
export interface UndoRequestPayload {
  roomId: string
}

// 悔棋响应 payload
export interface UndoResponsePayload {
  roomId: string
  accept: boolean
}

// 在线用户信息
export interface OnlineUserInfo {
  id: string
  alias: string
  score: number
  rank: string
  status: 'lobby' | 'playing' | 'spectating'
  roomId?: string
  roomName?: string
}

// 在线用户列表
export interface OnlineUsersPayload {
  count: number
  users: OnlineUserInfo[]
}

// ===== 消息类型守卫 =====

export function isRoomList(msg: WSMessage): msg is WSMessage<RoomListPayload> {
  return msg.type === 'room_list'
}

export function isRoomUpdate(msg: WSMessage): msg is WSMessage<RoomUpdatePayload> {
  return msg.type === 'room_update'
}

export function isMoveSync(msg: WSMessage): msg is WSMessage<MoveSyncPayload> {
  return msg.type === 'move'
}

export function isGameOver(msg: WSMessage): msg is WSMessage<GameOverPayload> {
  return msg.type === 'game_over'
}

export function isTimeUpdate(msg: WSMessage): msg is WSMessage<TimeUpdatePayload> {
  return msg.type === 'time_update'
}

export function isPlayerDisconnected(msg: WSMessage): msg is WSMessage<PlayerDisconnectedPayload> {
  return msg.type === 'player_disconnected'
}

export function isAliasCheckResult(msg: WSMessage): msg is WSMessage<AliasCheckResultPayload> {
  return msg.type === 'alias_check_result'
}

export function isEmojiBroadcast(msg: WSMessage): msg is WSMessage<EmojiBroadcastPayload> {
  return msg.type === 'emoji'
}

export function isAchievementUnlock(msg: WSMessage): msg is WSMessage<AchievementUnlockPayload> {
  return msg.type === 'achievement_unlock'
}

export function isLeaderboardUpdate(msg: WSMessage): msg is WSMessage<LeaderboardUpdatePayload> {
  return msg.type === 'leaderboard_update'
}

export function isError(msg: WSMessage): msg is WSMessage<ErrorPayload> {
  return msg.type === 'error'
}
