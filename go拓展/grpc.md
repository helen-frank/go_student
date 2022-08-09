https://colobu.com/2019/10/03/protobuf-ultimate-tutorial-in-go/

# 1. Quick Start

## 1.1 Install

### 1.1.1 Mac(m1)

```shell
brew install protobuf protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go

protoc --version
protoc-gen-go --version
protoc-gen-go-grpc --version
```

Vscode代码格式化

`vscode-proto`插件，插件自带高亮, 快捷键补全等等功能,

格式化功能依赖`clang-format`, `brew install clang-format`



## 1.2 example

项目结构

```shell
❯ tree
.
├── client
│   └── client.go
├── go.mod
├── go.sum
├── pb
│   ├── hello.pb.go
│   ├── hello.proto
│   └── hello_grpc.pb.go
└── server
    └── server.go

3 directories, 7 files
```

`pb/hello.pb`

```protobuf
syntax = "proto3"; // 版本声明，使用Protocol Buffers v3版本

option go_package = "pb/;pb"; // 指定go package名称

package pb; // 包名

// 定义服务
service Greeter {
  // SayHello 方法
  rpc SayHello(HelloRequest) returns (HelloResponse) {}
}

// 请求消息
message HelloRequest { string name = 1; }

// 响应消息
message HelloResponse { string reply = 1; }

```

生成go中间代码

```shell
protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./pb/hello.proto
```

`server/server.go`

```go
/*
	@file	server.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2022-07-04 16:04:42
*/

package main

import (
	"2/pb"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// hello server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println(in.Name, "call")
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	fmt.Println("start server, Listen: ", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
```

# 2. grpc流式

## 2.1 服务端流式rpc

`pb/hello.proto`

```protobuf
// 服务端返回流式数据
   rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
```

重新生成中间代码

`client/client.go`

```go
/*
	@file	client.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2022-07-04 16:04:38
*/

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"2/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// hello_client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func runLotsOfReplies(c pb.GreeterClient) {
	// server端流式RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v", err)
	}
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v", err)
		}
		log.Printf("got reply: %q\n", res.GetReply())
	}
}
func main() {
	flag.Parse()
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.GetReply())

	runLotsOfReplies(c)
}
```

`server/server.go`

```go
/*
	@file	server.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2022-07-04 16:04:42
*/

package main

import (
	"2/pb"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// hello server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println(in.Name, "call")
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

// LotsOfReplies 返回使用多种语言打招呼
func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好 ",
		"hello ",
		"こんにちは ",
		"여보세요 ",
	}
	fmt.Println(in.Name, "call")
	for i := range words {
		data := &pb.HelloResponse{
			Reply: words[i] + in.GetName(),
		}
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	fmt.Println("start server, Listen: ", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}

```

## 2.2 客户端流式rpc

`pb/hello.proto`

```protobuf
   // 客户端发送流式数据
   rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
```

`client/client.go`

```go
func runLotsOfGreeting(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端流式RPC
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed, err: %v", err)
	}
	names := []string{"helen", "1", "2", "3"}
	for _, name := range names {
		// 发送流式数据
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed: %v", err)
	}
	log.Printf("got reply: %v", res.GetReply())
}
```

`server/server.go`

```go
// LotsOfGreetings 接收流式数据
func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "你好："
	for {
		// 接收客户端发来的流式数据
		res, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			fmt.Println(reply)
			return stream.SendAndClose(&pb.HelloResponse{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += res.GetName() + " "
	}
}
```

## 2.3 双向流式RPC

`pb/hello.pb`

```protobuf
   // 双向流式数据
   rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
```

`client/client.go`

```go
func runBidiHello(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// 双向流模式
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed, err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI：%s\n", in.GetReply())
		}
	}()
	// 从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') // 读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
```



`server/server.go`

```go
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := magic(in.GetName()) // 对收到的数据做些处理

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "?", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}
```

# 3. metadata

类似于 HTTP 请求中的 Cookie 数据，元数据（[metadata](https://pkg.go.dev/google.golang.org/grpc/internal/metadata)）是关于特定 RPC 调用的信息（例如身份验证详细信息），采用键值对列表的形式，其中键是字符串，值通常是字符串，但也可以是二进制数据。 元数据对 gRPC 本身是不透明的——它允许客户端提供与服务器调用相关的信息，反之亦然





