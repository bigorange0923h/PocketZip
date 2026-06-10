# PocketUnzip MVP 设计文档

**日期**：2026-06-09  
**状态**：已批准  
**范围**：MVP 完整闭环

---

## 1. 技术选型

| 项目 | 选择 | 理由 |
|-----|------|------|
| 桌面框架 | Go + Wails v2 | 稳定版，文档完善，社区资源多 |
| 前端框架 | Vue + TypeScript | Wails 官方模板支持好 |
| 解压内核 | 7z.exe（命令行调用） | 支持格式多，开发复杂度低 |
| 本地数据库 | SQLite（modernc.org/sqlite） | 纯 Go 实现，无需 CGO，交叉编译方便 |
| 密码加密 | Windows DPAPI | 与用户账户绑定，无需主密码 |
| 7z.exe 分发 | 随程序分发 | 简单可靠，不依赖网络 |

---

## 2. 垂直切片计划

### 切片 1：核心解压流程

**目标**：选择文件 → 解压 → 看到日志

- 初始化 Wails + Vue 项目
- 实现 Archive Service（7z.exe 调用）
- 前端：文件选择 + 拖拽、解压按钮、实时日志面板

### 切片 2：密码管理

**目标**：加密压缩包自动尝试历史密码，失败则弹窗输入

- 实现 Security Service（DPAPI 加解密）
- 实现 Password Service（密码保存、匹配、自动尝试）
- 前端：密码输入弹窗、记住密码选项

### 切片 3：历史与配置

**目标**：查看解压历史、配置默认解压目录

- 实现 History Service（解压历史记录）
- 实现 Config Service（本地配置）
- 前端：历史记录页面、设置页面

### 切片 4：打包与优化

**目标**：Windows exe 可分发

- Windows exe 打包
- 7z.exe 随程序分发
- 错误处理完善
- UI 细节打磨（3D Elements 风格）

---

## 3. 项目结构

```text
PocketUnzip/
├── cmd/pocketunzip/
│   └── main.go                    # 程序入口
├── internal/
│   ├── app/
│   │   └── app.go                 # Wails App 绑定，编排层
│   ├── archive/
│   │   ├── archive.go             # 7z.exe 调用（已有）
│   │   └── archive_test.go
│   ├── password/
│   │   ├── password.go            # 密码保存、匹配、自动尝试
│   │   └── password_test.go
│   ├── history/
│   │   ├── history.go             # 解压历史 CRUD
│   │   └── history_test.go
│   ├── config/
│   │   ├── config.go              # 本地配置读写
│   │   └── config_test.go
│   ├── security/
│   │   ├── dpapi.go               # Windows DPAPI 加解密
│   │   ├── mask.go                # 日志脱敏（已有）
│   │   └── security_test.go
│   └── db/
│       └── db.go                  # SQLite 初始化、迁移
├── frontend/
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   ├── components/            # 通用组件
│   │   ├── views/                 # 页面
│   │   ├── composables/           # 组合式函数
│   │   └── assets/                # 静态资源
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── third_party/
│   └── 7zip/
│       ├── 7z.exe
│       └── 7z.dll
├── docs/
├── scripts/
├── go.mod
└── go.sum
```

**设计决策**：

1. **App Service 作为 Wails 绑定层** — 前端调用 App Service，App Service 编排各内部模块
2. **每个模块独立包** — 便于测试和维护
3. **前端使用 composables** — 封装与后端的交互逻辑

---

## 4. 核心模块设计

### 4.1 Archive Service

```go
type ExtractRequest struct {
    SevenZipPath string
    ArchivePath  string
    OutputDir    string
    Password     string
}

type ExtractResult struct {
    Success bool
    ExitErr error
}

func Extract(ctx context.Context, req ExtractRequest, onLog LogHandler) ExtractResult
```

### 4.2 Security Service

```go
// DPAPI 加解密（Windows 专用）
func Encrypt(plaintext []byte) ([]byte, error)
func Decrypt(ciphertext []byte) ([]byte, error)

// 日志脱敏（已有）
func MaskPasswordArg(line string) string
```

### 4.3 Password Service

```go
type PasswordRecord struct {
    ID                int64
    ArchivePath       string
    ArchiveName       string
    ArchiveHash       string
    EncryptedPassword []byte
    SuccessCount      int
    LastUsedAt        time.Time
    CreatedAt         time.Time
}

// 保存密码（加密后存入 SQLite）
func Save(db *sql.DB, archivePath, password string) error

// 匹配密码（按路径 → 文件名 → 常用密码顺序尝试）
func Match(db *sql.DB, archivePath string) ([]string, error)
```

### 4.4 History Service

```go
type ExtractHistory struct {
    ID           int64
    ArchivePath  string
    OutputDir    string
    Success      bool
    UsedPassword bool
    ErrorMessage string
    CreatedAt    time.Time
}

// 记录解压历史
func Record(db *sql.DB, history ExtractHistory) error

// 查询历史列表
func List(db *sql.DB, limit int) ([]ExtractHistory, error)
```

### 4.5 App Service（Wails 绑定）

```go
type App struct {
    ctx          context.Context
    db           *sql.DB
    sevenZipPath string
}

// Wails 绑定方法，前端可直接调用
func (a *App) SelectFile() (string, error)
func (a *App) SelectDirectory() (string, error)
func (a *App) Extract(archivePath, outputDir string) error
func (a *App) GetHistory() ([]ExtractHistory, error)
func (a *App) GetPasswordCandidates(archivePath string) ([]string, error)
func (a *App) SavePassword(archivePath, password string) error
```

---

## 5. 数据流

### 5.1 解压主流程

```
用户选择文件/拖拽
       ↓
前端调用 App.Extract(archivePath, outputDir)
       ↓
App Service 计算输出目录（默认同名文件夹）
       ↓
尝试无密码解压 ──→ 成功 ──→ 记录历史 ──→ 返回成功
       ↓ 失败
调用 Password.Match() 获取候选密码
       ↓
逐个尝试候选密码 ──→ 成功 ──→ 记录历史 ──→ 返回成功
       ↓ 全部失败
返回"需要密码"信号
       ↓
前端弹出密码输入框
       ↓
用户输入密码 ──→ 前端调用 App.Extract(archivePath, outputDir, password)
       ↓
解压成功 ──→ 询问是否保存密码 ──→ 记录历史 ──→ 返回成功
```

### 5.2 日志推送

```
7z.exe stdout/stderr
       ↓
Archive Service scanPipe()
       ↓
调用 onLog(line) 回调
       ↓
App Service 通过 Wails Events 推送前端
       ↓
前端实时显示日志
```

### 5.3 密码匹配优先级

```
1. archive_path 精确匹配
2. archive_name 精确匹配
3. archive_hash 匹配
4. success_count 最高的常用密码
5. 用户手动输入
```

---

## 6. 错误处理

### 6.1 7z.exe 错误判断

| 退出码 | 含义 | 处理方式 |
|-------|------|---------|
| 0 | 成功 | 正常返回 |
| 1 | 致命错误 | 显示错误信息 |
| 2 | 文件损坏 | 提示"压缩包已损坏" |
| 7 | 命令行错误 | 提示"参数错误" |
| 8 | 内存不足 | 提示"内存不足" |
| 255 | 用户中断 | 忽略 |

### 6.2 密码错误识别

```go
// 检查 7z 输出中是否包含密码错误关键词
func isPasswordError(output string) bool {
    keywords := []string{
        "Wrong password",
        "密码错误",
        "Cannot open encrypted archive",
    }
    // 检查输出中是否包含关键词
}
```

### 6.3 错误传递

```
7z.exe 错误
    ↓
Archive Service 解析退出码 + 输出
    ↓
返回结构化错误（类型 + 消息）
    ↓
App Service 根据错误类型决定行为
    ↓
前端显示友好提示
```

### 6.4 前端错误展示

- 密码错误 → 弹窗提示，允许重新输入
- 文件损坏 → 提示"压缩包已损坏，请检查文件"
- 路径不存在 → 提示"文件不存在"
- 权限不足 → 提示"没有写入权限"
- 其他错误 → 显示原始错误信息

---

## 7. 测试策略

### 7.1 单元测试

| 模块 | 测试内容 | 测试方式 |
|-----|---------|---------|
| Archive Service | 参数构造、错误解析 | Mock 7z.exe 输出 |
| Security Service | DPAPI 加解密、日志脱敏 | 本地加密解密验证 |
| Password Service | 保存、匹配、排序 | 内存 SQLite 测试 |
| History Service | 记录、查询、分页 | 内存 SQLite 测试 |
| Config Service | 读写配置 | 临时文件测试 |

### 7.2 集成测试

- 真实调用 7z.exe 解压测试压缩包
- 测试加密压缩包的密码匹配流程
- 测试历史记录的完整生命周期

### 7.3 测试数据

```
third_party/testdata/
├── test.zip                 # 普通压缩包
├── test_password.zip        # 密码为 "123456" 的压缩包
├── test_corrupt.zip         # 损坏的压缩包
└── test_chinese_path.zip    # 包含中文路径的压缩包
```

### 7.4 测试覆盖率目标

- Archive Service：80%+
- Security Service：90%+
- Password Service：80%+
- History Service：80%+
- App Service：集成测试覆盖核心流程

---

## 8. UI 设计

UI 采用 3D Elements 风格，详见 `docs/UI设计规范.md`。

关键设计要点：
- 深色渐变背景 + 冷色系主调
- 玻璃质感卡片 + 3D 悬停效果
- 霓虹风格进度条 + 发光动画
- 实时日志面板 + 扫描线效果

UI 设计需与用户讨论确认后实施。

---

## 9. 依赖清单

### Go 依赖

```
github.com/wailsapp/wails/v2
modernc.org/sqlite
```

### 前端依赖

```
vue
typescript
vite
```

### 外部工具

```
7z.exe / 7z.dll（随程序分发）
```

---

## 10. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|-----|------|---------|
| 7z.exe 输出格式不一致 | 日志解析失败 | 多格式测试，容错处理 |
| DPAPI 跨账户不可用 | 密码无法迁移 | 文档说明，后续支持导出 |
| Wails v2 限制 | 某些交互无法实现 | 查阅文档，社区求助 |
| 大文件解压卡顿 | UI 无响应 | 异步处理，进度回调 |

---

*本文档已通过 Spec 自检，无 TBD/TODO 项，无内部矛盾，范围聚焦于 MVP。*
