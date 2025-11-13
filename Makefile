# Go Ethereum Backend Makefile

.PHONY: build test clean run demo test-all abigen help

# 构建应用
build:
	@echo "Building Go Ethereum Backend..."
	go build -o bin/eth-backend cmd/main.go

# 运行应用
run:
	@echo "Starting Go Ethereum Backend..."
	go run cmd/main.go

# 运行演示程序
demo:
	@echo "Running blockchain demo..."
	go run cmd/demo.go

# 运行功能测试
test-all:
	@echo "Running all functional tests..."
	go run cmd/test_all.go

# 运行单元测试
test:
	@echo "Running unit tests..."
	cd tests && go test -v

# 安装 abigen 工具
abigen:
	@echo "Installing abigen tool..."
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 编译智能合约
compile-contracts:
	@echo "Compiling smart contracts..."
	node scripts/compile_sol.js

# 生成 Go 绑定代码
gen-bindings:
	@echo "Generating Go bindings..."
	abigen --abi=contracts/build/Counter.abi --bin=contracts/build/Counter.bin --pkg=counter --out=pkg/eth/counter.go

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf bin/
	rm -f go.sum

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy

# 显示帮助
help:
	@echo "Available commands:"
	@echo "  build              - 构建主应用程序"
	@echo "  run                - 运行主应用程序"
	@echo "  demo               - 运行区块链演示程序"
	@echo "  test-all           - 运行功能测试"
	@echo "  test               - 运行单元测试"
	@echo "  abigen             - 安装 abigen 工具"
	@echo "  compile-contracts  - 编译智能合约"
	@echo "  gen-bindings       - 生成 Go 绑定代码"
	@echo "  clean              - 清理构建文件"
	@echo "  deps               - 安装依赖"
	@echo "  help               - 显示此帮助信息"