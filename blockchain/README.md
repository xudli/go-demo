# 简单区块链实现

这是一个用 Go 语言实现的基础区块链项目。

## 核心功能

1. 区块结构 (block.go)
   - 区块头
     - 时间戳 (Timestamp)
     - 前一个区块的哈希值 (PrevHash)
     - 当前区块的哈希值 (Hash)
     - Nonce 值（用于工作量证明）
   - 区块数据 (Data)
   - 交易列表 (Transactions)

2. 区块链结构 (blockchain.go)
   - 使用 BoltDB 持久化存储区块数据
   - tip: 指向最新区块的哈希
   - 创世区块的实现
   - 添加新区块的方法

3. 工作量证明 (pow.go)
   - 实现 PoW (Proof of Work) 算法
   - 难度目标值设定 (targetBits)
   - 计算哈希值直到满足难度要求
   - 验证区块的有效性

## 技术实现细节

1. 区块创建流程
   - 设置区块基本信息（时间戳、前一区块哈希等）
   - 执行工作量证明获取有效的 Nonce 值
   - 使用 SHA256 计算区块哈希
   - 序列化区块数据并存储

2. 数据持久化
   - 使用 BoltDB 作为键值存储数据库
   - 使用 "blocks" bucket 存储所有区块数据
   - 使用 "l" 键存储最后一个区块的哈希值
   - 实现区块的序列化和反序列化

3. 工作量证明机制
   - 设置目标难度值 (targetBits = 24)
   - 组合区块数据（PrevHash + Data + Timestamp + targetBits + Nonce）
   - 计算 SHA256 哈希值
   - 验证哈希值是否满足难度要求

## 项目结构

``` yaml

blockchain/
├── block/
│   ├── block.go          // 区块结构定义
│   └── blockchain.go     // 区块链结构定义
├── pow/
│   └── pow.go           // 工作量证明实现
├── utils/
│   └── utils.go         // 工具函数
└── main.go              // 程序入口
```

## 使用方法

1. 创建新的区块链：

```go
bc := block.NewBlockchain()
```

2. 添加新的区块：

```go
bc.AddBlock([]byte("交易数据"))
```

## 待实现功能

1. 交易系统
   - 交易结构设计
   - UTXO 模型实现
   - 交易签名和验证

2. 钱包功能
   - 密钥对生成
   - 地址生成
   - 签名实现

3. 网络功能
   - P2P 网络实现
   - 节点间通信
   - 区块同步机制

4. 共识机制优化
   - 难度值动态调整
   - 分叉处理
   - 最长链选择

## 技术栈

- Go 1.20
- BoltDB (持久化存储)
- crypto/sha256 (哈希计算)
- encoding/gob (序列化)

这个更新后的 README.md 更好地反映了当前的实现细节，并清晰地列出了项目的结构和待实现的功能。是否需要我继续完善其他部分？