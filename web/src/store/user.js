import { defineStore } from 'pinia'
import { userApi } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => {
    // 从 localStorage 恢复状态
    const savedUserInfo = localStorage.getItem('userInfo')
    console.log('【Store】初始化状态:', {
      savedUserInfo,
      token: localStorage.getItem('token'),
      userId: localStorage.getItem('user_id')
    })

    return {
      userInfo: savedUserInfo ? JSON.parse(savedUserInfo) : null,
      token: localStorage.getItem('token'),
      userId: localStorage.getItem('user_id')
    }
  },

  getters: {
    isLoggedIn: (state) => !!state.token,
    getAvatar: (state) => {
      const avatar = state.userInfo?.avatar
      return avatar === 'http://localhost:8080/avatars/default.jpg' ? '/src/assets/images/default.jpg' : avatar || '/src/assets/images/default.jpg'
    },
    getUsername: (state) => state.userInfo?.username || '未知用户'
  },

  actions: {
    async login(data) {
      try {
        console.log('【Store】开始登录:', data)
        const loginRes = await userApi.login(data)
        console.log('【Store】登录响应:', loginRes)
        
        // 保存登录信息
        this.token = loginRes.data.token
        this.userId = loginRes.data.user_id
        localStorage.setItem('token', loginRes.data.token)
        localStorage.setItem('user_id', loginRes.data.user_id.toString())
        
        // 获取用户信息
        console.log('【Store】获取用户信息...')
        const userInfoRes = await userApi.getUserInfo({ user_id: loginRes.data.user_id })
        console.log('【Store】用户信息响应:', userInfoRes)
        
        // 保存用户信息
        if (userInfoRes && userInfoRes.data && userInfoRes.data.id) {
          this.userInfo = userInfoRes.data
          localStorage.setItem('userInfo', JSON.stringify(userInfoRes.data))
        }
        
        return loginRes
      } catch (error) {
        console.error('【Store】登录失败:', error)
        // 清理可能部分设置的状态
        this.logout()
        throw error
      }
    },

    setUserInfo(info) {
      console.log('【Store】更新用户信息:', info)
      // 直接使用返回的用户信息
      if (info && info.id) {
        this.userInfo = info
        localStorage.setItem('userInfo', JSON.stringify(info))
      } else {
        console.error('【Store】更新用户信息格式错误:', info)
      }
    },

    logout() {
      console.log('【Store】执行登出')
      this.token = null
      this.userId = null
      this.userInfo = null
      localStorage.removeItem('token')
      localStorage.removeItem('user_id')
      localStorage.removeItem('userInfo')
    }
  }
}) 