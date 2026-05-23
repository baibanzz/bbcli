package main

import (
	"bbcli/internal/server"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}

	// 创建并启动服务器
	srv, err := server.NewServer(dir, ":8080")
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 打开浏览器
	url := "http://localhost:8080"
	if err := openBrowser(url); err != nil {
		log.Printf("打开浏览器失败: %v", err)
	}

	log.Printf("bbcli 已在浏览器中打开: %s", url)
	log.Println("按 Ctrl+C 停止服务器")

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// openBrowser 打开浏览器
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
