<template>
  <div class="chat-window">
    <!-- 聊天头部 -->
    <div class="chat-header">
      <div class="friend-info">
        <span class="friend-name">{{ friendInfo?.username || `用户${toUserId}` }}</span>
        <span class="online-status">在线</span>
      </div>
    </div>

    <!-- 消息列表 -->
    <div class="message-list" ref="messageList">
      <el-scrollbar>
        <div v-if="loading" class="loading-wrapper">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        <template v-else>
          <div
            v-for="message in messages"
            :key="message.id"
            class="message-item"
            :class="{ 'message-self': message.sender_id === currentUserId }"
          >
            <el-avatar
              :src="message.sender?.avatar || ''"
              :size="32"
            >{{ getAvatarText(message.sender) }}</el-avatar>
            <div class="message-content">
              <div class="message-text">{{ message.content }}</div>
              <div class="message-time">{{ formatTime(message.created_at) }}</div>
            </div>
          </div>
        </template>
      </el-scrollbar>
    </div>

    <!-- 输入框 -->
    <div class="message-input">
      <el-input
        v-model="messageText"
        type="textarea"
        :rows="2"
        placeholder="输入消息..."
        @keyup.enter.native="sendMessage"
      ></el-input>
      <el-button type="primary" @click="sendMessage" :loading="sending">
        发送
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useUserStore } from '@/store/user'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'
import { Loading } from '@element-plus/icons-vue'
import { wsService, friendApi } from '@/api/social'

const route = useRoute()
const userStore = useUserStore()
const currentUserId = ref(userStore.userInfo?.id)
const toUserId = ref(parseInt(route.params.userId))
const friendInfo = ref(null)

const messages = ref([])
const messageText = ref('')
const sending = ref(false)
const messageList = ref(null)
const loading = ref(false)

// 获取好友信息
const getFriendInfo = async () => {
  try {
    const res = await friendApi.getFriends({
      user_id: currentUserId.value
    })
    if (res.code === 10000 && res.data?.FriendList) {
      const friend = res.data.FriendList.find(f => {
        const friendId = f.user_id === currentUserId.value ? f.friend_id : f.user_id
        return friendId === toUserId.value
      })
      if (friend) {
        friendInfo.value = {
          id: toUserId.value,
          username: friend.remark || `用户${toUserId.value}`,
          avatar: friend.avatar
        }
      }
    }
  } catch (error) {
    console.error('获取好友信息失败:', error)
  }
}

// 获取头像显示文本
const getAvatarText = (sender) => {
  if (!sender) return 'U'
  if (sender.username) return sender.username.charAt(0).toUpperCase()
  return `U${sender.id || ''}`
}

// 发送消息
const sendMessage = async () => {
  if (!messageText.value.trim()) return
  
  sending.value = true
  try {
    const content = messageText.value.trim()
    wsService.sendPrivateMessage(toUserId.value, content)
    // 立即在本地显示发送的消息
    messages.value.push({
      id: Date.now(),
      sender_id: currentUserId.value,
      receiver_id: toUserId.value,
      content: content,
      type: 0,
      created_at: Date.now(),
      sender: {
        id: currentUserId.value,
        username: userStore.userInfo?.username || '',
        avatar: userStore.userInfo?.avatar || ''
      }
    })
    messageText.value = ''
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送消息失败')
  }
  sending.value = false
}

// 处理收到的WebSocket消息
const handleWsMessage = (message) => {
  console.log('收到消息:', message)
  if (message.type === 'history') {
    // 处理历史消息
    const historyMessages = message.messages || []
    messages.value = historyMessages.map(msg => {
      // 尝试解析content字段中的实际消息内容
      let parsedContent = msg.content
      try {
        // 检查content是否包含map[content:xxx type:chat]格式
        if (msg.content.startsWith('map[')) {
          const contentMatch = msg.content.match(/content:([^\s\]]+)/)
          if (contentMatch) {
            parsedContent = contentMatch[1]
          }
        }
      } catch (error) {
        console.warn('解析消息内容失败:', error)
      }

      const isSelf = msg.sender_id === currentUserId.value
      return {
        id: msg.id,
        sender_id: msg.sender_id,
        receiver_id: msg.receiver_id,
        content: parsedContent,
        type: msg.type,
        created_at: msg.created_at || Date.now(),
        sender: {
          id: msg.sender_id,
          username: isSelf ? 
            userStore.userInfo?.username : friendInfo.value?.username || `用户${msg.sender_id}`,
          avatar: isSelf ? 
            userStore.userInfo?.avatar : friendInfo.value?.avatar || ''
        }
      }
    }).sort((a, b) => a.id - b.id) // 按消息ID升序排序

    nextTick(() => {
      scrollToBottom()
    })
  } else if (message.type === 'private') {
    // 如果不是自己发送的消息才添加到列表
    if (message.from !== currentUserId.value) {
      let parsedContent = message.content
      try {
        // 检查content是否包含map[content:xxx type:chat]格式
        if (message.content.startsWith('map[')) {
          const contentMatch = message.content.match(/content:([^\s\]]+)/)
          if (contentMatch) {
            parsedContent = contentMatch[1]
          }
        }
      } catch (error) {
        console.warn('解析消息内容失败:', error)
      }

      messages.value.push({
        id: Date.now(),
        sender_id: message.from,
        receiver_id: currentUserId.value,
        content: parsedContent,
        type: 0,
        created_at: message.timestamp * 1000 || Date.now(),
        sender: {
          id: message.from,
          username: friendInfo.value?.username || `用户${message.from}`,
          avatar: friendInfo.value?.avatar || ''
        }
      })
      nextTick(() => {
        scrollToBottom()
      })
    }
  } else if (message.type === 'system') {
    // 处理系统消息
    ElMessage.info(message.content || message.message)
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messageList.value) {
    const scrollbar = messageList.value.querySelector('.el-scrollbar__wrap')
    scrollbar.scrollTop = scrollbar.scrollHeight
  }
}

onMounted(async () => {
  loading.value = true
  try {
    // 获取好友信息
    await getFriendInfo()
    
    // 连接WebSocket
    const token = userStore.token
    if (token && toUserId.value) {
      // 连接私聊WebSocket
      wsService.connect(token, -1, toUserId.value)
      wsService.setMessageHandler(handleWsMessage)
    }
  } catch (error) {
    console.error('初始化聊天失败:', error)
    ElMessage.error('初始化聊天失败')
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  // 关闭WebSocket连接
  wsService.close()
})
</script>

<style scoped>
.chat-window {
  width: 800px;
  margin: 20px auto;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: calc(100vh - 120px);
}

.chat-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  background: #fff;
  border-radius: 8px 8px 0 0;
}

.friend-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.friend-name {
  font-size: 16px;
  font-weight: 500;
}

.online-status {
  font-size: 12px;
  color: #67c23a;
}

.message-list {
  flex: 1;
  overflow: hidden;
  padding: 20px;
  margin-bottom: 10px;
}

.message-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
}

.message-item.message-self {
  flex-direction: row-reverse;
}

.message-content {
  margin: 0 10px;
  max-width: 60%;
}

.message-text {
  background: #f4f4f5;
  padding: 10px 15px;
  border-radius: 4px;
  margin-bottom: 5px;
  word-break: break-all;
}

.message-self .message-text {
  background: #ecf5ff;
}

.message-time {
  font-size: 12px;
  color: #999;
  text-align: right;
}

.message-input {
  padding: 10px 20px;
  border-top: 1px solid #eee;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  background: #fff;
}

.message-input .el-input {
  flex: 1;
}

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #909399;
}

.loading-wrapper .el-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

:deep(.el-scrollbar__wrap) {
  overflow-x: hidden;
}
</style> 