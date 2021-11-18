package main

import (
	"flag"
	"fmt"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("i", "0.0.0.0", "ip地址")
	Port := flag.Int("p", 50051, "端口号")

	flag.Parse()
	fmt.Println("IP:", *IP)
	fmt.Println("Port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(listen)
	if err != nil {
		panic("failed to start gRPC:" + err.Error())
	}
}
