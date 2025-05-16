<template>
  <div class="video-player-container">
    <div class="video-main">
      <div class="video-wrapper">
        <video
          ref="videoPlayer"
          class="native-video-player"
          controls
          preload="auto"
          width="100%"
          height="100%"
          :poster="videoData.cover_url"
          :src="videoData.play_url"
          type="video/mp4"
          controlsList="nodownload"
          @play="handlePlay"
          @pause="handlePause"
        >
          您的浏览器不支持 HTML5 视频播放，请升级浏览器或使用其他浏览器
        </video>
      </div>
      <div class="video-info">
        <h1>{{ videoData.title }}</h1>
        <div class="video-stats">
          <span class="user-id">用户ID: {{ getUserId(videoData) }}</span>
          <span class="like-count">点赞数: {{ formatNumber(videoData.like_count || 0) }}</span>
          <span>{{ formatNumber(videoData.visit_count || 0) }}次观看</span>
          <span>{{ formatDate(videoData.created_at) }}</span>
        </div>
        <div class="author-info" v-if="videoData.author">
          <img :src="videoData.author.avatar || defaultAvatar" :alt="videoData.author.username" class="author-avatar">
          <span class="author-name">{{ videoData.author.username }}</span>
          <el-button 
            type="primary" 
            size="small"
            :loading="followLoading"
            @click="handleFollow"
          >
            {{ videoData.author.is_following ? '已关注' : '关注' }}
          </el-button>
        </div>
        <div class="video-actions">
          <el-button 
            :type="videoData.is_liked ? 'primary' : 'default'" 
            @click="handleLike"
            :loading="likeLoading"
          >
            <el-icon><CaretTop /></el-icon>
            {{ formatNumber(videoData.like_count || 0) }}
          </el-button>
          <el-button type="default" @click="handleShare">
            <el-icon><Share /></el-icon>
            分享
          </el-button>
          <el-button type="default" @click="handleCollect">
            <el-icon><Star /></el-icon>
            收藏
          </el-button>
        </div>
        <div class="video-description">
          {{ videoData.description }}
        </div>
      </div>
    </div>
    <div class="video-comments">
      <h2>评论 ({{ videoData.comment_count || 0 }})</h2>
      <div class="comment-input">
        <el-input
          type="textarea"
          :rows="2"
          placeholder="发表评论..."
          v-model="commentContent"
          :maxlength="200"
          show-word-limit
        >
        </el-input>
        <el-button type="primary" @click="submitComment" :loading="commentLoading">
          发表评论
        </el-button>
      </div>
      <div class="comment-list">
        <div v-if="loading" class="loading-wrapper">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        <template v-else>
          <div v-if="comments.length === 0" class="no-comments">
            暂无评论，快来抢沙发吧！
          </div>
          <div v-else v-for="comment in comments" :key="comment.id" class="comment-item">
            <img :src="comment.user.avatar || defaultAvatar" :alt="comment.user.username" class="comment-avatar">
            <div class="comment-content">
              <div class="comment-user">{{ comment.user.username }}</div>
              <div class="comment-text">{{ comment.content }}</div>
              <div class="comment-actions">
                <span>{{ formatDate(comment.created_at) }}</span>
                <span class="reply-btn" @click="replyToComment(comment)">回复</span>
                <span>
                  <el-icon><CaretTop /></el-icon>
                  {{ formatNumber(comment.like_count || 0) }}
                </span>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { formatDate, formatNumber } from '@/utils/format'
import { Loading, CaretTop, Share, Star } from '@element-plus/icons-vue'
import { videoApi } from '@/api/video'
import { userApi } from '@/api/user'
import { useUserStore } from '@/store/user'
import defaultAvatar from '@/assets/images/default.jpg'
import defaultCover from '@/assets/images/default.jpg'

const route = useRoute()
const userStore = useUserStore()
const videoId = ref(parseInt(route.params.id))
const videoPlayer = ref(null)
const videoData = ref({})
const comments = ref([])
const commentContent = ref('')
const loading = ref(false)
const commentLoading = ref(false)
const likeLoading = ref(false)
const followLoading = ref(false)

// 获取用户信息
const fetchAuthorInfo = async (userId) => {
  try {
    const response = await userApi.getUserInfo({ user_id: userId })
    if (response?.code === '10000' && response?.data) {
      const authorData = response.data
      return {
        id: authorData.id,
        username: authorData.username,
        avatar: authorData.avatar || defaultAvatar,
        is_following: authorData.isFollow || false,
        followCount: authorData.followCount || 0,
        followerCount: authorData.followerCount || 0,
        videoCount: authorData.videoCount || 0,
        likeCount: authorData.likeCount || 0
      }
    }
    return null
  } catch (error) {
    console.error('获取作者信息失败:', error)
    return null
  }
}

// 获取视频详情
const fetchVideoDetail = async () => {
  loading.value = true
  try {
    const response = await videoApi.getVideoDetail(videoId.value)
    if (response?.code === '10000' && response?.data) {
      const data = response.data
      
      // 调试输出原始数据
      console.log('原始视频详情数据:', data)
      
      // 处理封面URL
      let coverUrl = data.coverUrl || defaultCover
      if (coverUrl && coverUrl.includes('/videos/')) {
        coverUrl = coverUrl.replace('/videos/', '/covers/')
      }
      
      // 处理视频URL
      let playUrl = data.playUrl || ''

      // 获取作者ID
      const authorId = data.authorId || (data.author && data.author.id)
      
      // 获取作者详细信息
      let authorInfo = null
      if (authorId) {
        authorInfo = await fetchAuthorInfo(authorId)
      }

      // 如果获取不到作者信息，使用视频中的基本信息
      const author = authorInfo || {
        id: authorId,
        username: data.author?.username || '未知用户',
        avatar: data.author?.avatarUrl || defaultAvatar,
        is_following: data.author?.isFollow || false
      }

      videoData.value = {
        id: data.id,
        title: data.title || '无标题',
        description: data.description || '',
        cover_url: coverUrl,
        play_url: playUrl,
        visit_count: parseInt(data.visitCount || 0),
        like_count: parseInt(data.favoriteCount || 0),
        comment_count: parseInt(data.commentCount || 0),
        created_at: data.createdAt || Date.now(),
        is_liked: data.isFavorite || false,
        authorId: authorId,
        author: author
      }

      console.log('处理后的视频数据:', videoData.value)

      // 增加访问计数
      await videoApi.incrementVisitCount(data.id)
      
      // 延迟加载视频
      setTimeout(() => {
        if (videoPlayer.value) {
          videoPlayer.value.load();
        }
      }, 100);
    } else {
      ElMessage.error(response?.message || '获取视频详情失败')
    }
  } catch (error) {
    console.error('获取视频详情失败:', error)
    ElMessage.error(error.message || '获取视频详情失败')
  }
  loading.value = false
}

// 获取评论列表
const fetchComments = async () => {
  try {
    // TODO: 实现评论列表获取
    comments.value = []
  } catch (error) {
    console.error('获取评论列表失败:', error)
  }
}

// 提交评论
const submitComment = async () => {
  if (!commentContent.value.trim()) {
    return
  }

  commentLoading.value = true
  try {
    // TODO: 实现评论提交
    commentContent.value = ''
    await fetchComments() // 重新获取评论列表
  } catch (error) {
    console.error('提交评论失败:', error)
    ElMessage.error('提交评论失败')
  }
  commentLoading.value = false
}

// 回复评论
const replyToComment = (comment) => {
  commentContent.value = `@${comment.user.username} `
}

// 点赞视频
const handleLike = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  likeLoading.value = true
  try {
    // TODO: 实现视频点赞
    videoData.value.is_liked = !videoData.value.is_liked
    videoData.value.like_count += videoData.value.is_liked ? 1 : -1
  } catch (error) {
    console.error('点赞失败:', error)
    ElMessage.error('点赞失败')
  }
  likeLoading.value = false
}

// 关注作者
const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  followLoading.value = true
  try {
    // TODO: 实现关注功能
    videoData.value.author.is_following = !videoData.value.author.is_following
  } catch (error) {
    console.error('关注失败:', error)
    ElMessage.error('关注失败')
  }
  followLoading.value = false
}

// 分享视频
const handleShare = () => {
  // TODO: 实现分享功能
  ElMessage.info('分享功能开发中')
}

// 收藏视频
const handleCollect = () => {
  // TODO: 实现收藏功能
  ElMessage.info('收藏功能开发中')
}

// 修改视频点击处理逻辑，避免冲突
const handleVideoClick = () => {
  // 移除点击视频区域自动播放/暂停的功能，避免与原生控件冲突
  // 让浏览器原生控件处理播放/暂停
}

const handlePlay = () => {
  console.log('视频开始播放')
  // 视频播放时确保封面被隐藏
  if (videoPlayer.value) {
    videoPlayer.value.removeAttribute('poster')
    videoPlayer.value.setAttribute('data-playing', 'true')
  }
}

const handlePause = () => {
  console.log('视频暂停')
  if (videoPlayer.value) {
    videoPlayer.value.setAttribute('data-playing', 'false')
  }
}

// 获取用户ID
const getUserId = (video) => {
  // 首先检查作者对象中的ID
  if (video.author && video.author.id != null) {
    return video.author.id;
  }
  // 如果author.id不存在，尝试直接获取authorId
  if (video.authorId != null) {
    return video.authorId;
  }
  // 如果以上都失败，检查是否有poster_id字段
  if (video.poster_id != null) {
    return video.poster_id;
  }
  // 最后，检查是否有user_id字段
  if (video.user_id != null) {
    return video.user_id;
  }
  // 如果所有尝试都失败，返回未知
  return '未知';
};

onMounted(async () => {
  await fetchVideoDetail()
  await fetchComments()
  
  await nextTick()
  
  // 确保视频元素加载完成后正确初始化
  if (videoPlayer.value) {
    // 监听视频元素的错误事件
    videoPlayer.value.addEventListener('error', (e) => {
      console.error('视频加载错误:', e);
    });
    
    // 在视频元素准备好时触发
    videoPlayer.value.addEventListener('loadedmetadata', () => {
      console.log('视频元数据已加载');
    });
    
    // 设置视频元素的crossOrigin属性，解决可能的CORS问题
    videoPlayer.value.crossOrigin = 'anonymous';
    
    // 初始化视频元素
    videoPlayer.value.load();
  }
})

onBeforeUnmount(() => {
  // 清理工作
  if (videoPlayer.value) {
    videoPlayer.value.pause()
  }
})
</script>

<style scoped>
.video-player-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.video-main {
  margin-bottom: 20px;
  position: relative;
  z-index: 1;
}

.video-wrapper {
  position: relative;
  padding-top: 56.25%; /* 16:9 比例 */
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  z-index: 1;
}

.video-wrapper video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  outline: none; /* 移除视频焦点时的轮廓 */
  object-fit: contain; /* 确保视频内容完整显示 */
}

.video-info {
  position: relative;
  z-index: 1;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  margin-top: 20px;
}

.video-info h1 {
  margin: 0;
  font-size: 24px;
  line-height: 1.4;
}

.video-stats {
  color: #666;
  margin: 10px 0;
  display: flex;
  gap: 20px;
}

.author-info {
  display: flex;
  align-items: center;
  margin: 20px 0;
  padding: 15px 0;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
}

.author-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-right: 15px;
}

.author-name {
  font-size: 16px;
  font-weight: 500;
  margin-right: auto;
}

.video-actions {
  display: flex;
  gap: 15px;
  margin: 20px 0;
}

.video-description {
  white-space: pre-wrap;
  color: #666;
  line-height: 1.6;
}

.video-comments {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.comment-input {
  margin: 20px 0;
}

.comment-input .el-button {
  margin-top: 10px;
}

.comment-list {
  margin-top: 30px;
}

.comment-item {
  display: flex;
  margin: 20px 0;
  padding: 15px 0;
  border-bottom: 1px solid #eee;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 15px;
}

.comment-content {
  flex: 1;
}

.comment-user {
  font-weight: 500;
  margin-bottom: 5px;
}

.comment-text {
  line-height: 1.6;
  color: #333;
}

.comment-actions {
  margin-top: 8px;
  color: #666;
  font-size: 13px;
  display: flex;
  gap: 15px;
}

.reply-btn {
  color: #409EFF;
  cursor: pointer;
}

.reply-btn:hover {
  text-decoration: underline;
}

.loading-wrapper {
  text-align: center;
  padding: 40px 0;
  color: #909399;
}

.loading-wrapper .el-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.no-comments {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

/* 确保视频控件可见 */
.native-video-player {
  z-index: 10;
  background-color: #000;
}

/* 移除之前的 webkit 媒体控件样式，让浏览器使用默认控件样式 */
.native-video-player::-webkit-media-controls-enclosure {
  display: flex !important;
}

.native-video-player::-webkit-media-controls {
  display: flex !important;
  opacity: 1 !important;
  z-index: 20 !important;
}

/* 确保播放时隐藏封面 */
.native-video-player[data-playing="true"] {
  object-fit: contain;
}

.video-main {
  margin-bottom: 20px;
  position: relative;
  z-index: 1;
}

.video-info {
  position: relative;
  z-index: 1;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  margin-top: 20px;
}

.user-id, .like-count {
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 4px;
  margin-right: 10px;
}

.user-id {
  background-color: #e1f3ff;
  color: #3273dc;
}

.like-count {
  background-color: #fff2e8;
  color: #ff7d1a;
}
</style> 