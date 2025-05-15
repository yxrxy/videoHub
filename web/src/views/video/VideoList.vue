<template>
  <div class="video-list-container">
    <!-- 搜索框 -->
    <div class="search-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索视频"
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button @click="handleSearch">
            <el-icon><Search /></el-icon>
          </el-button>
        </template>
      </el-input>
    </div>

    <!-- 热门排行榜标题 -->
    <div class="rank-header">
      <h2 class="rank-title">热门排行榜</h2>
      <div class="rank-desc">实时更新最受欢迎的视频</div>
    </div>

    <!-- 视频列表 -->
    <div class="video-grid">
      <el-card
        v-for="video in videos"
        :key="video.id"
        class="video-card"
        shadow="hover"
        @click="goToVideo(video.id)"
      >
        <div class="video-cover">
          <img :src="video.cover_url || defaultCover" :alt="video.title">
          <div class="video-duration" v-if="video.duration">{{ formatDuration(video.duration) }}</div>
        </div>
        <div class="video-info">
          <h3 class="video-title" :title="video.title">{{ video.title }}</h3>
          <div class="video-meta">
            <span class="author">{{ video.author?.username || '未知用户' }}</span>
            <div class="stats">
              <span>{{ formatNumber(video.visit_count || 0) }}播放</span>
              <span>{{ formatNumber(video.like_count || 0) }}点赞</span>
            </div>
            <span class="upload-time">{{ formatDate(video.created_at) }}</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 加载更多 -->
    <div class="load-more" v-if="hasMore">
      <el-button
        :loading="loading"
        @click="loadMore"
      >
        加载更多
      </el-button>
    </div>

    <!-- 回到顶部 -->
    <el-backtop></el-backtop>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { formatDate, formatDuration, formatNumber } from '@/utils/format'
import { videoApi } from '@/api/video'
import defaultCover from '@/assets/images/default.jpg'

export default {
  name: 'VideoList',
  components: {
    Search
  },
  setup() {
    const router = useRouter()
    const searchQuery = ref('')
    const videos = ref([])
    const loading = ref(false)
    const hasMore = ref(true)
    const page = ref(1)
    const pageSize = ref(12)

    const fetchVideos = async (isLoadMore = false) => {
      loading.value = true
      try {
        const response = await videoApi.getHotVideos()
        
        if (response?.code === '10000' && response?.data) {
          const newVideos = Array.isArray(response.data) ? response.data : [response.data]
          // 处理每个视频对象，确保所有必要的字段都有默认值
          const processedVideos = newVideos.map(video => {
            // 处理封面URL，将videos替换为covers
            let coverUrl = video.coverUrl || defaultCover
            if (coverUrl && coverUrl.includes('/videos/')) {
              coverUrl = coverUrl.replace('/videos/', '/covers/')
            }

            return {
              id: video.id,
              title: video.title || '无标题',
              cover_url: coverUrl,
              duration: video.duration || 0,
              author: {
                id: video.authorId,
                username: '作者ID : ' + video.authorId
              },
              visit_count: parseInt(video.commentCount || 0),
              like_count: parseInt(video.favoriteCount || 0),
              created_at: Date.now()
            }
          })

          if (isLoadMore) {
            videos.value = [...videos.value, ...processedVideos]
          } else {
            videos.value = processedVideos
          }
          
          // 更新是否有更多数据的状态
          hasMore.value = processedVideos.length > 0
        } else {
          console.error('获取视频列表失败：', response)
          ElMessage.error(response?.message || '获取视频列表失败：数据格式错误')
        }
      } catch (error) {
        console.error('获取热门视频列表失败:', error)
        ElMessage.error(error.message || '获取视频列表失败')
      } finally {
        loading.value = false
      }
    }

    const handleSearch = () => {
      if (!searchQuery.value.trim()) return
      router.push({
        path: '/search',
        query: { q: searchQuery.value }
      })
    }

    const loadMore = () => {
      page.value++
      fetchVideos(true)
    }

    const goToVideo = (videoId) => {
      try {
        if (!videoId || isNaN(parseInt(videoId))) {
          ElMessage.error('无效的视频ID')
          return
        }
        router.push(`/video/${videoId}`)
      } catch (error) {
        console.error('跳转到视频详情页失败:', error)
        ElMessage.error('跳转失败')
      }
    }

    onMounted(() => {
      fetchVideos()
    })

    return {
      searchQuery,
      videos,
      loading,
      hasMore,
      handleSearch,
      loadMore,
      goToVideo,
      formatDate,
      formatDuration,
      formatNumber
    }
  }
}
</script>

<style scoped>
.video-list-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-header {
  margin-bottom: 20px;
}

.search-input {
  max-width: 600px;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.video-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.video-card:hover {
  transform: translateY(-5px);
}

.video-cover {
  position: relative;
  width: 100%;
  padding-top: 56.25%; /* 16:9 比例 */
}

.video-cover img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-duration {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 2px 4px;
  border-radius: 2px;
  font-size: 12px;
}

.video-info {
  padding: 12px;
}

.video-title {
  margin: 0;
  font-size: 14px;
  line-height: 1.4;
  height: 40px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.video-meta {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
}

.author {
  display: block;
  margin-bottom: 4px;
}

.stats {
  display: flex;
  gap: 12px;
  margin-bottom: 4px;
}

.upload-time {
  color: #999;
}

.load-more {
  text-align: center;
  margin-top: 40px;
}

.rank-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.rank-title {
  font-size: 24px;
  color: #303133;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.rank-desc {
  font-size: 14px;
  color: #909399;
}
</style> 