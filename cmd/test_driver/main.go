package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	// 必须引入 local 包以触发 init() 注册驱动
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"
	_ "github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers/local"
)

func main() {
	ctx := context.Background()

	// 1. 模拟从数据库取出的配置
	config := map[string]interface{}{
		"root_path": "./test_root", // 程序运行目录下的 test_root 文件夹
	}

	// 2. 通过工厂获取驱动 (解耦！)
	// 注意：我们不需要直接调用 local.NewLocalDriver()，而是通过字符串 "local" 获取
	driver, err := drivers.CreateInstance("local")
	if err != nil {
		log.Fatalf("无法创建驱动: %v", err)
	}

	// 3. 初始化
	if err := driver.Init(ctx, config); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	fmt.Println("Local Driver 初始化成功！")

	// 4. 测试：创建一个文件
	content := "Hello, Go File System Architecture!"
	err = driver.Create(ctx, "hello.txt", strings.NewReader(content), int64(len(content)))
	if err != nil {
		log.Fatalf("创建文件失败: %v", err)
	}
	fmt.Println("文件 hello.txt 创建成功")

	// 5. 测试：列出文件
	files, err := driver.List(ctx, "/")
	if err != nil {
		log.Fatalf("列出文件失败: %v", err)
	}

	fmt.Println("--- 文件列表 ---")
	for _, f := range files {
		fmt.Printf("[%s] %s (%d bytes)\n",
			map[bool]string{true: "DIR", false: "FILE"}[f.IsDir],
			f.Name,
			f.Size)
	}

	// 6. 安全测试 (可选)
	// err = driver.List(ctx, "../../")
	// 应该报错 access denied
}
