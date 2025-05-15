<template>
  <div class="social-container">
    <el-menu
      :default-active="activeMenu"
      class="social-menu"
      mode="horizontal"
      router
      @select="handleSelect"
    >
      <el-menu-item index="/social/friends">
        <el-icon><UserFilled /></el-icon>
        <span>好友列表</span>
      </el-menu-item>
      <el-menu-item index="/social/chatroom/list">
        <el-icon><ChatRound /></el-icon>
        <span>聊天室</span>
      </el-menu-item>
    </el-menu>
    
    <div class="social-content">
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onErrorCaptured } from 'vue'
import { useRoute } from 'vue-router'
import { UserFilled, ChatRound } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)

const handleSelect = (key) => {
  console.log('selected menu:', key)
}

// 添加错误处理
onErrorCaptured((err, instance, info) => {
  console.error('组件错误:', err)
  console.error('错误信息:', info)
  return false // 阻止错误继续传播
})
</script>

<style scoped>
.social-container {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.social-menu {
  background: transparent;
}

.social-content {
  flex: 1;
  background: var(--el-bg-color);
  border-radius: 8px;
  padding: 20px;
  box-shadow: var(--el-box-shadow-light);
}
</style> 