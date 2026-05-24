package main

import (
	"fmt"
	"strings"
	"time"
)

func simulateProgress(taskName, taskId string) {
	fmt.Println()
	fmt.Printf("  "+yellow("⏳")+" "+bold("任务 [%s] (%s)")+" "+yellow("正在备份...")+"\n", taskName, taskId)
	fmt.Println("  " + dim(strings.Repeat("─", 50)))

	total := 30
	barWidth := 30
	for i := 0; i <= total; i++ {
		percent := i * 100 / total
		filled := i * barWidth / total
		bar := green(strings.Repeat("█", filled)) + dim(strings.Repeat("░", barWidth-filled))

		fmt.Printf("\r  "+dim("进度:")+" [%s] "+bold("%3d%%"), bar, percent)

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
	fmt.Println("  " + green("✅") + " " + green("备份完成！"))
	fmt.Println("  " + dim(strings.Repeat("─", 50)))
}
