MODULE := github.com/yxrrxy/videoHub
IDL_PATH := idl
OUTPUT_DIR := output

CONFIG = $(shell go run scripts/config_tool.go $(1) $(2))

export MYSQL_USER = $(call CONFIG,mysql,username)
export MYSQL_PASSWORD = $(call CONFIG,mysql,password)
export MYSQL_ROOT_PASSWORD = $(call CONFIG,mysql,password) 
export MYSQL_DATABASE = $(call CONFIG,mysql,database)
export REDIS_PASSWORD = $(call CONFIG,redis,password)

env-up:
	@echo "Creating storage directories..."
	@mkdir -p src/storage/videos
	@mkdir -p src/storage/covers
	@mkdir -p src/storage/avatars
	@mkdir -p src/storage/chat
	@chmod -R 755 src/storage
#	@echo "Installing ffmpeg..."
#	@which ffmpeg || sudo apt-get update && sudo apt-get install -y ffmpeg
	@echo "Starting docker services..."
	@docker-compose -f docker/docker-compose.yml up -d mysql redis etcd

env-down:
	@docker-compose -f docker/docker-compose.yml down

.PHONY: all user video social interaction gateway env-up env-down clean kitex-gen-% start

# 构建服务的通用函数
define build_service
	@echo "Building $(1) RPC service..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=$(1)/rpc \
		--build-arg SERVICE_PORT=$(call CONFIG,$(1),rpc_port) \
		-t $(1)-rpc-service .
	@docker create --name temp-$(1)-rpc $(1)-rpc-service
	@docker cp temp-$(1)-rpc:/app/service_binary $(OUTPUT_DIR)/$(1)_rpc
	@docker rm temp-$(1)-rpc
	@echo "Starting $(1) RPC service..."
	@TMUX= tmux kill-session -t $(1) 2>/dev/null || true
	@TMUX= tmux new-session -d -s $(1)
	@TMUX= tmux send-keys -t $(1):0.0 '$(OUTPUT_DIR)/$(1)_rpc' C-m
endef

$(OUTPUT_DIR):
	@mkdir -p $(OUTPUT_DIR)

# 用户服务
user: $(OUTPUT_DIR)
	$(call build_service,user)

# 视频服务
video: $(OUTPUT_DIR)
	$(call build_service,video)

# 社交服务
social: $(OUTPUT_DIR)
	$(call build_service,social)

# 互动服务
interaction: $(OUTPUT_DIR)
	$(call build_service,interaction)

# 网关服务
gateway: $(OUTPUT_DIR)
	@echo "Building gateway service..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=gateway \
		--build-arg SERVICE_PORT=$(call CONFIG,gateway,port) \
		-t gateway-service .
	@docker create --name temp-gateway gateway-service
	@docker cp temp-gateway:/app/service_binary $(OUTPUT_DIR)/gateway
	@docker rm temp-gateway
	@echo "Starting gateway service..."
	@TMUX= tmux kill-session -t gateway 2>/dev/null || true
	@TMUX= tmux new-session -d -s gateway
	@TMUX= tmux send-keys -t gateway:0.0 '$(OUTPUT_DIR)/gateway' C-m

# 添加启动服务的命令
start:
	@echo "All services are running in tmux sessions:"
	@echo "- user service: tmux attach-session -t user"
	@echo "- video service: tmux attach-session -t video"
	@echo "- social service: tmux attach-session -t social"
	@echo "- interaction service: tmux attach-session -t interaction"
	@echo "- gateway service: tmux attach-session -t gateway"
	@echo "Use 'tmux attach-session -t <service_name>' to view logs"

clean:
	@echo "Cleaning build files and volumes..."
	@rm -rf $(OUTPUT_DIR)
	@docker-compose -f docker/docker-compose.yml down
	@docker volume rm docker_mysql_data docker_redis_data docker_etcd_data 2>/dev/null || true
	@echo "Cleaned build files and volumes"

kitex-gen-%:
	@kitex -module "${MODULE}" ${IDL_PATH}/$*.thrift
	@go mod tidy

# 添加启动所有服务的命令
all: env-up user video social interaction gateway start
