<template>
  <header class="header">
    <div class="logo">
      <router-link to="/">VideoHub</router-link>
    </div>
    <nav class="nav">
      <router-link to="/video/list">视频列表</router-link>
      <router-link to="/video/upload">上传视频</router-link>
      <router-link to="/search">搜索</router-link>
      <router-link to="/social/chat">社交</router-link>
    </nav>
    <div class="search-box">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索视频"
        @keyup.enter="handleSearch"
      >
        <template #suffix>
          <el-icon class="el-input__icon" @click="handleSearch">
            <Search />
          </el-icon>
        </template>
      </el-input>
    </div>
    <div class="user-info" v-if="userStore.isLoggedIn">
      <div class="avatar" @click="showUserMenu = !showUserMenu">
        <img :src="userStore.getAvatar" :alt="userStore.getUsername">
        <div class="user-menu" v-if="showUserMenu">
          <div class="menu-item" @click="$router.push('/profile')">个人资料</div>
          <div class="menu-item" @click="$router.push('/video/list?type=my')">我的视频</div>
          <div class="menu-item" @click="handleLogout">退出登录</div>
        </div>
      </div>
    </div>
    <div class="auth-buttons" v-else>
      <router-link to="/login" class="btn">登录</router-link>
      <router-link to="/register" class="btn btn-primary">注册</router-link>
    </div>
  </header>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { userApi } from '@/api/user'
import { Search } from '@element-plus/icons-vue'

export default {
  name: 'Header',
  components: {
    Search
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const showUserMenu = ref(false)
    const searchKeyword = ref('')

    const handleLogout = () => {
      userStore.logout()
      router.push('/login')
    }

    const handleSearch = () => {
      if (searchKeyword.value.trim()) {
        router.push({
          path: '/search',
          query: { q: searchKeyword.value.trim() }
        })
      }
    }

    onMounted(() => {
      console.log('【Header】组件挂载')
      console.log('【Header】当前登录状态:', {
        isLoggedIn: userStore.isLoggedIn,
        userId: userStore.userId,
        token: localStorage.getItem('token'),
        userInfo: userStore.userInfo
      })

      if (userStore.isLoggedIn && !userStore.userInfo) {
        console.log('【Header】开始获取用户信息...')
        const userId = parseInt(userStore.userId)
        if (isNaN(userId)) {
          console.error('【Header】无效的用户ID:', userStore.userId)
          userStore.logout()
          router.push('/login')
          return
        }

        userApi.getUserInfo({ user_id: userId })
          .then(res => {
            console.log('【Header】获取用户信息成功:', res)
            if (res && res.id) {
              userStore.setUserInfo(res)
            } else {
              console.error('【Header】用户信息响应格式错误:', res)
              throw new Error('获取用户信息失败')
            }
          })
          .catch(error => {
            console.error('【Header】获取用户信息失败:', error)
            if (error.response?.status === 401) {
              console.warn('【Header】Token失效，执行登出')
              userStore.logout()
              router.push('/login')
            }
          })
      }
    })

    return {
      userStore,
      showUserMenu,
      searchKeyword,
      handleLogout,
      handleSearch
    }
  }
}
</script>

<style scoped>
.header {
  background: white;
  padding: 15px 30px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logo a {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  text-decoration: none;
}

.nav {
  display: flex;
  gap: 20px;
}

.nav a {
  color: #666;
  text-decoration: none;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav a:hover {
  background-color: #f5f5f5;
}

.nav a.router-link-active {
  color: #409eff;
  font-weight: bold;
}

.search-box {
  flex: 1;
  max-width: 500px;
}

.user-info {
  position: relative;
  margin-left: auto;
}

.avatar {
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
}

.avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.user-menu {
  position: absolute;
  top: calc(100% + 5px);
  right: 0;
  background: white;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 5px 0;
  min-width: 150px;
  z-index: 1000;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.menu-item:hover {
  background-color: #f5f5f5;
}

.auth-buttons {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border-radius: 4px;
  text-decoration: none;
  transition: all 0.3s;
}

.btn:not(.btn-primary) {
  color: #666;
}

.btn-primary {
  background-color: #409eff;
  color: white;
}

.btn:hover {
  opacity: 0.8;
}
</style> 