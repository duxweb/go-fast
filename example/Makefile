.PHONY: build run generate clean

# 默认目标
all: generate build

# 生成注解
generate:
	@echo "Generating annotations..."
	@go generate ./...

# 构建项目（先生成）
build: generate
	@echo "Building..."
	@go build -o app .

# 运行项目（先生成）
run: generate
	@echo "Running..."
	@go run .

# 清理
clean:
	@rm -rf app
	@rm -rf runtime/annotations.go

# 仅构建，不生成
build-only:
	@go build -o app .

# 仅运行，不生成
run-only:
	@go run .