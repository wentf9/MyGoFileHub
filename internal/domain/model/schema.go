package model

// ConfigItem 描述存储源的一个配置选项
type ConfigItem struct {
	Name        string `json:"name"`        // 字段键名 (例如: root_path, host)
	Label       string `json:"label"`       // 显示名称 (例如: Root Path, 主机地址)
	Type        string `json:"type"`        // 类型: string, number, password, boolean
	Required    bool   `json:"required"`    // 是否必填
	Description string `json:"description"` // 帮助说明文字
	Default     any    `json:"default"`     // 默认值
}

// StorageDriverSchema 描述一个存储驱动及其需要的配置
type StorageDriverSchema struct {
	Type   string       `json:"type"`   // 驱动类型: local, smb
	Name   string       `json:"name"`   // 驱动显示名: Local Folder, SMB Share
	Config []ConfigItem `json:"config"` // 配置项数组
}
