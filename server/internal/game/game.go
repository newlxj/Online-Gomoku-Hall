package game

import (
	"gomoku-server/internal/message"
)

// Game 游戏逻辑
type Game struct {
	board       [15][15]int // 0=空, 1=黑, 2=白
	currentTurn int         // 当前回合: 1=黑, 2=白
	boardSize   int
	moveCount   int
}

// NewGame 创建新游戏
func NewGame(boardSize int) *Game {
	return &Game{
		board:       [15][15]int{},
		currentTurn: 1, // 黑棋先行
		boardSize:   boardSize,
		moveCount:   0,
	}
}

// GetBoard 获取棋盘
func (g *Game) GetBoard() [][]int {
	board := make([][]int, g.boardSize)
	for i := 0; i < g.boardSize; i++ {
		board[i] = make([]int, g.boardSize)
		for j := 0; j < g.boardSize; j++ {
			board[i][j] = g.board[i][j]
		}
	}
	return board
}

// GetCurrentTurn 获取当前回合
func (g *Game) GetCurrentTurn() int {
	return g.currentTurn
}

// MakeMove 落子
// 返回: 是否成功, 获胜连线(如果有)
func (g *Game) MakeMove(row, col, pieceType int) (bool, []message.Position) {
	// 检查是否轮到该玩家
	if pieceType != g.currentTurn {
		return false, nil
	}

	// 检查位置是否有效
	if row < 0 || row >= g.boardSize || col < 0 || col >= g.boardSize {
		return false, nil
	}

	// 检查位置是否为空
	if g.board[row][col] != 0 {
		return false, nil
	}

	// 落子
	g.board[row][col] = pieceType
	g.moveCount++

	// 检查胜利
	winningLine := g.checkWin(row, col, pieceType)

	// 切换回合
	if g.currentTurn == 1 {
		g.currentTurn = 2
	} else {
		g.currentTurn = 1
	}

	return true, winningLine
}

// checkWin 检查胜利
func (g *Game) checkWin(row, col, pieceType int) []message.Position {
	// 四个方向: 水平、垂直、两条对角线
	directions := [][2]int{
		{0, 1},  // 水平
		{1, 0},  // 垂直
		{1, 1},  // 对角线
		{1, -1}, // 反对角线
	}

	for _, dir := range directions {
		line := g.getLine(row, col, dir[0], dir[1], pieceType)
		if len(line) >= 5 {
			return line
		}
	}

	return nil
}

// getLine 获取某方向上的连续棋子
func (g *Game) getLine(row, col, dRow, dCol, pieceType int) []message.Position {
	line := []message.Position{{Row: row, Col: col}}

	// 正向
	r, c := row+dRow, col+dCol
	for r >= 0 && r < g.boardSize && c >= 0 && c < g.boardSize && g.board[r][c] == pieceType {
		line = append(line, message.Position{Row: r, Col: c})
		r += dRow
		c += dCol
	}

	// 反向
	r, c = row-dRow, col-dCol
	for r >= 0 && r < g.boardSize && c >= 0 && c < g.boardSize && g.board[r][c] == pieceType {
		line = append([]message.Position{{Row: r, Col: c}}, line...)
		r -= dRow
		c -= dCol
	}

	if len(line) >= 5 {
		return line[:5] // 只返回5个连续的
	}
	return nil
}

// IsBoardFull 检查棋盘是否已满
func (g *Game) IsBoardFull() bool {
	return g.moveCount >= g.boardSize*g.boardSize
}

// GetMoveCount 获取落子数量
func (g *Game) GetMoveCount() int {
	return g.moveCount
}

// IsValidMove 检查是否是有效落子
func (g *Game) IsValidMove(row, col int) bool {
	if row < 0 || row >= g.boardSize || col < 0 || col >= g.boardSize {
		return false
	}
	return g.board[row][col] == 0
}

// GetPiece 获取某位置的棋子
func (g *Game) GetPiece(row, col int) int {
	if row < 0 || row >= g.boardSize || col < 0 || col >= g.boardSize {
		return -1
	}
	return g.board[row][col]
}
