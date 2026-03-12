.PHONY: all build run test clean swag fmt lint

# 默认目标
all: fmt lint build

# 构建项目
build:
	go build -o bin/qzt-server cmd/server/main.go

# 运行项目
run:
	go run cmd/server/main.go

# 运行测试
test:
	go test -v -race -coverprofile=coverage.out ./...

# 查看测试覆盖率
coverage:
	go tool cover -html=coverage.out

# 清理构建产物
clean:
	rm -rf bin/
	rm -f coverage.out

# 生成 Swagger 文档
swag:
	swag init -g cmd/server/main.go -o docs/

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

# 安装依赖
deps:
	go mod download
	go mod tidy

# 安装开发工具
dev-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker 构建
docker-build:
	docker build -t qzt-go:latest .

# Docker 运行
docker-run:
	docker run -d --name qzt-go -p 8080:8080 qzt-go:latest
