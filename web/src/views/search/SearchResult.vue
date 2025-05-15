<template>
  <div class="search-result-container">
    <!-- 搜索框 -->
    <div class="search-header">
      <!-- 搜索模式切换 -->
      <div class="search-mode-switch">
        <el-radio-group v-model="searchMode" size="large">
          <el-radio-button value="exact">精确匹配</el-radio-button>
          <el-radio-button value="semantic">词意搜索</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 词意搜索 -->
      <div v-if="searchMode === 'semantic'" class="semantic-search">
        <el-input
          v-model="semanticQuery"
          placeholder="输入关键词进行词意搜索"
          class="search-input"
          clearable
          @keyup.enter="handleSemanticSearch"
        >
          <template #append>
            <el-button @click="handleSemanticSearch">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>

      <!-- 精确匹配搜索 -->
      <div v-else class="exact-search">
        <el-form :model="exactSearchForm" class="search-form">
          <el-form-item>
            <el-input
              v-model="exactSearchForm.keywords"
              placeholder="输入关键词进行精确搜索"
              class="search-input"
              clearable
              @keyup.enter="handleExactSearch"
            >
              <template #append>
                <el-button @click="handleExactSearch">
                  <el-icon><Search /></el-icon>
                </el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="用户名">
            <el-input
              v-model="exactSearchForm.username"
              placeholder="按用户名筛选"
              clearable
            />
          </el-form-item>

          <el-form-item label="上传时间">
            <el-date-picker
              v-model="exactSearchForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="x"
              :default-time="[
                new Date(2000, 1, 1, 0, 0, 0),
                new Date(2000, 1, 1, 23, 59, 59),
              ]"
            />
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div class="search-content" v-if="hasResults">
      <div class="result-stats">
        找到 {{ totalResults }} 个相关视频
      </div>

      <!-- 语义搜索额外信息 -->
      <div v-if="searchMode === 'semantic' && (searchSummary || relatedQueries.length > 0)" class="semantic-info">
        <!-- 搜索结果摘要 -->
        <div v-if="searchSummary" class="search-summary">
          <h3>搜索结果摘要</h3>
          <p>{{ searchSummary }}</p>
        </div>
        
        <!-- 相关搜索建议 -->
        <div v-if="relatedQueries.length > 0" class="related-queries">
          <h3>相关搜索</h3>
          <div class="query-tags">
            <el-tag
              v-for="query in relatedQueries"
              :key="query"
              class="query-tag"
              @click="handleRelatedQuery(query)"
            >
              {{ query }}
            </el-tag>
          </div>
        </div>
      </div>

      <div class="result-list">
        <el-card
          v-for="video in searchResults"
          :key="video.id"
          class="result-item"
          shadow="hover"
          @click="goToVideo(video.id)"
        >
          <div class="result-layout">
            <div class="video-thumbnail">
              <img :src="formatVideoCover(video.coverUrl)" :alt="video.title">
            </div>
            <div class="video-details">
              <h3 class="video-title">{{ video.title }}</h3>
              <div class="video-meta">
                <span class="like-count">点赞数: {{ formatNumber(video.favoriteCount || 0) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="totalPages > 1">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="totalResults"
          @current-change="handlePageChange"
          layout="prev, pager, next"
        >
        </el-pagination>
      </div>
    </div>

    <!-- 无结果提示 -->
    <div v-else-if="searched" class="no-results">
      <el-empty description="未找到相关视频">
        <template #description>
          <p>未找到与"{{ searchQuery }}"相关的视频</p>
          <p class="suggestions">建议：</p>
          <ul>
            <li>请检查输入是否正确</li>
            <li>尝试使用其他关键词</li>
            <li>使用更简单的关键词</li>
          </ul>
        </template>
      </el-empty>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatDate, formatDuration, formatNumber } from '@/utils/format'
import { videoApi } from '@/api/video'
import defaultCover from '@/assets/images/default.jpg'

export default {
  name: 'SearchResult',
  components: {
    Search
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const searchMode = ref('exact')
    const semanticQuery = ref('')
    const exactSearchForm = ref({
      keywords: '',
      username: '',
      dateRange: null
    })
    const currentPage = ref(1)
    const pageSize = ref(20)
    const totalResults = ref(0)
    const searchResults = ref([])
    const searched = ref(false)
    const loading = ref(false)
    const searchSummary = ref('')
    const relatedQueries = ref([])

    const hasResults = computed(() => searchResults.value.length > 0)
    const totalPages = computed(() => Math.ceil(totalResults.value / pageSize.value))

    const performSearch = async () => {
      if (searchMode.value === 'exact' && !exactSearchForm.value.keywords.trim()) return
      if (searchMode.value === 'semantic' && !semanticQuery.value.trim()) return

      searched.value = true
      loading.value = true
      try {
        let response
        if (searchMode.value === 'exact') {
          const [fromDate, toDate] = exactSearchForm.value.dateRange || [null, null]
          response = await videoApi.searchVideos({
            keywords: exactSearchForm.value.keywords,
            username: exactSearchForm.value.username,
            from_date: fromDate,
            to_date: toDate,
            page_num: currentPage.value,
            page_size: pageSize.value
          })
        } else {
          response = await videoApi.semanticSearchVideos({
            query: semanticQuery.value,
            page_num: currentPage.value,
            page_size: 5,
            threshold: 0.3
          })
        }

        if (response.code === '10000') {
          if (searchMode.value === 'semantic') {
            // 语义搜索返回的是results字段，包含多个语义结果项
            const results = response.data || [];
            console.log('语义搜索结果:', results);
            
            // 提取所有视频
            let allVideos = [];
            let summary = '';
            let queries = [];
            
            // 处理返回的语义搜索结果
            if (results && results.length > 0) {
              // 提取所有结果项中的视频
              results.forEach(item => {
                if (item.videos && item.videos.length > 0) {
                  // 确保每个视频对象都包含必要的字段
                  const processedVideos = item.videos.map(video => {
                    // 记录原始数据以便调试
                    console.log('原始视频数据:', video);
                    // 确保有authorId，如果没有则检查author.id
                    if (!video.authorId && video.author && video.author.id) {
                      video.authorId = video.author.id;
                    }
                    return video;
                  });
                  allVideos = [...allVideos, ...processedVideos];
                }
                // 使用第一个结果项的摘要和相关查询
                if (!summary && item.summary) {
                  summary = item.summary;
                }
                if (!queries.length && item.related_queries && item.related_queries.length) {
                  queries = item.related_queries;
                }
              });
            }
            
            console.log('处理后的视频列表:', allVideos);
            searchResults.value = allVideos;
            searchSummary.value = summary;
            relatedQueries.value = queries;
            totalResults.value = allVideos.length;
          } else {
            // 普通搜索返回的是videos字段
            const videos = response.data || [];
            console.log('普通搜索结果:', videos);
            
            // 确保每个视频对象都包含必要的字段
            searchResults.value = videos.map(video => {
              console.log('原始视频数据:', video);
              // 确保有authorId，如果没有则检查author.id
              if (!video.authorId && video.author && video.author.id) {
                video.authorId = video.author.id;
              }
              return video;
            });
            
            // 清空语义搜索相关数据
            searchSummary.value = '';
            relatedQueries.value = [];
            totalResults.value = response.total || 0;
          }
        } else {
          ElMessage.error(response.message || '搜索失败');
        }
      } catch (error) {
        console.error('搜索失败:', error);
        ElMessage.error('搜索失败');
      }
      loading.value = false;
    }

    const handleExactSearch = () => {
      currentPage.value = 1
      router.push({
        query: {
          ...getSearchQuery(),
          page: 1
        }
      })
      performSearch()
    }

    const handleSemanticSearch = () => {
      currentPage.value = 1
      router.push({
        query: {
          q: semanticQuery.value,
          mode: 'semantic',
          page: 1
        }
      })
      performSearch()
    }

    const handleRelatedQuery = (query) => {
      semanticQuery.value = query
      handleSemanticSearch()
    }

    const getSearchQuery = () => {
      const query = {
        mode: searchMode.value
      }

      if (searchMode.value === 'exact') {
        const [fromDate, toDate] = exactSearchForm.value.dateRange || [null, null]
        return {
          ...query,
          keywords: exactSearchForm.value.keywords,
          username: exactSearchForm.value.username,
          from_date: fromDate,
          to_date: toDate
        }
      } else {
        return {
          ...query,
          q: semanticQuery.value
        }
      }
    }

    const handlePageChange = (page) => {
      currentPage.value = page
      router.push({
        query: {
          ...getSearchQuery(),
          page
        }
      })
      performSearch()
    }

    const goToVideo = (videoId) => {
      router.push(`/video/${videoId}`)
    }

    // 获取作者名称
    const getAuthorName = (video) => {
      if (video.author && video.author.username) {
        return video.author.username;
      }
      return '未知用户';
    };

    // 格式化视频封面URL
    const formatVideoCover = (coverUrl) => {
      if (!coverUrl) return defaultCover;
      
      // 如果包含视频路径，替换为封面路径
      if (coverUrl.includes('/videos/')) {
        coverUrl = coverUrl.replace('/videos/', '/covers/');
      }
      
      // 确保URL格式正确
      if (!coverUrl.startsWith('http') && !coverUrl.startsWith('/')) {
        coverUrl = `/covers/${coverUrl}`;
      }
      
      return coverUrl;
    };

    // 获取用户ID
    const getUserId = (video) => {
      // 首先尝试直接获取authorId
      if (video.authorId != null && video.authorId !== undefined) {
        return video.authorId;
      }
      // 如果authorId不存在，尝试从author对象中获取id
      if (video.author && video.author.id != null) {
        return video.author.id;
      }
      // 如果以上都失败，检查是否有poster_id字段
      if (video.poster_id != null) {
        return video.poster_id;
      }
      // 最后，检查是否有 user_id 字段
      if (video.user_id != null) {
        return video.user_id;
      }
      // 如果所有尝试都失败，返回未知
      return '未知';
    };

    onMounted(() => {
      // 从URL参数中恢复搜索状态
      const { mode, keywords, username, from_date, to_date, page } = route.query
      if (mode) {
        searchMode.value = mode
        if (mode === 'exact') {
          exactSearchForm.value.keywords = keywords || ''
          exactSearchForm.value.username = username || ''
          exactSearchForm.value.dateRange = from_date && to_date ? [from_date, to_date] : null
        } else {
          semanticQuery.value = route.query.q || ''
        }
        currentPage.value = parseInt(page) || 1
        performSearch()
      }
    })

    return {
      searchMode,
      semanticQuery,
      exactSearchForm,
      currentPage,
      pageSize,
      totalResults,
      searchResults,
      searched,
      hasResults,
      totalPages,
      loading,
      searchSummary,
      relatedQueries,
      handleExactSearch,
      handleSemanticSearch,
      handlePageChange,
      goToVideo,
      formatDate,
      formatDuration,
      formatNumber,
      handleRelatedQuery,
      formatVideoCover,
      getAuthorName,
      getUserId
    }
  }
}
</script>

<style scoped>
.search-result-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-header {
  margin-bottom: 20px;
  width: 100%;
}

.search-mode-switch {
  margin-bottom: 20px;
  text-align: center;
}

.search-input {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.semantic-search {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
  width: 100%;
}

.exact-search {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
  width: 100%;
}

.search-filters {
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.filter-group {
  margin-bottom: 15px;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-label {
  margin-right: 10px;
  color: #606266;
}

.result-stats {
  margin-bottom: 20px;
  color: #606266;
}

.result-item {
  margin-bottom: 20px;
  cursor: pointer;
}

.result-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.result-layout {
  display: flex;
  gap: 20px;
}

.video-thumbnail {
  position: relative;
  width: 320px;
  height: 180px;
  flex-shrink: 0;
}

.video-thumbnail img {
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

.video-details {
  flex: 1;
}

.video-title {
  margin: 0 0 15px 0;
  font-size: 18px;
  line-height: 1.4;
  font-weight: 500;
  color: #303133;
}

.video-description {
  margin: 0 0 10px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.video-meta {
  color: #909399;
  font-size: 13px;
}

.video-meta > span {
  margin-right: 15px;
}

.video-meta .like-count {
  font-weight: bold;
  padding: 4px 10px;
  border-radius: 4px;
  margin-right: 10px;
  background-color: #fff2e8;
  color: #ff7d1a;
  display: inline-block;
}

.pagination-container {
  margin-top: 30px;
  text-align: center;
}

.no-results {
  margin-top: 60px;
  text-align: center;
}

.suggestions {
  margin: 20px 0 10px;
  color: #606266;
}

.suggestions ul {
  list-style: none;
  padding: 0;
  margin: 0;
  color: #909399;
}

.suggestions li {
  margin: 5px 0;
}

.search-form {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.search-form .el-form-item {
  margin-bottom: 20px;
  width: 100%;
}

.search-form .el-form-item:last-child {
  margin-bottom: 0;
}

.el-form-item__content {
  width: 100%;
}

.semantic-info {
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.search-summary {
  margin-bottom: 15px;
}

.search-summary h3 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #303133;
}

.search-summary p {
  color: #606266;
  line-height: 1.6;
}

.related-queries h3 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #303133;
}

.query-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.query-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.query-tag:hover {
  background-color: #409eff;
  color: white;
}
</style> 