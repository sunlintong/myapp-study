package main

import (
//	"os"
//	"fmt"
	"github.com/golang/glog"
	_ "myapp/routers"
	"github.com/astaxie/beego"

//	"github.com/golang/glog"
//	"log"
	"net"
	"time"
//	"os"
	"flag"

//	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "echotest/echotest"
	"google.golang.org/grpc/reflection"

)

const (
	port = ":50051"
)

type server struct{}


func(s *server) Echo(stream pb.Echo_EchoServer) error{
	
	//服务器并发线程主动地每隔5秒将时间加入流并发送
	go func(){
		ticker := time.NewTicker(5 * time.Second)
		for _ = range ticker.C{
			nowtime:=time.Now().Format("2006-01-02 15:04:05")
			echoreply := &pb.EchoReply{}
			echoreply.Nowtime = nowtime	
//			log.Println(echoreply)
			stream.Send(echoreply)
			glog.Infoln("send a time message:"+echoreply.Nowtime)
		}
}()
	//将echo回应加入流发送
	for true{
		echorequest,err := stream.Recv()
		//如果不判断echorequest是否为空，会出现如下bug：
		//当客户端断开与服务器的连接后，服务器会停止运行
		//因为如果没消息进来，echoreply中两个参数都是空字符串，echoreply也就是空指针了，不能发送空指针
		if echorequest == nil || err !=nil {
			break
		}
		echoreply := &pb.EchoReply{}
		echoreply.Output = "echo:"+echorequest.Input
		stream.Send(echoreply)
		glog.Infoln("send a echo message:"+echoreply.Output)
	}
	return nil	
}


func startEchoServer(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()
//	fmt.Println(os.Getenv("GO15VENDOREXPERIMENT"))
	//os.Getenv检索环境变量并返回值，如果变量是不存在的，这将是空的。
//	appname := os.Getenv("OEM")
//	version := os.Getenv("VER")
//	HOME := os.Getenv("HOME")
//	log.Println(appname,version,HOME,"wrg")
	//开启运行在端口：50051上的startEchoServer线程
	go startEchoServer()
	beego.Run()
}