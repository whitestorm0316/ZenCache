# ZenCache 缓存系统

## 概述
ZenCache 是一个用 Go 语言编写的高性能分布式缓存系统，旨在为应用程序提供高效、可靠的缓存解决方案。它支持并发访问、LRU 缓存淘汰策略，并具备通过一致性哈希算法实现的分布式缓存功能。

## 特性
- **并发安全**：采用读写锁机制，确保多线程环境下的缓存操作安全。
- **LRU 缓存淘汰**：当缓存达到最大容量时，自动淘汰最近最少使用的数据。
- **分布式缓存**：支持通过一致性哈希算法将缓存数据分布到多个节点。
- **可配置性**：可以通过配置文件加载缓存的相关参数。
- **HTTP 接口**：提供简单的 HTTP 接口用于存储、获取和删除缓存数据。

## 安装与依赖
### 依赖管理
ZenCache 使用 Go Modules 进行依赖管理。在项目根目录下，运行以下命令下载所有依赖：
```sh
go mod tidy
```

### 安装 Protocol Buffers 编译器
为了生成 Protocol Buffers 代码，需要安装 `protoc` 编译器和 `protoc-gen-go` 插件。可以使用以下命令进行安装：
```sh
# 检查 protoc 是否安装
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo protoc not found, please install Protocol Buffers compiler first
    echo You can download it from: https://github.com/protocolbuffers/protobuf/releases
    exit /b 1
)

# 检查 protoc-gen-go 是否安装
where protoc-gen-go >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    if %errorlevel% neq 0 (
        echo Failed to install protoc-gen-go
        exit /b 1
    )
)
```

### 生成 Protocol Buffers 代码
运行以下脚本生成 Protocol Buffers 代码：
```sh
./scripts/gen_proto.bat
```

## 配置文件
可以通过 `config.json` 文件配置缓存系统的相关参数。如果配置文件不存在，将使用默认配置。配置文件示例：
```json
{
    "hash": {
        "replicas": 3
    },
    // 其他配置项...
}
```

## 代码结构
- **`cmd`**：包含项目的入口文件 `main.go`。
- **`internal`**：
  - **`cache`**：实现了缓存的核心逻辑，包括 `ByteView`、`cache`、`Engine` 和 `Group` 等。
  - **`config`**：负责加载配置文件。
  - **`lru`**：实现了 LRU 缓存淘汰算法。
  - **`peers`**：定义了分布式缓存的节点选择接口。
  - **`transport`**：包含 HTTP 服务器的实现，提供缓存操作的 HTTP 接口。
  - **`consistenthash`**：实现了一致性哈希算法。
- **`scripts`**：包含生成 Protocol Buffers 代码的脚本。
- **`go.mod` 和 `go.sum`**：管理项目的依赖。

## 使用方法
### 启动服务器
在项目根目录下，运行以下命令启动缓存服务器：
```sh
go run cmd/main.go
```

### HTTP 接口
- **存储数据**：
  - **URL**：`/v1/store_key`
  - **方法**：`POST`
  - **请求体**：
```json
{
    "group": "test_group",
    "key": "test_key",
    "value": "test_value"
}
```
- **获取数据**：
  - **URL**：`/v1/get_key`
  - **方法**：`POST`
  - **请求体**：
```json
{
    "group": "test_group",
    "key": "test_key"
}
```
- **删除数据**：
  - **URL**：`/v1/delete_key`
  - **方法**：`POST`
  - **请求体**：
```json
{
    "group": "test_group",
    "key": "test_key"
}
```

## 测试
项目中包含了多个测试文件，用于验证各个模块的功能。可以使用以下命令运行所有测试：
```sh
go test ./...
```

## 性能测试
使用 `go test` 的 `-bench` 标志可以运行性能测试：
```sh
go test -bench=.
```
