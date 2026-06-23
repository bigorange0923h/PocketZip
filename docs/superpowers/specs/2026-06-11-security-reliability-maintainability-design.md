# PocketZip 安全、可靠性与可维护性改进设计

## 目标

在保留现有标签式导航和视觉风格的前提下，解决密码明文暴露、策略配置不完整、压缩包预览解析脆弱、路径穿越风险和 `App.vue` 职责过重问题。

## 范围

本轮包含：

- 历史密码仅在 Go 后端自动尝试，不再向前端返回明文。
- 手动密码成功后可按用户选择保存，默认选择保存。
- 密码库 API 只返回不含密码密文的摘要。
- 解压策略完整 JSON 持久化，并实际执行自动打开目录。
- 使用 `7z l -slt` 结构化输出解析压缩包条目。
- 解压前拒绝绝对路径、盘符路径、UNC 路径和 `..` 路径穿越。
- 保留现有标签导航，将解压工作区和工作流从 `App.vue` 拆出。
- 为上述行为补充自动化测试和构建验证。

本轮不包含：

- 重设计整体导航或视觉主题。
- 自动更新。
- 替换 7-Zip 解压内核。
- 在当前非 Windows 环境完成 Windows 实机验收。

## 架构

前端继续通过 Wails Bridge 调用 `internal/app.App`。`App` 仍作为业务编排层，但把配置序列化、输出目录计算和安全验证委托给可独立测试的函数。

```text
App.vue
  └─ ExtractWorkspace.vue
       └─ useExtractWorkflow.ts
            └─ Wails App API
                 ├─ archive: list / validate / extract
                 ├─ password: match / save / summaries
                 ├─ history
                 └─ app config / strategies
```

## 密码安全与保存流程

### 自动密码匹配

`App.Extract` 继续在 Go 后端读取并尝试已保存密码。前端不再调用或展示 `GetPasswordCandidates`，并从 Wails 绑定面删除该接口。

自动密码尝试失败后，`App.Extract` 返回 `ErrPasswordRequired`。前端弹窗只显示说明和手动输入框。

### 手动密码

密码弹窗提交：

```text
password: string
rememberPassword: boolean
```

“解压成功后记住此密码”默认勾选。`ExtractWithPassword` 仅在解压成功且 `rememberPassword` 为真时保存密码。

### 去重和摘要

保存密码时，后端先检查同一路径或文件名关联记录。若解密后密码相同，则更新使用时间而不是插入重复记录。

密码库接口返回专用摘要：

```text
id
archivePath
archiveName
successCount
lastUsedAt
createdAt
```

该接口不包含 `EncryptedPassword` 或任何可解密密码数据。

## 解压策略

策略使用 JSON 完整保存：

```json
{
  "name": "retry",
  "outputDir": "",
  "autoRetry": true,
  "maxRetries": 3,
  "autoOpen": false
}
```

读取时执行以下规则：

- 无存储记录时返回内置默认策略。
- JSON 无效时返回明确错误。
- `maxRetries` 限制为 `0..10`。
- 策略名称必须非空。

`ExtractWithStrategy` 根据策略计算实际输出目录，执行解压或重试，并在成功且 `autoOpen` 为真时打开该目录。

## 压缩包预览与路径安全

### 结构化列表

`archive.List` 调用：

```text
7z l -slt -y <archive>
```

解析 `Path =`、`Size =`、`Attributes =`、`Modified =` 等键值字段。解析器保留包含空格和中文的路径，不依赖列宽。

### 安全验证

每次实际解压前，后端列出压缩包条目并逐项验证路径。拒绝：

- 绝对路径。
- Windows 盘符路径。
- UNC 路径。
- 规范化后逃离目标目录的路径。
- 任意 `..` 路径段。

无法列出或验证条目时，解压终止并记录失败历史，不启动 `7z x`。

安全验证应用于普通解压、自动密码尝试、手动密码解压、批量解压和策略解压。

## 前端拆分

### `App.vue`

只负责：

- 页面标题和标签导航。
- 根据当前标签渲染工作区或管理组件。
- 保存当前选中的压缩包路径，供预览标签使用。

### `ExtractWorkspace.vue`

负责解压页面布局：

- 文件选择和批量选择。
- 输出目录和策略选择。
- 日志、结果、批量结果。
- 密码弹窗。

### `useExtractWorkflow.ts`

负责：

- 单文件和批量工作流状态。
- Wails API 调用。
- 日志订阅和解除订阅。
- 密码弹窗状态。
- 结果状态和输出目录计算。

### `PasswordDialog.vue`

删除历史密码候选展示，增加默认勾选的“解压成功后记住此密码”，提交密码与保存选择。

## 错误处理

- 安全校验失败返回可识别错误，不执行解压。
- 策略 JSON 无效返回错误，不静默回退为错误配置。
- 自动密码失败仍统一返回 `password required`，不泄漏候选值。
- 保存密码失败不应把已成功解压标记为失败，但需要向前端返回明确的保存失败信息。实现中采用组合错误，提示“解压成功但密码保存失败”。
- 数据库历史记录写入失败保持当前非阻塞行为，避免覆盖主要解压结果。

## 测试与验收

Go 自动化测试覆盖：

- `-slt` 条目解析，包括空格和中文路径。
- 路径穿越、绝对路径、盘符路径和安全路径。
- 解压在安全校验失败时不启动解压命令。
- 密码保存去重。
- 手动密码按 `rememberPassword` 保存。
- 密码摘要不含密文字段。
- 策略 JSON 保存、加载、校验和 `autoOpen` 执行。

前端验证：

- TypeScript 类型检查。
- Vite 生产构建。
- 密码弹窗不展示历史密码，并提交保存选择。
- `App.vue` 只保留导航职责。

最终验证命令：

```bash
GOCACHE=/tmp/pocketzip-go-cache go test ./...
cd frontend && npm run build
```

Windows 实机仍需验证：

- DPAPI 加解密。
- 内置 `7z.exe` 的 `-slt` 输出兼容性。
- 中文路径。
- Explorer 自动打开目录。
- 拖拽与 NSIS 安装包。
