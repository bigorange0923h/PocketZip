#!/bin/bash

echo "Building PocketZip..."

# 构建前端
cd frontend
npm install
npm run build
cd ..

# 构建 Wails 应用
wails build

# 复制 7z.exe 到输出目录（如果存在）
if [ -d "third_party/7zip" ]; then
    cp -r third_party/7zip/* build/bin/ 2>/dev/null || true
fi

echo "Build complete!"
echo "Output: build/bin/"
