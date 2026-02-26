package server

import (
	"encoding/json"
	"fmt"
	"gomoku-server/internal/leaderboard"
	"gomoku-server/internal/message"
	"gomoku-server/internal/player"
	"gomoku-server/internal/room"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Server WebSocket服务器
type Server struct {
	roomManager    *room.Manager
	aliasManager   *player.AliasManager
	leaderboardMgr *leaderboard.Manager

	// 客户端连接管理
	clients    map[string]*websocket.Conn // playerID -> conn
	clientMap  map[*websocket.Conn]string // conn -> playerID
	connMutexs map[*websocket.Conn]*sync.Mutex // 每个连接的写入互斥锁
	mu         sync.RWMutex

	// 在线用户管理
	onlineUsers   map[string]*OnlineUser // playerID -> user info
	onlineUsersMu sync.RWMutex

	// WebSocket升级器
	upgrader websocket.Upgrader

	// 心跳检测
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration
}

// OnlineUser 在线用户信息
type OnlineUser struct {
	ID     string
	Alias  string
	Status string // "lobby", "playing", "spectating"
	RoomID string
}

// NewServer 创建服务器
func NewServer(lbManager *leaderboard.Manager) *Server {
	return &Server{
		roomManager:       room.NewManager(),
		aliasManager:      player.NewAliasManager(),
		leaderboardMgr:    lbManager,
		clients:           make(map[string]*websocket.Conn),
		clientMap:         make(map[*websocket.Conn]string),
		connMutexs:        make(map[*websocket.Conn]*sync.Mutex),
		onlineUsers:       make(map[string]*OnlineUser),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  30 * time.Second,
	}
}

// HandleWebSocket 处理WebSocket连接
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// 生成玩家ID
	playerID := generatePlayerID()

	// 创建连接的写入互斥锁
	connMutex := &sync.Mutex{}

	// 注册连接
	s.mu.Lock()
	s.clients[playerID] = conn
	s.clientMap[conn] = playerID
	s.connMutexs[conn] = connMutex
	s.mu.Unlock()

	log.Printf("Player connected: %s", playerID)

	// 清理函数
	defer func() {
		s.mu.Lock()
		delete(s.clients, playerID)
		delete(s.clientMap, conn)
		delete(s.connMutexs, conn)
		s.mu.Unlock()

		s.handleDisconnect(playerID)
		log.Printf("Player disconnected: %s", playerID)
	}()

	// 启动心跳检测
	go s.heartbeatChecker(playerID, conn)

	// 消息循环
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}

		var msg message.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			s.sendError(conn, 400, "Invalid message format")
			continue
		}

		s.handleMessage(playerID, conn, msg)
	}
}

// handleMessage 处理消息
func (s *Server) handleMessage(playerID string, conn *websocket.Conn, msg message.Message) {
	// 使用 defer recover 防止单个消息处理 panic 影响整个服务器
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[handleMessage] Recovered from panic for player %s: %v", playerID, r)
			s.sendError(conn, 500, "Internal server error")
		}
	}()

	switch msg.Type {
	case message.MsgCheckAlias:
		s.handleCheckAlias(conn, msg)

	case message.MsgCreateRoom:
		s.handleCreateRoom(playerID, conn, msg)

	case message.MsgJoinRoom:
		s.handleJoinRoom(playerID, conn, msg)

	case message.MsgLeaveRoom:
		s.handleLeaveRoom(playerID, conn)

	case message.MsgSpectate:
		s.handleSpectate(playerID, conn, msg)

	case message.MsgReady:
		s.handleReady(playerID, conn, msg)

	case message.MsgMove:
		s.handleMove(playerID, conn, msg)

	case message.MsgEmoji:
		s.handleEmoji(playerID, conn, msg)

	case message.MsgGetRoomList:
		s.handleGetRoomList(conn)

	case message.MsgGetLeaderboard:
		s.handleGetLeaderboard(conn)

	case message.MsgGetOnlineUsers:
		s.handleGetOnlineUsers(conn)

	case message.MsgEnterLobby:
		s.handleEnterLobby(playerID, conn, msg)

	case message.MsgLeaveLobby:
		s.handleLeaveLobby(playerID)

	case message.MsgHeartbeat:
		// 更新活跃时间
		// TODO: 从玩家信息获取别名
	}
}

// handleCheckAlias 处理别名检查
func (s *Server) handleCheckAlias(conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	alias, _ := payload["alias"].(string)
	available := s.aliasManager.IsAvailable(alias)

	s.sendMessage(conn, message.NewMessage(message.MsgAliasCheckResult, message.AliasCheckResultPayload{
		Alias:     alias,
		Available: available,
	}))
}

// handleCreateRoom 处理创建房间
func (s *Server) handleCreateRoom(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	name, _ := payload["name"].(string)
	alias, _ := payload["alias"].(string)
	fmt.Printf("[handleCreateRoom] playerID=%s, name=%s, alias=%s\n", playerID, name, alias)

	// 注册别名
	if !s.aliasManager.Register(alias, playerID) {
		s.sendError(conn, 400, "Alias already in use")
		return
	}

	// 解析房间设置
	settings := s.parseRoomSettings(payload["settings"])

	// 创建房间
	newRoom := s.roomManager.CreateRoom(name, settings, playerID)
	fmt.Printf("[handleCreateRoom] Created room: %s\n", newRoom.ID)

	// 创建玩家
	p := player.NewPlayer(playerID, alias, conn)
	p.Score = s.leaderboardMgr.GetPlayerScore(alias)

	// 添加玩家到房间
	newRoom.AddPlayer(p)
	fmt.Printf("[handleCreateRoom] Added player to room, players count: %d\n", len(newRoom.Players))

	// 更新在线用户状态
	s.addOnlineUser(playerID, alias, "lobby", newRoom.ID)

	// 发送房间更新
	roomInfo := newRoom.GetFullInfo()
	fmt.Printf("[handleCreateRoom] Sending room_update to creator, status: %s\n", roomInfo.Status)
	s.sendMessage(conn, message.NewMessage(message.MsgRoomUpdate, roomInfo))

	// 广播房间列表给所有客户端
	s.broadcastRoomListToAll()
}

// handleJoinRoom 处理加入房间
func (s *Server) handleJoinRoom(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	roomID, _ := payload["roomId"].(string)
	alias, _ := payload["alias"].(string)
	fmt.Printf("[handleJoinRoom] playerID=%s, roomID=%s, alias=%s\n", playerID, roomID, alias)

	// 注册别名
	if !s.aliasManager.Register(alias, playerID) {
		s.sendError(conn, 400, "Alias already in use")
		return
	}

	// 获取房间
	rm := s.roomManager.GetRoom(roomID)
	if rm == nil {
		s.sendError(conn, 404, "Room not found")
		return
	}

	// 创建玩家
	p := player.NewPlayer(playerID, alias, conn)
	p.Score = s.leaderboardMgr.GetPlayerScore(alias)

	// 添加玩家到房间
	if _, err := s.roomManager.AddPlayerToRoom(roomID, p); err != nil {
		s.sendError(conn, 400, err.Error())
		return
	}

	// 更新在线用户状态
	s.addOnlineUser(playerID, alias, "lobby", roomID)

	// 重新获取房间以确保状态最新
	rm = s.roomManager.GetRoom(roomID)
	if rm != nil {
		fmt.Printf("[handleJoinRoom] Broadcasting room update, players: %d\n", len(rm.Players))
		// 广播房间更新给所有人
		s.broadcastRoomUpdate(rm)
	}

	// 广播房间列表给所有客户端
	s.broadcastRoomListToAll()
}

// handleLeaveRoom 处理离开房间
func (s *Server) handleLeaveRoom(playerID string, conn *websocket.Conn) {
	// 先检查是否是观战者
	spectatorRoom := s.findSpectatorRoom(playerID)
	if spectatorRoom != nil {
		// 从房间移除观战者
		s.roomManager.RemoveSpectatorFromRoom(spectatorRoom.ID, playerID)
		fmt.Printf("[handleLeaveRoom] Removed spectator %s from room %s\n", playerID, spectatorRoom.ID)

		// 更新在线用户状态
		s.updateOnlineUserStatus(playerID, "lobby", "")

		// 广播房间更新
		if updatedRoom := s.roomManager.GetRoom(spectatorRoom.ID); updatedRoom != nil {
			s.broadcastRoomUpdate(updatedRoom)
		}

		// 广播房间列表给所有客户端
		s.broadcastRoomListToAll()
		return
	}

	// 查找玩家所在房间
	rm := s.findPlayerRoom(playerID)
	if rm == nil {
		return
	}

	roomID := rm.ID
	wasPlaying := rm.Status == room.StatusPlaying

	// 如果游戏进行中离开，判对方获胜并扣分
	if wasPlaying {
		// 找到离开的玩家和留下的玩家
		var leaverAlias string
		var winnerAlias string
		var winnerPieceType int

		for _, p := range rm.Players {
			if p.ID == playerID {
				leaverAlias = p.Alias
			} else {
				winnerAlias = p.Alias
				winnerPieceType = p.PieceType
			}
		}

		// 停止计时器
		rm.StopTimer()

		// 设置游戏结束状态
		rm.SetGameFinished(winnerPieceType, "leave")

		// 广播游戏结束
		s.broadcastToRoom(rm, message.NewMessage(message.MsgGameOver, message.GameOverPayload{
			RoomID: rm.ID,
			Winner: winnerPieceType,
			Reason: "leave",
		}))

		// 更新排行榜：离开者扣分，留下者获胜
		if rm.Settings.RatedGame && leaverAlias != "" && winnerAlias != "" {
			// 记录结果：离开者判负，留下者获胜
			s.leaderboardMgr.RecordGameResult(winnerAlias, leaverAlias, false, 0)
			fmt.Printf("[handleLeaveRoom] Player %s left during game, %s wins by forfeit\n", leaverAlias, winnerAlias)
		}
	}

	// 移除玩家
	_, removedPlayer := s.roomManager.RemovePlayerFromRoom(roomID, playerID)

	// 注销别名并更新在线用户状态
	if removedPlayer != nil {
		s.aliasManager.Unregister(removedPlayer.Alias)
		s.updateOnlineUserStatus(playerID, "lobby", "")
	}

	// 广播房间更新
	if updatedRoom := s.roomManager.GetRoom(roomID); updatedRoom != nil {
		s.broadcastRoomUpdate(updatedRoom)
	}

	// 广播房间列表给所有客户端
	s.broadcastRoomListToAll()
}

// handleSpectate 处理观战
func (s *Server) handleSpectate(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	roomID, _ := payload["roomId"].(string)
	alias, _ := payload["alias"].(string)
	fmt.Printf("[handleSpectate] playerID=%s, roomID=%s, alias=%s\n", playerID, roomID, alias)

	// 获取房间
	rm := s.roomManager.GetRoom(roomID)
	if rm == nil {
		s.sendError(conn, 404, "Room not found")
		return
	}

	// 创建观战者
	spectator := player.NewSpectator(playerID, alias, conn)

	// 添加观战者
	s.roomManager.AddSpectatorToRoom(roomID, spectator)

	// 更新在线用户状态
	s.addOnlineUser(playerID, alias, "spectating", roomID)

	// 重新获取房间
	rm = s.roomManager.GetRoom(roomID)
	if rm != nil {
		// 发送房间状态给观战者
		s.sendMessage(conn, message.NewMessage(message.MsgRoomUpdate, rm.GetFullInfo()))
		// 广播房间更新
		s.broadcastRoomUpdate(rm)
	}
}

// handleReady 处理准备
func (s *Server) handleReady(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	roomID, _ := payload["roomId"].(string)
	fmt.Printf("[handleReady] playerID=%s, roomID=%s\n", playerID, roomID)

	rm := s.roomManager.GetRoom(roomID)
	if rm == nil {
		s.sendError(conn, 404, "Room not found")
		return
	}

	// 设置准备状态
	rm.SetPlayerReady(playerID, true)
	fmt.Printf("[handleReady] Set player %s ready\n", playerID)

	// 打印所有玩家的准备状态
	for _, p := range rm.Players {
		fmt.Printf("[handleReady] Player %s (alias=%s) IsReady=%v\n", p.ID, p.Alias, p.IsReady)
	}

	// 广播房间更新
	s.broadcastRoomUpdate(rm)

	// 广播房间列表给所有客户端（状态可能变化）
	s.broadcastRoomListToAll()

	// 检查是否可以开始
	canStart := rm.CanStart()
	fmt.Printf("[handleReady] CanStart=%v, Players count=%d\n", canStart, len(rm.Players))

	if canStart {
		fmt.Println("[handleReady] Starting game!")
		rm.StartGame()

		// 更新所有玩家的在线状态为 playing
		for _, p := range rm.Players {
			s.updateOnlineUserStatus(p.ID, "playing", roomID)
		}

		// 启动时间更新广播
		go s.broadcastTimeUpdates(rm)
		// 广播游戏开始
		s.broadcastRoomUpdate(rm)
		// 广播房间列表（游戏开始状态变化）
		s.broadcastRoomListToAll()
	}
}

// handleMove 处理落子
func (s *Server) handleMove(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	roomID, _ := payload["roomId"].(string)
	posData, _ := payload["position"].(map[string]interface{})
	row := int(posData["row"].(float64))
	col := int(posData["col"].(float64))

	fmt.Printf("[handleMove] playerID=%s, roomID=%s, position=(%d, %d)\n", playerID, roomID, row, col)

	rm := s.roomManager.GetRoom(roomID)
	if rm == nil {
		s.sendError(conn, 404, "Room not found")
		return
	}

	// 落子
	success, winner, winningLine := rm.MakeMove(playerID, message.Position{Row: row, Col: col})
	fmt.Printf("[handleMove] MakeMove result: success=%v, winner=%d\n", success, winner)

	if !success {
		fmt.Printf("[handleMove] Move failed, ignoring\n")
		return
	}

	// 广播房间更新（包含新棋盘状态）
	fmt.Printf("[handleMove] Broadcasting room update with new board state\n")
	s.broadcastRoomUpdate(rm)

	// 检查游戏结束
	if winner > 0 || len(winningLine) > 0 {
		s.handleGameOver(rm, winner, "win", winningLine)
	}
}

// handleEmoji 处理表情
func (s *Server) handleEmoji(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	roomID, _ := payload["roomId"].(string)
	emoji, _ := payload["emoji"].(string)

	rm := s.roomManager.GetRoom(roomID)
	if rm == nil {
		return
	}

	// 获取玩家别名
	var fromAlias string
	for _, p := range rm.Players {
		if p.ID == playerID {
			fromAlias = p.Alias
			break
		}
	}

	// 广播表情
	s.broadcastToRoom(rm, message.NewMessage(message.MsgEmojiBroadcast, message.EmojiBroadcastPayload{
		RoomID:     roomID,
		Emoji:      emoji,
		FromPlayer: playerID,
		FromAlias:  fromAlias,
	}))
}

// handleGetRoomList 处理获取房间列表
func (s *Server) handleGetRoomList(conn *websocket.Conn) {
	rooms := s.roomManager.GetAllRooms()
	s.sendMessage(conn, message.NewMessage(message.MsgRoomList, message.RoomListPayload{
		Rooms: rooms,
	}))
}

// handleGetLeaderboard 处理获取排行榜
func (s *Server) handleGetLeaderboard(conn *websocket.Conn) {
	entries := s.leaderboardMgr.GetLeaderboard(100)
	s.sendMessage(conn, message.NewMessage(message.MsgLeaderboardUpdate, message.LeaderboardUpdatePayload{
		Entries: entries,
	}))
}

// handleGameOver 处理游戏结束
func (s *Server) handleGameOver(rm *room.Room, winner int, reason string, winningLine []message.Position) {
	// 停止计时器
	rm.StopTimer()

	// 获取获胜者和失败者信息
	winnerAlias := ""
	loserAlias := ""
	scoreChanged := 0

	if rm.Settings.RatedGame && len(rm.Players) == 2 {
		for _, p := range rm.Players {
			if p.PieceType == winner {
				winnerAlias = p.Alias
			} else {
				loserAlias = p.Alias
			}
		}

		if winnerAlias != "" && loserAlias != "" {
			// 计算奖励分(击败高分玩家额外加分)
			winnerScore := s.leaderboardMgr.GetPlayerScore(winnerAlias)
			loserScore := s.leaderboardMgr.GetPlayerScore(loserAlias)
			bonus := 0
			if loserScore > winnerScore {
				bonus = (loserScore - winnerScore) / 10
			}
			scoreChanged = 25 + bonus // 基础25分 + 奖励分

			s.leaderboardMgr.RecordGameResult(winnerAlias, loserAlias, false, bonus)

			// 更新玩家对象中的分数
			for _, p := range rm.Players {
				if p.Alias == winnerAlias {
					p.Score = s.leaderboardMgr.GetPlayerScore(winnerAlias)
				} else if p.Alias == loserAlias {
					p.Score = s.leaderboardMgr.GetPlayerScore(loserAlias)
				}
			}
		}
	}

	// 先广播房间更新（更新状态为finished，包含更新后的分数）
	s.broadcastRoomUpdate(rm)

	// 广播游戏结束
	s.broadcastToRoom(rm, message.NewMessage(message.MsgGameOver, message.GameOverPayload{
		RoomID:       rm.ID,
		Winner:       winner,
		Reason:       reason,
		WinningLine:  winningLine,
		WinnerAlias:  winnerAlias,
		LoserAlias:   loserAlias,
		ScoreChanged: scoreChanged,
	}))

	// 广播房间列表（状态变为finished）
	s.broadcastRoomListToAll()

	// 广播排行榜更新
	s.broadcastLeaderboardToAll()
}

// handleDisconnect 处理断线
func (s *Server) handleDisconnect(playerID string) {
	// 移除在线用户
	s.removeOnlineUser(playerID)

	// 先检查是否是观战者
	spectatorRoom := s.findSpectatorRoom(playerID)
	if spectatorRoom != nil {
		// 从房间移除观战者
		s.roomManager.RemoveSpectatorFromRoom(spectatorRoom.ID, playerID)
		fmt.Printf("[handleDisconnect] Removed spectator %s from room %s\n", playerID, spectatorRoom.ID)

		// 广播房间更新
		if updatedRoom := s.roomManager.GetRoom(spectatorRoom.ID); updatedRoom != nil {
			s.broadcastRoomUpdate(updatedRoom)
		}

		// 广播房间列表给所有客户端
		s.broadcastRoomListToAll()
		return
	}

	// 查找玩家所在房间
	rm := s.findPlayerRoom(playerID)
	if rm == nil {
		return
	}

	roomID := rm.ID

	// 如果游戏进行中断线，先处理判负
	if rm.Status == room.StatusPlaying {
		gameOver, winner := rm.HandleDisconnect(playerID)
		if gameOver {
			s.handleGameOver(rm, winner, "disconnect", nil)
		}
	}

	// 从房间移除玩家
	_, removedPlayer := s.roomManager.RemovePlayerFromRoom(roomID, playerID)

	// 注销别名
	if removedPlayer != nil {
		s.aliasManager.Unregister(removedPlayer.Alias)
	}

	// 广播房间更新（如果房间还存在）
	if updatedRoom := s.roomManager.GetRoom(roomID); updatedRoom != nil {
		s.broadcastRoomUpdate(updatedRoom)
	}

	// 广播房间列表给所有客户端
	s.broadcastRoomListToAll()
}

// heartbeatChecker 心跳检测
func (s *Server) heartbeatChecker(playerID string, conn *websocket.Conn) {
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()

	lastPong := time.Now()

	conn.SetPongHandler(func(string) error {
		lastPong = time.Now()
		return nil
	})

	for {
		select {
		case <-ticker.C:
			if time.Since(lastPong) > s.heartbeatTimeout {
				conn.Close()
				return
			}
			conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

// 辅助方法
func (s *Server) sendMessage(conn *websocket.Conn, msg message.Message) {
	// 获取连接的写入锁
	s.mu.RLock()
	connMutex, ok := s.connMutexs[conn]
	s.mu.RUnlock()

	if !ok {
		return // 连接已关闭
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[sendMessage] Marshal error: %v", err)
		return
	}

	// 使用互斥锁保护写入
	connMutex.Lock()
	defer connMutex.Unlock()

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[sendMessage] Write error: %v", err)
	}
}

func (s *Server) sendError(conn *websocket.Conn, code int, errMsg string) {
	s.sendMessage(conn, message.NewMessage(message.MsgError, message.ErrorPayload{
		Code:    code,
		Message: errMsg,
	}))
}

// broadcastRoomListToAll 广播房间列表给所有连接的客户端
func (s *Server) broadcastRoomListToAll() {
	rooms := s.roomManager.GetAllRooms()
	msg := message.NewMessage(message.MsgRoomList, message.RoomListPayload{
		Rooms: rooms,
	})
	data, _ := json.Marshal(msg)

	s.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(s.clients))
	for _, conn := range s.clients {
		conns = append(conns, conn)
	}
	connMutexs := make(map[*websocket.Conn]*sync.Mutex)
	for conn, mutex := range s.connMutexs {
		connMutexs[conn] = mutex
	}
	s.mu.RUnlock()

	for _, conn := range conns {
		if mutex, ok := connMutexs[conn]; ok {
			mutex.Lock()
			conn.WriteMessage(websocket.TextMessage, data)
			mutex.Unlock()
		}
	}
}

func (s *Server) broadcastToRoom(rm *room.Room, msg message.Message) {
	data, _ := json.Marshal(msg)
	fmt.Printf("[broadcastToRoom] Broadcasting to room %s, msg type: %s\n", rm.ID, msg.Type)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 发送给玩家
	for _, p := range rm.Players {
		conn, ok := s.clients[p.ID]
		connMutex, mutexOk := s.connMutexs[conn]
		fmt.Printf("[broadcastToRoom] Player %s (alias=%s): conn exists=%v, IsConnected=%v\n",
			p.ID, p.Alias, ok, p.IsConnected)
		if ok && p.IsConnected && mutexOk {
			connMutex.Lock()
			err := conn.WriteMessage(websocket.TextMessage, data)
			connMutex.Unlock()
			if err != nil {
				fmt.Printf("[broadcastToRoom] Error sending to player %s: %v\n", p.ID, err)
			}
		}
	}

	// 发送给观战者
	for _, spec := range rm.Spectators {
		if conn, ok := s.clients[spec.ID]; ok {
			if connMutex, mutexOk := s.connMutexs[conn]; mutexOk {
				connMutex.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				connMutex.Unlock()
			}
		}
	}
}

func (s *Server) broadcastRoomUpdate(rm *room.Room) {
	msg := message.NewMessage(message.MsgRoomUpdate, rm.GetFullInfo())
	s.broadcastToRoom(rm, msg)
}

// broadcastLeaderboardToAll 广播排行榜给所有客户端
func (s *Server) broadcastLeaderboardToAll() {
	entries := s.leaderboardMgr.GetLeaderboard(100)
	msg := message.NewMessage(message.MsgLeaderboardUpdate, message.LeaderboardUpdatePayload{
		Entries: entries,
	})
	data, _ := json.Marshal(msg)

	s.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(s.clients))
	for _, conn := range s.clients {
		conns = append(conns, conn)
	}
	connMutexs := make(map[*websocket.Conn]*sync.Mutex)
	for conn, mutex := range s.connMutexs {
		connMutexs[conn] = mutex
	}
	s.mu.RUnlock()

	for _, conn := range conns {
		if mutex, ok := connMutexs[conn]; ok {
			mutex.Lock()
			conn.WriteMessage(websocket.TextMessage, data)
			mutex.Unlock()
		}
	}
}

// broadcastTimeUpdates 定期广播时间更新
func (s *Server) broadcastTimeUpdates(rm *room.Room) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 检查游戏是否还在进行中
			if rm.Status != room.StatusPlaying {
				// 检查是否是超时导致的游戏结束
				if rm.Status == room.StatusFinished && rm.WinReason == "timeout" {
					fmt.Printf("[broadcastTimeUpdates] Game ended by timeout in room %s, winner=%d\n", rm.ID, rm.Winner)
					s.handleGameOver(rm, rm.Winner, "timeout", nil)
				}
				fmt.Printf("[broadcastTimeUpdates] Game no longer playing, stopping time updates for room %s\n", rm.ID)
				return
			}

			// 广播时间更新
			timeUpdate := rm.GetTimeUpdate()
			s.broadcastToRoom(rm, message.NewMessage(message.MsgTimeUpdate, timeUpdate))
		}
	}
}

func (s *Server) findPlayerRoom(playerID string) *room.Room {
	// 简单遍历查找
	for _, rm := range s.roomManager.GetAllRooms() {
		room := s.roomManager.GetRoom(rm.ID)
		if room == nil {
			continue
		}
		for _, p := range room.Players {
			if p.ID == playerID {
				return room
			}
		}
	}
	return nil
}

// findSpectatorRoom 查找观战者所在房间
func (s *Server) findSpectatorRoom(spectatorID string) *room.Room {
	for _, rm := range s.roomManager.GetAllRooms() {
		room := s.roomManager.GetRoom(rm.ID)
		if room == nil {
			continue
		}
		for _, spec := range room.Spectators {
			if spec.ID == spectatorID {
				return room
			}
		}
	}
	return nil
}

func (s *Server) parseRoomSettings(settings interface{}) message.RoomSettings {
	result := message.RoomSettings{
		TimeLimit:     600,
		MoveTimeLimit: 30,
		FirstMove:     "random",
		RatedGame:     true,
	}

	if settings == nil {
		return result
	}

	settingsMap, ok := settings.(map[string]interface{})
	if !ok {
		return result
	}

	if v, ok := settingsMap["timeLimit"].(float64); ok {
		result.TimeLimit = int(v)
	}
	if v, ok := settingsMap["moveTimeLimit"].(float64); ok {
		result.MoveTimeLimit = int(v)
	}
	if v, ok := settingsMap["firstMove"].(string); ok {
		result.FirstMove = v
	}
	if v, ok := settingsMap["ratedGame"].(bool); ok {
		result.RatedGame = v
	}

	return result
}

func generatePlayerID() string {
	return fmt.Sprintf("player_%d", time.Now().UnixNano())
}

// ========== 在线用户管理 ==========

// addOnlineUser 添加在线用户
func (s *Server) addOnlineUser(playerID, alias, status, roomID string) {
	s.onlineUsersMu.Lock()
	defer s.onlineUsersMu.Unlock()

	s.onlineUsers[playerID] = &OnlineUser{
		ID:     playerID,
		Alias:  alias,
		Status: status,
		RoomID: roomID,
	}

	// 广播在线用户更新
	go s.broadcastOnlineUsersToAll()
}

// updateOnlineUserStatus 更新在线用户状态
func (s *Server) updateOnlineUserStatus(playerID, status, roomID string) {
	s.onlineUsersMu.Lock()
	defer s.onlineUsersMu.Unlock()

	if user, ok := s.onlineUsers[playerID]; ok {
		user.Status = status
		user.RoomID = roomID
	}

	// 广播在线用户更新
	go s.broadcastOnlineUsersToAll()
}

// removeOnlineUser 移除在线用户
func (s *Server) removeOnlineUser(playerID string) {
	s.onlineUsersMu.Lock()
	defer s.onlineUsersMu.Unlock()

	delete(s.onlineUsers, playerID)

	// 广播在线用户更新
	go s.broadcastOnlineUsersToAll()
}

// getOnlineUsersList 获取在线用户列表
func (s *Server) getOnlineUsersList() message.OnlineUsersPayload {
	s.onlineUsersMu.RLock()
	defer s.onlineUsersMu.RUnlock()

	users := make([]message.OnlineUserInfo, 0, len(s.onlineUsers))

	for _, user := range s.onlineUsers {
		if user.Alias == "" {
			continue // 跳过未设置别名的用户
		}

		info := message.OnlineUserInfo{
			ID:     user.ID,
			Alias:  user.Alias,
			Status: user.Status,
		}

		// 获取分数和段位
		info.Score = s.leaderboardMgr.GetPlayerScore(user.Alias)
		info.Rank = s.getRankByScore(info.Score)

		// 获取房间名称
		if user.RoomID != "" {
			info.RoomID = user.RoomID
			if rm := s.roomManager.GetRoom(user.RoomID); rm != nil {
				info.RoomName = rm.Name
			}
		}

		users = append(users, info)
	}

	return message.OnlineUsersPayload{
		Count: len(users),
		Users: users,
	}
}

// getRankByScore 根据分数获取段位
func (s *Server) getRankByScore(score int) string {
	switch {
	case score >= 1000:
		return "大师"
	case score >= 600:
		return "钻石"
	case score >= 300:
		return "黄金"
	case score >= 100:
		return "白银"
	default:
		return "青铜"
	}
}

// broadcastOnlineUsersToAll 广播在线用户列表给所有客户端
func (s *Server) broadcastOnlineUsersToAll() {
	// 使用 defer recover 防止 panic 导致服务器崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[broadcastOnlineUsersToAll] Recovered from panic: %v", r)
		}
	}()

	payload := s.getOnlineUsersList()
	msg := message.NewMessage(message.MsgOnlineUsers, payload)
	data, _ := json.Marshal(msg)

	s.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(s.clients))
	for _, conn := range s.clients {
		conns = append(conns, conn)
	}
	connMutexs := make(map[*websocket.Conn]*sync.Mutex)
	for conn, mutex := range s.connMutexs {
		connMutexs[conn] = mutex
	}
	s.mu.RUnlock()

	for _, conn := range conns {
		if mutex, ok := connMutexs[conn]; ok {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[broadcastOnlineUsersToAll] Write panic recovered: %v", r)
					}
				}()
				mutex.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mutex.Unlock()
			}()
		}
	}
}

// handleGetOnlineUsers 处理获取在线用户列表
func (s *Server) handleGetOnlineUsers(conn *websocket.Conn) {
	payload := s.getOnlineUsersList()
	s.sendMessage(conn, message.NewMessage(message.MsgOnlineUsers, payload))
}

// handleEnterLobby 处理进入大厅
func (s *Server) handleEnterLobby(playerID string, conn *websocket.Conn, msg message.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		s.sendError(conn, 400, "Invalid payload")
		return
	}

	alias, _ := payload["alias"].(string)
	if alias == "" {
		// 如果没有别名，使用默认格式
		alias = "Player_" + playerID[:6]
	}

	// 检查别名是否已被使用
	if !s.aliasManager.IsAvailable(alias) {
		// 别名已被使用，生成一个新的
		alias = s.aliasManager.GenerateUniqueAlias(alias)
	}

	// 注册别名
	s.aliasManager.Register(alias, playerID)

	// 添加到在线用户列表
	s.addOnlineUser(playerID, alias, "lobby", "")

	fmt.Printf("[handleEnterLobby] Player %s entered lobby with alias: %s\n", playerID, alias)

	// 发送确认消息给客户端
	s.sendMessage(conn, message.NewMessage("lobby_entered", map[string]interface{}{
		"alias":     alias,
		"playerId":  playerID,
	}))
}

// handleLeaveLobby 处理离开大厅（返回主菜单）
func (s *Server) handleLeaveLobby(playerID string) {
	// 先获取别名用于注销
	s.onlineUsersMu.RLock()
	var alias string
	if user, ok := s.onlineUsers[playerID]; ok {
		alias = user.Alias
	}
	s.onlineUsersMu.RUnlock()

	// 注销别名
	if alias != "" {
		s.aliasManager.Unregister(alias)
	}

	// 移除在线用户
	s.removeOnlineUser(playerID)

	fmt.Printf("[handleLeaveLobby] Player %s left lobby\n", playerID)
}
