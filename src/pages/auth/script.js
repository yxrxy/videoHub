// 表单切换函数
function toggleForm(type) {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    if (type === 'login') {
        registerForm.classList.add('hidden');
        loginForm.classList.remove('hidden');
    } else {
        loginForm.classList.add('hidden');
        registerForm.classList.remove('hidden');
    }
}

// 登录处理
async function handleLogin(event) {
    event.preventDefault();
    
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch('/user/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();
        
        if (data.base.code === 0) {
            // 登录成功
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('refresh_token', data.data.refresh_token);
            showToast('登录成功', 'success');
            // 跳转到主页
            setTimeout(() => {
                window.location.href = '/';
            }, 1500);
        } else {
            showToast(data.base.msg, 'error');
        }
    } catch (error) {
        showToast('网络错误，请稍后重试', 'error');
    }
}

// 注册处理
async function handleRegister(event) {
    event.preventDefault();
    
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const confirm = document.getElementById('register-confirm').value;

    if (password !== confirm) {
        showToast('两次输入的密码不一致', 'error');
        return;
    }

    try {
        const response = await fetch('/user/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();
        
        if (data.base.code === 0) {
            // 注册成功
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('refresh_token', data.data.refresh_token);
            showToast('注册成功', 'success');
            // 跳转到主页
            setTimeout(() => {
                window.location.href = '/';
            }, 1500);
        } else {
            showToast(data.base.msg, 'error');
        }
    } catch (error) {
        showToast('网络错误，请稍后重试', 'error');
    }
}

// 显示提示信息
function showToast(message, type = 'info') {
    // 检查是否已存在 toast
    let toast = document.querySelector('.toast');
    if (toast) {
        toast.remove();
    }

    // 创建新的 toast
    toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;

    document.body.appendChild(toast);

    // 添加动画
    setTimeout(() => {
        toast.classList.add('show');
    }, 100);

    // 3秒后移除
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
            toast.remove();
        }, 300);
    }, 3000);
}

// 添加 Toast 样式
const style = document.createElement('style');
style.textContent = `
    .toast {
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 24px;
        border-radius: 4px;
        color: white;
        font-size: 14px;
        opacity: 0;
        transform: translateY(-20px);
        transition: all 0.3s ease;
    }

    .toast.show {
        opacity: 1;
        transform: translateY(0);
    }

    .toast-success {
        background-color: var(--success-color);
    }

    .toast-error {
        background-color: var(--error-color);
    }

    .toast-info {
        background-color: var(--primary-color);
    }
`;
document.head.appendChild(style); 