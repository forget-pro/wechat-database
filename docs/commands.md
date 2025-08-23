# Chatlog 命令行指令详细文档

本文档详细介绍了 Chatlog 工具的所有可用命令行指令。

## 目录

- [基础命令](#基础命令)
- [密钥管理](#密钥管理)
- [数据解密](#数据解密)
- [服务管理](#服务管理)
- [配置管理](#配置管理)
- [系统信息](#系统信息)
- [实用工具](#实用工具)
- [HTTP API 接口](#http-api-接口)
- [使用示例](#使用示例)

## 基础命令

### `chatlog` (默认命令)

启动 Terminal UI 界面，提供图形化操作界面。

```bash
chatlog [flags]
```

**参数：**

- `--debug`: 启用调试模式
- `-h, --help`: 显示帮助信息

**功能：**

- 启动交互式 Terminal UI 界面
- 提供菜单式操作，包括解密数据、启动服务等
- 适合不熟悉命令行的用户使用

**使用示例：**

```bash
# 启动 Terminal UI
chatlog

# 启动并开启调试模式
chatlog --debug
```

---

## 密钥管理

### `chatlog key`

获取微信数据库加密密钥。

```bash
chatlog key [flags]
```

**参数：**

- `-p, --pid <进程ID>`: 指定微信进程 ID（可选）
- `-h, --help`: 显示帮助信息

**功能：**

- 从正在运行的微信进程中提取数据库加密密钥
- 需要微信程序正在运行
- 在 macOS 上可能需要临时关闭 SIP

**使用示例：**

```bash
# 自动检测微信进程并获取密钥
chatlog key

# 指定微信进程ID获取密钥
chatlog key -p 12345
```

### `chatlog get-secret`

获取数据加密密钥（支持 JSON 输出）。

```bash
chatlog get-secret [flags]
```

**参数：**

- `-p, --pid <进程ID>`: 指定微信进程 ID
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 与 `key` 命令功能相同，但支持 JSON 格式输出
- 适合脚本自动化使用

**使用示例：**

```bash
# 获取密钥并以JSON格式输出
chatlog get-secret -j

# 指定进程ID并获取密钥
chatlog get-secret -p 12345 -j
```

**JSON 输出格式：**

```json
{
  "success": true,
  "key": "your_encryption_key_here",
  "message": ""
}
```

### `chatlog set-secret`

设置数据加密密钥。

```bash
chatlog set-secret [flags]
```

**参数：**

- `-k, --key <密钥>`: 要设置的加密密钥
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 手动设置数据库加密密钥
- 用于已知密钥的情况下直接设置

**使用示例：**

```bash
# 设置加密密钥
chatlog set-secret -k "your_encryption_key"

# 设置密钥并以JSON格式输出结果
chatlog set-secret -k "your_key" -j
```

---

## 数据解密

### `chatlog decrypt`

解密微信数据库文件。

```bash
chatlog decrypt [flags]
```

**参数：**

- `-d, --data-dir <目录>`: 微信数据目录路径
- `-w, --work-dir <目录>`: 解密后文件的工作目录
- `-k, --key <密钥>`: 数据库加密密钥
- `-p, --platform <平台>`: 平台类型 (windows/darwin，默认为当前系统)
- `-v, --version <版本>`: 微信版本 (3 或 4，默认为 3)
- `-h, --help`: 显示帮助信息

**功能：**

- 解密微信数据库文件到指定工作目录
- 支持微信 3.x 和 4.0 版本
- 支持 Windows 和 macOS 平台

**使用示例：**

```bash
# 基本解密命令
chatlog decrypt -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "your_key"

# 指定微信4.0版本解密
chatlog decrypt -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "your_key" -v 4

# 跨平台解密（在macOS上解密Windows数据）
chatlog decrypt -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "your_key" -p windows
```

### `chatlog decrypt-data`

解密微信数据库文件（支持 JSON 输出）。

```bash
chatlog decrypt-data [flags]
```

**参数：**

- `-d, --data-dir <目录>`: 微信数据目录路径
- `-w, --work-dir <目录>`: 解密后文件的工作目录
- `-k, --key <密钥>`: 数据库加密密钥
- `-p, --platform <平台>`: 平台类型 (windows/darwin)
- `-v, --version <版本>`: 微信版本 (3 或 4)
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 与 `decrypt` 命令功能相同
- 额外支持 JSON 格式输出，适合脚本自动化

**使用示例：**

```bash
# 解密并以JSON格式输出结果
chatlog decrypt-data -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "your_key" -j
```

---

## 服务管理

### `chatlog server`

启动 HTTP API 服务器。

```bash
chatlog server [flags]
```

**参数：**

- `-a, --addr <地址>`: 服务器监听地址 (默认: 127.0.0.1:5030)
- `-d, --data-dir <目录>`: 数据目录路径
- `-w, --work-dir <目录>`: 工作目录路径
- `-p, --platform <平台>`: 平台类型 (默认为当前系统)
- `-v, --version <版本>`: 微信版本 (默认为 3)
- `-h, --help`: 显示帮助信息

**功能：**

- 启动 HTTP API 服务器
- 提供 RESTful API 接口访问聊天数据
- 支持 MCP 协议集成

**使用示例：**

```bash
# 使用默认设置启动服务器
chatlog server

# 指定监听地址和端口
chatlog server -a "0.0.0.0:8080"

# 指定数据和工作目录
chatlog server -d "/path/to/data" -w "/path/to/work"
```

### `chatlog stop`

停止正在运行的聊天日志服务。

```bash
chatlog stop
```

**功能：**

- 查找并停止正在运行的 chatlog 进程
- 支持 Windows 和 Unix 系统
- 自动排除当前进程

**使用示例：**

```bash
# 停止服务
chatlog stop
```

### `chatlog exit`

安全退出聊天日志程序。

```bash
chatlog exit
```

**功能：**

- 停止所有正在运行的服务
- 清理资源后安全退出
- 比直接 kill 进程更安全

**使用示例：**

```bash
# 安全退出程序
chatlog exit
```

### `chatlog auto-decrypt`

启动或停止自动解密功能。

```bash
chatlog auto-decrypt [flags]
```

**参数：**

- `-e, --enable`: 启用自动解密
- `-d, --disable`: 禁用自动解密
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 监控微信数据变化并自动解密
- 简化数据更新流程

**使用示例：**

```bash
# 启用自动解密
chatlog auto-decrypt -e

# 禁用自动解密
chatlog auto-decrypt -d

# 启用自动解密并以JSON格式输出
chatlog auto-decrypt -e -j
```

---

## 配置管理

### `chatlog data-dir`

设置数据目录路径。

```bash
chatlog data-dir [flags]
```

**参数：**

- `-d, --dir <目录>`: 要设置的数据目录路径
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 设置微信数据目录路径
- 持久化保存配置

**使用示例：**

```bash
# 设置数据目录
chatlog data-dir -d "/path/to/wechat/data"

# 设置数据目录并以JSON格式输出
chatlog data-dir -d "/path/to/wechat/data" -j
```

### `chatlog work-dir`

设置工作目录路径。

```bash
chatlog work-dir [flags]
```

**参数：**

- `-d, --dir <目录>`: 要设置的工作目录路径
- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 设置解密后文件的工作目录
- 持久化保存配置

**使用示例：**

```bash
# 设置工作目录
chatlog work-dir -d "/path/to/work/dir"

# 设置工作目录并以JSON格式输出
chatlog work-dir -d "/path/to/work/dir" -j
```

---

## 系统信息

### `chatlog status`

获取微信进程信息。

```bash
chatlog status [flags]
```

**参数：**

- `-j, --json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

**功能：**

- 显示当前运行的微信进程信息
- 包括进程 ID、账号、版本等信息

**使用示例：**

```bash
# 查看微信进程状态
chatlog status

# 以JSON格式输出状态信息
chatlog status -j
```

**JSON 输出示例：**

```json
{
  "success": true,
  "processes": [
    {
      "pid": 12345,
      "account": "wxid_example",
      "version": "3.9.10.19",
      "platform": "windows"
    }
  ]
}
```

### `chatlog version`

显示程序版本信息。

```bash
chatlog version [flags]
```

**参数：**

- `-m`: 显示详细版本信息
- `-h, --help`: 显示帮助信息

**功能：**

- 显示 chatlog 程序版本
- 可显示详细的构建信息

**使用示例：**

```bash
# 显示版本
chatlog version

# 显示详细版本信息
chatlog version -m
```

---

## 实用工具

### `chatlog dumpmemory` (仅 macOS)

导出微信内存数据。

```bash
chatlog dumpmemory
```

**功能：**

- 导出微信进程的内存数据
- 仅支持 macOS 系统
- 用于调试和分析

**使用示例：**

```bash
# 导出内存数据（仅macOS）
chatlog dumpmemory
```

---

---

## HTTP API 接口

启动 HTTP 服务器后，可以通过 RESTful API 访问微信数据。所有 API 都支持分页功能。

### 基础信息

- **默认地址**: `http://127.0.0.1:5030`
- **内容类型**: `application/json` 或 `text/plain` 或 `text/csv`
- **分页参数**: `limit`（每页数量）、`offset`（偏移量）

### API 端点

#### 1. 获取联系人列表

**端点**: `GET /api/contact`

**参数**:

- `keyword` (string, 可选): 搜索关键字，支持模糊匹配用户名、别名、备注、昵称
- `limit` (int, 可选): 每页返回的联系人数量，默认返回所有
- `offset` (int, 可选): 偏移量，用于分页，默认为 0
- `format` (string, 可选): 返回格式，支持 `json`、`csv`，默认为 csv

**响应格式 (JSON)**:

```json
{
  "items": [
    {
      "user_name": "wxid_example",
      "alias": "alias_example",
      "remark": "备注名称",
      "nick_name": "昵称",
      "avatar": "头像URL",
      "avatar_hd": "高清头像URL"
    }
  ],
  "total": 150,
  "page": {
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

**使用示例**:

```bash
# 获取所有联系人
curl "http://127.0.0.1:5030/api/contact?format=json"

# 分页获取联系人（每页20个，第2页）
curl "http://127.0.0.1:5030/api/contact?format=json&limit=20&offset=20"

# 搜索包含"张三"的联系人
curl "http://127.0.0.1:5030/api/contact?format=json&keyword=张三"

# 搜索并分页
curl "http://127.0.0.1:5030/api/contact?format=json&keyword=张三&limit=10&offset=0"
```

#### 2. 获取群聊列表

**端点**: `GET /api/chatroom`

**参数**:

- `keyword` (string, 可选): 搜索关键字，支持模糊匹配群名称、群 ID
- `limit` (int, 可选): 每页返回的群聊数量
- `offset` (int, 可选): 偏移量，用于分页
- `format` (string, 可选): 返回格式，支持 `json`、`csv`

**响应格式 (JSON)**:

```json
{
  "items": [
    {
      "name": "12345678901@chatroom",
      "nick_name": "群聊名称",
      "remark": "群聊备注",
      "member_count": 25,
      "members": ["wxid_member1", "wxid_member2"]
    }
  ],
  "total": 50,
  "page": {
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}
```

**使用示例**:

```bash
# 获取所有群聊
curl "http://127.0.0.1:5030/api/chatroom?format=json"

# 分页获取群聊
curl "http://127.0.0.1:5030/api/chatroom?format=json&limit=10&offset=0"

# 搜索包含"工作"的群聊
curl "http://127.0.0.1:5030/api/chatroom?format=json&keyword=工作"
```

#### 3. 获取聊天记录

**端点**: `GET /api/chatlog`

**参数**:

- `time` (string, 必需): 时间范围，格式如 `2024-01-01` 或 `2024-01-01:2024-01-31`
- `talker` (string, 可选): 对话者 ID，可以是联系人 ID 或群聊 ID
- `sender` (string, 可选): 发送者 ID，在群聊中用于过滤特定发送者
- `keyword` (string, 可选): 消息内容关键字搜索
- `limit` (int, 可选): 每页返回的消息数量
- `offset` (int, 可选): 偏移量，用于分页
- `format` (string, 可选): 返回格式，支持 `json`、`csv`、`plain`（纯文本）
- `pagination` (bool, 可选): 是否返回分页信息，仅在 format=json 时有效

**响应格式 (JSON, 带分页)**:

```json
{
  "items": [
    {
      "id": 12345,
      "talker": "wxid_example",
      "sender": "wxid_sender",
      "content": "消息内容",
      "timestamp": 1640995200,
      "type": 1,
      "sender_avatar": "发送者头像URL",
      "sender_avatar_hd": "发送者高清头像URL"
    }
  ],
  "total": 1500,
  "page": {
    "limit": 50,
    "offset": 0,
    "has_more": true
  }
}
```

**响应格式 (JSON, 不带分页)**:

```json
[
  {
    "id": 12345,
    "talker": "wxid_example",
    "sender": "wxid_sender",
    "content": "消息内容",
    "timestamp": 1640995200,
    "type": 1,
    "sender_avatar": "发送者头像URL",
    "sender_avatar_hd": "发送者高清头像URL"
  }
]
```

**使用示例**:

```bash
# 获取2024年1月的所有消息（带分页）
curl "http://127.0.0.1:5030/api/chatlog?time=2024-01&format=json&pagination=true&limit=50&offset=0"

# 获取与特定联系人的聊天记录
curl "http://127.0.0.1:5030/api/chatlog?time=2024-01&talker=wxid_example&format=json"

# 在群聊中搜索特定发送者的消息
curl "http://127.0.0.1:5030/api/chatlog?time=2024-01&talker=12345@chatroom&sender=wxid_sender&format=json"

# 关键字搜索消息
curl "http://127.0.0.1:5030/api/chatlog?time=2024-01&keyword=重要&format=json&pagination=true"

# 获取纯文本格式的聊天记录
curl "http://127.0.0.1:5030/api/chatlog?time=2024-01&talker=wxid_example&format=plain"
```

### 分页说明

所有列表类 API 都支持分页功能：

1. **`limit`**: 指定每页返回的记录数量
2. **`offset`**: 指定从第几条记录开始返回（从 0 开始）
3. **`total`**: 响应中包含符合条件的记录总数
4. **`has_more`**: 表示是否还有更多记录

**分页计算示例**:

```bash
# 第1页（前20条）
curl "http://127.0.0.1:5030/api/contact?limit=20&offset=0"

# 第2页（第21-40条）
curl "http://127.0.0.1:5030/api/contact?limit=20&offset=20"

# 第3页（第41-60条）
curl "http://127.0.0.1:5030/api/contact?limit=20&offset=40"
```

### 错误处理

API 返回的错误格式：

```json
{
  "error": "错误信息",
  "code": "ERROR_CODE"
}
```

常见 HTTP 状态码：

- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

### JavaScript/Python 调用示例

**JavaScript (fetch)**:

```javascript
// 获取联系人列表（分页）
async function getContacts(page = 0, pageSize = 20, keyword = "") {
  const params = new URLSearchParams({
    format: "json",
    limit: pageSize,
    offset: page * pageSize,
  });

  if (keyword) {
    params.append("keyword", keyword);
  }

  const response = await fetch(`http://127.0.0.1:5030/api/contact?${params}`);
  const data = await response.json();

  return {
    contacts: data.items,
    total: data.total,
    hasMore: data.page?.has_more || false,
  };
}

// 获取聊天记录（带分页）
async function getChatMessages(time, talker, page = 0, pageSize = 50) {
  const params = new URLSearchParams({
    time: time,
    format: "json",
    pagination: "true",
    limit: pageSize,
    offset: page * pageSize,
  });

  if (talker) {
    params.append("talker", talker);
  }

  const response = await fetch(`http://127.0.0.1:5030/api/chatlog?${params}`);
  return await response.json();
}
```

**Python (requests)**:

```python
import requests

def get_contacts(page=0, page_size=20, keyword=''):
    """获取联系人列表"""
    params = {
        'format': 'json',
        'limit': page_size,
        'offset': page * page_size
    }

    if keyword:
        params['keyword'] = keyword

    response = requests.get('http://127.0.0.1:5030/api/contact', params=params)
    data = response.json()

    return {
        'contacts': data['items'],
        'total': data['total'],
        'has_more': data.get('page', {}).get('has_more', False)
    }

def get_chat_messages(time_range, talker=None, page=0, page_size=50):
    """获取聊天记录"""
    params = {
        'time': time_range,
        'format': 'json',
        'pagination': 'true',
        'limit': page_size,
        'offset': page * page_size
    }

    if talker:
        params['talker'] = talker

    response = requests.get('http://127.0.0.1:5030/api/chatlog', params=params)
    return response.json()

# 使用示例
contacts_result = get_contacts(page=0, page_size=20, keyword='张三')
messages_result = get_chat_messages('2024-01', talker='wxid_example')
```

---

## 使用示例

### 完整工作流程

1. **查看微信进程状态**

```bash
chatlog status
```

2. **获取加密密钥**

```bash
chatlog key
```

3. **解密数据库**

```bash
chatlog decrypt -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "obtained_key"
```

4. **启动 HTTP 服务**

```bash
chatlog server -a "127.0.0.1:5030"
```

5. **停止服务**

```bash
chatlog stop
```

### 自动化脚本示例

```bash
#!/bin/bash

# 获取密钥
KEY=$(chatlog get-secret -j | jq -r '.key')

if [ "$KEY" != "null" ] && [ "$KEY" != "" ]; then
    echo "成功获取密钥: $KEY"

    # 解密数据
    chatlog decrypt-data -d "/path/to/wechat/data" -w "/path/to/work/dir" -k "$KEY" -j

    # 启动服务
    chatlog server &

    echo "服务已启动"
else
    echo "获取密钥失败"
    exit 1
fi
```

### JSON API 集成示例

```python
import subprocess
import json

def get_wechat_key():
    """获取微信加密密钥"""
    result = subprocess.run(['chatlog', 'get-secret', '-j'],
                          capture_output=True, text=True)
    data = json.loads(result.stdout)
    return data['key'] if data['success'] else None

def decrypt_wechat_data(data_dir, work_dir, key):
    """解密微信数据"""
    result = subprocess.run(['chatlog', 'decrypt-data',
                           '-d', data_dir, '-w', work_dir,
                           '-k', key, '-j'],
                          capture_output=True, text=True)
    data = json.loads(result.stdout)
    return data['success']

# 使用示例
key = get_wechat_key()
if key:
    success = decrypt_wechat_data('/path/to/data', '/path/to/work', key)
    print(f"解密{'成功' if success else '失败'}")
```

---

## 注意事项

1. **权限要求**：

   - Windows: 可能需要管理员权限
   - macOS: 获取密钥前需要临时关闭 SIP

2. **微信状态**：

   - 获取密钥时微信必须正在运行
   - 解密数据时微信可以关闭

3. **数据安全**：

   - 密钥信息敏感，请妥善保管
   - 仅处理自己合法拥有的数据

4. **平台兼容性**：

   - 大部分命令支持跨平台使用
   - `dumpmemory` 仅支持 macOS

5. **版本兼容性**：
   - 支持微信 3.x 和 4.0 版本
   - 不同版本数据库结构可能不同

---

## 故障排除

### 常见问题

1. **无法获取密钥**

   - 确保微信正在运行
   - 检查是否有足够权限
   - macOS 用户检查 SIP 状态

2. **解密失败**

   - 验证密钥是否正确
   - 检查数据目录路径
   - 确认微信版本设置

3. **服务启动失败**
   - 检查端口是否被占用
   - 验证工作目录权限
   - 查看调试日志信息

### 调试方法

```bash
# 启用调试模式
chatlog --debug

# 查看详细错误信息
chatlog command --help
```

有关更多帮助信息，请访问项目文档或提交 Issue。
