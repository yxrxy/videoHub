MODULE := github.com/yxrrxy/videoHub
IDL_PATH := idl
OUTPUT_DIR := output

CONFIG = $(shell go run scripts/config_tool.go $(1) $(2))

# 导出 MySQL 环境变量
export MYSQL_USER = $(call CONFIG,mysql,username)
export MYSQL_PASSWORD = $(call CONFIG,mysql,password)
export MYSQL_DATABASE = $(call CONFIG,mysql,database)

# 在 Makefile 中添加调试命令
env-up:
	@echo "MYSQL_USER: ${MYSQL_USER}"
	@echo "MYSQL_PASSWORD: ${MYSQL_PASSWORD}"
	@echo "MYSQL_DATABASE: ${MYSQL_DATABASE}"
	@docker-compose -f docker/docker-compose.yml up -d mysql

env-down:
	@docker-compose -f docker/docker-compose.yml down

.PHONY: all user user-http user-rpc env-up env-down clean kitex-gen-%

# 添加默认目标
user: user-http user-rpc

# 添加构建目录检查
$(OUTPUT_DIR):
	@mkdir -p $(OUTPUT_DIR)

user-http: $(OUTPUT_DIR)
	@echo "Building user HTTP service..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=user \
		--build-arg SERVICE_PORT=$(call CONFIG,user,http_port) \
		-t user-http-service .
	@docker create --name temp-user-http user-http-service
	@docker cp temp-user-http:/app/service_binary $(OUTPUT_DIR)/user_http
	@docker rm temp-user-http
	@echo "Starting user HTTP service from $(OUTPUT_DIR)/user_http..."
	@$(OUTPUT_DIR)/user_http

user-rpc: $(OUTPUT_DIR)
	@echo "Building user RPC service..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=user/rpc \
		--build-arg SERVICE_PORT=$(call CONFIG,user,rpc_port) \
		-t user-rpc-service .
	@docker create --name temp-user-rpc user-rpc-service
	@docker cp temp-user-rpc:/app/service_binary $(OUTPUT_DIR)/user_rpc
	@docker rm temp-user-rpc
	@echo "Starting user RPC service from $(OUTPUT_DIR)/user_rpc..."
	@$(OUTPUT_DIR)/user_rpc

clean:
	@echo "Cleaning build files and volumes..."
	@rm -rf $(OUTPUT_DIR)
	@docker-compose -f docker/docker-compose.yml down
	@docker volume rm docker_mysql_data 2>/dev/null || true
	@echo "Cleaned build files and volumes"

kitex-gen-%:
	@kitex -module "${MODULE}" ${IDL_PATH}/$*.thrift
	@go mod tidy
