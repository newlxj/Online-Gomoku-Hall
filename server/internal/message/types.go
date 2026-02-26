package message

import "time"

// MessageType 定义WebSocket消息类型
type MessageType string

const (
	// 客户端 -> 服务器
	MsgCreateRoom   MessageType = "create_room"
	MsgJoinRoom     MessageType = "join_room"
	MsgLeaveRoom    MessageType = "leave_room"
	MsgSpectate     MessageType = "spectate"
	MsgReady        MessageType = "ready"
	MsgMove         MessageType = "move"
	MsgEmoji        MessageType = "emoji"
	MsgCheckAlias   MessageType = "check_alias"
	MsgHeartbeat    MessageType = "heartbeat"
	MsgGetRoomList  MessageType = "get_room_list"
	MsgGetLeaderboard MessageType = "get_leaderboard"
	MsgGetOnlineUsers  MessageType = "get_online_users"
	MsgEnterLobby   MessageType = "enter_lobby"
	MsgLeaveLobby   MessageType = "leave_lobby"
	MsgSurrenderRequest  MessageType = "surrender_request"
	MsgSurrenderResponse MessageType = "surrender_response"
	MsgUndoRequest       MessageType = "undo_request"
	MsgUndoResponse      MessageType = "undo_response"

	// 服务器 -> 客户端
	MsgRoomList           MessageType = "room_list"
	MsgRoomUpdate         MessageType = "room_update"
	MsgGameStart          MessageType = "game_start"
	MsgMoveSync           MessageType = "move"
	MsgGameOver           MessageType = "game_over"
	MsgTimeUpdate         MessageType = "time_update"
	MsgPlayerDisconnected MessageType = "player_disconnected"
	MsgAliasCheckResult   MessageType = "alias_check_result"
	MsgEmojiBroadcast     MessageType = "emoji"
	MsgAchievementUnlock  MessageType = "achievement_unlock"
	MsgLeaderboardUpdate  MessageType = "leaderboard_update"
	MsgOnlineUsers        MessageType = "online_users"
	MsgError              MessageType = "error"
)

// Message WebSocket消息结构
type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, payload interface{}) Message {
	return Message{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now().UnixMilli(),
	}
}

// CreateRoomPayload 创建房间消息
type CreateRoomPayload struct {
	Name     string       `json:"name"`
	Settings RoomSettings `json:"settings"`
	Alias    string       `json:"alias"`
}

// JoinRoomPayload 加入房间消息
type JoinRoomPayload struct {
	RoomID string `json:"roomId"`
	Alias  string `json:"alias"`
}

// SpectatePayload 观战消息
type SpectatePayload struct {
	RoomID string `json:"roomId"`
	Alias  string `json:"alias"`
}

// ReadyPayload 准备消息
type ReadyPayload struct {
	RoomID string `json:"roomId"`
}

// MovePayload 落子消息
type MovePayload struct {
	RoomID   string `json:"roomId"`
	Position Position `json:"position"`
}

// EmojiPayload 表情消息
type EmojiPayload struct {
	RoomID string `json:"roomId"`
	Emoji  string `json:"emoji"`
}

// CheckAliasPayload 检查别名消息
type CheckAliasPayload struct {
	Alias string `json:"alias"`
}

// Position 位置
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// RoomSettings 房间设置
type RoomSettings struct {
	TimeLimit     int    `json:"timeLimit"`     // 每人总时间(秒), 默认600秒
	MoveTimeLimit int    `json:"moveTimeLimit"` // 每步限时(秒), 默认30秒
	FirstMove     string `json:"firstMove"`     // host/guest/random
	RatedGame     bool   `json:"ratedGame"`     // 是否计入排行榜
}

// RoomInfo 房间信息(用于列表显示)
type RoomInfo struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	HostID       string        `json:"hostId"`
	HostName     string        `json:"hostName"`
	PlayerCount  int           `json:"playerCount"`
	SpectatorCount int         `json:"spectatorCount"`
	Status       string        `json:"status"`
	CreatedAt    int64         `json:"createdAt"`
	Settings     RoomSettings  `json:"settings"`
}

// RoomListPayload 房间列表
type RoomListPayload struct {
	Rooms []RoomInfo `json:"rooms"`
}

// RoomUpdatePayload 房间更新
type RoomUpdatePayload struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	HostID       string              `json:"hostId"`
	Players      []PlayerInfo        `json:"players"`
	Spectators   []SpectatorInfo     `json:"spectators"`
	Status       string              `json:"status"`
	CurrentTurn  int                 `json:"currentTurn"`
	Board        [][]int             `json:"board"`
	Settings     RoomSettings        `json:"settings"`
	MoveHistory  []MoveRecord        `json:"moveHistory"`
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	ID            string `json:"id"`
	Alias         string `json:"alias"`
	PieceType     int    `json:"pieceType"`
	IsReady       bool   `json:"isReady"`
	RemainingTime int    `json:"remainingTime"`
	MoveTimeLeft  int    `json:"moveTimeLeft"`
	IsConnected   bool   `json:"isConnected"`
	Score         int    `json:"score"`
	Rank          string `json:"rank"`
}

// SpectatorInfo 观战者信息
type SpectatorInfo struct {
	ID        string `json:"id"`
	Alias     string `json:"alias"`
	JoinedAt  int64  `json:"joinedAt"`
}

// MoveRecord 落子记录
type MoveRecord struct {
	Position  Position `json:"position"`
	PieceType int      `json:"pieceType"`
	Timestamp int64    `json:"timestamp"`
	PlayerID  string   `json:"playerId"`
}

// MoveSyncPayload 落子同步
type MoveSyncPayload struct {
	RoomID   string   `json:"roomId"`
	Position Position `json:"position"`
	PlayerID string   `json:"playerId"`
	PieceType int     `json:"pieceType"`
}

// GameOverPayload 游戏结束
type GameOverPayload struct {
	RoomID       string     `json:"roomId"`
	Winner       int        `json:"winner"` // 0=平局, 1=黑, 2=白
	Reason       string     `json:"reason"` // win/disconnect/timeout/leave
	WinningLine  []Position `json:"winningLine,omitempty"`
	WinnerAlias  string     `json:"winnerAlias,omitempty"`
	LoserAlias   string     `json:"loserAlias,omitempty"`
	ScoreChanged int        `json:"scoreChanged,omitempty"` // 获胜者得分变化
}

// TimeUpdatePayload 时间更新
type TimeUpdatePayload struct {
	RoomID  string               `json:"roomId"`
	Players []PlayerTimeInfo     `json:"players"`
}

// PlayerTimeInfo 玩家时间信息
type PlayerTimeInfo struct {
	PlayerID      string `json:"playerId"`
	RemainingTime int    `json:"remainingTime"`
	MoveTimeLeft  int    `json:"moveTimeLeft"`
}

// AliasCheckResultPayload 别名检查结果
type AliasCheckResultPayload struct {
	Alias     string `json:"alias"`
	Available bool   `json:"available"`
}

// EmojiBroadcastPayload 表情广播
type EmojiBroadcastPayload struct {
	RoomID   string `json:"roomId"`
	Emoji    string `json:"emoji"`
	FromPlayer string `json:"fromPlayer"`
	FromAlias  string `json:"fromAlias"`
}

// AchievementUnlockPayload 成就解锁
type AchievementUnlockPayload struct {
	AchievementID   string `json:"achievementId"`
	AchievementName string `json:"achievementName"`
	Description     string `json:"description"`
	Reward          int    `json:"reward"`
}

// LeaderboardEntry 排行榜条目
type LeaderboardEntry struct {
	Rank      int     `json:"rank"`
	Alias     string  `json:"alias"`
	Wins      int     `json:"wins"`
	Losses    int     `json:"losses"`
	Draws     int     `json:"draws"`
	WinRate   float64 `json:"winRate"`
	Score     int     `json:"score"`
	RankTitle string  `json:"rankTitle"`
}

// LeaderboardUpdatePayload 排行榜更新
type LeaderboardUpdatePayload struct {
	Entries []LeaderboardEntry `json:"entries"`
}

// ErrorPayload 错误消息
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// OnlineUserInfo 在线用户信息
type OnlineUserInfo struct {
	ID        string `json:"id"`
	Alias     string `json:"alias"`
	Score     int    `json:"score"`
	Rank      string `json:"rank"`
	Status    string `json:"status"` // "lobby", "playing", "spectating"
	RoomID    string `json:"roomId,omitempty"`
	RoomName  string `json:"roomName,omitempty"`
}

// OnlineUsersPayload 在线用户列表
type OnlineUsersPayload struct {
	Count   int              `json:"count"`
	Users   []OnlineUserInfo `json:"users"`
}

// SurrenderRequestPayload 认输请求
type SurrenderRequestPayload struct {
	RoomID    string `json:"roomId"`
	PlayerID  string `json:"playerId"`
	PlayerAlias string `json:"playerAlias"`
}

// SurrenderResponsePayload 认输响应
type SurrenderResponsePayload struct {
	RoomID    string `json:"roomId"`
	FromPlayer string `json:"fromPlayer"`
	Accept    bool   `json:"accept"`
}

// UndoRequestPayload 悔棋请求
type UndoRequestPayload struct {
	RoomID      string   `json:"roomId"`
	PlayerID    string   `json:"playerId"`
	PlayerAlias string   `json:"playerAlias"`
	MovePosition Position `json:"movePosition"` // 要悔的棋的位置
}

// UndoResponsePayload 悔棋响应
type UndoResponsePayload struct {
	RoomID     string `json:"roomId"`
	FromPlayer string `json:"fromPlayer"`
	Accept     bool   `json:"accept"`
}

// GameActionRequestPayload 游戏动作请求（服务器转发给对方）
type GameActionRequestPayload struct {
	RoomID      string `json:"roomId"`
	Action      string `json:"action"` // "surrender" or "undo"
	FromPlayer  string `json:"fromPlayer"`
	FromAlias   string `json:"fromAlias"`
	MovePosition *Position `json:"movePosition,omitempty"` // 悔棋时使用
}

// GameActionResponsePayload 游戏动作响应（服务器转发给请求者）
type GameActionResponsePayload struct {
	RoomID     string `json:"roomId"`
	Action     string `json:"action"` // "surrender" or "undo"
	FromPlayer string `json:"fromPlayer"`
	FromAlias  string `json:"fromAlias"`
	Accept     bool   `json:"accept"`
}
