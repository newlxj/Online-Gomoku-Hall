package room

import (
	"fmt"
	"gomoku-server/internal/message"
	"gomoku-server/internal/player"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager 房间管理器
type Manager struct {
	rooms    map[string]*Room
	mu       sync.RWMutex
}

// NewManager 创建房间管理器
func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// CreateRoom 创建房间
func (m *Manager) CreateRoom(name string, settings message.RoomSettings, hostID string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	roomID := generateRoomID()
	room := NewRoom(roomID, name, hostID, settings)
	m.rooms[roomID] = room

	return room
}

// GetRoom 获取房间
func (m *Manager) GetRoom(roomID string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.rooms[roomID]
}

// DeleteRoom 删除房间
func (m *Manager) DeleteRoom(roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, roomID)
}

// GetWaitingRooms 获取等待中的房间列表
func (m *Manager) GetWaitingRooms() []message.RoomInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rooms := make([]message.RoomInfo, 0)
	for _, room := range m.rooms {
		if room.IsWaiting() {
			rooms = append(rooms, room.GetInfo())
		}
	}
	return rooms
}

// GetAllRooms 获取所有房间列表
func (m *Manager) GetAllRooms() []message.RoomInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rooms := make([]message.RoomInfo, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, room.GetInfo())
	}
	return rooms
}

// GetRandomWaitingRoom 随机获取一个等待中的房间
func (m *Manager) GetRandomWaitingRoom() *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, room := range m.rooms {
		if room.IsWaiting() {
			return room
		}
	}
	return nil
}

// AddPlayerToRoom 添加玩家到房间
func (m *Manager) AddPlayerToRoom(roomID string, p *player.Player) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}

	if !room.AddPlayer(p) {
		return nil, fmt.Errorf("room is full")
	}

	return room, nil
}

// RemovePlayerFromRoom 从房间移除玩家
func (m *Manager) RemovePlayerFromRoom(roomID, playerID string) (*Room, *player.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, nil
	}

	removedPlayer := room.RemovePlayer(playerID)

	// 如果房间空了，删除房间
	if len(room.Players) == 0 {
		delete(m.rooms, roomID)
	}

	return room, removedPlayer
}

// AddSpectatorToRoom 添加观战者到房间
func (m *Manager) AddSpectatorToRoom(roomID string, s *player.Spectator) (*Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}

	room.AddSpectator(s)
	return room, nil
}

// RemoveSpectatorFromRoom 从房间移除观战者
func (m *Manager) RemoveSpectatorFromRoom(roomID, spectatorID string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return
	}

	room.RemoveSpectator(spectatorID)
}

// GetRoomCount 获取房间数量
func (m *Manager) GetRoomCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.rooms)
}

// GetPlayingRoomCount 获取进行中的房间数量
func (m *Manager) GetPlayingRoomCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, room := range m.rooms {
		if room.Status == StatusPlaying {
			count++
		}
	}
	return count
}

// CleanupEmptyRooms 清理空房间
func (m *Manager) CleanupEmptyRooms() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for id, room := range m.rooms {
		if len(room.Players) == 0 {
			delete(m.rooms, id)
			count++
		}
	}
	return count
}

// CleanupOldFinishedRooms 清理已结束超过一定时间的房间
func (m *Manager) CleanupOldFinishedRooms(maxAge time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	now := time.Now()
	for id, room := range m.rooms {
		if room.Status == StatusFinished && now.Sub(room.FinishedAt) > maxAge {
			delete(m.rooms, id)
			count++
		}
	}
	return count
}

// generateRoomID 生成房间ID
func generateRoomID() string {
	id := uuid.New()
	return id.String()[:8] // 使用前8位作为房间ID
}
