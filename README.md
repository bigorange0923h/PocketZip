# PocketUnzip

PocketUnzip 是一款面向 Windows 的轻量解压工具，第一阶段使用 `7z.exe` 作为解压内核，重点提供：常见格式解压、密码记忆、自动尝试历史密码、解压历史和实时日志。

## 技术路线

- 桌面框架：Go + Wails
- 解压内核：7z.exe
- 本地数据库：SQLite
- 密码加密：Windows DPAPI
- 前端：Vue / React 均可，建议 MVP 使用 Vue + TypeScript

## MVP 功能

- 选择或拖拽压缩包
- 自动创建同名输出目录
- 调用 7z.exe 解压
- 加密压缩包密码输入
- 成功密码加密保存
- 下次自动尝试历史密码
- 解压历史记录
- 实时解压日志

## 项目目录

```text
PocketUnzip/
├─ cmd/pocketunzip/          # 程序入口
├─ internal/app/             # 应用编排层
├─ internal/archive/         # 7z 解压适配层
├─ internal/password/        # 密码匹配与保存
├─ internal/history/         # 解压历史
├─ internal/config/          # 本地配置
├─ internal/security/        # DPAPI、安全工具
├─ frontend/                 # 桌面前端
├─ third_party/7zip/         # 7z.exe、7z.dll 放置目录
├─ docs/                     # 需求与架构文档
└─ scripts/                  # 构建脚本
```

## 下一步

建议先完成最小闭环：选择文件 → 自动创建目录 → 调用 7z.exe → 展示日志 → 记录历史。
