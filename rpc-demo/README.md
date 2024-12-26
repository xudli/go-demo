# 自定义 RPC 框架实现

一个基于 TCP 协议的轻量级 RPC 框架实现，使用 JSON 作为序列化方案。

## 项目结构

``` yaml
rpc-demo/
├── main.go          # 服务端实现
├── client/
│   └── client.go    # 客户端实现
└── go.mod          # Go 模块文件
```

## 核心组件

### 1. 消息协议

``` golang
// 请求消息
type Request struct {
    ServiceMethod string      `json:"method"`    // 服务方法名
    Args         interface{} `json:"args"`      // 参数
    RequestID    uint64      `json:"request_id"` // 请求ID
}

// 响应消息
type Response struct {
    RequestID uint64      `json:"request_id"` // 对应请求ID
    Error    string      `json:"error"`      // 错误信息
    Result   interface{} `json:"result"`     // 返回结果
}
```

### 2. 服务端

- **Server 结构**：管理服务注册和请求处理
- **服务注册**：支持动态服务注册
- **并发处理**：每个客户端连接使用独立的 goroutine 处理

### 3. 客户端

- **Client 结构**：管理与服务端的连接
- **请求发送**：处理请求的编码和发送
- **响应接收**：处理响应的接收和解码

## 实现原理

1. **服务注册**

   ``` golang
   server := NewServer()
   server.Register("ServiceName.MethodName", service)
   ```

2. **请求处理流程**
   - 接收客户端连接
   - 解码请求消息
   - 查找对应服务
   - 执行方法调用
   - 返回响应结果

3. **远程调用流程**
   - 创建请求对象
   - 发送到服务端
   - 等待响应
   - 解析返回结果

## 使用示例

### 服务端

```golang
func main() {
    server := NewServer()
    arith := new(ArithService)
    server.Register("Arith.Multiply", arith)
    
    listener, _ := net.Listen("tcp", ":1234")
    // 处理连接...
}
```

### 客户端

```golang
func main() {
    client, _ := NewClient("localhost:1234")
    args := Args{A: 7, B: 8}
    result, _ := client.Call("Arith.Multiply", args)
    log.Printf("结果: %v", result)
}
```

## 特点

- ✅ 简单轻量：核心实现代码简洁
- ✅ 易扩展：支持注册任意服务
- ✅ 类型安全：基于 Go 的类型系统
- ✅ 并发处理：支持多客户端同时连接

## 局限性

- 仅支持 JSON 序列化
- 无连接池机制
- 缺少超时处理
- 无服务发现功能

## 后续优化方向

1. **性能优化**
   - 实现连���池
   - 使用更高效的序列化方式
   - 添加压缩支持

2. **功能增强**
   - 添加超时机制
   - 实现重试逻辑
   - 支持异步调用

3. **可靠性**
   - 添加心跳检测
   - 实现负载均衡
   - 支持服务发现

4. **安全性**
   - 添加认证机制
   - 实现加密传输

## 运行说明

1. 启动服务端：

```bash
go run main.go
```

2. 运行客户端：

```bash
go run client/client.go
```

## 许可证

MIT License

## 为什么选择 Go 实现

相比 Java，Go 语言在实现 RPC 框架时具有以下优势：

### 1. 并发处理

```golang
// Go 的 goroutine 非常轻量级
go server.handleConnection(conn)  // 几行代码实现并发
```

```java
// Java 需要显式创建线程池
ExecutorService executor = Executors.newCachedThreadPool();
executor.submit(() -> handleConnection(conn));
```

### 2. 接口实现

```golang
// Go 隐式实现接口
type Service interface {
    Method(args Args) Result
}

type MyService struct{}
func (s *MyService) Method(args Args) Result { ... }  // 自动实现接口
```

```java
// Java 需要显式声明接口实现
public interface Service {
    Result method(Args args);
}

public class MyService implements Service {
    @Override
    public Result method(Args args) { ... }
}
```

### 3. JSON 处理

```golang
// Go 的 JSON 处理简洁
json.NewEncoder(conn).Encode(request)
json.NewDecoder(conn).Decode(&response)
```

```java
// Java 需要更多样板代码
ObjectMapper mapper = new ObjectMapper();
mapper.writeValue(outputStream, request);
Response response = mapper.readValue(inputStream, Response.class);
```

### 4. 错误处理

```golang
// Go 的错误处理直接
if err != nil {
    return err
}
```

```java
// Java 的异常处理冗长
try {
    // 业务逻辑
} catch (IOException e) {
    throw new RPCException("Failed to process", e);
} finally {
    // 清理资源
}
```

### 5. 网络编程

```golang
// Go 的网络编程API简洁
listener, _ := net.Listen("tcp", ":1234")
conn, _ := listener.Accept()
```

```java
// Java 需要更多代码
ServerSocket serverSocket = new ServerSocket(1234);
Socket socket = serverSocket.accept();
InputStream in = socket.getInputStream();
OutputStream out = socket.getOutputStream();
```

Go 语言的这些特性使得 RPC 框架的实现更加简洁：
- 内置并发支持
- 隐式接口实现
- 简洁的标准库
- 直接的错误处理
- 轻量级的反射机制

虽然 Java 提供了更严格的类型系统和更好的可扩展性，但也带来了更多的样板代码。Go 的设计理念更适合构建这样的网络服务框架。