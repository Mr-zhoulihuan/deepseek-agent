package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      AI 备份系统管理终端                  ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println("所有操作由 DeepSeek AI 驱动")

	for {
		fmt.Println("\n请选择操作:")
		fmt.Println("  1. 📋 查询文件备份任务")
		fmt.Println("  2. 🗄️  查询数据库备份任务")
		fmt.Println("  3. ➕ 创建文件备份任务")
		fmt.Println("  4. 🗄️  创建数据库备份任务")
		fmt.Println("  5. 💬 自由对话")
		fmt.Println("  6. ❌ 退出")
		fmt.Print("\n请输入选项 (1-6): ")

		scanned := scanner.Scan()
		if !scanned {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			handleQueryFileBackup()
		case "2":
			handleQueryDatabaseBackup()
		case "3":
			handleCreateFileBackup()
		case "4":
			handleCreateDatabaseBackup()
		case "5":
			handleFreeChat(scanner)
		case "6":
			fmt.Println("再见！")
			return
		default:
			fmt.Println("❌ 无效选项，请输入 1-6")
		}
	}
}

func handleQueryFileBackup() {
	fmt.Println("\n── 查询文件备份任务 ──")
	fmt.Println("⏳ 正在向 DeepSeek 查询...")
	if err := QueryFileBackup(1, 3); err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
	}
}

func handleQueryDatabaseBackup() {
	fmt.Println("\n── 查询数据库备份任务 ──")
	fmt.Println("⏳ 正在向 DeepSeek 查询...")
	if err := QueryDatabaseBackup(1, 3); err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
	}
}

func handleCreateFileBackup() {
	fmt.Println("\n── 创建文件备份任务 ──")
	fmt.Println("⏳ 正在向 DeepSeek 提交创建请求...")
	if err := CreateFileBackup(
		"Web 服务器数据备份",
		"web-prod-01",
		"192.168.1.10",
		"/var/www/html",
		"/backup/web-prod-01",
		"每天 02:00 执行",
		"gzip",
		"AES-256",
		30,
	); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
	}
}

func handleCreateDatabaseBackup() {
	fmt.Println("\n── 创建数据库备份任务 ──")
	fmt.Println("⏳ 正在向 DeepSeek 提交创建请求...")
	if err := CreateDatabaseBackup(
		"MySQL 生产库全量备份",
		"MySQL",
		"192.168.1.30",
		3306,
		"production_db",
		"全量备份",
		"/backup/mysql/production_db",
		"每天 03:00 执行",
		"gzip",
		"AES-256",
		30,
	); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
	}
}
