# README

## go mod init

`go mod init github.com/1005281342/trpc-demo`

## 环境安装

安装 [官方命令行工具](https://github.com/trpc-group/trpc-cmdline/blob/main/README.zh_CN.md) 文档

通过 `go install trpc.group/trpc-go/trpc-cmdline/trpc@latest` 安装 `trpc` 命令行工具

[使用 trpc setup 一键安装所有依赖](https://github.com/trpc-group/trpc-cmdline/blob/main/README.zh_CN.md#%E4%BD%BF%E7%94%A8-trpc-setup-%E4%B8%80%E9%94%AE%E5%AE%89%E8%A3%85%E6%89%80%E6%9C%89%E4%BE%9D%E8%B5%96) 
可能会因为一些网络问题导致设置失败，可以通过如下方式手动配置
- 下载官方提供的二进制文件，如 mac 系统的 https://github.com/trpc-group/trpc-cmdline/releases/tag/v0.0.1-darwin
- 添加可执行权限 `chmod +x 工具二进制文件`
- 将这些工具拷贝到 `go/bin` 目录下，并通过 `which 工具二进制文件` 确定这些工具可以被找到

## 开始编写一个 tRPC 服务

### 协议设计

参考官方提供的示例协议 https://github.com/trpc-group/trpc-cmdline/blob/main/docs/helloworld/helloworld.proto

创建 `helloworld.proto` 文件并添加如下协议内容

```protobuf
syntax = "proto3";
package helloworld;

option go_package = "github.com/1005281342/trpc-demo/helloworld";

// HelloRequest is hello request.
message HelloRequest {
    string msg = 1;
}

// HelloResponse is hello response.
message HelloResponse {
    string msg = 1;
}

// HelloWorldService handles hello request and echo message.
service HelloWorldService {
    // Hello says hello.
    rpc Hello(HelloRequest) returns(HelloResponse);
}
```

### 生成项目 

`trpc create -p helloworld.proto -o .`

### 启动服务

`go run .`

```
2023-10-21 15:52:32.359 DEBUG   maxprocs/maxprocs.go:47 maxprocs: Leaving GOMAXPROCS=12: CPU quota undefined
2023-10-21 15:52:32.360 INFO    server/service.go:167   process:8575, trpc service:helloworld.HelloWorldService launch success, tcp:127.0.0.1:8000, serving ...
```

### 发起请求

`go run cmd/client/main.go`

客户端日志
```
2023-10-21 15:54:57.034 DEBUG   debuglog@v1.0.0/log.go:236      client request:/helloworld.HelloWorldService/Hello, cost:1.880357ms, to:127.0.0.1:8000
2023-10-21 15:54:57.034 DEBUG   client/main.go:25       simple  rpc   receive: 
```

服务端日志
```
2023-10-21 15:54:57.033 DEBUG   debuglog@v1.0.0/log.go:196      server request:/helloworld.HelloWorldService/Hello, cost:9.418µs, from:127.0.0.1:60328
```

## 业务编码

### 修改 Hello 接口

- 响应用户输入的 msg

    ```go
    // Hello Hello says hello.
    func (s *helloWorldServiceImpl) Hello(
        ctx context.Context,
        req *pb.HelloRequest,
    ) (*pb.HelloResponse, error) {
        rsp := &pb.HelloResponse{
            Msg: req.GetMsg(),
        }
        return rsp, nil
    }
    ```

- 重启服务

- 修改客户端 `cmd/client/main.go` 中 `callHelloWorldServiceHello()` 的请求内容
    
    ```go
    // Example usage of unary client.
    reply, err := proxy.Hello(ctx, &pb.HelloRequest{
        Msg: "hello tRPC",
    })
    ```
- 发送请求进行测试

    客户端日志
    ```
    2023-10-21 16:13:10.556 DEBUG   debuglog@v1.0.0/log.go:236      client request:/helloworld.HelloWorldService/Hello, cost:2.020233ms, to:127.0.0.1:8000
    2023-10-21 16:13:10.556 DEBUG   client/main.go:27       simple  rpc   receive: msg:"hello tRPC"
    ```        

    服务端日志
    ```
    2023-10-21 16:13:10.556 DEBUG   debuglog@v1.0.0/log.go:196      server request:/helloworld.HelloWorldService/Hello, cost:9.456µs, from:127.0.0.1:63248
    ```
  

### 新增 SayHi 功能

- 接口协议

  ```protobuf
  ...
  
  // SayHiReq is says hi request
  message SayHiReq {}
  
  // SayHiRsp is says hi response
  message SayHiRsp {
      string msg = 1;
  }
  
  // HelloWorldService handles hello request and echo message.
  service HelloWorldService {
      ...
  
      // SayHi says hi
      rpc SayHi(SayHiReq) returns (SayHiRsp);
  }
  ```
  
- 更新桩代码 `trpc create -p helloworld.proto -o ./stub/github.com/1005281342/trpc-demo/helloworld --rpconly`

- 实现 SayHi 接口

  ```go
  func (s *helloWorldServiceImpl) SayHi(ctx context.Context, req *pb.SayHiReq) (*pb.SayHiRsp, error) {
      return &pb.SayHiRsp{Msg: "hi"}, nil
  }
  ```

- 启动服务 `go mod tidy && go run .`

- 添加测试代码

  ```go
  func main() {
      // Load configuration following the logic in trpc.NewServer.
      cfg, err := trpc.LoadConfig(trpc.ServerConfigPath)
      if err != nil {
          panic("load config fail: " + err.Error())
      }
      trpc.SetGlobalConfig(cfg)
      if err := trpc.Setup(cfg); err != nil {
          panic("setup plugin fail: " + err.Error())
      }
      //callHelloWorldServiceHello()
      callHelloWorldServiceSayHi()
  }
  
  func callHelloWorldServiceSayHi() {
      proxy := pb.NewHelloWorldServiceClientProxy(
          client.WithTarget("ip://127.0.0.1:8000"),
          client.WithProtocol("trpc"),
      )
      ctx := trpc.BackgroundContext()
      // Example usage of unary client.
      reply, err := proxy.SayHi(ctx, &pb.SayHiReq{})
      if err != nil {
          log.Fatalf("err: %v", err)
      }
      log.Debugf("[SayHi] -- simple  rpc   receive: %+v", reply)
  }
  ```
  
- 发送请求进行测试 `go run cmd/client/main.go`

  客户端日志
  ```
  2023-10-21 16:47:08.698 DEBUG   debuglog@v1.0.0/log.go:236      client request:/helloworld.HelloWorldService/SayHi, cost:2.969017ms, to:127.0.0.1:8000
  2023-10-21 16:47:08.698 DEBUG   client/main.go:39       [SayHi] -- simple  rpc   receive: msg:"hi"
  ```
  
  服务端日志
  ```
  2023-10-21 16:47:08.697 DEBUG   debuglog@v1.0.0/log.go:196      server request:/helloworld.HelloWorldService/SayHi, cost:3.751µs, from:127.0.0.1:52401
  ```