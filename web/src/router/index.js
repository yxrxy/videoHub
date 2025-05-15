import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/Register.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/video/list',
    name: 'VideoList',
    component: () => import('@/views/video/VideoList.vue')
  },
  {
    path: '/video/upload',
    name: 'VideoUpload',
    component: () => import('@/views/video/Upload.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/video/:id',
    name: 'VideoPlayer',
    component: () => import('@/views/video/VideoPlayer.vue')
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/search/SearchResult.vue')
  },
  {
    path: '/social',
    name: 'Social',
    component: () => import('@/views/social/Social.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'friends',
        name: 'Friends',
        component: () => import('@/views/social/Friends.vue')
      },
      {
        path: 'chat/:userId',
        name: 'PrivateChat',
        component: () => import('@/views/social/PrivateChat.vue')
      },
      {
        path: 'chatroom/list',
        name: 'ChatRoomList',
        component: () => import('@/views/social/ChatRoomList.vue')
      },
      {
        path: 'chatroom/:roomId',
        name: 'ChatRoom', 
        component: () => import('@/views/social/ChatRoom.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 检查路由是否需要认证
  if (to.meta.requiresAuth) {
    // 检查用户是否已登录
    const isAuthenticated = localStorage.getItem('token')
    if (!isAuthenticated) {
      // 如果没有登录，重定向到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router 