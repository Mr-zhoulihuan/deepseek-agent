package main

import (
	"fmt"
	"time"
)

func simulateProgress(taskName, taskId string) {
	fmt.Printf("\n⏳ 任务 [%s] (%s) 正在备份...\n", taskName, taskId)
	fmt.Println("──────────────────────────────────────────────────────")

	total := 30
	for i := 0; i <= total; i++ {
		percent := i * 100 / total
		bar := ""
		for j := 0; j < total; j++ {
			if j < i {
				bar += "█"
			} else {
				bar += "░"
			}
		}

		fmt.Printf("\r  进度: [%s] %3d%%", bar, percent)

		rate := 100
		if percent < 30 {
			rate = 200
		} else if percent < 60 {
			rate = 150
		} else if percent < 80 {
			rate = 100
		} else if percent < 95 {
			rate = 80
		} else {
			rate = 150
		}
		time.Sleep(time.Duration(rate) * time.Millisecond)
	}

	fmt.Println()
	fmt.Println("✅ 备份完成！")
	fmt.Println("──────────────────────────────────────────────────────")
}
