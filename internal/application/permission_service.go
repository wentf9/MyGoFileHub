package application

import (
	"context"
	"strings"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
)

type PermissionService struct {
	permRepo repository.PermissionRepository
	userRepo repository.UserRepository
}

func NewPermissionService(pRepo repository.PermissionRepository, uRepo repository.UserRepository) *PermissionService {
	return &PermissionService{permRepo: pRepo, userRepo: uRepo}
}

// CheckPermission 检查细粒度权限
// action: "read" 或 "write"
func (s *PermissionService) CheckPermission(ctx context.Context, userID, sourceID uint, path string, action string) bool {
	// 1. 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false
	}

	// 2. 超级管理员拥有所有权限
	if user.Role == "admin" {
		return true
	}

	// 3. 获取该用户在该源下的所有权限规则
	// 优化点：这里可以在内存做缓存 (Redis/MemoryCache)，避免每次文件操作都查库
	perms, err := s.permRepo.FindByUserAndSource(ctx, userID, sourceID)
	if err != nil || len(perms) == 0 {
		return false // 默认拒绝 (白名单模式)
	}

	// 4. 最长前缀匹配算法
	var bestMatch *model.UserPermission
	maxLen := -1

	// 确保路径以 / 开头，且统一格式
	checkPath := cleanPath(path)

	for _, p := range perms {
		prefix := cleanPath(p.PathPrefix)

		// 检查 path 是否以 prefix 开头
		if strings.HasPrefix(checkPath, prefix) {
			// 找到匹配，看是否是更长的匹配
			if len(prefix) > maxLen {
				maxLen = len(prefix)
				bestMatch = p
			}
		}
	}

	// 5. 如果没有匹配的规则，默认拒绝
	if bestMatch == nil {
		return false
	}

	// 6. 根据动作检查字段
	if action == "write" {
		return bestMatch.AllowWrite
	}
	return bestMatch.AllowRead
}

// cleanPath 简单的路径标准化
func cleanPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	// 确保目录匹配的准确性，比如 /work 不应该匹配 /working
	// 这里的逻辑可以根据需求搞更复杂，暂时简单处理
	return p
}
