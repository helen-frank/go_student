# 1. quick start

## 1.1 install

Make sure that your `$GOBIN` is in your `$PATH`.

Protobuf&go插件

```shell
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Grpc-gatway

```shell
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```

## 1.2 demo

```shell
❯ tree
.
├── client
│   └── main.go
├── go.mod
├── go.sum
├── proto
│   └── service
│       ├── service.pb.go
│       ├── service.proto
│       └── service_grpc.pb.go
└── server
    └── main.go

5 directories, 7 files
```

`service.proto`

```protobuf
syntax = "proto3";

option go_package = ".;service";

service Call {
  rpc Call(CallRequest) returns (CallResponse) {}
}

message CallRequest {
  string name = 1;
}

message CallResponse {
  string response_msg = 1;
}
```

生成代码

```shell
protoc --go_out=./proto/service/ ./proto/service/service.proto
protoc --go-grpc_out=./proto/service/ ./proto/service/service.proto
```

`server/main.go`

```go
/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2023-02-01 14:50:48
*/

package main

import (
	"context"
	"net"

	"testproto/proto/service"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

type server struct {
	service.UnimplementedCallServer
}

func (s *server) Call(ctx context.Context, req *service.CallRequest) (*service.CallResponse, error) {
	return &service.CallResponse{ResponseMsg: "hello " + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	klog.Info("listen start")

	grpcServer := grpc.NewServer()

	service.RegisterCallServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		klog.Error("run grpc server failed, err: ", err)
		panic(err)
	}
}

```

`client/main.go`

```go
/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2023-02-01 14:50:41
*/

package main

import (
	"context"
	"testproto/proto/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := service.NewCallClient(conn)

	resp, err := client.Call(context.TODO(), &service.CallRequest{Name: "helen"})
	if err != nil {
		panic(err)
	}

	klog.Info(resp.GetResponseMsg())
}

```

## 1.3 ssl

```shell
brew install openssl
```

|      |                                                              |
| ---- | ------------------------------------------------------------ |
| key  | 服务器上的私钥文件，用于对发送给客户端数据的加密，以及对从客户端接收到数据的解密 |
| csr  | 证书签名请求文件，用于提交给证书颁发机构(CA)对证书签名       |
| crt  | 由证书颁发机构(CA)签名后的证书，或者是开发者自签的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息 |
| pem  | 是基于Base64编码的证书格式，扩展名包括PEM，CRT和CER          |

### 1.3.1 生成证书

```shell
# 生成私钥
openssl genrsa -out ./key/server.key 2048

# 生成证书
openssl req -new -x509 -key ./key/server.key -out ./key/server.crt -days 36500
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:CN
State or Province Name (full name) []:ShangHai
Locality Name (eg, city) []:ShangHai
Organization Name (eg, company) []:Helen
Organizational Unit Name (eg, section) []:Go
Common Name (eg, fully qualified host name) []:helen
Email Address []:helenfrank@protonmail.com

# 生成 csr
openssl req -new -key ./key/server.key -out ./key/server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:CN
State or Province Name (full name) []:ShangHai
Locality Name (eg, city) []:ShangHai
Organization Name (eg, company) []:Helen 
Organizational Unit Name (eg, section) []:go
Common Name (eg, fully qualified host name) []:helen
Email Address []:helenfrank@protonmail.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:

```

# 2. grpc

## 2.1 quick start



