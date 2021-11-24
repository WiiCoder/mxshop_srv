package main

import (
	"flag"
	"fmt"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"net"
)

func main() {
	IP := flag.String("i", "127.0.0.1", "ip地址")
	Port := flag.Int("p", 50000, "端口号")
	Env := flag.String("e", "dev", "运行环境")

	initialize.InitLogger()
	initialize.InitConfig(*Env)
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("IP:", *IP)
	zap.S().Info("Port:", *Port)
	zap.S().Info("Env:", *Env)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	logPanic(err)

	// 生成健康检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", *IP, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",     // 15s 服务不可达时，注销服务
		Status:                         "passing", // 服务启动时，默认正常
	}

	// 生成注册对象
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceID,
		Port:    *Port,
		Tags:    []string{"user", "srv"},
		Address: *IP,
		Check:   check,
	}

	err = client.Agent().ServiceRegister(registration)
	logPanic(err)

	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("failed to start gRPC:" + err.Error())
		}
	}()

	// 接收终止服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Infof("[consul] 注销失败: %s", serviceID)
	}
	zap.S().Infof("[consul] 注销成功: %s", serviceID)

}

func logPanic(err error) {
	if err != nil {
		zap.S().Panic(err)
	}
}
