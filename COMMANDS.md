# ChatLog 命令行指令说明

本项目已经增加了命令行支持，所有命令都集成在 `cmd_unified.go` 文件中。每个命令都支持 JSON 格式输出。

## 可用命令

### 1. 获取微信进程信息

```bash
go run . status [--json]
```

返回当前微信实例的完整状态信息，包括：

- Account: 微信账号
- PID: 进程 ID
- Status: 状态（online/offline）
- ExePath: 可执行文件路径
- Platform: 平台（darwin/windows）
- Version: 微信版本
- Session: 会话时间
- Data Key: 数据加密密钥
- Data Usage: 数据目录使用量
- Data Dir: 数据目录路径
- Work Usage: 工作目录使用量
- Work Dir: 工作目录路径
- HTTP Server: HTTP 服务状态
- Auto Decrypt: 自动解密状态

**示例：**

```bash
# 普通输出（格式与终端UI一致）
go run . status

# JSON 输出
go run . status --json
```

### 2. 获取数据密钥

```bash
go run . get-secret [--pid PID] [--json]
```

从微信进程中提取数据加密密钥。

**参数：**

- `--pid, -p`: 指定微信进程 ID（可选，如果有多个进程时需要指定）
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
# 自动获取密钥
go run . get-secret

# 指定进程ID获取密钥
go run . get-secret --pid 12345

# JSON 输出
go run . get-secret --json
```

### 3. 解密数据

```bash
go run . decrypt-data --data-dir DIR --key KEY [--work-dir DIR] [--platform PLATFORM] [--version VERSION] [--json]
```

解密微信数据库文件。

**参数：**

- `--data-dir, -d`: 微信数据目录（必需）
- `--key, -k`: 加密密钥（必需）
- `--work-dir, -w`: 解密后文件存储目录（可选）
- `--platform, -p`: 平台类型，默认为当前系统
- `--version, -v`: 微信版本，默认为 3
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
go run . decrypt-data \
  --data-dir "/path/to/wechat/data" \
  --key "your_encryption_key" \
  --work-dir "/path/to/output"
```

### 4. 开启/停止自动解密

```bash
go run . auto-decrypt [--stop] [--json]
```

启动或停止自动解密监控。

**参数：**

- `--stop, -s`: 停止自动解密
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
# 启动自动解密（会阻塞，按Ctrl+C停止）
go run . auto-decrypt

# 停止自动解密
go run . auto-decrypt --stop

# JSON 输出
go run . auto-decrypt --json
```

### 5. 设置工作目录

```bash
go run . work-dir --path PATH [--json]
```

设置解密文件的工作目录。

**参数：**

- `--path, -p`: 工作目录路径（必需）
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
go run . work-dir --path "/path/to/work/directory"
```

### 6. 设置数据密钥

```bash
go run . set-secret --key KEY [--json]
```

手动设置数据加密密钥。

**参数：**

- `--key, -k`: 加密密钥（必需）
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
go run . set-secret --key "your_encryption_key"
```

### 7. 设置数据目录

```bash
go run . data-dir --path PATH [--json]
```

设置微信数据目录。

**参数：**

- `--path, -p`: 数据目录路径（必需）
- `--json, -j`: 以 JSON 格式输出

**示例：**

```bash
go run . data-dir --path "/path/to/wechat/data"
```

## JSON 输出格式

所有命令都支持 `--json` 参数，输出格式统一为：

```json
{
  "success": true,
  "message": "操作描述",
  "data": "具体数据（根据命令不同而不同）"
}
```

## 保留原功能

原有的终端 UI 交互方式依然保留，直接运行 `go run .` 即可进入交互模式。

## 常用工作流程

1. 首先获取微信进程信息：

   ```bash
   go run . status --json
   ```

2. 获取加密密钥：

   ```bash
   go run . get-secret --json
   ```

3. 解密数据：

   ```bash
   go run . decrypt-data --data-dir "数据目录" --key "密钥" --json
   ```

4. 可选：开启自动解密：
   ```bash
   go run . auto-decrypt --json
   ```
