package player

import (
	"fmt"
	"sync"
	"time"
)

// AliasManager 别名管理器 - 管理在线玩家别名唯一性
type AliasManager struct {
	aliases map[string]*AliasInfo // alias -> info
	mu      sync.RWMutex
}

// AliasInfo 别名信息
type AliasInfo struct {
	PlayerID    string
	ConnectedAt time.Time
	LastActive  time.Time
}

// NewAliasManager 创建别名管理器
func NewAliasManager() *AliasManager {
	return &AliasManager{
		aliases: make(map[string]*AliasInfo),
	}
}

// Register 注册别名
func (am *AliasManager) Register(alias, playerID string) bool {
	am.mu.Lock()
	defer am.mu.Unlock()

	if info, exists := am.aliases[alias]; exists {
		// 检查是否是同一个玩家
		if info.PlayerID == playerID {
			info.LastActive = time.Now()
			return true
		}
		// 别名已被其他玩家占用
		return false
	}

	am.aliases[alias] = &AliasInfo{
		PlayerID:    playerID,
		ConnectedAt: time.Now(),
		LastActive:  time.Now(),
	}
	return true
}

// Unregister 注销别名
func (am *AliasManager) Unregister(alias string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.aliases, alias)
}

// IsAvailable 检查别名是否可用
func (am *AliasManager) IsAvailable(alias string) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	_, exists := am.aliases[alias]
	return !exists
}

// UpdateActivity 更新活跃时间
func (am *AliasManager) UpdateActivity(alias string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if info, exists := am.aliases[alias]; exists {
		info.LastActive = time.Now()
	}
}

// GetOnlineAliases 获取所有在线别名
func (am *AliasManager) GetOnlineAliases() []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	aliases := make([]string, 0, len(am.aliases))
	for alias := range am.aliases {
		aliases = append(aliases, alias)
	}
	return aliases
}

// GetOnlineCount 获取在线人数
func (am *AliasManager) GetOnlineCount() int {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return len(am.aliases)
}

// CleanupInactive 清理不活跃的别名(超时5分钟)
func (am *AliasManager) CleanupInactive(timeout time.Duration) int {
	am.mu.Lock()
	defer am.mu.Unlock()

	count := 0
	now := time.Now()
	for alias, info := range am.aliases {
		if now.Sub(info.LastActive) > timeout {
			delete(am.aliases, alias)
			count++
		}
	}
	return count
}

// 随机别名生成的形容词列表
var adjectives = []string{
	"疾风", "暗影", "烈焰", "冰霜", "雷霆",
	"星空", "月影", "日光", "云霄", "海浪",
	"山岳", "森林", "沙漠", "极光", "彩虹",
	"神秘", "勇敢", "智慧", "沉稳", "迅猛",
	"精准", "冷静", "热情", "坚韧", "灵动",
}

// 随机别名生成的名词列表
var nouns = []string{
	"棋手", "大师", "王者", "侠客", "剑士",
	"法师", "猎手", "骑士", "战士", "智者",
	"圣者", "仙人", "英雄", "传奇", "宗师",
	"高手", "达人", "专家", "新星", "霸主",
}

// GenerateRandomAlias 生成随机别名
func GenerateRandomAlias() string {
	// 使用时间戳作为简单的随机种子
	seed := time.Now().UnixNano()
	adjIndex := int(seed) % len(adjectives)
	nounIndex := int(seed / 1000) % len(nouns)
	return adjectives[adjIndex] + nouns[nounIndex]
}

// GenerateUniqueAlias 生成唯一的别名
// 如果提供的别名已被占用，会在后面添加数字后缀
func (am *AliasManager) GenerateUniqueAlias(baseAlias string) string {
	if am.IsAvailable(baseAlias) {
		return baseAlias
	}

	// 如果基础别名已被占用，尝试添加数字后缀
	for i := 2; i <= 100; i++ {
		newAlias := fmt.Sprintf("%s%d", baseAlias, i)
		if am.IsAvailable(newAlias) {
			return newAlias
		}
	}

	// 如果还是冲突，生成随机别名
	for i := 0; i < 10; i++ {
		randomAlias := GenerateRandomAlias()
		if am.IsAvailable(randomAlias) {
			return randomAlias
		}
	}

	// 最后使用时间戳
	return fmt.Sprintf("Player_%d", time.Now().UnixNano()%100000)
}
