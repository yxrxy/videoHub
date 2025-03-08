MODULE := github.com/yxrrxy/videoHub
IDL_PATH := idl
OUTPUT_DIR := output

CONFIG = $(shell go run scripts/config_tool.go $(1) $(2))

# 导出 MySQL 环境变量
export MYSQL_PASSWORD = $(call CONFIG,mysql,password)
export MYSQL_DATABASE = $(call CONFIG,mysql,database)

.PHONY: user env-up env-down clean kitex-gen-%

user:
	@echo "Building user service..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_NAME=user \
		--build-arg SERVICE_PORT=$(call CONFIG,user,Port) \
		-t user-service .
	@docker create --name temp-user user-service
	@docker cp temp-user:/app/user_service $(OUTPUT_DIR)/
	@docker rm temp-user
	@echo "Starting user service from $(OUTPUT_DIR)/user_service..."
	@$(OUTPUT_DIR)/user_service

env-up:
	@docker-compose -f docker/docker-compose.yml up -d mysql

env-down:
	@docker-compose -f docker/docker-compose.yml down

clean:
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleaned build files"

kitex-gen-%:
	@kitex -module "${MODULE}" ${IDL_PATH}/$*.thrift
	@go mod tidy
