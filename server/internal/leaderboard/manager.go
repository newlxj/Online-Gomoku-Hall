package leaderboard

import (
	"gomoku-server/internal/message"
	"time"
)

// Manager 排行榜管理器
type Manager struct {
	storage   *Storage
	flushTick *time.Ticker
	stopChan  chan struct{}
}

// NewManager 创建排行榜管理器
func NewManager(filePath string) (*Manager, error) {
	storage, err := NewStorage(filePath)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		storage:   storage,
		flushTick: time.NewTicker(10 * time.Second),
		stopChan:  make(chan struct{}),
	}

	// 启动定时保存
	go m.autoFlush()

	return m, nil
}

// autoFlush 自动保存
func (m *Manager) autoFlush() {
	for {
		select {
		case <-m.flushTick.C:
			if m.storage.IsDirty() {
				m.storage.Flush()
			}
		case <-m.stopChan:
			return
		}
	}
}

// Stop 停止管理器
func (m *Manager) Stop() {
	close(m.stopChan)
	m.flushTick.Stop()
	// 最后保存一次
	m.storage.ForceFlush()
}

// RecordGameResult 记录游戏结果
func (m *Manager) RecordGameResult(winner, loser string, isDraw bool, winnerBonus int) {
	if isDraw {
		m.storage.RecordDraw(winner)
		m.storage.RecordDraw(loser)
	} else {
		m.storage.RecordWin(winner, winnerBonus)
		m.storage.RecordLoss(loser)
	}
}

// GetLeaderboard 获取排行榜
func (m *Manager) GetLeaderboard(topN int) []message.LeaderboardEntry {
	entries := m.storage.GetTopEntries(topN)
	result := make([]message.LeaderboardEntry, 0, len(entries))

	for i, entry := range entries {
		totalGames := entry.Wins + entry.Losses + entry.Draws
		winRate := 0.0
		if totalGames > 0 {
			winRate = float64(entry.Wins) / float64(totalGames) * 100
		}

		result = append(result, message.LeaderboardEntry{
			Rank:      i + 1,
			Alias:     entry.Alias,
			Wins:      entry.Wins,
			Losses:    entry.Losses,
			Draws:     entry.Draws,
			WinRate:   winRate,
			Score:     entry.Score,
			RankTitle: GetRankTitle(entry.Score),
		})
	}

	return result
}

// GetPlayerRank 获取玩家排名
func (m *Manager) GetPlayerRank(alias string) int {
	entries := m.storage.GetTopEntries(1000) // 获取足够多的条目
	for i, entry := range entries {
		if entry.Alias == alias {
			return i + 1
		}
	}
	return -1
}

// GetPlayerScore 获取玩家分数
func (m *Manager) GetPlayerScore(alias string) int {
	entry := m.storage.GetEntry(alias)
	if entry == nil {
		return 0
	}
	return entry.Score
}

// ForceFlush 强制保存
func (m *Manager) ForceFlush() error {
	return m.storage.ForceFlush()
}

// Rank 定义段位
type Rank struct {
	Title string
	Min   int
	Max   int
}

// 段位列表
var ranks = []Rank{
	{"大师", 1000, 999999},
	{"钻石", 600, 999},
	{"黄金", 300, 599},
	{"白银", 100, 299},
	{"青铜", 0, 99},
}

// GetRankTitle 获取段位称号
func GetRankTitle(score int) string {
	for _, rank := range ranks {
		if score >= rank.Min && score <= rank.Max {
			return rank.Title
		}
	}
	return "青铜"
}
