import request from '@/utils/request'

// WebSocket服务
const wsService = {
  // WebSocket实例
  ws: null,
  messageHandler: null,

  // 连接WebSocket
  connect(token, roomId = -1, toUser = null) {
    if (this.ws) {
      this.close()
    }

    let wsUrl = `ws://localhost:8080/api/v1/social/ws/connect?token=${token}&room_id=${roomId}`
    if (toUser) {
      wsUrl += `&to_user=${toUser}`
    }

    this.ws = new WebSocket(wsUrl)
    
    this.ws.onopen = () => {
      console.log('WebSocket连接已建立')
    }
    
    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data)
      if (this.messageHandler) {
        this.messageHandler(message)
      }
    }
    
    this.ws.onclose = () => {
      console.log('WebSocket连接已关闭')
      this.ws = null
    }
    
    this.ws.onerror = (error) => {
      console.error('WebSocket错误:', error)
      this.ws = null
    }
  },

  // 发送消息
  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket未连接')
    }
  },

  // 发送私聊消息
  sendPrivateMessage(toUser, content) {
    if (!this.ws) return
    this.ws.send(JSON.stringify({
      type: 'private',
      to: toUser,
      content: content
    }))
  },

  // 发送群聊消息
  sendChatMessage(content) {
    if (!this.ws) return
    this.ws.send(JSON.stringify({
      type: 'chat',
      content: content
    }))
  },

  // 关闭连接
  close() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  },

  // 设置消息处理回调
  setMessageHandler(handler) {
    this.messageHandler = handler
  }
}

export { wsService }

// 私聊相关接口
export const messageApi = {
  // 获取私聊消息历史
  getMessageHistory(params) {
    return request({
      url: '/api/v1/social/private/messages',
      method: 'get',
      params
    })
  },

  // 标记消息已读
  markMessageRead(data) {
    return request({
      url: '/api/v1/social/message/read',
      method: 'post',
      data
    })
  },

  // 获取未读消息数
  getUnreadCount() {
    return request({
      url: '/api/v1/social/message/unread/count',
      method: 'get'
    })
  }
}

// 聊天室相关接口
export const chatRoomApi = {
  // 创建聊天室
  createChatRoom(data) {
    return request({
      url: '/api/v1/social/chatroom',
      method: 'post',
      data
    })
  },

  // 获取聊天室列表
  getChatRooms(params) {
    return request({
      url: '/api/v1/social/chatrooms',
      method: 'get',
      params
    })
  },

  // 获取聊天室详情
  getChatRoomDetail(roomId) {
    return request({
      url: `/api/v1/social/chatroom/${roomId}`,
      method: 'get'
    })
  },

  // 获取聊天室消息历史
  getMessageHistory(params) {
    return request({
      url: `/api/v1/social/chatroom/${params.roomId}/messages`,
      method: 'get',
      params
    })
  },

  // 设置成员角色
  setMemberRole(data) {
    return request({
      url: '/api/v1/social/chatroom/member/role',
      method: 'put',
      data
    })
  },

  // 移除成员
  removeMember(data) {
    return request({
      url: '/api/v1/social/chatroom/member/remove',
      method: 'post',
      data
    })
  }
}

// 好友相关接口
export const friendApi = {
  // 获取好友列表
  getFriends(params) {
    return request({
      url: '/api/v1/social/friends',
      method: 'get',
      params
    })
  },

  // 发送好友申请
  sendFriendRequest(data) {
    return request({
      url: '/api/v1/social/friend',
      method: 'post',
      data
    })
  },

  // 获取好友申请列表
  getFriendRequests(params) {
    return request({
      url: '/api/v1/social/friend/requests',
      method: 'get',
      params
    })
  },

  // 处理好友申请
  handleFriendRequest(data) {
    return request({
      url: '/api/v1/social/friend/request/handle',
      method: 'post',
      data
    })
  },

  // 修改好友备注
  updateFriendRemark(data) {
    return request({
      url: '/api/v1/social/friend/remark',
      method: 'put',
      data
    })
  },

  // 删除好友
  deleteFriend(friendId) {
    return request({
      url: `/api/v1/social/friend/${friendId}`,
      method: 'delete'
    })
  }
} 