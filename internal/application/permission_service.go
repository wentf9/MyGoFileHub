package application

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
)

type PermissionService struct {
	permRepo repository.PermissionRepository
	userRepo repository.UserRepository
}

var permissionCache = sync.Map{} // 用于存储用户的权限数据

func NewPermissionService(pRepo repository.PermissionRepository, uRepo repository.UserRepository) *PermissionService {
	return &PermissionService{permRepo: pRepo, userRepo: uRepo}
}

// CheckPermission 检查细粒度权限
// action: "read" 或 "write"
func (s *PermissionService) CheckPermission(ctx context.Context, username string, sourceID uint, path string, action string) bool {
	// 1. 获取用户信息
	var user *model.User
	var err error
	if vaule, ok := userCache.Load(username); ok {
		user = vaule.(*model.User)
	} else {
		// 查询用户
		user, err = s.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return false
		}
		value, loaded := userCache.LoadOrStore(username, user)
		if loaded {
			user = value.(*model.User)
		}
	}

	// 2. 超级管理员拥有所有权限
	if user.Role == "admin" {
		return true
	}

	// 3. 获取该用户在该源下的所有权限规则

	var perms []*model.UserPermission
	if vaule, ok := permissionCache.Load(username + "_" + strconv.FormatUint(uint64(sourceID), 10)); ok {
		perms = vaule.([]*model.UserPermission)
	} else {
		perms, err = s.permRepo.FindByUserAndSource(ctx, user.ID, sourceID)
		if err != nil || len(perms) == 0 {
			permissionCache.Store(username+"_"+strconv.FormatUint(uint64(sourceID), 10), []*model.UserPermission{})
			return false
		}
		value, loaded := permissionCache.LoadOrStore(username+"_"+strconv.FormatUint(uint64(sourceID), 10), perms)
		if loaded {
			perms = value.([]*model.UserPermission)
		}
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
