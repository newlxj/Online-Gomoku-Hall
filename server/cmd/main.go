package main

import (
	"embed"
	"flag"
	"fmt"
	"gomoku-server/internal/leaderboard"
	"gomoku-server/internal/server"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 命令行参数
	port := flag.Int("port", 8080, "服务器监听端口")
	host := flag.String("host", "0.0.0.0", "服务器监听地址")
	dataDir := flag.String("data", "./data", "数据存储目录")
	flag.Parse()

	// 环境变量可以覆盖命令行参数
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envHost := os.Getenv("HOST"); envHost != "" {
		*host = envHost
	}

	// 初始化排行榜管理器
	lbManager, err := leaderboard.NewManager(*dataDir + "/scores.json")
	if err != nil {
		log.Fatalf("Failed to initialize leaderboard: %v", err)
	}
	defer lbManager.Stop()

	// 创建服务器
	srv := server.NewServer(lbManager)

	// WebSocket 路由
	http.HandleFunc("/ws", srv.HandleWebSocket)

	// 静态文件服务 - 从嵌入的文件系统
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to load static files: %v", err)
	}

	// 处理静态文件请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// WebSocket 路径不处理
		if path == "/ws" {
			http.NotFound(w, r)
			return
		}

		// 根路径直接返回 index.html 内容
		if path == "/" || path == "" {
			data, err := fs.ReadFile(staticFS, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		// 移除前导斜杠，因为 staticFS 的根已经是 static/ 目录
		filePath := strings.TrimPrefix(path, "/")

		// 尝试打开文件
		file, err := staticFS.Open(filePath)
		if err != nil {
			// 文件不存在
			// 对于静态资源请求（assets/、.js、.css等），返回404而不是index.html
			if strings.HasPrefix(path, "/assets/") ||
				strings.HasSuffix(path, ".js") ||
				strings.HasSuffix(path, ".css") ||
				strings.HasSuffix(path, ".ico") ||
				strings.HasSuffix(path, ".svg") ||
				strings.HasSuffix(path, ".png") ||
				strings.HasSuffix(path, ".jpg") ||
				strings.HasSuffix(path, ".woff") ||
				strings.HasSuffix(path, ".woff2") {
				http.NotFound(w, r)
				return
			}
			// 对于其他路径（SPA路由），返回 index.html
			data, err := fs.ReadFile(staticFS, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}
		file.Close()

		// 设置正确的 MIME 类型
		ext := filepath.Ext(path)
		switch ext {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".woff":
			w.Header().Set("Content-Type", "font/woff")
		case ".woff2":
			w.Header().Set("Content-Type", "font/woff2")
		}

		// 直接读取并返回文件内容，避免重定向问题
		data, err := fs.ReadFile(staticFS, filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Write(data)
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Server starting on %s...", addr)
	log.Printf("Access the game at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
