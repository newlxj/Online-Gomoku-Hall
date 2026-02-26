package player

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Player 玩家实体
type Player struct {
	ID            string          `json:"id"`
	Alias         string          `json:"alias"`
	PieceType     int             `json:"pieceType"` // 1=黑, 2=白
	IsReady       bool            `json:"isReady"`
	RemainingTime int             `json:"remainingTime"` // 剩余总时间(秒)
	MoveTimeLeft  int             `json:"moveTimeLeft"`  // 当前步剩余时间(秒)
	IsConnected   bool            `json:"isConnected"`
	Conn          *websocket.Conn `json:"-"`
	LastActive    time.Time       `json:"-"`
	RoomID        string          `json:"roomId"`
	Score         int             `json:"score"`
	mu            sync.RWMutex
}

// Spectator 观战者实体
type Spectator struct {
	ID        string          `json:"id"`
	Alias     string          `json:"alias"`
	JoinedAt  time.Time       `json:"joinedAt"`
	Conn      *websocket.Conn `json:"-"`
	RoomID    string          `json:"roomId"`
}

// NewPlayer 创建新玩家
func NewPlayer(id, alias string, conn *websocket.Conn) *Player {
	return &Player{
		ID:          id,
		Alias:       alias,
		PieceType:   0, // 未分配
		IsReady:     false,
		RemainingTime: 600, // 默认10分钟
		MoveTimeLeft:  30,  // 默认30秒
		IsConnected:  true,
		Conn:         conn,
		LastActive:   time.Now(),
		RoomID:       "",
		Score:        0,
	}
}

// NewSpectator 创建新观战者
func NewSpectator(id, alias string, conn *websocket.Conn) *Spectator {
	return &Spectator{
		ID:       id,
		Alias:    alias,
		JoinedAt: time.Now(),
		Conn:     conn,
		RoomID:   "",
	}
}

// SetPieceType 设置棋子类型
func (p *Player) SetPieceType(pieceType int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.PieceType = pieceType
}

// SetReady 设置准备状态
func (p *Player) SetReady(ready bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsReady = ready
}

// UpdateTime 更新时间
func (p *Player) UpdateTime(totalTime, moveTime int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.RemainingTime = totalTime
	p.MoveTimeLeft = moveTime
}

// DecrementMoveTime 减少步时
func (p *Player) DecrementMoveTime() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.MoveTimeLeft > 0 {
		p.MoveTimeLeft--
	}
	return p.MoveTimeLeft
}

// DecrementTotalTime 减少总时间
func (p *Player) DecrementTotalTime() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.RemainingTime > 0 {
		p.RemainingTime--
	}
	return p.RemainingTime
}

// ResetMoveTime 重置步时
func (p *Player) ResetMoveTime(moveTimeLimit int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.MoveTimeLeft = moveTimeLimit
}

// SetConnected 设置连接状态
func (p *Player) SetConnected(connected bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsConnected = connected
}

// UpdateLastActive 更新最后活跃时间
func (p *Player) UpdateLastActive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LastActive = time.Now()
}

// GetInfo 获取玩家信息(用于消息)
func (p *Player) GetInfo() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return map[string]interface{}{
		"id":            p.ID,
		"alias":         p.Alias,
		"pieceType":     p.PieceType,
		"isReady":       p.IsReady,
		"remainingTime": p.RemainingTime,
		"moveTimeLeft":  p.MoveTimeLeft,
		"isConnected":   p.IsConnected,
		"score":         p.Score,
	}
}
