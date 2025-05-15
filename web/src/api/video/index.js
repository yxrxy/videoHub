import request from '@/utils/request'

export const videoApi = {
  // 获取视频列表
  getVideoList(params) {
    return request({
      url: '/api/v1/video/list',
      method: 'get',
      params: {
        user_id: parseInt(params.user_id),
        page: params.page || 1,
        size: params.size || 10,
        category: params.category || ''
      }
    })
  },

  // 获取视频详情
  getVideoDetail(id) {
    return request({
      url: `/api/v1/video/${id}`,
      method: 'get'
    })
  },

  // 上传视频
  uploadVideo(data) {
    return request({
      url: '/api/v1/video/publish',
      method: 'post',
      data: {
        user_id: parseInt(localStorage.getItem('user_id')),
        video_data: data.file,
        content_type: data.file.type,
        title: data.title,
        description: data.description,
        category: data.category,
        tags: data.tags,
        is_private: data.isPrivate
      },
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取热门视频
  getHotVideos(params = {}) {
    return request({
      url: '/api/v1/video/hot',
      method: 'get',
      params: {
        limit: params.limit || 10,
        category: params.category,
        last_visit: params.lastVisit,
        last_like: params.lastLike,
        last_id: params.lastId
      }
    })
  },

  // 增加访问量
  incrementVisitCount(videoId) {
    return request({
      url: `/api/v1/video/${videoId}/visit`,
      method: 'post',
      data: {
        video_id: videoId
      }
    })
  },

  // 点赞视频
  likeVideo(videoId) {
    return request({
      url: `/api/v1/video/${videoId}/like`,
      method: 'post',
      data: {
        video_id: videoId
      }
    })
  },

  // 搜索视频
  searchVideos(params) {
    return request({
      url: '/api/v1/video/search',
      method: 'post',
      data: {
        keywords: params.keywords,
        page_size: params.pageSize || 10,
        page_num: params.pageNum || 1,
        from_date: params.fromDate,
        to_date: params.toDate,
        username: params.username
      }
    })
  },

  // 语义搜索视频
  semanticSearchVideos(params) {
    return request({
      url: '/api/v1/video/semantic',
      method: 'post',
      data: {
        query: params.query,
        page_size: params.page_size || 10,
        page_num: params.page_num || 1,
        threshold: params.threshold || 0.3
      }
    })
  }
} 