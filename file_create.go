package main

import "fmt"

const createSystemPrompt = `你是一个企业备份系统管理助手。用户会向你提交备份任务创建请求。
你需要根据用户提供的任务信息，确认任务创建成功，并以友好的中文格式输出任务详情。
输出应包含：任务ID(根据当前时间生成)、任务名称、服务器/数据库信息、备份路径、
调度策略、压缩/加密配置、保留天数等信息。
格式清晰易读，用分隔线装饰。`

func CreateFileBackup(taskName, serverName, serverIp, sourcePath, destPath, schedule, compressAlg, encryptAlg string, retentionDays int) error {
	prompt := fmt.Sprintf(
		`请帮我创建一个文件备份任务，要求如下：
- 任务名称: %s
- 服务器: %s (%s)
- 源路径: %s
- 目标路径: %s
- 调度策略: %s
- 压缩算法: %s
- 加密算法: %s
- 保留天数: %d天

请确认任务创建成功，输出任务详情。`,
		taskName, serverName, serverIp, sourcePath, destPath, schedule, compressAlg, encryptAlg, retentionDays,
	)

	result, err := chatWithDeepSeek(createSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)
	return nil
}

func CreateDatabaseBackup(taskName, dbType, host string, port int, dbName, backupType, destPath, schedule, compressAlg, encryptAlg string, retentionDays int) error {
	prompt := fmt.Sprintf(
		`请帮我创建一个数据库备份任务，要求如下：
- 任务名称: %s
- 数据库类型: %s
- 连接地址: %s:%d
- 数据库名: %s
- 备份类型: %s
- 目标路径: %s
- 调度策略: %s
- 压缩算法: %s
- 加密算法: %s
- 保留天数: %d天

请确认任务创建成功，输出任务详情。`,
		taskName, dbType, host, port, dbName, backupType, destPath, schedule, compressAlg, encryptAlg, retentionDays,
	)

	result2, err := chatWithDeepSeek(createSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result2)
	return nil
}
