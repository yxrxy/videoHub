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

### 3. 互动服务 (InteractionService)
- Like
- GetLikes
- Comment
- GetComments
- DeleteComment
- LikeComment

### 4. 社交服务 (SocialService)
- SendPrivateMessage
- GetPrivateMessages
- CreateChatRoom
- WsChat

## 技术栈

- **语言与框架**: Go 1.21, Kitex(RPC), Hertz(HTTP), Thrift(IDL)
- **数据存储**: MySQL(GORM), Redis
- **中间件**: Kafka, etcd, JWT
- **监控**: OpenTelemetry
- **部署**: Docker, Docker Compose

## 系统架构

基于微服务架构，划分为用户、视频、互动和社交四个核心服务，通过Gateway统一入口处理请求。服务间通过Kitex RPC框架通信，使用etcd进行服务发现。

## 核心实现

### 缓存策略
- 热门视频使用Redis缓存
- 用户令牌缓存
- 点赞、访问量等计数器使用Redis实现
- 热门排行榜使用Redis Sorted Set

### 异步处理
- 视频上传后通过Kafka消息队列异步处理
- 视频封面生成、元数据提取(时长、分辨率)
- 异步数据写入与统计更新

### 社交功能
- 基于WebSocket的实时聊天
- 私信和群聊支持

### 数据处理
- 视频分类与标签系统

## 开发与部署

项目使用Makefile自动化构建，支持单独编译和部署各个服务。环境依赖通过Docker Compose管理，包括MySQL、Redis、etcd和Kafka等。

测试用例请参考Apifox链接。