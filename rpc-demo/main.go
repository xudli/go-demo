package main

import (
    "encoding/json"
    "log"
    "net"
)

// Request RPC请求结构
type Request struct {
    ServiceMethod string      `json:"method"`
    Args         interface{} `json:"args"`
    RequestID    uint64      `json:"request_id"`
}

// Response RPC响应结构
type Response struct {
    RequestID uint64      `json:"request_id"`
    Error    string      `json:"error"`
    Result   interface{} `json:"result"`
}

// ArithService 算术服务
type ArithService struct{}

// Args 计算参数
type Args struct {
    A, B int
}

// Multiply 乘法实现
func (a *ArithService) Multiply(args Args) int {
    return args.A * args.B
}

// Server RPC服务器
type Server struct {
    services map[string]interface{}
}

// NewServer 创建新的RPC服务器
func NewServer() *Server {
    return &Server{
        services: make(map[string]interface{}),
    }
}

// Register 注册服务
func (s *Server) Register(name string, service interface{}) {
    s.services[name] = service
}

// handleConnection 处理单个连接
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    decoder := json.NewDecoder(conn)
    encoder := json.NewEncoder(conn)
    
    for {
        // 解码请求
        var req Request
        if err := decoder.Decode(&req); err != nil {
            log.Printf("解码请求失败: %v", err)
            return
        }

        // 处理请求
        var resp Response
        resp.RequestID = req.RequestID

        // 查找服务
        if service, ok := s.services[req.ServiceMethod]; ok {
            switch req.ServiceMethod {
            case "Arith.Multiply":
                // 将请求参数转换为Args类型
                var args Args
                argsBytes, _ := json.Marshal(req.Args)
                json.Unmarshal(argsBytes, &args)
                
                // 调用服务方法
                result := service.(*ArithService).Multiply(args)
                resp.Result = result
            }
        } else {
            resp.Error = "服务方法不存在"
        }

        // 发送响应
        if err := encoder.Encode(resp); err != nil {
            log.Printf("编码响应失败: %v", err)
            return
        }
    }
}

func main() {
    // 创建服务器
    server := NewServer()
    
    // 注册服务
    arith := new(ArithService)
    server.Register("Arith.Multiply", arith)

    // 监听端口
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("启动服务器失败:", err)
    }
    defer listener.Close()
    
    log.Println("RPC 服务器正在监听端口 1234...")

    // 处理连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
        go server.handleConnection(conn)
    }
}