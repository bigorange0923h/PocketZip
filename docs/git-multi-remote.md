# Git 多远程仓库推送配置

> **策略**：GitHub 为主仓库，Gitee 作为国内备份镜像。

## 一、添加两个远程源

```bash
# 添加 origin（主仓库 - GitHub）
git remote add origin https://github.com/BigOrange/PocketZip.git

# 添加 backup（国内备份 - Gitee）
git remote add backup https://gitee.com/bigorange_gitee/PocketZip.git
```

## 二、查看远程源配置

```bash
git remote -v
```

预期输出：

```
backup  https://gitee.com/bigorange_gitee/PocketZip.git (fetch)
backup  https://gitee.com/bigorange_gitee/PocketZip.git (push)
origin  https://github.com/BigOrange/PocketZip.git (fetch)
origin  https://github.com/BigOrange/PocketZip.git (push)
```

## 三、同时推送到两个远程源

### 方式一：逐个推送

```bash
git push origin main
git push backup main
```

### 方式二：配置一个 origin 同时推送两个地址

```bash
# 给 origin 添加第二个 push 地址
git remote set-url --add --push origin https://github.com/BigOrange/PocketZip.git
git remote set-url --add --push origin https://gitee.com/bigorange_gitee/PocketZip.git
```

验证配置：

```bash
git remote -v
```

预期输出：

```
origin  https://github.com/BigOrange/PocketZip.git (fetch)
origin  https://github.com/BigOrange/PocketZip.git (push)
origin  https://gitee.com/bigorange_gitee/PocketZip.git (push)
```

之后只需：

```bash
git push origin main
```

即可同时推送到两个远程仓库。

## 四、批量推送所有远程源

```bash
git push --all origin
git push --all backup
```

或者使用通配符：

```bash
git push --all
```

## 五、注意事项

1. **首次推送**需要添加 `-u` 参数关联分支：
   ```bash
   git push -u origin main
   git push -u backup main
   ```

2. **删除远程分支**时需要分别操作：
   ```bash
   git push origin --delete <分支名>
   git push backup --delete <分支名>
   ```

3. **拉取代码**始终从 GitHub（origin）拉取，Gitee 仅作备份同步：
   ```bash
   git pull origin main
   ```
