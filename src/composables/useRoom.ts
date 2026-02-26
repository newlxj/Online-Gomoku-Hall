import { ref, computed, readonly, type Ref, type ComputedRef } from 'vue'
import type { Room, RoomInfo, RoomSettings, Player } from '@/types/multiplayer'
import type { Position } from '@/types/game'
import { useWebSocket } from './useWebSocket'
import { useAlias } from './useAlias'

export interface UseRoomReturn {
  // 状态
  currentRoom: Ref<Room | null>
  roomList: Ref<RoomInfo[]>
  isInRoom: ComputedRef<boolean>
  isHost: ComputedRef<boolean>
  isPlaying: ComputedRef<boolean>
  myPlayer: ComputedRef<Player | null>
  opponentPlayer: ComputedRef<Player | null>
  isMyTurn: ComputedRef<boolean>

  // 方法
  createRoom: (name: string, settings: RoomSettings) => void
  joinRoom: (roomId: string) => void
  leaveRoom: () => void
  spectate: (roomId: string) => void
  setReady: () => void
  makeMove: (position: Position) => void
  sendEmoji: (emoji: string) => void
  refreshRoomList: () => void
  requestSurrender: () => void
  requestUndo: () => void
}

// 单例状态 - 避免每次调用都创建新的状态
const currentRoom = ref<Room | null>(null) as Ref<Room | null>
const roomList = ref<RoomInfo[]>([]) as Ref<RoomInfo[]>
const myPlayerId = ref<string>('')
let subscriptionsSetup = false

export function useRoom(): UseRoomReturn {
  console.log('[useRoom] Called')
  const { send, subscribe, isConnected } = useWebSocket()
  const { currentAlias } = useAlias()

  console.log('[useRoom] WebSocket state:', { isConnected: isConnected.value })

  // 计算属性
  const isInRoom = computed(() => currentRoom.value !== null)
  const isHost = computed(() => {
    if (!currentRoom.value) return false
    return currentRoom.value.hostId === myPlayerId.value
  })
  const isPlaying = computed(() => currentRoom.value?.status === 'playing')
  const myPlayer = computed((): Player | null => {
    if (!currentRoom.value) return null
    return currentRoom.value.players.find(p => p.id === myPlayerId.value) || null
  })
  const opponentPlayer = computed((): Player | null => {
    if (!currentRoom.value) return null
    return currentRoom.value.players.find(p => p.id !== myPlayerId.value) || null
  })
  const isMyTurn = computed(() => {
    if (!currentRoom.value || !myPlayer.value) {
      console.log('[useRoom] isMyTurn: false (no room or no myPlayer)')
      return false
    }
    const result = currentRoom.value.currentTurn === myPlayer.value.pieceType
    console.log('[useRoom] isMyTurn:', result, '- currentTurn:', currentRoom.value.currentTurn, 'myPieceType:', myPlayer.value.pieceType)
    return result
  })

  // 只设置一次订阅
  if (!subscriptionsSetup) {
    console.log('[useRoom] Setting up subscriptions...')

    // 订阅消息 - 注意payload结构与后端保持一致
    subscribe<{ rooms: RoomInfo[] }>('room_list', (payload) => {
      console.log('[useRoom] Received room_list:', payload)
      roomList.value = payload.rooms || []
      console.log('[useRoom] Updated roomList:', roomList.value.length, 'rooms')
    })

    subscribe<Room>('room_update', (room) => {
      console.log('[useRoom] Received room_update:', JSON.stringify(room, null, 2))
      currentRoom.value = room

      // 设置 myPlayerId - 直接使用模块级别的 currentAlias
      // 注意：这里需要从 useAlias 获取当前值
      const aliasModule = useAlias()
      const alias = aliasModule.currentAlias.value
      console.log('[useRoom] Current alias:', alias)
      console.log('[useRoom] Room players:', room.players.map(p => ({ id: p.id, alias: p.alias, isReady: p.isReady })))

      if (alias) {
        const myPlayer = room.players.find(p => p.alias === alias)
        if (myPlayer) {
          myPlayerId.value = myPlayer.id
          console.log('[useRoom] Set myPlayerId:', myPlayerId.value, 'for alias:', alias)
        } else {
          console.warn('[useRoom] Could not find player with alias:', alias)
        }
      } else {
        console.warn('[useRoom] No current alias set')
      }
    })

    subscribe<{ roomId: string; position: Position; playerId: string; pieceType: number }>(
      'move',
      (payload) => {
        console.log('[useRoom] Received move:', payload)
      }
    )

    // 订阅时间更新
    subscribe<{ roomId: string; players: Array<{ playerId: string; remainingTime: number; moveTimeLeft: number }> }>(
      'time_update',
      (payload) => {
        console.log('[useRoom] Received time_update:', payload)
        if (currentRoom.value && currentRoom.value.id === payload.roomId) {
          // 更新玩家时间
          for (const playerTime of payload.players) {
            const player = currentRoom.value.players.find(p => p.id === playerTime.playerId)
            if (player) {
              player.remainingTime = playerTime.remainingTime
              player.moveTimeLeft = playerTime.moveTimeLeft
            }
          }
        }
      }
    )

    subscribe<{ code: number; message: string }>('error', (payload) => {
      console.error('[useRoom] Received error:', payload.message)
      alert('错误: ' + payload.message)
    })

    subscriptionsSetup = true
    console.log('[useRoom] Subscriptions setup complete')
  }

  // 创建房间
  function createRoom(name: string, settings: RoomSettings) {
    if (!isConnected.value || !currentAlias.value) {
      console.warn('[useRoom] Cannot create room: not connected or no alias')
      return
    }

    console.log('[useRoom] Creating room:', name, settings)
    send('create_room', {
      name,
      settings,
      alias: currentAlias.value,
    })
  }

  // 加入房间
  function joinRoom(roomId: string) {
    if (!isConnected.value || !currentAlias.value) {
      console.warn('[useRoom] Cannot join room: not connected or no alias')
      return
    }

    console.log('[useRoom] Joining room:', roomId)
    send('join_room', {
      roomId,
      alias: currentAlias.value,
    })
  }

  // 离开房间
  function leaveRoom() {
    if (!currentRoom.value) return

    console.log('[useRoom] Leaving room')
    send('leave_room', {})
    currentRoom.value = null
    myPlayerId.value = ''
  }

  // 观战
  function spectate(roomId: string) {
    if (!isConnected.value || !currentAlias.value) {
      console.warn('[useRoom] Cannot spectate: not connected or no alias')
      return
    }

    console.log('[useRoom] Spectating room:', roomId)
    send('spectate', {
      roomId,
      alias: currentAlias.value,
    })
  }

  // 准备
  function setReady() {
    if (!currentRoom.value) return

    console.log('[useRoom] Setting ready')
    send('ready', {
      roomId: currentRoom.value.id,
    })
  }

  // 落子
  function makeMove(position: Position) {
    if (!currentRoom.value || !isMyTurn.value) return

    console.log('[useRoom] Making move:', position)
    send('move', {
      roomId: currentRoom.value.id,
      position,
    })
  }

  // 发送表情
  function sendEmoji(emoji: string) {
    if (!currentRoom.value) return

    console.log('[useRoom] Sending emoji:', emoji)
    send('emoji', {
      roomId: currentRoom.value.id,
      emoji,
    })
  }

  // 刷新房间列表
  function refreshRoomList() {
    console.log('[useRoom] Refreshing room list, isConnected:', isConnected.value)
    send('get_room_list', {})
  }

  // 请求认输
  function requestSurrender() {
    if (!currentRoom.value) return

    console.log('[useRoom] Requesting surrender')
    send('surrender_request', {
      roomId: currentRoom.value.id,
    })
  }

  // 请求悔棋
  function requestUndo() {
    if (!currentRoom.value) return

    console.log('[useRoom] Requesting undo')
    send('undo_request', {
      roomId: currentRoom.value.id,
    })
  }

  return {
    currentRoom: currentRoom,
    roomList: roomList,
    isInRoom,
    isHost,
    isPlaying,
    myPlayer,
    opponentPlayer,
    isMyTurn,
    createRoom,
    joinRoom,
    leaveRoom,
    spectate,
    setReady,
    makeMove,
    sendEmoji,
    refreshRoomList,
    requestSurrender,
    requestUndo,
  }
}
