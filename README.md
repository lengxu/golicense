# GoLicense - 授权控制系统

一个基于RSA公私钥的Go语言授权控制系统，支持硬件指纹绑定和模块化授权管理。

## 系统架构

```
客户端启动 -> 自动检查license.dat
          -> 如无授权则自动生成req.dat 
          -> 等待管理员处理
管理员    -> 收到req.dat
          -> 使用licgen生成license.dat (包含客户信息)
          -> 发回给客户端
客户端    -> 放置license.dat到bin/目录
          -> 重启程序，自动验证授权
          -> 各模块启动时自动检查授权
```

## 核心特性

✅ **自动化授权管理**: 客户端启动时自动生成授权请求  
✅ **硬件指纹绑定**: 基于CPU、主板、系统UUID等硬件特征  
✅ **RSA 4096位加密**: 私钥签名，公钥验证，内置密钥无需外部文件  
✅ **客户信息管理**: 支持客户名称、组织等信息记录  
✅ **模块化授权**: 可控制goscan、gopasswd、goweb等特定模块  
✅ **跨平台支持**: Windows/Linux兼容的硬件指纹获取  
✅ **时间控制**: 支持试用期、年度授权等时间限制  
✅ **友好提示**: 授权状态、过期提醒等用户友好界面  
✅ **Base58字符串格式**: req.dat和license.dat都是Base58编码的单行字符串，避免易混淆字符，易于复制粘贴

## 目录结构

```
golicense/
├── server/           # 服务端代码（授权生成）
├── client/           # 客户端代码（验证库）
├── cmd/
│   ├── reqgen/      # 生成req.dat工具
│   ├── licgen/      # 生成license.dat工具
│   └── liccheck/    # 检查license.dat工具
└── README.md
```

## 快速开始

### 方案一：自动化流程（推荐）

**客户端操作：**
1. 启动任何需要授权的程序（goscan/gopasswd/goweb）
2. 程序自动检测到无授权，生成 `bin/req.dat` (单行字符串文件)
3. 复制req.dat内容发送给管理员（或直接发送文件）

**管理员操作：**
```bash
# 根据req.dat生成授权文件
cd cmd/licgen
go run main.go -i req.dat -c "客户名称" -org "客户公司" -d 365
# 生成包含客户信息的1年期license.dat
```

**客户端部署：**
```bash
# 将license.dat放入bin目录（或直接复制文件内容到license.dat）
cp license.dat goweb/bin/
# 重新启动程序，自动验证通过
```

### 方案二：手动生成（可选）

**1. 手动生成授权请求**
```bash
cd cmd/reqgen
go run main.go -o req.dat
```

**2. 生成授权文件**
```bash
cd cmd/licgen  
go run main.go -i req.dat -c "张三" -org "ABC公司" -d 365
```

**3. 验证授权文件**
```bash
cd cmd/liccheck
go run main.go -l license.dat
```

## 集成到现有模块

### 推荐集成方式（自动化）

**goweb集成**
```go
// goweb/main.go
import "golicense/client"

func main() {
    // 自动授权检查（包含自动生成req.dat）
    if err := client.AutoLicenseCheck("goweb"); err != nil {
        log.Fatal("Authorization required:", err)
    }
    // 继续原有逻辑...
}
```

**goscan集成**
```go
// goscan/v2/cmd/cmd.go
import "golicense/client"

func Execute() {
    // 自动授权检查
    if err := client.AutoLicenseCheck("goscan"); err != nil {
        fmt.Println("Authorization required:", err)
        os.Exit(1)
    }
    // 继续扫描逻辑...
}
```

**gopasswd集成**
```go
// gopasswd/cmd/cmd.go
import "golicense/client"

func Execute() {
    // 自动授权检查
    if err := client.AutoLicenseCheck("gopasswd"); err != nil {
        fmt.Println("Authorization required:", err)
        os.Exit(1)
    }
    // 继续攻击逻辑...
}
```

### 简单集成方式（仅验证）

如果只需要验证已有授权，可以使用：
```go
// 快速验证授权
if err := client.QuickLicenseCheck("module_name"); err != nil {
    log.Fatal("License validation failed:", err)
}
```

## 安全特性

- **RSA 4096位密钥**: 服务端私钥签名，客户端公钥验证
- **硬件指纹绑定**: 基于CPU ID、主板序列号、系统UUID
- **AES-256加密**: 敏感数据全程加密存储
- **时间限制**: 支持试用期、年度授权等时间控制
- **模块化授权**: 可控制具体模块的使用权限
- **跨平台支持**: Windows/Linux兼容的硬件指纹获取

## 命令行工具

### reqgen - 授权请求生成工具
```bash
reqgen [选项]
  -o string    输出文件路径 (默认 "req.dat")
  -h          显示帮助信息
```

### licgen - 授权文件生成工具
```bash
licgen -i <req.dat> [选项]
  -i string    输入的req.dat文件路径 (必需)
  -o string    输出的license.dat文件路径 (默认 "license.dat")
  -d int       授权有效期天数 (默认 365)
  -h          显示帮助信息
```

### liccheck - 授权文件检查工具
```bash
liccheck [选项]
  -l string    license.dat文件路径 (默认 "license.dat")
  -m string    检查特定模块授权
  -h          显示帮助信息
```

## 硬件指纹获取

### Windows平台
- CPU ProcessorId
- 主板序列号
- BIOS UUID

### Linux平台
- Machine ID (/etc/machine-id)
- CPU信息hash
- DMI产品UUID

## 文件格式

### req.dat (授权请求文件)
```
REQ:8S6CuJCr7RnCoMyEK8AHnByBD6FaE4HCs1R6C9aYSPYM55Evg5w1iYyszwQ3tgfUhpNP...
```
- 格式：`REQ:` + Base58编码的Gzip压缩JSON数据
- 优势：单行字符串，避免易混淆字符(0,O,I,l)，易于复制粘贴

### license.dat (授权文件)
```
LIC:aeNtLDXfwLSaYbNhn2fm9jyynwxb9ZgxTEVyD7hEsxK3T5FQD1B9iP1v7ZQcKVoB1PPe...
```
- 格式：`LIC:` + Base58编码的Gzip压缩JSON数据
- 优势：Base58编码无易混淆字符，压缩后体积适中，文本格式便于传输

## 构建和部署

```bash
# 构建所有工具
go mod tidy

# 构建客户端工具
cd cmd/reqgen && go build
cd cmd/liccheck && go build

# 构建服务端工具  
cd cmd/licgen && go build
```

## 注意事项

1. **私钥安全**: 服务端私钥需要妥善保管，不可泄露
2. **硬件变更**: 硬件更换后需要重新申请授权
3. **时间同步**: 确保系统时间准确，避免授权时间判断错误
4. **文件权限**: license.dat建议设置适当的文件权限
5. **备份恢复**: 重要授权文件建议备份