<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <div class="profile-header">
        <div class="profile-avatar">
          <el-avatar 
            :size="100" 
            :src="userStore.getAvatar || defaultAvatar"
            @error="handleAvatarError"
          ></el-avatar>
          <el-upload
            class="avatar-uploader"
            :action="`${baseURL}/api/v1/user/avatar`"
            :headers="uploadHeaders"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            name="avatar_data"
          >
            <el-button size="small">更换头像</el-button>
          </el-upload>
        </div>
        <div class="profile-info">
          <h2>{{ userStore.getUsername }}</h2>
          <p>{{ userStore.userInfo?.bio || '这个人很懒，什么都没写~' }}</p>
          <div class="profile-stats">
            <div class="stat-item">
              <div class="stat-value">{{ userStore.userInfo?.followCount || 0 }}</div>
              <div class="stat-label">关注</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ userStore.userInfo?.followerCount || 0 }}</div>
              <div class="stat-label">粉丝</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ userStore.userInfo?.videoCount || 0 }}</div>
              <div class="stat-label">视频</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ userStore.userInfo?.likeCount || 0 }}</div>
              <div class="stat-label">获赞</div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-tabs v-model="activeTab" class="profile-tabs">
      <el-tab-pane label="我的视频" name="videos">
        <div class="video-grid" v-if="videos && videos.length > 0">
          <div v-for="video in videos" :key="video.id" class="video-item" @click="$router.push(`/video/${video.id}`)">
            <el-card :body-style="{ padding: '0px' }">
              <img 
                :src="video.cover_url" 
                class="video-cover"
                @error="handleVideoCoverError"
              >
              <div class="video-info">
                <h3>{{ video.title }}</h3>
                <p>{{ formatNumber(video.visit_count || 0) }}次观看 · {{ formatDate(video.created_at) }}</p>
              </div>
            </el-card>
          </div>
        </div>
        <el-empty v-else description="暂无视频"></el-empty>
      </el-tab-pane>

      <el-tab-pane label="账号设置" name="settings">
        <el-form :model="userForm" label-width="100px">
          <el-form-item label="用户名">
            <el-input v-model="userForm.username" disabled></el-input>
          </el-form-item>
          <el-form-item label="个人简介">
            <el-input type="textarea" v-model="userForm.bio" disabled></el-input>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { videoApi } from '@/api/video'
import { userApi } from '@/api/user'
import { formatDate, formatNumber } from '@/utils/format'
import { ElMessage } from 'element-plus'

export default {
  name: 'UserProfile',
  setup() {
    const router = useRouter()
    const baseURL = import.meta.env.VITE_API_BASE_URL || ''
    const userStore = useUserStore()
    const activeTab = ref('videos')
    const videos = ref([])
    const defaultAvatar = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIiB2aWV3Qm94PSIwIDAgMjQgMjQiIGZpbGw9Im5vbmUiIHN0cm9rZT0iIzQwOWVmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTIwIDIxdi0yYTQgNCAwIDAgMC00LTRIOGE0IDQgMCAwIDAtNCA0djIiPjwvcGF0aD48Y2lyY2xlIGN4PSIxMiIgY3k9IjciIHI9IjQiPjwvY2lyY2xlPjwvc3ZnPg=='
    const userForm = ref({
      username: userStore.getUsername,
      bio: userStore.userInfo?.bio || ''
    })

    const uploadHeaders = computed(() => ({
      Authorization: `Bearer ${userStore.token}`
    }))

    const fetchUserVideos = async () => {
      try {
        const res = await videoApi.getVideoList({
          user_id: parseInt(userStore.userId),
          page: 1,
          size: 10
        })
        console.log('获取到的视频列表:', res)
        const videoList = res?.data || []
        videos.value = videoList.map(video => ({
          ...video,
          cover_url: video.coverUrl ? video.coverUrl.replace('/videos/', '/covers/') : '',
          play_url: video.playUrl ? `${baseURL}${video.playUrl}` : ''
        }))
        console.log('处理后的视频列表:', videos.value[0].cover_url)
      } catch (error) {
        console.error('获取用户视频失败:', error)
        ElMessage.error('获取视频列表失败')
      }
    }

    const handleAvatarSuccess = (res) => {
      if (res && res.data) {
        const avatarUrl = res.data.startsWith('http') ? res.data : `${baseURL}${res.data}`
        userStore.setUserInfo({
          ...userStore.userInfo,
          avatar: avatarUrl
        })
        ElMessage.success('头像更新成功')
      } else {
        console.error('头像上传响应格式错误:', res)
        ElMessage.error('头像更新失败')
      }
    }

    const handleAvatarError = (e) => {
      console.log('头像加载失败，使用默认头像')
      e.target.src = defaultAvatar
    }

    const handleVideoCoverError = (e) => {
      console.log('视频封面加载失败，使用默认封面')
      e.target.src = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIiB2aWV3Qm94PSIwIDAgMjQgMjQiIGZpbGw9Im5vbmUiIHN0cm9rZT0iIzQwOWVmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiPjxyZWN0IHg9IjIiIHk9IjIiIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgcng9IjIuMTgiIHJ5PSIyLjE4Ii8+PGxpbmUgeDE9IjciIHkxPSIyIiB4Mj0iNyIgeTI9IjIyIi8+PHBhdGggZD0iTTIgN2gyMCIvPjwvc3ZnPg=='
    }

    onMounted(async () => {
      console.log('Profile 组件挂载')
      console.log('当前登录状态:', { 
        isLoggedIn: userStore.isLoggedIn,
        userId: userStore.userId,
        userInfo: userStore.userInfo,
        baseURL
      })
      
      if (!userStore.isLoggedIn) {
        console.warn('用户未登录，跳转到登录页')
        router.push('/login')
        return
      }

      try {
        const res = await userApi.getUserInfo({ user_id: parseInt(userStore.userId) })
        console.log('获取到用户信息:', res)
        if (res && res.data) {
          const userData = {
            ...res.data,
            avatar: res.data.avatar ? (res.data.avatar.startsWith('http') ? res.data.avatar : `${baseURL}${res.data.avatar}`) : ''
          }
          userStore.setUserInfo(userData)
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        ElMessage.error('获取用户信息失败')
      }

      await fetchUserVideos()
    })

    return {
      baseURL,
      userStore,
      activeTab,
      videos,
      userForm,
      uploadHeaders,
      formatDate,
      formatNumber,
      handleAvatarSuccess,
      defaultAvatar,
      handleAvatarError,
      handleVideoCoverError
    }
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: flex-start;
  padding: 20px;
}

.profile-avatar {
  margin-right: 40px;
  text-align: center;
}

.avatar-uploader {
  margin-top: 10px;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  margin: 0 0 10px 0;
}

.profile-tabs {
  margin-top: 20px;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.video-item {
  cursor: pointer;
}

.video-cover {
  width: 100%;
  height: 140px;
  object-fit: cover;
}

.video-info {
  padding: 10px;
}

.video-info h3 {
  margin: 0;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.video-info p {
  margin: 5px 0 0;
  font-size: 12px;
  color: #666;
}

.profile-stats {
  display: flex;
  margin-top: 20px;
}

.stat-item {
  margin-right: 40px;
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
}

.stat-label {
  color: #666;
}
</style> 