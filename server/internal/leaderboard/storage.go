package leaderboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage 排行榜存储
type Storage struct {
	filePath string
	entries  map[string]*Entry // alias -> entry
	mu       sync.RWMutex
	dirty    bool
	lastSave time.Time
}

// Entry 排行榜条目
type Entry struct {
	Alias      string    `json:"alias"`
	Wins       int       `json:"wins"`
	Losses     int       `json:"losses"`
	Draws      int       `json:"draws"`
	Score      int       `json:"score"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// StorageData 存储文件数据结构
type StorageData struct {
	Entries    []*Entry `json:"entries"`
	LastUpdate string   `json:"lastUpdate"`
}

// NewStorage 创建存储
func NewStorage(filePath string) (*Storage, error) {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	s := &Storage{
		filePath: filePath,
		entries:  make(map[string]*Entry),
		dirty:    false,
		lastSave: time.Time{},
	}

	// 加载现有数据
	if err := s.load(); err != nil {
		// 文件不存在是正常的，忽略
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return s, nil
}

// load 从文件加载
func (s *Storage) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var storageData StorageData
	if err := json.Unmarshal(data, &storageData); err != nil {
		return err
	}

	for _, entry := range storageData.Entries {
		s.entries[entry.Alias] = entry
	}

	return nil
}

// save 保存到文件
func (s *Storage) save() error {
	entries := make([]*Entry, 0, len(s.entries))
	for _, entry := range s.entries {
		entries = append(entries, entry)
	}

	storageData := StorageData{
		Entries:    entries,
		LastUpdate: time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(storageData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// GetOrCreateEntry 获取或创建条目
func (s *Storage) GetOrCreateEntry(alias string) *Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[alias]
	if !exists {
		entry = &Entry{
			Alias:     alias,
			Wins:      0,
			Losses:    0,
			Draws:     0,
			Score:     0,
			UpdatedAt: time.Now(),
		}
		s.entries[alias] = entry
	}

	return entry
}

// RecordWin 记录胜利
func (s *Storage) RecordWin(alias string, bonusScore int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[alias]
	if !exists {
		entry = &Entry{
			Alias: alias,
		}
		s.entries[alias] = entry
	}

	entry.Wins++
	entry.Score += 25 + bonusScore // 基础25分 + 奖励分
	if entry.Score < 0 {
		entry.Score = 0
	}
	entry.UpdatedAt = time.Now()
	s.dirty = true
}

// RecordLoss 记录失败
func (s *Storage) RecordLoss(alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[alias]
	if !exists {
		entry = &Entry{
			Alias: alias,
		}
		s.entries[alias] = entry
	}

	entry.Losses++
	entry.Score -= 10
	if entry.Score < 0 {
		entry.Score = 0
	}
	entry.UpdatedAt = time.Now()
	s.dirty = true
}

// RecordDraw 记录平局
func (s *Storage) RecordDraw(alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[alias]
	if !exists {
		entry = &Entry{
			Alias: alias,
		}
		s.entries[alias] = entry
	}

	entry.Draws++
	entry.Score += 5
	entry.UpdatedAt = time.Now()
	s.dirty = true
}

// GetEntry 获取条目
func (s *Storage) GetEntry(alias string) *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entries[alias]
}

// GetTopEntries 获取前N名
func (s *Storage) GetTopEntries(n int) []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 转换为切片
	entries := make([]*Entry, 0, len(s.entries))
	for _, entry := range s.entries {
		entries = append(entries, entry)
	}

	// 按分数排序(简单冒泡，数据量小)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Score > entries[i].Score {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

// Flush 如果有变化则保存
func (s *Storage) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.dirty {
		return nil
	}

	if err := s.save(); err != nil {
		return err
	}

	s.dirty = false
	s.lastSave = time.Now()
	return nil
}

// ForceFlush 强制保存
func (s *Storage) ForceFlush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.save()
}

// IsDirty 是否有未保存的变化
func (s *Storage) IsDirty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

// GetEntryCount 获取条目数量
func (s *Storage) GetEntryCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
