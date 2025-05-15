/**
 * 格式化时间
 * @param {number|string|Date} time - 时间戳、时间字符串或Date对象
 * @returns {string} 格式化后的时间字符串
 */
export const formatTime = (time) => {
  if (!time) return ''
  
  const date = new Date(time)
  const now = new Date()
  
  // 如果是无效日期，返回空字符串
  if (isNaN(date.getTime())) return ''
  
  // 获取时间差（毫秒）
  const diff = now - date
  
  // 小于1分钟
  if (diff < 60 * 1000) {
    return '刚刚'
  }
  
  // 小于1小时
  if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))}分钟前`
  }
  
  // 小于24小时
  if (diff < 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 60 * 1000))}小时前`
  }
  
  // 小于30天
  if (diff < 30 * 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (24 * 60 * 60 * 1000))}天前`
  }
  
  // 其他情况显示完整日期
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  // 如果是今年，不显示年份
  if (year === now.getFullYear()) {
    return `${month}-${day} ${hours}:${minutes}`
  }
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

/**
 * 格式化日期
 * @param {number|string|Date} date - 时间戳、时间字符串或Date对象
 * @returns {string} 格式化后的日期字符串 (YYYY-MM-DD)
 */
export const formatDate = (date) => {
  if (!date) return ''
  
  const d = new Date(date)
  
  // 如果是无效日期，返回空字符串
  if (isNaN(d.getTime())) return ''
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
}

/**
 * 格式化数字
 * @param {number} num - 要格式化的数字
 * @returns {string} 格式化后的字符串
 */
export const formatNumber = (num) => {
  if (typeof num !== 'number') return '0'
  
  if (num >= 100000000) {
    return `${(num / 100000000).toFixed(1)}亿`
  }
  
  if (num >= 10000) {
    return `${(num / 10000).toFixed(1)}万`
  }
  
  return num.toString()
}

/**
 * 格式化视频时长
 * @param {number} seconds 视频时长（秒）
 * @returns {string} 格式化后的时长字符串
 */
export function formatDuration(seconds) {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
  }
  return `${minutes}:${String(secs).padStart(2, '0')}`
}

/**
 * 格式化文件大小
 * @param {number} bytes 文件大小（字节）
 * @returns {string} 格式化后的文件大小字符串
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
} 