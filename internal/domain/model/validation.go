package model

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrEmptyStorageKey 存储源 key 为空
	ErrEmptyStorageKey = errors.New("storage key cannot be empty")
	// ErrReservedStorageKey 存储源 key 使用保留字
	ErrReservedStorageKey = errors.New("storage key cannot be a reserved word")
	// ErrInvalidStorageKeyFormat 存储源 key 格式错误
	ErrInvalidStorageKeyFormat = errors.New("storage key must start with a lowercase letter, followed by lowercase letters, digits, underscores, or hyphens (1-32 chars)")
	// ErrEmptyFileName 文件名为空
	ErrEmptyFileName = errors.New("file name cannot be empty")
	// ErrInvalidFileName 文件名包含非法字符
	ErrInvalidFileName = errors.New("file name cannot contain path separators, null bytes, or be '.' or '..'")
	// ErrFileNameTooLong 文件名过长
	ErrFileNameTooLong = errors.New("file name too long (max 255 characters)")
	// ErrInvalidPath 路径格式错误
	ErrInvalidPath = errors.New("path cannot contain '..' or must be absolute")
)

// 合法 key 模式：小写字母开头，后跟小写字母/数字/下划线/中划线，1-32 字符
var storageKeyRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)

// 保留字黑名单（与系统路由冲突的 key）
var reservedKeys = map[string]bool{
	// API 路由
	"api": true, "@api": true,
	// 文件操作
	"cp": true, "@cp": true,
	"mv": true, "@mv": true,
	// 服务路由
	"webdav": true,
	// 前端静态文件
	"ui": true,
	// 其他系统路径
	"static": true,
}

// ValidateStorageKey 验证存储源 key 的合法性
// 规则：
// 1. 只能包含小写字母、数字、下划线、中划线
// 2. 必须以字母开头
// 3. 长度 1-32 字符
// 4. 不能使用保留字
func ValidateStorageKey(key string) error {
	if key == "" {
		return ErrEmptyStorageKey
	}
	// 检查保留字（包括带 @ 前缀的）
	if reservedKeys[key] || reservedKeys["@"+key] || reservedKeys[strings.TrimPrefix(key, "@")] {
		return ErrReservedStorageKey
	}
	// 检查格式
	if !storageKeyRegex.MatchString(key) {
		return ErrInvalidStorageKeyFormat
	}
	return nil
}

// ValidateFileName 验证文件名的合法性
// 规则：
// 1. 不能为空、. 或 ..
// 2. 不能包含路径分隔符（/ \）
// 3. 不能包含 \0 等控制字符
// 4. 长度不超过 255
func ValidateFileName(name string) error {
	if name == "" || name == "." || name == ".." {
		return ErrEmptyFileName
	}
	// 检查路径分隔符
	if strings.ContainsAny(name, `/\`) {
		return ErrInvalidFileName
	}
	// 检查控制字符（包括 \0）
	for _, r := range name {
		if r < 32 {
			return ErrInvalidFileName
		}
	}
	// 检查长度
	if len(name) > 255 {
		return ErrFileNameTooLong
	}
	return nil
}

// ValidatePath 验证文件路径的合法性
// 规则：
// 1. 必须以 / 开头（绝对路径）
// 2. 不能包含 .. （路径遍历）
func ValidatePath(path string) error {
	if path == "" {
		return nil // 空路径表示根目录，允许
	}
	// 检查是否以 / 开头
	if !strings.HasPrefix(path, "/") {
		return ErrInvalidPath
	}
	// 检查路径遍历
	if strings.Contains(path, "..") {
		return ErrInvalidPath
	}
	return nil
}
