<template>
  <div class="chatroom-list-container">
    <div class="list-header">
      <h2>聊天室列表</h2>
      <el-button type="primary" @click="showCreateRoomDialog">
        创建聊天室
      </el-button>
    </div>

    <div class="room-list">
      <el-card v-for="room in rooms" :key="room.id" class="room-card">
        <div class="room-info">
          <el-avatar :size="50">{{ room.name[0] }}</el-avatar>
          <div class="room-details">
            <h3>{{ room.name }}</h3>
            <p class="room-meta">
              <span>创建者: {{ room.creator_id }}</span>
              <span>{{ room.type === 1 ? '私聊' : '群聊' }}</span>
            </p>
          </div>
        </div>
        <div class="room-actions">
          <el-button type="primary" @click="enterRoom(room)">进入聊天室</el-button>
          <el-button v-if="isRoomOwner(room)" type="danger" plain @click="deleteRoom(room)">
            删除
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 创建聊天室对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建聊天室"
      width="500px"
    >
      <el-form :model="roomForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="roomForm.name" placeholder="请输入聊天室名称"></el-input>
        </el-form-item>
        <el-form-item label="类型" required>
          <el-radio-group v-model="roomForm.type">
            <el-radio :label="1">私聊</el-radio>
            <el-radio :label="2">群聊</el-radio>
          </el-radio-group>
          <div class="type-info">
            <small v-if="roomForm.type === 1">私聊仅限您和一个好友之间的交流</small>
            <small v-else>群聊可以添加多个好友一起聊天</small>
          </div>
        </el-form-item>
        <el-form-item label="成员" v-if="roomForm.type === 2">
          <el-select
            v-model="roomForm.members"
            multiple
            filterable
            placeholder="请选择成员"
          >
            <el-option
              v-for="friend in friends"
              :key="friend.id"
              :label="friend.username"
              :value="friend.id"
            >
              <span>{{ friend.username }}</span>
            </el-option>
          </el-select>
          <div class="friend-info" v-if="!friends.some(f => f.id > 3)">
            <el-alert type="info" show-icon :closable="false">
              未找到好友，已自动添加默认成员
            </el-alert>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="createRoom" :loading="creating">
            创建
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { chatRoomApi, friendApi } from '@/api/social'

const router = useRouter()
const userStore = useUserStore()
const rooms = ref([])
const friends = ref([])
const createDialogVisible = ref(false)
const creating = ref(false)

const roomForm = ref({
  name: '',
  type: 2,
  members: []
})

// 获取聊天室列表
const fetchRooms = async () => {
  try {
    const res = await chatRoomApi.getChatRooms(
      {
        user_id: userStore.userInfo.id,
        page: 1,
        size: 50
      }
    )
    rooms.value = res.data
  } catch (error) {
    console.error('获取聊天室列表失败:', error)
    ElMessage.error('获取聊天室列表失败')
  }
}

// 获取好友列表
const fetchFriends = async () => {
  if (!userStore.userInfo?.id) {
    console.error('用户未登录')
    return
  }

  try {
    const res = await friendApi.getFriends({
      user_id: userStore.userInfo.id,
      page: 1,
      size: 50
    })
    console.log('获取到好友列表:', res)
    if (res.code === 10000) {
      // 处理好友列表数据
      if (res.data && res.data.FriendList && res.data.FriendList.length > 0) {
        // 解析好友数据，格式化为下拉框所需格式
        friends.value = res.data.FriendList.map(friend => {
          // 这里的属性名需要根据实际API返回结构调整
          const friendId = friend.user_id === userStore.userInfo.id ? friend.friend_id : friend.user_id
          return {
            id: friendId,
            username: friend.remark || `用户${friendId}` // 如果有备注则显示备注，否则显示用户ID
          }
        })
      } else {
        // 没有好友时，创建默认的1,2,3成员
        friends.value = [
          { id: 1, username: '用户1' },
          { id: 2, username: '用户2' },
          { id: 3, username: '用户3' }
        ]
        // 默认选中这些成员
        if (roomForm.value.type === 2) {
          roomForm.value.members = [1, 2, 3]
        }
      }
    } else {
      console.error('获取好友列表失败:', res.message)
      // 获取失败时也使用默认成员
      friends.value = [
        { id: 1, username: '用户1' },
        { id: 2, username: '用户2' },
        { id: 3, username: '用户3' }
      ]
      if (roomForm.value.type === 2) {
        roomForm.value.members = [1, 2, 3]
      }
    }
  } catch (error) {
    console.error('获取好友列表失败:', error)
    ElMessage.error('获取好友列表失败')
    // 出错时也使用默认成员
    friends.value = [
      { id: 1, username: '用户1' },
      { id: 2, username: '用户2' },
      { id: 3, username: '用户3' }
    ]
    if (roomForm.value.type === 2) {
      roomForm.value.members = [1, 2, 3]
    }
  }
}

// 显示创建聊天室对话框
const showCreateRoomDialog = () => {
  createDialogVisible.value = true
  roomForm.value = {
    name: '',
    type: 2,
    members: []
  }
  fetchFriends()
}

// 创建聊天室
const createRoom = async () => {
  if (!roomForm.value.name.trim()) {
    ElMessage.warning('请输入聊天室名称')
    return
  }

  if (roomForm.value.type === 2 && roomForm.value.members.length === 0) {
    ElMessage.warning('请选择群聊成员')
    return
  }

  creating.value = true
  try {
    console.log('创建聊天室参数:', {
      name: roomForm.value.name,
      creator_id: userStore.userInfo.id,
      type: roomForm.value.type,
      member_ids: roomForm.value.members
    })
    
    await chatRoomApi.createChatRoom({
      name: roomForm.value.name,
      creator_id: userStore.userInfo.id,
      type: roomForm.value.type,
      member_ids: roomForm.value.members
    })
    
    ElMessage.success('创建聊天室成功')
    createDialogVisible.value = false
    fetchRooms()
  } catch (error) {
    console.error('创建聊天室失败:', error)
    ElMessage.error('创建聊天室失败')
  }
  creating.value = false
}

// 进入聊天室
const enterRoom = (room) => {
  router.push(`/social/chatroom/${room.id}`)
}

// 删除聊天室
const deleteRoom = async (room) => {
  try {
    await ElMessageBox.confirm('确定要删除该聊天室吗？', '提示', {
      type: 'warning'
    })
    
    await chatRoomApi.deleteChatRoom(room.id)
    ElMessage.success('删除聊天室成功')
    fetchRooms()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除聊天室失败:', error)
      ElMessage.error('删除聊天室失败')
    }
  }
}

// 判断是否是聊天室所有者
const isRoomOwner = (room) => {
  return room.owner_id === userStore.userInfo?.id
}

onMounted(() => {
  fetchRooms()
})
</script>

<style scoped>
.chatroom-list-container {
  padding: 20px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.list-header h2 {
  margin: 0;
}

.room-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.room-card {
  border-radius: 8px;
}

.room-info {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.room-details {
  flex: 1;
}

.room-details h3 {
  margin: 0 0 5px 0;
}

.room-meta {
  margin: 0;
  color: #666;
  font-size: 14px;
  display: flex;
  gap: 10px;
}

.room-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 新增样式 */
.type-info {
  margin-top: 5px;
  color: #909399;
}

.friend-info {
  margin-top: 10px;
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-dialog__body) {
  padding-top: 10px;
}
</style> 