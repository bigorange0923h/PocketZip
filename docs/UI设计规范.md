# PocketUnzip UI 设计规范

## 1. 设计风格：3D Elements

PocketUnzip 采用 **3D Elements** 设计风格，通过 CSS 3D、透视、分层阴影和光晕，在二维屏幕上营造「屏幕里有实体装置」的视觉错觉。

### 1.1 风格定位

3D Elements 风格适用于需要突出科技感、创新力和高价值感的场景。它不适合信息极度密集、强调纯效率的后台，而更适合「第一印象很重要」「需要强烈视觉记忆点」的页面。

PocketUnzip 作为一款轻量解压工具，采用此风格可以：

- 传递「技术领先」「注重细节」「体验高端」的品牌印象
- 在众多解压工具中形成强烈的视觉记忆点
- 为桌面工具赋予现代感和品质感

### 1.2 视觉设计理念

3D Elements 的核心是通过「分层 + 透视 + 光影」在二维屏幕上构建一个伪 3D 空间。

页面中的模块（卡片、按钮、图标、图表）不再只是扁平方块，而是像浮在舞台上的小方盒：

- 它们有厚度、有投影、有旋转角度
- 背景通常使用深色渐变或带网格的星空式底板
- 前景元素通过明亮渐变、光带和发光边缘凸显出来
- 用户感觉自己在浏览一组「实体组件」的陈列架

---

## 2. 材质与质感

### 2.1 常用材质

| 材质类型 | 实现方式 | 应用场景 |
|---------|---------|---------|
| 玻璃（Glassmorphism） | 半透明背景、模糊（blur）、内外阴影、细边框 | 主内容卡片、弹窗、侧边栏 |
| 金属 | 高光渐变、大范围柔和阴影 | 工具栏、状态栏 |
| 亚克力 | 半透明 + 模糊 + 细微纹理 | 次级面板、悬浮提示 |
| 霓虹塑料 | 饱和渐变 + 发光边缘 | 按钮、图标、进度条 |

### 2.2 质感实现原则

所有质感均由 CSS 渐变、阴影和 3D 变换组合而成，尽量减少位图纹理的依赖，以保证在不同分辨率下的可控性。

重点元素（如主 CTA、核心功能区）叠加体积光、反射光或彩色边缘高光，让它们在画面中显得更「重」也更「近」。

---

## 3. 主题配色

### 3.1 主色调

以冷色渐变为主，营造「未来实验室 / 创意工作台」的氛围。

```css
:root {
  /* 主背景色 */
  --bg-primary: #0a0e1a;
  --bg-secondary: #111827;
  --bg-gradient: linear-gradient(135deg, #0a0e1a 0%, #1a1040 50%, #0d1f2d 100%);

  /* 主强调色 */
  --accent-primary: #6366f1;
  --accent-secondary: #8b5cf6;
  --accent-gradient: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a78bfa 100%);

  /* 辅助强调色 */
  --accent-cyan: #22d3ee;
  --accent-emerald: #34d399;
  --accent-amber: #fbbf24;

  /* 玻璃质感 */
  --glass-bg: rgba(255, 255, 255, 0.05);
  --glass-border: rgba(255, 255, 255, 0.1);
  --glass-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);

  /* 文字色 */
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --text-muted: #64748b;

  /* 发光效果 */
  --glow-primary: 0 0 20px rgba(99, 102, 241, 0.3);
  --glow-cyan: 0 0 20px rgba(34, 211, 238, 0.3);
  --glow-success: 0 0 20px rgba(52, 211, 153, 0.3);
}
```

### 3.2 色彩使用规则

| 元素类型 | 色彩方案 |
|---------|---------|
| 页面背景 | 深色渐变（深蓝、紫色、青色） |
| 主要操作按钮 | 主强调色渐变 + 发光边缘 |
| 成功状态 | 翡翠绿 + 绿色光晕 |
| 警告状态 | 琥珀色 + 暖色光晕 |
| 错误状态 | 玫瑰红 + 红色光晕 |
| 次要元素 | 玻璃质感（半透明 + 模糊） |
| 文字 | 浅灰白为主，深灰为辅 |

### 3.3 暖色点缀

少量暖色用于关键区域，形成视觉焦点：

- 解压成功提示
- 密码保存确认
- 重要操作的 CTA 按钮

---

## 4. 交互体验

### 4.1 深度反馈

交互反馈强调「深度」和「视角变化」。

**卡片悬停效果：**

```css
.card {
  transform: perspective(1000px) rotateX(0) rotateY(0) translateZ(0);
  transition: transform 0.4s ease, box-shadow 0.4s ease;
}

.card:hover {
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg) translateZ(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), var(--glow-primary);
}
```

**按钮悬停效果：**

```css
.btn-primary {
  position: relative;
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.btn-primary::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(99, 102, 241, 0.4);
}

.btn-primary:hover::before {
  left: 100%;
}
```

### 4.2 3D 模型交互

3D 模型（如立方体、图标）可以持续缓慢旋转，在悬停时暂停或加速，强调「这是一个可以被探索的物件」。

```css
.cube {
  animation: rotate 20s linear infinite;
}

.cube:hover {
  animation-play-state: paused;
  /* 或加速：animation-duration: 5s; */
}

@keyframes rotate {
  from { transform: rotateY(0) rotateX(15deg); }
  to { transform: rotateY(360deg) rotateX(15deg); }
}
```

### 4.3 动效节奏

动效节奏控制在 **0.3–0.6 秒**之间，既有重量感，又不会太拖沓。

| 交互类型 | 持续时间 | 缓动函数 |
|---------|---------|---------|
| 卡片悬停 | 0.4s | ease |
| 按钮点击 | 0.3s | ease-out |
| 弹窗打开 | 0.5s | cubic-bezier(0.16, 1, 0.3, 1) |
| 弹窗关闭 | 0.3s | ease-in |
| 3D 旋转 | 持续 | linear |
| 光带滑过 | 0.5s | ease |

---

## 5. 整体氛围

### 5.1 氛围定位

3D Elements 风格营造的是一种「未来实验室 / 创意工作台」的氛围：

- 背景像昏暗的工作室或数据机房
- 前景是一排排悬浮的设备、卡片和模块
- 影视级的光影与微妙的透视让用户觉得自己在操作一个真实的控制台或产品展示台

### 5.2 应用到 PocketUnzip

对于 PocketUnzip，这种氛围可以这样应用：

| 界面区域 | 设计手法 |
|---------|---------|
| 主窗口背景 | 深色渐变 + 微妙网格线 |
| 文件选择区 | 玻璃质感卡片 + 拖拽时的光晕反馈 |
| 解压进度 | 霓虹风格进度条 + 发光动画 |
| 密码输入 | 金属质感输入框 + 聚焦时的光边 |
| 历史记录 | 悬浮卡片列表 + 悬停时的深度变化 |
| 日志面板 | 半透明终端风格 + 扫描线效果 |
| 成功/失败状态 | 绿色/红色光晕 + 图标动画 |

---

## 6. 组件示例

### 6.1 主卡片组件

```css
.main-card {
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  border-radius: 16px;
  padding: 24px;
  box-shadow: var(--glass-shadow);
  transform: perspective(1000px) rotateX(0) rotateY(0) translateZ(0);
  transition: all 0.4s ease;
}

.main-card:hover {
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg) translateZ(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), var(--glow-primary);
  border-color: rgba(99, 102, 241, 0.3);
}
```

### 6.2 主要按钮

```css
.btn-primary {
  background: var(--accent-gradient);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.3);
  transform: translateY(0);
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.5);
}
```

### 6.3 输入框

```css
.input-field {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 10px;
  padding: 12px 16px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.input-field:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2), var(--glow-primary);
}
```

### 6.4 状态指示器

```css
.status-success {
  color: var(--accent-emerald);
  text-shadow: var(--glow-success);
}

.status-error {
  color: #f43f5e;
  text-shadow: 0 0 20px rgba(244, 63, 94, 0.3);
}

.status-processing {
  color: var(--accent-cyan);
  text-shadow: var(--glow-cyan);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
```

---

## 7. 注意事项

### 7.1 性能优化

- 使用 `transform` 和 `opacity` 进行动画，避免触发重排
- `backdrop-filter` 在某些场景下性能较差，谨慎使用
- 3D 变换使用 `will-change: transform` 提示浏览器优化

### 7.2 可访问性

- 确保文字与背景的对比度符合 WCAG AA 标准
- 为动画提供 `prefers-reduced-motion` 的降级方案
- 重要状态不仅依赖颜色，还需配合图标或文字

### 7.3 响应式适配

- 在小屏幕上减少 3D 变换的幅度
- 移动端禁用复杂的悬停效果
- 确保玻璃质感在低性能设备上仍可正常显示

---

## 8. 设计资源

### 8.1 参考色板

- 主色：#6366f1（Indigo）
- 辅助：#8b5cf6（Violet）
- 点缀：#22d3ee（Cyan）、#34d399（Emerald）、#fbbf24（Amber）

### 8.2 参考项目

- Linear.app（科技感仪表盘）
- Vercel Dashboard（简洁 3D 感）
- Raycast 官网（深色主题 + 光影）

---

*本规范基于 3D Elements 设计风格制定，作为 PocketUnzip 前端开发的视觉指导。*
