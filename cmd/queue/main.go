package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"app/internal/core/services"

	"github.com/spf13/viper"
)

var (
	// 命令行参数
	configFile string
	queueName  string
	listQueues bool
	clearQueue bool
	stopQueue  bool
	startQueue bool
)

func init() {
	// 注册命令行参数
	flag.StringVar(&configFile, "config", "configs/config.yaml", "配置文件路径")
	flag.StringVar(&queueName, "queue", "", "队列名称")
	flag.BoolVar(&listQueues, "list", false, "列出所有队列")
	flag.BoolVar(&clearQueue, "clear", false, "清空队列")
	flag.BoolVar(&stopQueue, "stop", false, "停止队列")
	flag.BoolVar(&startQueue, "start", false, "启动队列")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 加载配置
	if err := loadConfig(configFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建队列服务
	queueService, err := services.NewQueueService()
	if err != nil {
		log.Fatalf("Failed to create queue service: %v", err)
	}

	// 处理命令
	switch {
	case listQueues:
		// 列出所有队列
		queues := queueService.GetActiveQueues()
		if len(queues) == 0 {
			fmt.Println("No active queues")
			return
		}

		fmt.Println("Active queues:")
		for _, name := range queues {
			fmt.Printf("- %s\n", name)
		}

	case clearQueue:
		// 清空队列
		if queueName == "" {
			log.Fatal("Queue name is required")
		}

		ctx := context.Background()
		if err := queueService.Clear(ctx, queueName); err != nil {
			log.Fatalf("Failed to clear queue %s: %v", queueName, err)
		}
		fmt.Printf("Queue %s cleared\n", queueName)

	case stopQueue:
		// 停止队列
		if queueName == "" {
			// 停止所有队列
			queueService.Stop()
			fmt.Println("All queues stopped")
			return
		}

		// TODO: 实现停止指定队列的功能
		fmt.Printf("Stopping queue %s...\n", queueName)

	case startQueue:
		// 启动队列
		if err := queueService.Start(); err != nil {
			log.Fatalf("Failed to start queue service: %v", err)
		}

		// 等待信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		fmt.Println("Queue service started, press Ctrl+C to stop")
		<-sigChan

		// 停止服务
		queueService.Stop()
		fmt.Println("Queue service stopped")

	default:
		// 显示帮助信息
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
}

// loadConfig 加载配置
func loadConfig(configFile string) error {
	// 设置配置文件
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	return nil
}
