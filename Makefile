MODULE := github.com/yxrrxy/videoHub
IDL_PATH := idl
OUTPUT_DIR := output

CONFIG = $(shell go run scripts/config_tool.go $(1) $(2))

export MYSQL_USER = $(call CONFIG,mysql,username)
export MYSQL_PASSWORD = $(call CONFIG,mysql,password)
export MYSQL_DATABASE = $(call CONFIG,mysql,database)

env-up:
	@echo "MYSQL_USER: ${MYSQL_USER}"
	@echo "MYSQL_PASSWORD: ${MYSQL_PASSWORD}"
	@echo "MYSQL_DATABASE: ${MYSQL_DATABASE}"
	@docker-compose -f docker/docker-compose.yml up -d mysql

env-down:
	@docker-compose -f docker/docker-compose.yml down

.PHONY: all user video env-up env-down clean kitex-gen-%

define build_service
	@echo "Building $(1) services..."
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=$(1) \
		--build-arg SERVICE_PORT=$(call CONFIG,$(1),http_port) \
		-t $(1)-http-service .
	@docker create --name temp-$(1)-http $(1)-http-service
	@docker cp temp-$(1)-http:/app/service_binary $(OUTPUT_DIR)/$(1)_http
	@docker rm temp-$(1)-http
	@docker build \
		-f docker/baseBuild/Dockerfile \
		--build-arg SERVICE_PATH=$(1)/rpc \
		--build-arg SERVICE_PORT=$(call CONFIG,$(1),rpc_port) \
		-t $(1)-rpc-service .
	@docker create --name temp-$(1)-rpc $(1)-rpc-service
	@docker cp temp-$(1)-rpc:/app/service_binary $(OUTPUT_DIR)/$(1)_rpc
	@docker rm temp-$(1)-rpc
	@echo "Starting $(1) services..."
	@TMUX= tmux kill-session -t $(1) 2>/dev/null || true
	@TMUX= tmux new-session -d -s $(1)
	@TMUX= tmux select-window -t $(1):0
	@TMUX= tmux split-window -h -t $(1):0
	@TMUX= tmux select-pane -t $(1):0.0
	@TMUX= tmux send-keys -t $(1):0.0 '$(OUTPUT_DIR)/$(1)_http' C-m
	@TMUX= tmux select-pane -t $(1):0.1
	@TMUX= tmux send-keys -t $(1):0.1 '$(OUTPUT_DIR)/$(1)_rpc' C-m
	@TMUX= tmux attach-session -t $(1)
endef

$(OUTPUT_DIR):
	@mkdir -p $(OUTPUT_DIR)

user: $(OUTPUT_DIR)
	$(call build_service,user)

video: $(OUTPUT_DIR)
	$(call build_service,video)

clean:
	@echo "Cleaning build files and volumes..."
	@rm -rf $(OUTPUT_DIR)
	@docker-compose -f docker/docker-compose.yml down
	@docker volume rm docker_mysql_data 2>/dev/null || true
	@echo "Cleaned build files and volumes"

kitex-gen-%:
	@kitex -module "${MODULE}" ${IDL_PATH}/$*.thrift
	@go mod tidy
