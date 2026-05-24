package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DbBackupTask struct {
	TaskID        string `json:"任务ID"`
	TaskName      string `json:"任务名称"`
	DbType        string `json:"数据库类型"`
	Host          string `json:"连接地址"`
	DbName        string `json:"数据库名称"`
	BackupType    string `json:"备份类型"`
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
	DbSize        string `json:"数据库大小"`
	BackupSize    string `json:"备份文件大小"`
	RetentionDays string `json:"保留天数"`
}

const dbQuerySystemPrompt = `你是一个企业备份系统管理助手。用户会查询数据库备份任务信息。
你必须严格按照以下 JSON 数组格式输出，不要使用 markdown 代码块，只输出纯 JSON：

[
  {
    "任务ID": "DB-20240501-001",
    "任务名称": "MySQL 生产库全量备份",
    "数据库类型": "MySQL",
    "连接地址": "192.168.1.30:3306",
    "数据库名称": "production_db",
    "备份类型": "全量备份",
    "目标路径": "/backup/mysql/production_db",
    "状态": "已启用",
    "运行状态": "空闲",
    "调度策略": "每天 03:00 执行",
    "创建时间": "2024-01-15 10:30:00",
    "上次执行时间": "2024-05-01 03:00:05",
    "上次执行结果": "成功",
    "总执行次数": "120",
    "成功次数": "118",
    "失败次数": "2",
    "数据库大小": "256 GB",
    "备份文件大小": "85 GB",
    "保留天数": "30天"
  }
]

要求：
- 必须使用上述字段名（中文），不得更改
- 生成 3 个任务，包含不同的数据库类型(MySQL、PostgreSQL、SQL Server)
- 包含不同的备份类型(全量备份、增量备份)
- 至少有一个任务是"备份中"运行状态
- 运行状态只能是: 空闲、备份中、已失效、已暂停
- 上次执行结果只能是: 成功、失败
- 状态只能是: 已启用、正常
- 只输出 JSON 数组，不要有任何其他文字`

func QueryDatabaseBackup(page, pageSize int) error {
	prompt := fmt.Sprintf(
		`请生成 %d 个数据库备份任务信息（第 %d 页，每页 %d 条）。`,
		pageSize, page, pageSize,
	)

	result, err := chatWithDeepSeek(dbQuerySystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)

	tasks, err := parseDbBackupTasks(result.Content)
	if err != nil {
		fmt.Printf("⚠️  数据解析失败，显示原始输出:\n%s\n", result.Content)
		return nil
	}

	renderDbBackupTable(tasks)
	return nil
}

func parseDbBackupTasks(content string) ([]DbBackupTask, error) {
	content = strings.TrimSpace(content)
	if idx := strings.Index(content, "["); idx != -1 {
		content = content[idx:]
	}
	if idx := strings.LastIndex(content, "]"); idx != -1 {
		content = content[:idx+1]
	}

	var tasks []DbBackupTask
	if err := json.Unmarshal([]byte(content), &tasks); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}
	return tasks, nil
}

func renderDbBackupTable(tasks []DbBackupTask) {
	headers := []string{"任务ID", "任务名称", "数据库", "连接地址", "备份类型", "运行状态", "上次结果", "调度策略"}

	var rows [][]string
	for _, t := range tasks {
		rows = append(rows, []string{
			t.TaskID,
			t.TaskName,
			t.DbType + "/" + t.DbName,
			t.Host,
			t.BackupType,
			statusColor(t.RunStatus),
			statusColor(t.LastResult),
			t.Schedule,
		})
	}

	PrintSubTitle("📊 数据库备份任务汇总")
	PrintTable(headers, rows)

	for _, t := range tasks {
		renderDbBackupDetail(t)
	}
}

func renderDbBackupDetail(t DbBackupTask) {
	fields := map[string]string{
		"任务ID":     t.TaskID,
		"任务名称":     t.TaskName,
		"数据库类型":    t.DbType,
		"连接地址":     t.Host,
		"数据库名称":    t.DbName,
		"备份类型":     t.BackupType,
		"目标路径":     t.DestPath,
		"状态":       statusColor(t.Status),
		"运行状态":     statusColor(t.RunStatus),
		"调度策略":     t.Schedule,
		"创建时间":     t.CreateTime,
		"上次执行时间":   t.LastRunTime,
		"上次执行结果":   statusColor(t.LastResult),
		"总执行/成功/失败": t.TotalRuns + " / " + t.SuccessRuns + " / " + t.FailedRuns,
		"数据库大小":    t.DbSize,
		"备份文件大小":   t.BackupSize,
		"保留天数":     t.RetentionDays,
	}

	PrintSubTitle("📋 " + t.TaskName)
	PrintDetailTable(fields)
}
