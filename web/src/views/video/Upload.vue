<template>
  <div class="upload-container">
    <el-card class="upload-card">
      <template #header>
        <div class="card-header">
          <h2>上传视频</h2>
        </div>
      </template>
      
      <el-steps :active="currentStep" finish-status="success" align-center>
        <el-step title="上传视频文件"></el-step>
        <el-step title="填写信息"></el-step>
      </el-steps>

      <!-- 步骤1：上传视频文件 -->
      <div v-if="currentStep === 0" class="upload-step">
        <div class="video-upload-container">
          <el-upload
            class="video-uploader"
            drag
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleVideoChange"
            accept="video/mp4,video/quicktime"
          >
            <div v-if="!videoUrl">
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                将视频拖到此处，或<em>点击上传</em>
              </div>
              <div class="upload-tip">
                支持mp4、mov格式，最大500MB
              </div>
            </div>
            <div v-else class="upload-preview">
              <video :src="videoUrl" controls></video>
            </div>
          </el-upload>
        </div>
        <div class="step-actions">
          <el-button type="primary" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>

      <!-- 步骤2：填写视频信息 -->
      <div v-if="currentStep === 1" class="upload-step">
        <el-form :model="videoForm" label-width="100px">
          <el-form-item label="视频标题" required>
            <el-input v-model="videoForm.title" placeholder="请输入视频标题"></el-input>
          </el-form-item>

          <el-form-item label="视频描述">
            <el-input
              type="textarea"
              v-model="videoForm.description"
              :rows="4"
              placeholder="请输入视频描述"
            ></el-input>
          </el-form-item>

          <el-form-item label="视频分类">
            <el-select v-model="videoForm.category" placeholder="请选择视频分类">
              <el-option
                v-for="item in categories"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="标签">
            <el-tag
              v-for="tag in videoForm.tags"
              :key="tag"
              closable
              @close="handleTagClose(tag)"
              class="tag-item"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-if="tagInputVisible"
              ref="tagInput"
              v-model="tagInputValue"
              class="tag-input"
              size="small"
              @keyup.enter="handleTagConfirm"
              @blur="handleTagConfirm"
            ></el-input>
            <el-button v-else size="small" @click="showTagInput">
              + 添加标签
            </el-button>
          </el-form-item>
        </el-form>
        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="publishVideo" :loading="publishing">
            发布视频
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted, nextTick, computed } from 'vue'
import { Plus, UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { videoApi } from '@/api/video'
import { useUserStore } from '@/store/user'

export default {
  name: 'VideoUpload',
  components: {
    Plus,
    UploadFilled
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const currentStep = ref(0)
    const videoUrl = ref('')
    const videoFile = ref(null)
    const uploadProgress = ref(0)
    const publishing = ref(false)
    const tagInputVisible = ref(false)
    const tagInputValue = ref('')
    const tagInput = ref(null)

    const videoForm = ref({
      title: '',
      description: '',
      category: '',
      tags: []
    })

    // 获取上传用的请求头
    const uploadHeaders = computed(() => ({
      Authorization: `Bearer ${userStore.token}`
    }))

    const categories = ref([
      { label: '游戏', value: 'game' },
      { label: '音乐', value: 'music' },
      { label: '电影', value: 'movie' },
      { label: '动画', value: 'animation' },
      { label: '科技', value: 'tech' },
      { label: '生活', value: 'life' }
    ])

    const handleVideoChange = (file) => {
      console.log('选择视频文件:', {
        fileName: file.name,
        fileType: file.type,
        fileSize: `${(file.size / 1024 / 1024).toFixed(2)}MB`
      })

      const isMP4 = file.raw.type === 'video/mp4'
      const isMOV = file.raw.type === 'video/quicktime'
      const isLt500M = file.raw.size / 1024 / 1024 < 500

      if (!isMP4 && !isMOV) {
        console.warn('视频格式不符合要求:', file.raw.type)
        ElMessage.error('视频格式必须是 mp4 或 mov!')
        return false
      }
      if (!isLt500M) {
        console.warn('视频文件过大:', `${(file.raw.size / 1024 / 1024).toFixed(2)}MB`)
        ElMessage.error('视频大小不能超过 500MB!')
        return false
      }

      videoFile.value = file.raw
      videoUrl.value = URL.createObjectURL(file.raw)
      return true
    }

    const handleCoverSuccess = (res) => {
      console.log('封面上传成功响应:', res)
      videoForm.value.cover = res.url
    }

    const beforeCoverUpload = (file) => {
      console.log('准备上传封面:', {
        fileName: file.name,
        fileType: file.type,
        fileSize: `${(file.size / 1024 / 1024).toFixed(2)}MB`
      })
      const isImage = file.type.startsWith('image/')
      const isLt2M = file.size / 1024 / 1024 < 2

      if (!isImage) {
        console.warn('封面格式不符合要求:', file.type)
        ElMessage.error('封面必须是图片格式!')
      }
      if (!isLt2M) {
        console.warn('封面文件过大:', `${(file.size / 1024 / 1024).toFixed(2)}MB`)
        ElMessage.error('封面图片大小不能超过 2MB!')
      }
      return isImage && isLt2M
    }

    const handleError = (error) => {
      console.error('上传出错:', error)
      ElMessage.error('上传失败，请重试')
      uploadProgress.value = 0
    }

    const handleTagClose = (tag) => {
      videoForm.value.tags = videoForm.value.tags.filter(t => t !== tag)
    }

    const showTagInput = () => {
      tagInputVisible.value = true
      nextTick(() => {
        tagInput.value.focus()
      })
    }

    const handleTagConfirm = () => {
      if (tagInputValue.value) {
        if (!videoForm.value.tags.includes(tagInputValue.value)) {
          videoForm.value.tags.push(tagInputValue.value)
        }
      }
      tagInputVisible.value = false
      tagInputValue.value = ''
    }

    const nextStep = () => {
      // 验证当前步骤
      if (currentStep.value === 0) {
        // 第一步：验证视频文件
        if (!videoFile.value) {
          ElMessage.error('请先选择视频文件')
          return
        }
        // 如果验证通过，进入下一步
        currentStep.value++
      }
    }

    const prevStep = () => {
      if (currentStep.value > 0) {
        currentStep.value--
      }
    }

    const publishVideo = async () => {
      if (!videoFile.value) {
        ElMessage.error('请先上传视频文件')
        return
      }

      if (!videoForm.value.title.trim()) {
        ElMessage.error('请输入视频标题')
        return
      }

      console.log('开始发布视频:', {
        formData: videoForm.value,
        videoFile: videoFile.value
      })

      publishing.value = true
      try {
        // 创建 FormData
        const formData = new FormData()
        formData.append('video_data', videoFile.value)
        formData.append('title', videoForm.value.title)
        if (videoForm.value.description) {
          formData.append('description', videoForm.value.description)
        }

        const response = await fetch('/api/v1/video/publish', {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${userStore.token}`
          },
          body: formData
        })

        const result = await response.json()
        console.log('发布视频响应:', result)
        
        if (result.code === '10000') {
          ElMessage.success('视频发布成功！')
          router.push('/video/list')
        } else {
          console.error('发布视频失败:', result)
          ElMessage.error(result.message || '发布失败，请重试')
        }
      } catch (error) {
        console.error('发布视频出错:', error)
        ElMessage.error('发布失败，请重试')
      }
      publishing.value = false
    }

    const percentageFormat = (percentage) => {
      return `${percentage}%`
    }

    return {
      currentStep,
      videoUrl,
      videoFile,
      uploadProgress,
      videoForm,
      categories,
      publishing,
      tagInputVisible,
      tagInputValue,
      tagInput,
      uploadHeaders,
      handleVideoChange,
      handleCoverSuccess,
      beforeCoverUpload,
      handleError,
      handleTagClose,
      showTagInput,
      handleTagConfirm,
      nextStep,
      prevStep,
      publishVideo,
      percentageFormat
    }
  }
}
</script>

<style scoped>
.upload-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.upload-card {
  margin-bottom: 20px;
}

.card-header {
  text-align: center;
}

.upload-step {
  margin-top: 30px;
  display: flex;
  flex-direction: column;
  min-height: 400px;
  justify-content: space-between;
}

.video-upload-container {
  flex: 1;
  margin-bottom: 40px;
  position: relative;
  z-index: 1;
}

.video-uploader {
  width: 100%;
  height: 300px;
}

.upload-preview {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-preview video {
  width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.upload-tip {
  font-size: 12px;
  color: #666;
  margin-top: 10px;
}

.step-actions {
  padding: 30px 0;
  text-align: center;
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: auto;
  position: relative;
  z-index: 10;
  background: white;
}

.cover-uploader {
  width: 200px;
  height: 112px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.cover-uploader:hover {
  border-color: #409EFF;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 200px;
  height: 112px;
  line-height: 112px;
  text-align: center;
}

.cover-preview {
  width: 200px;
  height: 112px;
  object-fit: cover;
}

.tag-item {
  margin-right: 10px;
  margin-bottom: 10px;
}

.tag-input {
  width: 100px;
  margin-right: 10px;
  vertical-align: bottom;
}

.publish-preview {
  text-align: center;
}

.preview-card {
  max-width: 500px;
  margin: 20px auto;
}

.preview-cover {
  width: 100%;
  height: 280px;
  object-fit: cover;
}

.preview-info {
  padding: 20px;
}

.preview-info h4 {
  margin: 0 0 10px 0;
}

.preview-tags {
  margin-top: 10px;
}

.preview-tags .el-tag {
  margin-right: 5px;
}
</style> 