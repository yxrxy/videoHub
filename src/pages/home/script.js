// 获取用户信息
async function fetchUserInfo() {
    try {
        const response = await fetch('/user/info', {
            headers: {
                'Authorization': localStorage.getItem('token')
            }
        });
        const data = await response.json();
        
        if (data.base.code === 0) {
            updateUserInfo(data.data);
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
    }
}

// 更新用户信息显示
function updateUserInfo(user) {
    document.getElementById('user-avatar').src = user.avatar_url;
    document.getElementById('menu-avatar').src = user.avatar_url;
    document.getElementById('username').textContent = user.username;
    document.getElementById('user-id').textContent = `ID: ${user.id}`;
}

// 加载视频列表
async function loadVideos() {
    try {
        const response = await fetch('/video/list', {
            headers: {
                'Authorization': localStorage.getItem('token')
            }
        });
        const data = await response.json();
        
        if (data.base.code === 0) {
            renderVideos(data.data.videos);
        }
    } catch (error) {
        console.error('加载视频列表失败:', error);
    }
}

// 渲染视频列表
function renderVideos(videos) {
    const grid = document.getElementById('video-grid');
    grid.innerHTML = videos.map(video => `
        <div class="video-card">
            <img src="${video.cover_url}" alt="${video.title}" class="video-thumbnail">
            <div class="video-info">
                <h3 class="video-title">${video.title}</h3>
                <div class="video-meta">
                    <div class="video-author">
                        <img src="${video.author.avatar_url}" alt="${video.author.username}" class="author-avatar">
                        <span>${video.author.username}</span>
                    </div>
                    <span>•</span>
                    <span>${formatViews(video.views)}次观看</span>
                </div>
            </div>
        </div>
    `).join('');
}

// 加载好友列表
async function loadFriends() {
    try {
        const response = await fetch('/user/friends', {
            headers: {
                'Authorization': localStorage.getItem('token')
            }
        });
        const data = await response.json();
        
        if (data.base.code === 0) {
            renderFriends(data.data.friends);
        }
    } catch (error) {
        console.error('加载好友列表失败:', error);
    }
}

// 渲染好友列表
function renderFriends(friends) {
    const list = document.getElementById('friends-list');
    list.innerHTML = friends.map(friend => `
        <div class="friend-item">
            <img src="${friend.avatar_url}" alt="${friend.username}" class="friend-avatar">
            <div class="friend-info">
                <div class="friend-name">${friend.username}</div>
                <div class="friend-status ${friend.online ? 'online' : ''}">${friend.online ? '在线' : '离线'}</div>
            </div>
        </div>
    `).join('');
}

// 显示上传模态框
function showUploadModal() {
    document.getElementById('upload-modal').classList.remove('hidden');
}

// 隐藏上传模态框
function hideUploadModal() {
    document.getElementById('upload-modal').classList.add('hidden');
}

// 处理视频上传
async function handleUpload(event) {
    event.preventDefault();
    
    const formData = new FormData();
    formData.append('file', document.getElementById('video-file').files[0]);
    formData.append('title', document.getElementById('video-title').value);
    formData.append('description', document.getElementById('video-description').value);

    try {
        const response = await fetch('/video/publish', {
            method: 'POST',
            headers: {
                'Authorization': localStorage.getItem('token')
            },
            body: formData
        });
        
        const data = await response.json();
        if (data.base.code === 0) {
            showToast('视频上传成功 ʕ•ᴥ•ʔ', 'success');
            hideUploadModal();
            loadVideos();  // 重新加载视频列表
        } else {
            showToast(data.base.msg, 'error');
        }
    } catch (error) {
        showToast('上传失败，请重试 (｡•́︿•̀｡)', 'error');
    }
}

// 处理登出
function handleLogout() {
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/auth/index.html';
}

// 切换用户菜单
function toggleUserMenu() {
    document.getElementById('user-menu-dropdown').classList.toggle('hidden');
}

// 格式化观看次数
function formatViews(views) {
    if (views >= 10000) {
        return (views / 10000).toFixed(1) + '万';
    }
    return views;
}

// 显示提示信息
function showToast(message, type = 'info') {
    // ... 与登录页面相同的 Toast 实现 ...
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    fetchUserInfo();
    loadVideos();
    loadFriends();
    
    // 设置上传区域的拖拽功能
    const uploadArea = document.getElementById('upload-area');
    const fileInput = document.getElementById('video-file');
    
    uploadArea.addEventListener('click', () => fileInput.click());
    
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = 'var(--primary-color)';
    });
    
    uploadArea.addEventListener('dragleave', () => {
        uploadArea.style.borderColor = '';
    });
    
    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = '';
        fileInput.files = e.dataTransfer.files;
    });
    
    // 点击页面其他地方关闭用户菜单
    document.addEventListener('click', (e) => {
        if (!e.target.closest('.user-menu')) {
            document.getElementById('user-menu-dropdown').classList.add('hidden');
        }
    });
}); 