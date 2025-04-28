# VideoHub

API文档: https://2pr5l6um26.apifox.cn

## 服务接口列表

### 1. 用户服务 (UserService)
- Register
- Login
- GetUserInfo
- UploadAvatar
- RefreshToken

### 2. 视频服务 (VideoService)
- Publish
- List
- Detail
- GetHotVideos
- Delete
- IncrementVisitCount
- IncrementLikeCount
- Search
- SemanticSearch

### 3. 互动服务 (InteractionService)
- Like
- GetLikes
- Comment
- GetComments
- DeleteComment
- LikeComment

### 4. 社交服务 (SocialService)
- PrivateChat
- RoomChat
- CreateChatRoom

## 技术栈

- **语言与框架**: Go 1.21, Kitex(RPC), Hertz(HTTP), Thrift(IDL)
- **数据存储**: MySQL(GORM), Redis, Elasticsearch
- **中间件**: Kafka, etcd, JWT
- **AI与搜索**: OpenAI API, 向量检索
- **监控**: OpenTelemetry
- **部署**: Docker, Docker Compose

## 系统架构

基于微服务架构，划分为用户、视频、互动和社交四个核心服务，通过Gateway统一入口处理请求。服务间通过Kitex RPC框架通信，使用etcd进行服务发现。

## 核心实现

### 视频服务
- 基于OpenAI的语义搜索功能
- 视频元数据提取和标签管理
- 热门视频排行榜（访问量、点赞数）

### 缓存策略
- 热门视频使用Redis缓存
- 用户令牌缓存
- 点赞、访问量等计数器使用Redis实现
- 热门排行榜使用Redis Sorted Set
- 分服务的Redis DB隔离（用户服务DB0、视频服务DB1、社交服务DB2）

### 异步处理
- 视频上传后通过Kafka消息队列异步处理
- 视频封面生成、元数据提取(时长、分辨率)
- 异步数据写入与统计更新
- 异步更新Elasticsearch索引和向量数据库

### 社交功能
- 基于WebSocket的实时聊天
- 支持私聊消息和聊天室
- 聊天记录持久化存储

### 搜索功能
- 基于Elasticsearch的全文搜索
- 集成OpenAI API的语义搜索
- 支持多字段复合查询
- 搜索结果缓存优化

### 安全性
- JWT用户认证
- API访问频率限制
- 视频内容安全检查
- 敏感信息加密存储

## 开发与部署

项目使用Makefile自动化构建，支持单独编译和部署各个服务。环境依赖通过Docker Compose管理，包括MySQL、Redis、etcd、Elasticsearch和Kafka等。

### 环境要求
- Docker 20.10+
- Docker Compose v2.0+
- Go 1.21+
- Make

### 快速开始

1. 克隆项目
```bash
git clone https://github.com/yxrxy/videoHub.git
cd videoHub
```

2. 配置环境变量
```bash
cp docker/env/example.env docker/env/.env
# 修改配置，包括：
# - openai keys, upyun keys (将config.ymal.example重命名为cinfig.yaml后修改)
```

3. 启动基础服务
```bash
make env-up 
# 启动 MySQL、Redis、Elasticsearch、Kafka、etcd
```

4. 编译并启动服务
```
make user        # 启动用户服务
make video       # 启动视频服务
make interaction #启动评论，点赞(视频交互)服务
make social      # 启动社交服务
make gateway     # 启动网关服务
```

### 目录结构
```
.
├── app                     # 应用服务目录
│   ├── gateway            # API 网关服务
│   │   ├── handler       # 请求处理器
│   │   ├── model         # 数据模型
│   │   ├── mw            # 中间件
│   │   ├── pack         # 数据打包工具
│   │   ├── router       # 路由配置
│   │   ├── rpc          # RPC 客户端
│   │   └── ws           # WebSocket 处理
│   ├── interaction       # 互动服务
│   │   ├── controllers  # 控制器
│   │   ├── domain       # 领域模型
│   │   ├── infrastructure # 基础设施
│   │   └── usecase      # 用例实现
│   ├── social           # 社交服务
│   │   ├── controllers  # 控制器
│   │   ├── domain       # 领域模型
│   │   ├── infrastructure # 基础设施
│   │   └── usecase      # 用例实现
│   ├── user             # 用户服务
│   │   ├── controllers  # 控制器
│   │   ├── domain       # 领域模型
│   │   ├── infrastructure # 基础设施
│   │   └── usecase      # 用例实现
│   └── video            # 视频服务
│       ├── controllers  # 控制器
│       ├── domain       # 领域模型
│       ├── infrastructure # 基础设施
│       └── usecase      # 用例实现
├── cmd                    # 主程序入口
│   ├── gateway           # 网关服务入口
│   ├── interaction       # 互动服务入口
│   ├── social            # 社交服务入口
│   ├── user              # 用户服务入口
│   └── video             # 视频服务入口
├── config                # 配置文件目录
│   ├── elasticsearch    # ES 配置
│   ├── kibana          # Kibana 配置
│   ├── otel-collector  # OpenTelemetry 配置
│   └── sql             # SQL 初始化脚本
├── docker               # Docker 相关配置
│   ├── data            # 数据持久化目录
│   ├── env             # 环境变量配置
│   └── script          # 部署和维护脚本
├── idl                  # 接口定义文件
│   ├── api             # API 接口定义
│   └── model.thrift    # 公共模型定义
├── kitex_gen           # Kitex 生成的代码
├── pkg                 # 公共包
│   ├── base           # 基础组件
│   ├── constants      # 常量定义
│   ├── errno          # 错误码
│   ├── jwt            # JWT 工具
│   ├── kafka          # Kafka 工具
│   ├── middleware     # 中间件
│   ├── storage        # 存储工具
│   ├── upyun          # 又拍云工具
│   └── utils          # 工具函数
└── src                # 资源文件
    └── storage        # 存储目录
        ├── avatars    # 头像存储
        ├── chat       # 聊天文件
        ├── covers     # 视频封面
        └── videos     # 视频文件
```

测试用例请参考Apifox链接。