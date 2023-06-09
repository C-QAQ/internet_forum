package main

import (
	"fmt"
	"internet_forum/controller"
	"internet_forum/dao/mysql"
	"internet_forum/dao/redis"
	"internet_forum/logger"
	"internet_forum/pkg/snowflake"
	"internet_forum/router"
	"internet_forum/setting"
)

// @title GoWeb论坛
// @version 1.0
// @description 使用gin mysql redis 实现的web论坛

// @contact.name ShaoChong
// @contact.url coder.cc

// @host 127.0.0.1:8081
// @BasePath /
func main() {
	// 加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err: %v\n", err)
		return
	}

	// 注册路由
	r := router.Setup(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
