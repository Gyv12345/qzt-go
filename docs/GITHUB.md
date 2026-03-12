# 推送到 GitHub 指南

## 前提条件
- ✅ Git 已安装
- ✅ 已有 GitHub 账号
- ✅ 本地仓库已提交

## 方式一：一键推送（推荐）

### 1. 创建 GitHub 仓库
1. 访问 https://github.com/new
2. 仓库名称: `qzt-go`
3. 描述: `企智通 CRM 系统 Go 语言后端`
4. 选择 **Public** 或 **Private**
5. **不要**勾选 "Add a README file"
6. 点击 "Create repository"

### 2. 复制仓库地址
创建后会看到类似这样的地址：
```
https://github.com/你的用户名/qzt-go.git
```

### 3. 运行推送脚本
```powershell
cd D:\ccproject\qzt-go
.\scripts\push-to-github.ps1
```
然后粘贴仓库地址。

## 方式二：手动推送

```bash
# 1. 添加远程仓库
git remote add origin https://github.com/你的用户名/qzt-go.git

# 2. 推送
git push -u origin master
```

## 方式三：SSH（推荐长期使用）

### 1. 生成 SSH 密钥
```bash
ssh-keygen -t ed25519 -C "your_email@example.com"
```

### 2. 添加到 GitHub
1. 复制 `~/.ssh/id_ed25519.pub` 内容
2. 访问 https://github.com/settings/keys
3. 点击 "New SSH key"，粘贴内容

### 3. 使用 SSH 地址
```bash
git remote set-url origin git@github.com:你的用户名/qzt-go.git
git push -u origin master
```

## ✅ 验证

推送成功后访问：
```
https://github.com/你的用户名/qzt-go
```

应该能看到所有代码！
