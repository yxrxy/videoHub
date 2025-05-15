<template>
  <div class="friends-container">
    <!-- 好友列表 -->
    <div class="friends-list">
      <div class="friends-header">
        <h3>我的好友</h3>
        <el-button type="primary" size="small" @click="showAddFriendDialog">
          添加好友
        </el-button>
      </div>

      <!-- 好友分组 -->
      <div class="friends-content">
        <el-scrollbar>
          <el-collapse v-model="activeGroups">
            <!-- 在线好友 -->
            <el-collapse-item name="online">
              <template #title>
                <div class="group-title">
                  <span>在线好友</span>
                  <span class="friend-count">{{ onlineFriends.length }}</span>
                </div>
              </template>
              <div
                v-for="friend in onlineFriends"
                :key="friend.id"
                class="friend-item"
              >
                <el-avatar :src="friend.avatar" :size="40">
                  {{ getAvatarText(friend) }}
                </el-avatar>
                <div class="friend-info">
                  <div class="friend-name">
                    {{ friend.username }}
                    <span class="friend-remark" v-if="friend.remark">
                      ({{ friend.remark }})
                    </span>
                  </div>
                  <div class="friend-status online">在线</div>
                </div>
                <div class="friend-actions">
                  <el-button type="primary" link @click="startChat(friend)">
                    发消息
                  </el-button>
                  <el-dropdown>
                    <el-button link>
                      <el-icon><More /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="editRemark(friend)">
                          修改备注
                        </el-dropdown-item>
                        <el-dropdown-item @click="deleteFriend(friend)" class="text-danger">
                          删除好友
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </el-collapse-item>

            <!-- 离线好友 -->
            <el-collapse-item name="offline">
              <template #title>
                <div class="group-title">
                  <span>离线好友</span>
                  <span class="friend-count">{{ offlineFriends.length }}</span>
                </div>
              </template>
              <div
                v-for="friend in offlineFriends"
                :key="friend.id"
                class="friend-item"
              >
                <el-avatar :src="friend.avatar" :size="40">
                  {{ getAvatarText(friend) }}
                </el-avatar>
                <div class="friend-info">
                  <div class="friend-name">
                    {{ friend.username }}
                    <span class="friend-remark" v-if="friend.remark">
                      ({{ friend.remark }})
                    </span>
                  </div>
                  <div class="friend-status offline">离线</div>
                </div>
                <div class="friend-actions">
                  <el-button type="primary" link @click="startChat(friend)">
                    发消息
                  </el-button>
                  <el-dropdown>
                    <el-button link>
                      <el-icon><More /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="editRemark(friend)">
                          修改备注
                        </el-dropdown-item>
                        <el-dropdown-item @click="deleteFriend(friend)" class="text-danger">
                          删除好友
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </el-collapse-item>
          </el-collapse>
        </el-scrollbar>
      </div>
    </div>

    <!-- 好友申请列表 -->
    <div class="friend-requests">
      <div class="requests-header">
        <h3>好友申请</h3>
        <el-badge :value="pendingRequests.length" :hidden="!pendingRequests.length">
          待处理
        </el-badge>
      </div>
      <div class="requests-content">
        <el-scrollbar>
          <div v-if="friendRequests.length === 0" class="empty-requests">
            暂无好友申请
          </div>
          <div
            v-for="request in friendRequests"
            :key="request.id"
            class="request-item"
          >
            <el-avatar :src="request.sender?.avatar" :size="40">
              {{ request.sender?.username[0] }}
            </el-avatar>
            <div class="request-info">
              <div class="request-name">{{ request.sender?.username }}</div>
              <div class="request-message">{{ request.message || '请求添加你为好友' }}</div>
            </div>
            <div class="request-actions" v-if="request.status === 0">
              <el-button type="primary" size="small" @click="handleRequest(request, 1)">
                接受
              </el-button>
              <el-button size="small" @click="handleRequest(request, 2)">
                拒绝
              </el-button>
            </div>
            <div class="request-status" v-else>
              {{ request.status === 1 ? '已接受' : '已拒绝' }}
            </div>
          </div>
        </el-scrollbar>
      </div>
    </div>

    <!-- 添加好友对话框 -->
    <el-dialog
      v-model="addFriendDialogVisible"
      title="添加好友"
      width="500px"
    >
      <el-form :model="searchForm">
        <el-form-item>
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入用户ID"
            type="number"
            clearable
            @keyup.enter="searchUsers"
          >
            <template #append>
              <el-button @click="searchUsers">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>

      <div class="search-results" v-if="searchResults.length">
        <div
          v-for="user in searchResults"
          :key="user.id"
          class="user-item"
        >
          <el-avatar :src="user.avatar" :size="40">
            {{ user.username[0] }}
          </el-avatar>
          <div class="user-info">
            <div class="user-name">{{ user.username }}</div>
            <div class="user-id">ID: {{ user.id }}</div>
          </div>
          <el-button
            type="primary"
            size="small"
            @click="sendFriendRequest(user)"
            :disabled="friends.value.some(f => f.id === user.id)"
          >
            {{ friends.value.some(f => f.id === user.id) ? '已是好友' : '添加好友' }}
          </el-button>
        </div>
      </div>
      <div v-else-if="searched" class="empty-results">
        未找到该用户
      </div>
    </el-dialog>

    <!-- 修改备注对话框 -->
    <el-dialog
      v-model="remarkDialogVisible"
      title="修改备注"
      width="400px"
    >
      <el-form :model="remarkForm">
        <el-form-item label="备注">
          <el-input v-model="remarkForm.remark" placeholder="请输入备注"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="remarkDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveRemark">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Search } from '@element-plus/icons-vue'
import { friendApi } from '@/api/social'
import { useUserStore } from '@/store/user'
import { userApi } from '@/api/user'

const userStore = useUserStore()
const router = useRouter()

// 组件状态
const friends = ref([])
const friendRequests = ref([])
const addFriendDialogVisible = ref(false)
const remarkDialogVisible = ref(false)
const searchKeyword = ref('')
const searchResults = ref([])
const searching = ref(false)
const activeGroups = ref(['online', 'offline'])

// 清理函数
const cleanup = () => {
  friends.value = []
  friendRequests.value = []
  searchResults.value = []
  addFriendDialogVisible.value = false
  remarkDialogVisible.value = false
}

onBeforeUnmount(() => {
  cleanup()
})

// 计算在线好友列表
const onlineFriends = computed(() => {
  return friends.value?.filter(f => f.online) || []
})

// 计算离线好友列表
const offlineFriends = computed(() => {
  return friends.value?.filter(f => !f.online) || []
})

// 计算待处理的好友申请
const pendingRequests = computed(() => {
  return friendRequests.value?.filter(r => r.status === 0) || []
})

// 获取好友列表
const getFriendsList = async () => {
  try {
    const res = await friendApi.getFriends({
      user_id: userStore.userInfo.id
    })
    friends.value = res.data || []
  } catch (error) {
    console.error('获取好友列表失败:', error)
    ElMessage.error('获取好友列表失败')
  }
}

// 获取好友申请列表
const getFriendRequests = async () => {
  try {
    const res = await friendApi.getFriendRequests({
      user_id: userStore.userInfo.id,
      type: 0, // 获取收到的申请
      page: 1,
      size: 50
    })
    if (res.code === '10000') {
      friendRequests.value = res.data || []
    } else {
      console.error('获取好友申请列表失败:', res.message)
      friendRequests.value = []
    }
  } catch (error) {
    console.error('获取好友申请列表失败:', error)
    ElMessage.error('获取好友申请列表失败')
    friendRequests.value = []
  }
}

// 判断是否是好友
const isFriend = (userId) => {
  return friends.value?.some(f => f.id === userId) || false
}

// 搜索用户
const searchUsers = async () => {
  if (!searchKeyword.value.trim()) return
  
  try {
    // 尝试将输入转换为用户ID
    const userId = parseInt(searchKeyword.value.trim())
    if (isNaN(userId)) {
      ElMessage.warning('请输入正确的用户ID')
      return
    }

    const res = await userApi.getUserInfo({
      user_id: userId
    })

    if (res.code === '10000') {
      searchResults.value = [res.User]
    } else {
      searchResults.value = []
    }
    searching.value = true
  } catch (error) {
    console.error('搜索用户失败:', error)
    ElMessage.error('搜索用户失败')
    searchResults.value = []
    searching.value = true
  }
}

// 发送好友申请
const sendFriendRequest = async (user) => {
  // 检查是否已经是好友
  if (friends.value.some(f => f.id === user.id)) {
    ElMessage.warning('该用户已经是你的好友')
    return
  }

  try {
    await friendApi.sendFriendRequest({
      user_id: userStore.userInfo.id,
      friend_id: user.id,
    })
    ElMessage.success('好友申请已发送')
    addFriendDialogVisible.value = false
    searchKeyword.value = ''
    searchResults.value = []
    searching.value = false
  } catch (error) {
    console.error('发送好友申请失败:', error)
    ElMessage.error('发送好友申请失败',error.response.data.message)
  }
}

// 处理好友申请
const handleRequest = async (request, action) => {
  try {
    const res = await friendApi.handleFriendRequest({
      request_id: request.id,
      user_id: userStore.userInfo.id,
      action: action
    })
    
    if (res.code === '10000') {
      ElMessage.success(action === 1 ? '已接受好友申请' : '已拒绝好友申请')
      getFriendRequests()
      if (action === 1) {
        getFriendsList()
      }
    } else {
      throw new Error(res.message)
    }
  } catch (error) {
    console.error('处理好友申请失败:', error)
    ElMessage.error('处理好友申请失败')
  }
}

// 开始聊天
const startChat = (friend) => {
  router.push(`/social/chat/${friend.id}`)
}

// 修改备注
const editRemark = (friend) => {
  currentFriend.value = friend
  remarkForm.value = {
    remark: friend.remark || '',
    friendId: friend.id
  }
  remarkDialogVisible.value = true
}

// 保存备注
const saveRemark = async () => {
  if (!currentFriend.value) return
  
  try {
    await friendApi.updateFriend({
      user_id: userStore.userInfo.id,
      friend_id: remarkForm.value.friendId,
      remark: remarkForm.value.remark
    })
    ElMessage.success('修改备注成功')
    remarkDialogVisible.value = false
    getFriendsList()
  } catch (error) {
    console.error('保存备注失败:', error)
    ElMessage.error('保存备注失败')
  }
}

// 删除好友
const deleteFriend = async (friend) => {
  try {
    await ElMessageBox.confirm('确定要删除该好友吗？', '提示', {
      type: 'warning'
    })
    
    await friendApi.deleteFriend({
      user_id: userStore.userInfo.id,
      friend_id: friend.id
    })
    ElMessage.success('删除好友成功')
    getFriendsList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除好友失败:', error)
      ElMessage.error('删除好友失败')
    }
  }
}

// 显示添加好友对话框
const showAddFriendDialog = () => {
  addFriendDialogVisible.value = true
  searchKeyword.value = ''
  searchResults.value = []
  searching.value = false
}

// 在script setup部分添加getAvatarText函数
const getAvatarText = (user) => {
  if (!user || !user.username) return 'U'
  return user.username.charAt(0).toUpperCase()
}

onMounted(() => {
  getFriendsList(
    {
      user_id: userStore.userInfo.id
    }
  )
  getFriendRequests(
    {
      user_id: userStore.userInfo.id,
      type : 0,
      page: 1,
      size: 50
    }
  )
})
</script>

<style scoped>
.friends-container {
  display: flex;
  height: calc(100vh - 60px);
  background: #fff;
}

.friends-list {
  flex: 2;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
}

.friends-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.friends-header h3 {
  margin: 0;
}

.friends-content {
  flex: 1;
  overflow: hidden;
}

.group-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.friend-count {
  font-size: 12px;
  color: #999;
}

.friend-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f5f7fa;
}

.friend-info {
  flex: 1;
  margin: 0 10px;
}

.friend-name {
  font-weight: 500;
  margin-bottom: 5px;
}

.friend-remark {
  font-size: 12px;
  color: #999;
}

.friend-status {
  font-size: 12px;
}

.friend-status.online {
  color: #67c23a;
}

.friend-status.offline {
  color: #909399;
}

.friend-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.friend-requests {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.requests-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.requests-header h3 {
  margin: 0;
}

.requests-content {
  flex: 1;
  overflow: hidden;
  padding: 20px;
}

.empty-requests {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

.request-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f5f7fa;
}

.request-info {
  flex: 1;
  margin: 0 10px;
}

.request-name {
  font-weight: 500;
  margin-bottom: 5px;
}

.request-message {
  font-size: 12px;
  color: #666;
}

.request-actions {
  display: flex;
  gap: 10px;
}

.request-status {
  font-size: 12px;
  color: #999;
}

.search-results {
  max-height: 400px;
  overflow-y: auto;
  margin-top: 20px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f5f7fa;
}

.user-info {
  flex: 1;
  margin: 0 10px;
}

.user-name {
  font-weight: 500;
  margin-bottom: 5px;
}

.user-id {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.empty-results {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

.text-danger {
  color: #f56c6c;
}
</style> 