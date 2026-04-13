# ProKit

基于 Go 的 Protocol Buffers 工具库，提供序列化/反序列化封装以及带长度前缀的帧协议读写。

## 项目结构

```
prokit/
├── codec/
│   └── codec.go          # 序列化/反序列化封装 + Frame 帧协议
├── proto/
│   ├── user.proto         # 用户消息定义
│   └── user.pb.go         # protoc 生成的 Go 代码
├── tools/
│   └── genproto.sh        # protoc 代码生成脚本
├── go.mod
└── README.md
```

## 依赖

- Go 1.26+
- [protoc](https://github.com/protocolbuffers/protobuf) (Protocol Buffers 编译器)
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) (Go 代码生成插件)

安装 protoc-gen-go：

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Proto 消息定义

`proto/user.proto` 定义了以下类型：

### 枚举

| 枚举 | 值 | 说明 |
|------|-----|------|
| `Gender` | `GENDER_UNSPECIFIED(0)`, `GENDER_MALE(1)`, `GENDER_FEMALE(2)` | 用户性别 |
| `UserStatus` | `USER_STATUS_UNSPECIFIED(0)`, `USER_STATUS_ACTIVE(1)`, `USER_STATUS_INACTIVE(2)`, `USER_STATUS_BANNED(3)` | 用户状态 |

### 消息

#### Address - 联系地址

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `province` | `string` | 1 | 省份 |
| `city` | `string` | 2 | 城市 |
| `district` | `string` | 3 | 区/县 |
| `street` | `string` | 4 | 街道 |
| `zip_code` | `string` | 5 | 邮编 |

#### User - 用户信息

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `id` | `int64` | 1 | 用户 ID |
| `username` | `string` | 2 | 用户名 |
| `email` | `string` | 3 | 邮箱 |
| `phone` | `string` | 4 | 电话 |
| `gender` | `Gender` | 5 | 性别 |
| `age` | `int32` | 6 | 年龄 |
| `address` | `Address` | 7 | 地址 |
| `tags` | `repeated string` | 8 | 标签列表 |
| `extra` | `map<string, string>` | 9 | 扩展键值对 |
| `status` | `UserStatus` | 10 | 用户状态 |
| `created_at` | `google.protobuf.Timestamp` | 11 | 创建时间 |
| `updated_at` | `google.protobuf.Timestamp` | 12 | 更新时间 |

#### UserList - 用户列表

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `users` | `repeated User` | 1 | 用户集合 |
| `total_count` | `int32` | 2 | 总数 |

## 生成 Proto 代码

```bash
# 使用脚本
bash tools/genproto.sh

# 或手动执行
protoc --go_out=. --go_opt=paths=source_relative proto/user.proto
```

## codec 包使用

### 基础序列化 / 反序列化

```go
import (
    "github.com/tattoo1880/protkit/codec"
    pb "github.com/tattoo1880/protkit/proto"
)

user := &pb.User{
    Id:       1,
    Username: "alice",
    Email:    "alice@example.com",
    Status:   pb.UserStatus_USER_STATUS_ACTIVE,
}

// 序列化
data, err := codec.Marshal(user)

// 反序列化
got := &pb.User{}
err = codec.Unmarshal(data, got)
```

### Frame 帧协议

适用于 TCP 连接、文件流、消息队列等需要消息边界分隔的场景。
帧格式：`[4 字节大端长度] + [protobuf payload]`。

```go
import (
    "bufio"
    "bytes"

    "github.com/tattoo1880/protkit/codec"
    pb "github.com/tattoo1880/protkit/proto"
)

var buf bytes.Buffer

// 写入帧
err := codec.WriteFrame(&buf, user)

// 读取帧
reader := bufio.NewReader(&buf)
got := &pb.User{}
err = codec.ReadFrame(reader, got)
```

## 运行测试

```bash
go test ./...
```