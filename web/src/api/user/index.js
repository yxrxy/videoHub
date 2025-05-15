import request from '@/utils/request'

export const userApi = {
  // 用户登录
  login(data) {
    return request({
      url: '/api/v1/user/login',
      method: 'post',
      data: {
        username: data.username,
        password: data.password
      }
    })
  },

  // 用户注册
  register(data) {
    return request({
      url: '/api/v1/user/register',
      method: 'post',
      data: {
        username: data.username,
        password: data.password
      }
    })
  },

  // 获取用户信息
  getUserInfo(data) {
    return request({
      url: '/api/v1/user/info',
      method: 'get',
      params: {
        user_id: parseInt(data.user_id)
      }
    })
  },

  // 上传头像
  uploadAvatar(data) {
    const formData = new FormData()
    formData.append('avatar_data', data.file)
    formData.append('content_type', data.file.type)

    return request({
      url: '/api/v1/user/avatar',
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
} 