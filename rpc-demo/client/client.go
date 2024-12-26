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

// Args 计算参数
type Args struct {
    A, B int
}

// Client RPC客户端
type Client struct {
    conn net.Conn
    requestID uint64
}

// NewClient 创建新的RPC客户端
func NewClient(addr string) (*Client, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }
    return &Client{conn: conn}, nil
}

// Call 调用远程服务
func (c *Client) Call(serviceMethod string, args interface{}) (interface{}, error) {
    c.requestID++
    
    // 创建请求
    req := &Request{
        ServiceMethod: serviceMethod,
        Args:         args,
        RequestID:    c.requestID,
    }

    // 发送请求
    encoder := json.NewEncoder(c.conn)
    if err := encoder.Encode(req); err != nil {
        return nil, err
    }

    // 接收响应
    decoder := json.NewDecoder(c.conn)
    var resp Response
    if err := decoder.Decode(&resp); err != nil {
        return nil, err
    }

    if resp.Error != "" {
        return nil, error(nil)
    }
    return resp.Result, nil
}

func main() {
    client, err := NewClient("localhost:1234")
    if err != nil {
        log.Fatal("连接服务器失败:", err)
    }

    args := Args{A: 7, B: 8}
    result, err := client.Call("Arith.Multiply", args)
    if err != nil {
        log.Fatal("调用失败:", err)
    }

    log.Printf("计算结果: %d * %d = %v", args.A, args.B, result)
} 