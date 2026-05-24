package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CreatedTask struct {
	TaskID        string `json:"任务ID"`
	TaskName      string `json:"任务名称"`
	ServerInfo    string `json:"服务器信息"`
	DbInfo        string `json:"数据库信息"`
	SourcePath    string `json:"源路径"`
	DestPath      string `json:"目标路径"`
	Schedule      string `json:"调度策略"`
	CompressAlg   string `json:"压缩算法"`
	EncryptAlg    string `json:"加密算法"`
	RetentionDays string `json:"保留天数"`
	Status        string `json:"状态"`
	CreateTime    string `json:"创建时间"`
}

const createSystemPrompt = `你是一个企业备份系统管理助手。用户会提交备份任务创建请求。
你必须确认任务创建成功，并严格按照以下 JSON 格式输出任务详情，不要使用 markdown 代码块，只输出纯 JSON：

{
  "任务ID": "TASK-20240501-001",
  "任务名称": "任务名称",
  "服务器信息": "服务器名 (IP)",
  "数据库信息": "数据库类型 (host:port) / 数据库名",
  "源路径": "/source/path",
  "目标路径": "/backup/path",
  "调度策略": "每天 02:00 执行",
  "压缩算法": "gzip",
  "加密算法": "AES-256",
  "保留天数": "30天",
  "状态": "已创建",
  "创建时间": "2024-05-01 10:30:00"
}

要求：
- 必须使用上述字段名（中文），不得更改
- 任务ID 根据当前日期时间生成
- 状态必须是"已创建"
- 确认任务创建成功
- 只输出 JSON 对象，不要有任何其他文字`

func CreateFileBackup(taskName, serverName, serverIp, sourcePath, destPath, schedule, compressAlg, encryptAlg string, retentionDays int) error {
	prompt := fmt.Sprintf(
		`请创建一个文件备份任务，参数如下：
- 任务名称: %s
- 服务器: %s (%s)
- 源路径: %s
- 目标路径: %s
- 调度策略: %s
- 压缩算法: %s
- 加密算法: %s
- 保留天数: %d天`,
		taskName, serverName, serverIp, sourcePath, destPath, schedule, compressAlg, encryptAlg, retentionDays,
	)

	result, err := chatWithDeepSeek(createSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)

	task, err := parseCreatedTask(result.Content)
	if err != nil {
		fmt.Printf("⚠️  数据解析失败，显示原始输出:\n%s\n", result.Content)
		return nil
	}

	renderCreatedTask(task)
	return nil
}

func CreateDatabaseBackup(taskName, dbType, host string, port int, dbName, backupType, destPath, schedule, compressAlg, encryptAlg string, retentionDays int) error {
	prompt := fmt.Sprintf(
		`请创建一个数据库备份任务，参数如下：
- 任务名称: %s
- 数据库类型: %s
- 连接地址: %s:%d
- 数据库名: %s
- 备份类型: %s
- 目标路径: %s
- 调度策略: %s
- 压缩算法: %s
- 加密算法: %s
- 保留天数: %d天`,
		taskName, dbType, host, port, dbName, backupType, destPath, schedule, compressAlg, encryptAlg, retentionDays,
	)

	result, err := chatWithDeepSeek(createSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)

	task, err := parseCreatedTask(result.Content)
	if err != nil {
		fmt.Printf("⚠️  数据解析失败，显示原始输出:\n%s\n", result.Content)
		return nil
	}

	renderCreatedTask(task)
	return nil
}

func parseCreatedTask(content string) (*CreatedTask, error) {
	content = strings.TrimSpace(content)
	if idx := strings.Index(content, "{"); idx != -1 {
		content = content[idx:]
	}
	if idx := strings.LastIndex(content, "}"); idx != -1 {
		content = content[:idx+1]
	}

	var task CreatedTask
	if err := json.Unmarshal([]byte(content), &task); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}
	return &task, nil
}

func renderCreatedTask(t *CreatedTask) {
	PrintSuccess("任务创建成功！")

	fields := map[string]string{
		"任务ID": t.TaskID,
		"任务名称": t.TaskName,
		"创建时间": t.CreateTime,
		"状态":   green(t.Status),
		"调度策略": t.Schedule,
		"压缩算法": t.CompressAlg,
		"加密算法": t.EncryptAlg,
		"保留天数": t.RetentionDays,
	}

	if t.ServerInfo != "" {
		fields["服务器信息"] = t.ServerInfo
		fields["源路径"] = t.SourcePath
	}
	if t.DbInfo != "" {
		fields["数据库信息"] = t.DbInfo
	}
	if t.DestPath != "" {
		fields["目标路径"] = t.DestPath
	}

	PrintDetailTable(fields)
}
