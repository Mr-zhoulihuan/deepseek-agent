package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FileBackupTask struct {
	TaskID        string `json:"任务ID"`
	TaskName      string `json:"任务名称"`
	ServerName    string `json:"服务器名称"`
	ServerIP      string `json:"服务器IP"`
	SourcePath    string `json:"源路径"`
	DestPath      string `json:"目标路径"`
	Status        string `json:"状态"`
	RunStatus     string `json:"运行状态"`
	Schedule      string `json:"调度策略"`
	CreateTime    string `json:"创建时间"`
	LastRunTime   string `json:"上次执行时间"`
	LastResult    string `json:"上次执行结果"`
	TotalRuns     string `json:"总执行次数"`
	SuccessRuns   string `json:"成功次数"`
	FailedRuns    string `json:"失败次数"`
	TotalSize     string `json:"总数据大小"`
	CompressSize  string `json:"压缩后大小"`
	FileCount     string `json:"文件数量"`
	RetentionDays string `json:"保留天数"`
}

const fileQuerySystemPrompt = `你是一个企业备份系统管理助手。用户会查询文件备份任务信息。
你必须严格按照以下 JSON 数组格式输出，不要使用 markdown 代码块，只输出纯 JSON：

[
  {
    "任务ID": "FILE-20240501-001",
    "任务名称": "Web 服务器数据备份",
    "服务器名称": "web-prod-01",
    "服务器IP": "192.168.1.10",
    "源路径": "/var/www/html",
    "目标路径": "/backup/web-prod-01",
    "状态": "已启用",
    "运行状态": "空闲",
    "调度策略": "每天 02:00 执行",
    "创建时间": "2024-01-10 09:00:00",
    "上次执行时间": "2024-05-01 02:00:03",
    "上次执行结果": "成功",
    "总执行次数": "150",
    "成功次数": "148",
    "失败次数": "2",
    "总数据大小": "50 GB",
    "压缩后大小": "12 GB",
    "文件数量": "15000",
    "保留天数": "30天"
  }
]

要求：
- 必须使用上述字段名（中文），不得更改
- 生成 3 个任务，包含不同场景：网站数据备份、日志归档、配置文件备份等
- 至少有一个任务是"备份中"运行状态
- 运行状态只能是: 空闲、备份中、已失效、已暂停
- 上次执行结果只能是: 成功、失败
- 状态只能是: 已启用、正常
- 只输出 JSON 数组，不要有任何其他文字`

func QueryFileBackup(page, pageSize int) error {
	prompt := fmt.Sprintf(
		`请生成 %d 个文件备份任务信息（第 %d 页，每页 %d 条）。`,
		pageSize, page, pageSize,
	)

	result, err := chatWithDeepSeek(fileQuerySystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)

	tasks, err := parseFileBackupTasks(result.Content)
	if err != nil {
		fmt.Printf("⚠️  数据解析失败，显示原始输出:\n%s\n", result.Content)
		return nil
	}

	renderFileBackupTable(tasks)
	return nil
}

func parseFileBackupTasks(content string) ([]FileBackupTask, error) {
	content = strings.TrimSpace(content)
	if idx := strings.Index(content, "["); idx != -1 {
		content = content[idx:]
	}
	if idx := strings.LastIndex(content, "]"); idx != -1 {
		content = content[:idx+1]
	}

	var tasks []FileBackupTask
	if err := json.Unmarshal([]byte(content), &tasks); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}
	return tasks, nil
}

func renderFileBackupTable(tasks []FileBackupTask) {
	headers := []string{"任务ID", "任务名称", "服务器", "源路径", "运行状态", "上次结果", "调度策略"}

	var rows [][]string
	for _, t := range tasks {
		rows = append(rows, []string{
			t.TaskID,
			t.TaskName,
			t.ServerName + "/" + t.ServerIP,
			t.SourcePath,
			statusColor(t.RunStatus),
			statusColor(t.LastResult),
			t.Schedule,
		})
	}

	PrintSubTitle("📊 文件备份任务汇总")
	PrintTable(headers, rows)

	for _, t := range tasks {
		renderFileBackupDetail(t)
	}
}

func renderFileBackupDetail(t FileBackupTask) {
	fields := map[string]string{
		"任务ID":         t.TaskID,
		"任务名称":         t.TaskName,
		"服务器名称":        t.ServerName,
		"服务器IP":        t.ServerIP,
		"源路径":          t.SourcePath,
		"目标路径":         t.DestPath,
		"状态":           statusColor(t.Status),
		"运行状态":         statusColor(t.RunStatus),
		"调度策略":         t.Schedule,
		"创建时间":         t.CreateTime,
		"上次执行时间":       t.LastRunTime,
		"上次执行结果":       statusColor(t.LastResult),
		"总执行/成功/失败":     t.TotalRuns + " / " + t.SuccessRuns + " / " + t.FailedRuns,
		"总数据大小":        t.TotalSize,
		"压缩后大小":        t.CompressSize,
		"文件数量":         t.FileCount,
		"保留天数":         t.RetentionDays,
	}

	PrintSubTitle("📋 " + t.TaskName)
	PrintDetailTable(fields)
}
