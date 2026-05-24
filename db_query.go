package main

import "fmt"

const dbQuerySystemPrompt = `你是一个企业备份系统管理助手。用户会向你查询数据库备份任务信息。
请生成合理的数据库备份任务数据，包含以下字段：
- 任务ID、任务名称、数据库类型(MySQL/PostgreSQL/SQL Server)
- 连接地址(host:port)、数据库名称
- 备份类型(全量备份/增量备份)、目标路径
- 状态(已启用/正常)、运行状态(空闲/备份中/已失效/已暂停)
- 调度策略、创建时间、上次执行时间、上次执行结果(成功/失败)
- 总执行次数、成功次数、失败次数
- 数据库大小、备份文件大小、保留天数

请以友好的中文格式输出，每个任务用分隔线隔开，信息清晰易读。` + "\n\n生成3个数据库备份任务，不要使用markdown代码块，纯文本输出。"

func QueryDatabaseBackup(page, pageSize int) error {
	prompt := fmt.Sprintf(
		`请生成 %d 个数据库备份任务信息（第 %d 页，每页 %d 条）。
要求包含不同的数据库类型(MySQL、PostgreSQL、SQL Server)，
不同的备份类型(全量备份、增量备份)。
至少有一个任务是"备份中"状态。`,
		pageSize, page, pageSize,
	)

	result, err := chatWithDeepSeek(dbQuerySystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("调用 DeepSeek 失败: %w", err)
	}

	fmt.Println()
	printWithReasoning(result)
	return nil
}
