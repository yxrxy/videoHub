<template>
  <div class="home">
    <Header />
    <div class="main-content">
      <aside class="sidebar">
        <nav>
          <router-link to="/video/list">视频列表</router-link>
          <router-link to="/video/upload">上传视频</router-link>
          <router-link to="/profile">个人中心</router-link>
          <router-link to="/social/friends">社交中心</router-link>
        </nav>
      </aside>
      <main class="content">
        <h1>欢迎来到VideoHub</h1>
        <div class="quick-actions">
          <div class="action-card" @click="$router.push('/video/upload')">
            <el-icon><Upload /></el-icon>
            <span>上传视频</span>
          </div>
          <div class="action-card" @click="$router.push('/video/list')">
            <el-icon><VideoCamera /></el-icon>
            <span>浏览热门视频</span>
          </div>
          <div class="action-card" @click="$router.push('/profile')">
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </div>
          <div class="action-card" @click="$router.push('/social/friends')">
            <el-icon><ChatDotRound /></el-icon>
            <span>社交中心</span>
          </div>
        </div>
        <div class="recent-videos" v-if="recentVideos.length">
          <h2>热门视频</h2>
          <div class="video-grid">
            <div v-for="video in recentVideos" :key="video.id" class="video-card" @click="$router.push(`/video/${video.id}`)">
              <img :src="video.cover_url" :alt="video.title">
              <div class="video-info">
                <h3>{{ video.title }}</h3>
                <p class="author">{{ video.author_name }}</p>
                <p class="stats">
                  <span>{{ formatNumber(video.visit_count) }}次观看</span>
                  <span>{{ formatDate(video.created_at) }}</span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import Header from '@/components/layout/Header.vue'
import { videoApi } from '@/api/video'
import { formatDate, formatNumber } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { Upload, VideoCamera, User, ChatDotRound } from '@element-plus/icons-vue'

export default {
  name: 'Home',
  components: {
    Header,
    Upload,
    VideoCamera,
    User,
    ChatDotRound
  },
  setup() {
    const recentVideos = ref([])

    const fetchRecentVideos = async () => {
      try {
        const res = await videoApi.getHotVideos({ 
          limit: 12,
          category: '',
          last_id: 0,
          last_like: 0,
          last_visit: 0
        })
        recentVideos.value = res.VideoList || []
      } catch (error) {
        console.error('获取热门视频失败:', error)
        ElMessage.error('获取视频列表失败')
      }
    }

    onMounted(() => {
      fetchRecentVideos()
    })

    return {
      recentVideos,
      formatDate,
      formatNumber
    }
  }
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.main-content {
  display: flex;
  padding: 20px;
  gap: 20px;
  max-width: 1440px;
  margin: 0 auto;
}

.sidebar {
  width: 200px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: fit-content;
}

.sidebar nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sidebar a {
  padding: 10px;
  color: #333;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.sidebar a:hover {
  background-color: #f0f0f0;
}

.content {
  flex: 1;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

h1 {
  margin-bottom: 30px;
  color: #333;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.action-card {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
  text-align: center;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-card i {
  font-size: 24px;
  margin-bottom: 10px;
  display: block;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.video-card {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.3s;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.video-card:hover {
  transform: translateY(-5px);
}

.video-card img {
  width: 100%;
  aspect-ratio: 16/9;
  object-fit: cover;
}

.video-info {
  padding: 12px;
}

.video-info h3 {
  margin: 0;
  font-size: 16px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.video-info .author {
  color: #666;
  margin: 8px 0;
  font-size: 14px;
}

.video-info .stats {
  display: flex;
  justify-content: space-between;
  color: #999;
  font-size: 12px;
  margin: 0;
}
</style> 