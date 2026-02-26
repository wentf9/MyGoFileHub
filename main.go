package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/config"
	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	_ "github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers/local"
	_ "github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers/smb"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/persistence"
	"github.com/wentf9/MyGoFileHub/internal/interface/api"
)

// 嵌入前端构建产物 (生产模式)
// 路径相对于项目根目录
//go:embed all:frontend/dist
var embeddedFS embed.FS

func main() {
	// 检查 --version 或 -v 参数
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			application.PrintVersion()
			return
		}
	}

	// 1. 准备测试目录
	_ = os.MkdirAll(config.AppConfig.DataDir+"/test_storage", 0755)

	// 2. 初始化数据库
	db, err := persistence.InitDB(config.AppConfig.DataDir + "/gofile.db")
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
			Key:  "test",
			Type: "local",
			Config: model.JSONMap{
				"root_path": config.AppConfig.DataDir + "./test_storage",
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
	userService := application.NewUserService(userRepo)
	sourceService := application.NewSourceService(sourceRepo, fileService)

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

	// 5. 根据运行模式初始化 Router
	var r *gin.Engine
	var staticFS fs.FS

	if config.AppConfig.Mode == "dev" {
		// 开发模式：不嵌入静态文件，仅 API
		fmt.Println("[INFO] Running in development mode (API only)")
		fmt.Printf("[INFO] Frontend directory: %s\n", config.AppConfig.FrontendDir)
		staticFS = nil // nil 表示不服务静态文件
	} else {
		// 生产模式：使用嵌入式静态文件
		fmt.Println("[INFO] Running in production mode (embedded frontend)")
		// 获取 dist 目录的嵌入文件系统
		staticFS, _ = fs.Sub(embeddedFS, "frontend/dist")
	}

	// 初始化路由
	r = api.InitRouter(fileService, authService, userService, sourceService, staticFS)

	// 6. 启动
	fmt.Printf("Server starting on %s:%s\n", config.AppConfig.Listen, config.AppConfig.ServerPort)
	if config.AppConfig.Mode == "dev" {
		fmt.Println("\nDevelopment mode:")
		fmt.Println("  - API:  http://" + config.AppConfig.Listen + ":" + config.AppConfig.ServerPort)
		fmt.Println("  - Frontend: Run 'npm run dev' in frontend/my-go-file-hub-ui directory")
	} else {
		fmt.Println("\nProduction mode:")
		fmt.Println("  - Access: http://" + config.AppConfig.Listen + ":" + config.AppConfig.ServerPort)
	}
	fmt.Println()

	if err := r.Run(config.AppConfig.Listen + ":" + config.AppConfig.ServerPort); err != nil {
		panic(err)
	}
}
