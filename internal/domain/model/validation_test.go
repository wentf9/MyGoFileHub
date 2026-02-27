package model

import (
	"strings"
	"testing"
)

// genStr 生成指定长度的字符串（用于测试）
func genStr(n int) string {
	return strings.Repeat("a", n)
}

func TestValidateStorageKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr error
	}{
		// 有效的 key
		{"valid_simple", "my-storage", nil},
		{"valid_with_underscore", "my_storage", nil},
		{"valid_with_number", "storage123", nil},
		{"valid_single_char", "a", nil},
		{"valid_max_length", "abcdefghijklmnopqrstuvwxyz123456", nil}, // 32 chars

		// 无效的 key - 空值
		{"empty", "", ErrEmptyStorageKey},

		// 无效的 key - 保留字
		{"reserved_api", "api", ErrReservedStorageKey},
		{"reserved_webdav", "webdav", ErrReservedStorageKey},
		{"reserved_ui", "ui", ErrReservedStorageKey},
		{"reserved_cp", "cp", ErrReservedStorageKey},
		{"reserved_mv", "mv", ErrReservedStorageKey},
		{"reserved_static", "static", ErrReservedStorageKey},
		{"reserved_with_at", "@api", ErrReservedStorageKey},

		// 无效的 key - 格式错误
		{"uppercase", "MyStorage", ErrInvalidStorageKeyFormat},
		{"number_start", "123storage", ErrInvalidStorageKeyFormat},
		{"underscore_start", "_storage", ErrInvalidStorageKeyFormat},
		{"hyphen_start", "-storage", ErrInvalidStorageKeyFormat},
		{"special_chars", "storage@123", ErrInvalidStorageKeyFormat},
		{"space", "my storage", ErrInvalidStorageKeyFormat},
		{"too_long", "abcdefghijklmnopqrstuvwxyz1234567", ErrInvalidStorageKeyFormat}, // 33 chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStorageKey(tt.key)
			if err != tt.wantErr {
				t.Errorf("ValidateStorageKey(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileName(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr error
	}{
		// 有效的文件名
		{"valid_simple", "file.txt", nil},
		{"valid_with_space", "my file.txt", nil},
		{"valid_chinese", "文件.txt", nil},
		{"valid_long", genStr(255), nil},

		// 无效的文件名 - 空值/特殊值
		{"empty", "", ErrEmptyFileName},
		{"dot", ".", ErrEmptyFileName},
		{"double_dot", "..", ErrEmptyFileName},

		// 无效的文件名 - 包含非法字符
		{"slash", "file/name", ErrInvalidFileName},
		{"backslash", "file\\name", ErrInvalidFileName},

		// 无效的文件名 - 过长
		{"too_long", genStr(256), ErrFileNameTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileName(tt.file)
			if err != tt.wantErr {
				t.Errorf("ValidateFileName(%q) error = %v, wantErr %v", tt.file, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		// 有效的路径
		{"empty", "", nil},
		{"root", "/", nil},
		{"valid_simple", "/files", nil},
		{"valid_nested", "/files/photos/image.jpg", nil},

		// 无效的路径 - 路径遍历
		{"parent_dir", "../etc/passwd", ErrInvalidPath},
		{"parent_dir_in_path", "/files/../etc/passwd", ErrInvalidPath},
		{"double_parent", "../../etc/passwd", ErrInvalidPath},

		// 无效的路径 - 不是绝对路径
		{"relative", "files/photos", ErrInvalidPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if err != tt.wantErr {
				t.Errorf("ValidatePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
