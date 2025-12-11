package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/persistence"
	"github.com/wentf9/MyGoFileHub/internal/interface/api"

	// !!! 极其重要：必须引入 driver 包以触发 init() 注册
	_ "github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers/local"
)

func main() {
	// 1. 准备测试目录
	_ = os.MkdirAll("./test_storage", 0755)

	// 2. 初始化数据库
	db, err := persistence.InitDB("gofile.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// --- 自动插入测试数据 (种子数据) ---
	// 检查是否有数据，没有则插入一个本地源
	var count int64
	db.Model(&model.StorageSource{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding test data...")
		testSource := model.StorageSource{
			Name: "Local Test Disk",
			Type: "local",
			Config: model.JSONMap{
				"root_path": "./test_storage",
			},
		}
		db.Create(&testSource)
		fmt.Printf("Created seed storage source with ID: %d\n", testSource.ID)
	}
	// --------------------------------

	// 3. 初始化 Repository
	sourceRepo := persistence.NewSourceRepository(db)
	userRepo := persistence.NewUserRepository(db)
	permRepo := persistence.NewPermissionRepository(db)

	// 4. 初始化 Service (注入 Repo)
	permService := application.NewPermissionService(permRepo, userRepo)
	fileService := application.NewFileService(sourceRepo, permService)
	authService := application.NewAuthService(userRepo)

	// --- Seeding: 创建默认管理员 ---
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		fmt.Println("Creating default admin user...")
		// 这里的 Register 方法会自动 Hash 密码
		err := authService.Register(context.Background(), "admin", "admin123", "admin")
		if err != nil {
			log.Printf("Failed to create admin: %v", err)
		} else {
			fmt.Println("Admin created. User: admin, Pass: admin123")
		}
	}

	// 5. 初始化 Router
	r := api.InitRouter(fileService, authService)

	// 6. 启动
	fmt.Println("Server starting on :8080...")
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
