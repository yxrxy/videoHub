import axios from 'axios'
import router from '@/router'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 15000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    console.log('发送请求:', config.url, config)
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    console.log('收到响应:', response.config.url, res)
    
    // 检查响应状态
    if (res.Base && res.Base.code !== 0) {
      ElMessage.error(res.Base.msg || '请求失败')
      return Promise.reject(new Error(res.Base.msg || '请求失败'))
    }
    
    return res
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      console.error('响应错误:', {
        url: error.config.url,
        status,
        data
      })
      
      if (status === 401) {
        // token 过期或无效
        localStorage.removeItem('token')
        localStorage.removeItem('user_id')
        router.push('/login')
        ElMessage.error('登录已过期，请重新登录')
      } else {
        ElMessage.error(data.Base?.msg || '请求失败')
      }
      return Promise.reject(data)
    }
    console.error('网络错误:', error)
    ElMessage.error('网络错误，请检查网络连接')
    return Promise.reject(error)
  }
)

export default request