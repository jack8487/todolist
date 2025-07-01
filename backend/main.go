package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"todolist/config"
	"todolist/docs"
	"todolist/internal/api"
	"todolist/internal/repository"
	"todolist/internal/service"
)

// @title TodoList API
// @version 1.0
// @description TodoList 项目的 API 文档
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// 加载配置
	err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	err = repository.InitDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建路由
	r := gin.Default()

	// 设置 Swagger 配置
	docs.SwaggerInfo.Title = "TodoList API"
	docs.SwaggerInfo.Description = "TodoList 项目的 API 文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 添加 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建仓储实例
	userRepo := repository.NewUserRepository(repository.DB)
	taskRepo := repository.NewTaskRepository(repository.DB)

	// 创建服务实例
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// 创建处理器实例
	userHandler := api.NewUserHandler(userService)
	taskHandler := api.NewTaskHandler(taskService)

	// 注册路由
	userHandler.RegisterRoutes(r)
	taskHandler.RegisterRoutes(r)

	// 启动服务器
	r.Run(":8080")
}
