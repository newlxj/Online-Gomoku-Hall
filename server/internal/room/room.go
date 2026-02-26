package room

import (
	"fmt"
	"gomoku-server/internal/game"
	"gomoku-server/internal/message"
	"gomoku-server/internal/player"
	"sync"
	"time"
)

// RoomStatus 房间状态
type RoomStatus string

const (
	StatusWaiting  RoomStatus = "waiting"
	StatusReady    RoomStatus = "ready"
	StatusPlaying  RoomStatus = "playing"
	StatusFinished RoomStatus = "finished"
)

// Room 房间实体
type Room struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	HostID       string                 `json:"hostId"`
	Players      []*player.Player       `json:"players"`
	Spectators   []*player.Spectator    `json:"spectators"`
	Status       RoomStatus             `json:"status"`
	Game         *game.Game             `json:"-"`
	Settings     message.RoomSettings   `json:"settings"`
	CreatedAt    time.Time              `json:"createdAt"`
	StartedAt    time.Time              `json:"startedAt,omitempty"`
	FinishedAt   time.Time              `json:"finishedAt,omitempty"`
	Winner       int                    `json:"winner"` // 0=未结束, 1=黑胜, 2=白胜, 3=平局
	WinReason    string                 `json:"winReason"`
	MoveHistory  []message.MoveRecord   `json:"moveHistory"`
	mu           sync.RWMutex

	// 计时器相关
	timerStopChan chan struct{} `json:"-"`
}

// NewRoom 创建新房间
func NewRoom(id, name, hostID string, settings message.RoomSettings) *Room {
	// 设置默认值
	if settings.TimeLimit == 0 {
		settings.TimeLimit = 600 // 默认10分钟
	}
	if settings.MoveTimeLimit == 0 {
		settings.MoveTimeLimit = 30 // 默认30秒
	}
	if settings.FirstMove == "" {
		settings.FirstMove = "random"
	}

	return &Room{
		ID:          id,
		Name:        name,
		HostID:      hostID,
		Players:     make([]*player.Player, 0, 2),
		Spectators:  make([]*player.Spectator, 0),
		Status:      StatusWaiting,
		Game:        nil,
		Settings:    settings,
		CreatedAt:   time.Now(),
		MoveHistory: make([]message.MoveRecord, 0),
	}
}

// AddPlayer 添加玩家
func (r *Room) AddPlayer(p *player.Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= 2 {
		return false
	}

	p.RoomID = r.ID
	p.RemainingTime = r.Settings.TimeLimit
	p.MoveTimeLeft = r.Settings.MoveTimeLimit
	p.IsConnected = true // 确保连接状态为true
	r.Players = append(r.Players, p)

	// 如果是房主，设置为主机
	if p.ID == r.HostID && len(r.Players) == 1 {
		p.PieceType = 1 // 房主默认黑棋
	}

	// 当有2名玩家时，更新状态为ready
	if len(r.Players) == 2 {
		r.Status = StatusReady
		fmt.Printf("[AddPlayer] Room now has 2 players, status changed to ready\n")
	}

	return true
}

// RemovePlayer 移除玩家
func (r *Room) RemovePlayer(playerID string) *player.Player {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.Players {
		if p.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			p.RoomID = ""
			p.PieceType = 0
			p.IsReady = false

			// 更新房间状态
			if len(r.Players) < 2 {
				r.Status = StatusWaiting
			}

			return p
		}
	}
	return nil
}

// AddSpectator 添加观战者
func (r *Room) AddSpectator(s *player.Spectator) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s.RoomID = r.ID
	r.Spectators = append(r.Spectators, s)
}

// RemoveSpectator 移除观战者
func (r *Room) RemoveSpectator(spectatorID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.Spectators {
		if s.ID == spectatorID {
			r.Spectators = append(r.Spectators[:i], r.Spectators[i+1:]...)
			s.RoomID = ""
			return
		}
	}
}

// SetPlayerReady 设置玩家准备状态
func (r *Room) SetPlayerReady(playerID string, ready bool) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.Players {
		if p.ID == playerID {
			// 先更新房间状态（如果在 finished 状态需要重置）
			r.updateStatusAfterReadyLocked(playerID)
			// 然后设置当前玩家的准备状态
			p.SetReady(ready)
			return true
		}
	}
	return false
}

// updateStatusAfterReadyLocked 在玩家准备后更新房间状态（需要持有锁）
// currentPlayerID 是当前点击准备的玩家ID，用于在重置后恢复其准备状态
func (r *Room) updateStatusAfterReadyLocked(currentPlayerID string) {
	// 如果游戏结束了，允许重新准备开始新游戏
	if r.Status == StatusFinished {
		// 重置游戏状态
		r.Game = nil
		r.Winner = 0
		r.WinReason = ""
		r.MoveHistory = make([]message.MoveRecord, 0)
		// 重置所有玩家的准备状态（要求双方重新准备）
		for _, p := range r.Players {
			p.IsReady = false
		}
		// 检查玩家数量后设置状态
		if len(r.Players) < 2 {
			r.Status = StatusWaiting
			return
		}
		r.Status = StatusReady
		return
	}

	// 只有在waiting或ready状态下才更新
	if r.Status != StatusWaiting && r.Status != StatusReady {
		return
	}

	// 如果只有一人或无人，状态为waiting
	if len(r.Players) < 2 {
		r.Status = StatusWaiting
		return
	}

	// 有2人时，状态变为ready（等待双方准备）
	r.Status = StatusReady
}

// CanStart 检查是否可以开始游戏
func (r *Room) CanStart() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.Players) != 2 {
		return false
	}

	for _, p := range r.Players {
		if !p.IsReady {
			return false
		}
	}

	return true
}

// StartGame 开始游戏
func (r *Room) StartGame() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 内联检查，避免死锁（不能在持有写锁时调用CanStart因为它会尝试获取读锁）
	if len(r.Players) != 2 {
		return false
	}
	for _, p := range r.Players {
		if !p.IsReady {
			return false
		}
	}

	fmt.Printf("[StartGame] Starting game with %d players\n", len(r.Players))

	// 根据先手规则分配棋子颜色
	switch r.Settings.FirstMove {
	case "host":
		r.Players[0].PieceType = 1 // 房主黑棋
		r.Players[1].PieceType = 2 // 访客白棋
	case "guest":
		r.Players[0].PieceType = 2 // 房主白棋
		r.Players[1].PieceType = 1 // 访客黑棋
	default: // random
		if time.Now().UnixNano()%2 == 0 {
			r.Players[0].PieceType = 1
			r.Players[1].PieceType = 2
		} else {
			r.Players[0].PieceType = 2
			r.Players[1].PieceType = 1
		}
	}

	// 初始化游戏
	r.Game = game.NewGame(15)
	r.Status = StatusPlaying
	r.StartedAt = time.Now()
	r.MoveHistory = make([]message.MoveRecord, 0)

	// 重置时间
	for _, p := range r.Players {
		p.RemainingTime = r.Settings.TimeLimit
		p.MoveTimeLeft = r.Settings.MoveTimeLimit
	}

	// 启动计时器
	r.timerStopChan = make(chan struct{})
	go r.gameTimer()

	return true
}

// gameTimer 游戏计时器
func (r *Room) gameTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.decrementTime()
		case <-r.timerStopChan:
			return
		}
	}
}

// decrementTime 减少时间
// 规则：当自己下棋时，自己的总时间和步时同时倒计时
// 步时结束或总时间结束都自动判输
func (r *Room) decrementTime() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != StatusPlaying || r.Game == nil {
		return
	}

	currentTurn := r.Game.GetCurrentTurn()

	// 减少当前玩家的时间
	for _, p := range r.Players {
		if p.PieceType == currentTurn {
			// 同时减少总时间和步时
			p.RemainingTime--
			p.MoveTimeLeft--

			// 检查步时超时（优先判断步时）
			if p.MoveTimeLeft <= 0 {
				r.Status = StatusFinished
				r.FinishedAt = time.Now()
				// 设置对手获胜
				for _, op := range r.Players {
					if op.ID != p.ID {
						r.Winner = op.PieceType
					}
				}
				r.WinReason = "timeout"
				return
			}

			// 检查总时间超时
			if p.RemainingTime <= 0 {
				r.Status = StatusFinished
				r.FinishedAt = time.Now()
				// 设置对手获胜
				for _, op := range r.Players {
					if op.ID != p.ID {
						r.Winner = op.PieceType
					}
				}
				r.WinReason = "timeout"
			}
			break
		}
	}
}

// StopTimer 停止计时器
func (r *Room) StopTimer() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.timerStopChan != nil {
		close(r.timerStopChan)
		r.timerStopChan = nil
	}
}

// GetTimeUpdate 获取时间更新消息的payload
func (r *Room) GetTimeUpdate() message.TimeUpdatePayload {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]message.PlayerTimeInfo, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, message.PlayerTimeInfo{
			PlayerID:      p.ID,
			RemainingTime: p.RemainingTime,
			MoveTimeLeft:  p.MoveTimeLeft,
		})
	}

	return message.TimeUpdatePayload{
		RoomID:  r.ID,
		Players: players,
	}
}

// MakeMove 落子
func (r *Room) MakeMove(playerID string, pos message.Position) (bool, int, []message.Position) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != StatusPlaying || r.Game == nil {
		return false, 0, nil
	}

	// 找到玩家
	var currentPlayer *player.Player
	for _, p := range r.Players {
		if p.ID == playerID {
			currentPlayer = p
			break
		}
	}

	if currentPlayer == nil {
		return false, 0, nil
	}

	// 检查是否轮到该玩家
	if r.Game.GetCurrentTurn() != currentPlayer.PieceType {
		return false, 0, nil
	}

	// 落子
	success, winningLine := r.Game.MakeMove(pos.Row, pos.Col, currentPlayer.PieceType)
	if !success {
		return false, 0, nil
	}

	// 记录落子
	r.MoveHistory = append(r.MoveHistory, message.MoveRecord{
		Position:  pos,
		PieceType: currentPlayer.PieceType,
		Timestamp: time.Now().UnixMilli(),
		PlayerID:  playerID,
	})

	// 检查胜利
	if len(winningLine) > 0 {
		r.Status = StatusFinished
		r.FinishedAt = time.Now()
		r.Winner = currentPlayer.PieceType
		r.WinReason = "win"
		return true, currentPlayer.PieceType, winningLine
	}

	// 检查平局
	if r.Game.IsBoardFull() {
		r.Status = StatusFinished
		r.FinishedAt = time.Now()
		r.Winner = 0
		r.WinReason = "draw"
		return true, 0, nil
	}

	// 重置下一个玩家的步时（回合已经切换）
	nextTurn := r.Game.GetCurrentTurn()
	for _, p := range r.Players {
		if p.PieceType == nextTurn {
			p.MoveTimeLeft = r.Settings.MoveTimeLimit
			break
		}
	}

	return true, 0, nil
}

// GetCurrentTurnPlayer 获取当前回合玩家
func (r *Room) GetCurrentTurnPlayer() *player.Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.Game == nil {
		return nil
	}

	currentTurn := r.Game.GetCurrentTurn()
	for _, p := range r.Players {
		if p.PieceType == currentTurn {
			return p
		}
	}
	return nil
}

// GetPlayerByAlias 根据别名获取玩家
func (r *Room) GetPlayerByAlias(alias string) *player.Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.Players {
		if p.Alias == alias {
			return p
		}
	}
	return nil
}

// GetInfo 获取房间信息(用于列表)
func (r *Room) GetInfo() message.RoomInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hostName := ""
	for _, p := range r.Players {
		if p.ID == r.HostID {
			hostName = p.Alias
			break
		}
	}

	return message.RoomInfo{
		ID:             r.ID,
		Name:           r.Name,
		HostID:         r.HostID,
		HostName:       hostName,
		PlayerCount:    len(r.Players),
		SpectatorCount: len(r.Spectators),
		Status:         string(r.Status),
		CreatedAt:      r.CreatedAt.UnixMilli(),
		Settings:       r.Settings,
	}
}

// GetFullInfo 获取完整房间信息(用于更新)
func (r *Room) GetFullInfo() message.RoomUpdatePayload {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]message.PlayerInfo, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, message.PlayerInfo{
			ID:            p.ID,
			Alias:         p.Alias,
			PieceType:     p.PieceType,
			IsReady:       p.IsReady,
			RemainingTime: p.RemainingTime,
			MoveTimeLeft:  p.MoveTimeLeft,
			IsConnected:   p.IsConnected,
			Score:         p.Score,
		})
	}

	spectators := make([]message.SpectatorInfo, 0, len(r.Spectators))
	for _, s := range r.Spectators {
		spectators = append(spectators, message.SpectatorInfo{
			ID:       s.ID,
			Alias:    s.Alias,
			JoinedAt: s.JoinedAt.UnixMilli(),
		})
	}

	var board [][]int
	if r.Game != nil {
		board = r.Game.GetBoard()
	} else {
		board = make([][]int, 15)
		for i := range board {
			board[i] = make([]int, 15)
		}
	}

	currentTurn := 0
	if r.Game != nil {
		currentTurn = r.Game.GetCurrentTurn()
	}

	return message.RoomUpdatePayload{
		ID:          r.ID,
		Name:        r.Name,
		HostID:      r.HostID,
		Players:     players,
		Spectators:  spectators,
		Status:      string(r.Status),
		CurrentTurn: currentTurn,
		Board:       board,
		Settings:    r.Settings,
		MoveHistory: r.MoveHistory,
	}
}

// IsWaiting 是否等待中
func (r *Room) IsWaiting() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Status == StatusWaiting && len(r.Players) < 2
}

// HandleTimeout 处理超时
func (r *Room) HandleTimeout(playerID string) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != StatusPlaying {
		return false, 0
	}

	// 找到超时玩家
	var timeoutPlayer *player.Player
	var winner int
	for _, p := range r.Players {
		if p.ID == playerID {
			timeoutPlayer = p
		} else {
			winner = p.PieceType
		}
	}

	if timeoutPlayer == nil {
		return false, 0
	}

	r.Status = StatusFinished
	r.FinishedAt = time.Now()
	r.Winner = winner
	r.WinReason = "timeout"

	return true, winner
}

// HandleDisconnect 处理断线
func (r *Room) HandleDisconnect(playerID string) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 标记玩家断线
	for _, p := range r.Players {
		if p.ID == playerID {
			p.IsConnected = false
			break
		}
	}

	// 如果游戏中，判对方获胜
	if r.Status == StatusPlaying {
		var winner int
		for _, p := range r.Players {
			if p.ID != playerID {
				winner = p.PieceType
			}
		}

		r.Status = StatusFinished
		r.FinishedAt = time.Now()
		r.Winner = winner
		r.WinReason = "disconnect"
		return true, winner
	}

	return false, 0
}

// SetGameFinished 设置游戏结束状态（用于中途离开等情况）
func (r *Room) SetGameFinished(winner int, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Status = StatusFinished
	r.FinishedAt = time.Now()
	r.Winner = winner
	r.WinReason = reason
}

// UndoMove 撤销最后一步棋（需要请求者的棋子颜色）
func (r *Room) UndoMove(requesterPieceType int) (bool, message.MoveRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != StatusPlaying || r.Game == nil {
		return false, message.MoveRecord{}
	}

	// 检查是否有棋可悔
	if len(r.MoveHistory) == 0 {
		return false, message.MoveRecord{}
	}

	// 获取最后一步棋
	lastMove := r.MoveHistory[len(r.MoveHistory)-1]

	// 规则：只能悔对手刚下的那一步
	// 最后一步必须是对手下的棋（棋子颜色与请求者不同）
	if lastMove.PieceType == requesterPieceType {
		// 不能悔自己下的棋
		return false, message.MoveRecord{}
	}

	// 从历史中移除
	r.MoveHistory = r.MoveHistory[:len(r.MoveHistory)-1]

	// 撤销棋盘上的棋子
	r.Game.UndoMove(lastMove.Position.Row, lastMove.Position.Col)

	// 回合切换回去
	r.Game.SwitchTurn()

	// 重置请求者的步时（因为现在轮到请求者了）
	for _, p := range r.Players {
		if p.PieceType == requesterPieceType {
			p.MoveTimeLeft = r.Settings.MoveTimeLimit
			break
		}
	}

	return true, lastMove
}

// GetLastMove 获取最后一步棋
func (r *Room) GetLastMove() *message.MoveRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.MoveHistory) == 0 {
		return nil
	}

	lastMove := r.MoveHistory[len(r.MoveHistory)-1]
	return &lastMove
}
