import { ref, computed, onUnmounted, type Ref, type ComputedRef } from 'vue'
import type { WSMessage, MessageType } from '@/types/websocket'

export type ConnectionStatus = 'connecting' | 'connected' | 'disconnected' | 'error'

export interface UseWebSocketReturn {
  status: Ref<ConnectionStatus>
  isConnected: ComputedRef<boolean>
  error: Ref<string | null>
  connect: () => void
  disconnect: () => void
  send: <T>(type: MessageType, payload: T) => void
  subscribe: <T>(type: MessageType, handler: (payload: T) => void) => () => void
  lastMessage: Ref<WSMessage | null>
}

// ========== 单例状态 ==========
// 获取WebSocket URL
const getWebSocketUrl = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'

  // 开发环境下，前端运行在5173端口，后端在8080端口
  // 生产环境下，前后端同端口
  if (import.meta.env.DEV) {
    // 开发模式：固定连接到后端8080端口
    const hostname = window.location.hostname || 'localhost'
    return `ws://${hostname}:8080/ws`
  } else {
    // 生产模式：使用相对路径
    const host = window.location.host
    return `${protocol}//${host}/ws`
  }
}

console.log('[WebSocket Module] Loaded')

// 共享状态
const status = ref<ConnectionStatus>('disconnected')
const isConnected = computed(() => status.value === 'connected')
const error = ref<string | null>(null)
const lastMessage = ref<WSMessage | null>(null)

// 共享WebSocket实例
let ws: WebSocket | null = null
let reconnectAttempts = 0
let heartbeatTimer: ReturnType<typeof setInterval> | null = null
let connectionCount = 0
let isConnecting = false

// 共享配置
const config = {
  getUrl: getWebSocketUrl, // 改为函数，在连接时获取
  reconnect: true,
  reconnectInterval: 3000,
  maxReconnectAttempts: 5,
  heartbeatInterval: 10000,
}

// 消息处理器映射
const handlers = new Map<MessageType, Set<(payload: unknown) => void>>()

// ========== 连接函数 ==========
function connect() {
  console.log('[WebSocket] connect() called')
  console.log('[WebSocket] Current state:', {
    wsState: ws?.readyState,
    status: status.value,
    isConnecting,
    connectionCount
  })

  if (ws?.readyState === WebSocket.OPEN) {
    console.log('[WebSocket] Already connected')
    return
  }

  if (ws?.readyState === WebSocket.CONNECTING || isConnecting) {
    console.log('[WebSocket] Already connecting')
    return
  }

  isConnecting = true
  status.value = 'connecting'
  error.value = null

  try {
    const url = config.getUrl()
    console.log('[WebSocket] Creating WebSocket connection to:', url)
    ws = new WebSocket(url)

    ws.onopen = () => {
      isConnecting = false
      status.value = 'connected'
      reconnectAttempts = 0
      console.log('[WebSocket] ✅ Connected successfully!')
      startHeartbeat()

      // 连接成功后请求房间列表
      console.log('[WebSocket] Requesting room list...')
      send('get_room_list', {})
    }

    ws.onclose = (event) => {
      isConnecting = false
      status.value = 'disconnected'
      stopHeartbeat()
      console.log('[WebSocket] ❌ Disconnected, code:', event.code, 'reason:', event.reason)

      // 自动重连
      if (config.reconnect && reconnectAttempts < config.maxReconnectAttempts) {
        reconnectAttempts++
        console.log(`[WebSocket] 🔄 Reconnecting... (${reconnectAttempts}/${config.maxReconnectAttempts})`)
        setTimeout(connect, config.reconnectInterval)
      }
    }

    ws.onerror = (event) => {
      isConnecting = false
      status.value = 'error'
      error.value = 'WebSocket连接错误'
      console.error('[WebSocket] ⚠️ Connection error!', event)
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        lastMessage.value = msg
        console.log('[WebSocket] 📩 Received:', msg.type, msg.payload)

        // 分发消息到订阅者
        const typeHandlers = handlers.get(msg.type)
        if (typeHandlers) {
          console.log(`[WebSocket] Dispatching ${msg.type} to ${typeHandlers.size} handlers`)
          typeHandlers.forEach((handler) => {
            try {
              handler(msg.payload)
            } catch (e) {
              console.error(`[WebSocket] Handler error for ${msg.type}:`, e)
            }
          })
        } else {
          console.log(`[WebSocket] No handlers registered for: ${msg.type}`)
        }
      } catch (e) {
        console.error('[WebSocket] Parse message error:', e)
      }
    }
  } catch (e) {
    isConnecting = false
    status.value = 'error'
    error.value = '无法创建WebSocket连接'
    console.error('[WebSocket] Create connection error:', e)
  }
}

function disconnect() {
  console.log('[WebSocket] disconnect() called')
  stopHeartbeat()
  if (ws) {
    ws.close(1000, 'Client disconnect')
    ws = null
  }
  status.value = 'disconnected'
  isConnecting = false
}

function send<T>(type: MessageType, payload: T) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.warn('[WebSocket] Cannot send, not connected. State:', ws?.readyState, 'Message type:', type)
    return
  }

  const msg: WSMessage<T> = {
    type,
    payload,
    timestamp: Date.now(),
  }

  console.log('[WebSocket] 📤 Sending:', type, payload)
  ws.send(JSON.stringify(msg))
}

function subscribe<T>(type: MessageType, handler: (payload: T) => void): () => void {
  console.log('[WebSocket] Subscribing to:', type)

  if (!handlers.has(type)) {
    handlers.set(type, new Set())
  }
  handlers.get(type)!.add(handler as (payload: unknown) => void)

  console.log('[WebSocket] Total handlers for', type, ':', handlers.get(type)!.size)

  // 返回取消订阅函数
  return () => {
    handlers.get(type)?.delete(handler as (payload: unknown) => void)
    console.log('[WebSocket] Unsubscribed from:', type)
  }
}

function startHeartbeat() {
  stopHeartbeat()
  heartbeatTimer = setInterval(() => {
    if (ws?.readyState === WebSocket.OPEN) {
      send('heartbeat', {})
    }
  }, config.heartbeatInterval)
}

function stopHeartbeat() {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

// ========== 导出的composable函数 ==========
export function useWebSocket(): UseWebSocketReturn {
  console.log('[useWebSocket] Called')
  console.log('[useWebSocket] Current state:', {
    connectionCount,
    status: status.value,
    isConnected: isConnected.value
  })

  // 引用计数
  connectionCount++
  console.log('[useWebSocket] connectionCount incremented to:', connectionCount)

  // 首次使用时自动连接
  if (connectionCount === 1 && status.value === 'disconnected' && !isConnecting) {
    console.log('[useWebSocket] First usage, scheduling connection...')
    setTimeout(() => {
      console.log('[useWebSocket] Executing delayed connect()')
      connect()
    }, 50)
  }

  // 组件卸载时
  onUnmounted(() => {
    connectionCount--
    console.log('[useWebSocket] Component unmounted, connectionCount:', connectionCount)
    // 所有组件都卸载后才断开连接
    if (connectionCount <= 0) {
      connectionCount = 0
      console.log('[useWebSocket] All components unmounted, disconnecting')
      disconnect()
    }
  })

  return {
    status,
    isConnected,
    error,
    connect,
    disconnect,
    send,
    subscribe,
    lastMessage,
  }
}
