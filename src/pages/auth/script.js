// 表情变化函数
function changeMascotFace(face) {
    document.getElementById('mascot-face').textContent = face;
}

// 表单切换函数
function toggleForm(type) {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const mascotFace = document.getElementById('mascot-face');
    
    if (type === 'login') {
        registerForm.classList.add('hidden');
        loginForm.classList.remove('hidden');
        mascotFace.textContent = 'ʕ•ᴥ•ʔ';
    } else {
        loginForm.classList.add('hidden');
        registerForm.classList.remove('hidden');
        mascotFace.textContent = 'ʕ◕ᴥ◕ʔ';
    }
}

// 登录处理
async function handleLogin(event) {
    event.preventDefault();
    changeMascotFace('ʕ •`ᴥ´• ʔ');
    
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
        console.log('登录响应:', data);  // 添加日志
        
        if (data.base.code === 0) {
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('refresh_token', data.data.refresh_token);
            showToast('登录成功 ʕ•ᴥ•ʔ ✨', 'success');
            changeMascotFace('ʕ ᵔᴥᵔ ʔ');
            setTimeout(() => {
                window.location.href = '/home/index.html';
            }, 1500);
        } else {
            showToast(data.base.msg + ' (´･_･`)', 'error');
            changeMascotFace('ʕ ´•̥̥̥ ᴥ•̥̥̥` ʔ');
        }
    } catch (error) {
        showToast('网络错误，请稍后重试 (｡•́︿•̀｡)', 'error');
        changeMascotFace('ʕ ≧ᴥ≦ ʔ');
    }
}

// 注册处理
async function handleRegister(event) {
    event.preventDefault();
    changeMascotFace('ʕ •`ᴥ´• ʔ');
    
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const confirm = document.getElementById('register-confirm').value;

    if (password !== confirm) {
        showToast('两次输入的密码不一致 (｡•́︿•̀｡)', 'error');
        changeMascotFace('ʕ ´•̥̥̥ ᴥ•̥̥̥` ʔ');
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
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('refresh_token', data.data.refresh_token);
            showToast('注册成功 ʕ•ᴥ•ʔ ✨', 'success');
            changeMascotFace('ʕ ᵔᴥᵔ ʔ');
            setTimeout(() => {
                window.location.href = '/auth/index.html';
            }, 1500);
        } else {
            showToast(data.base.msg + ' (´･_･`)', 'error');
            changeMascotFace('ʕ ´•̥̥̥ ᴥ•̥̥̥` ʔ');
        }
    } catch (error) {
        showToast('网络错误，请稍后重试 (｡•́︿•̀｡)', 'error');
        changeMascotFace('ʕ ≧ᴥ≦ ʔ');
    }
}

// 显示提示信息
function showToast(message, type = 'info') {
    let toast = document.querySelector('.toast');
    if (toast) {
        toast.remove();
    }

    toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    // 添加图标
    const icon = document.createElement('span');
    icon.textContent = type === 'success' ? '🌟' : '💫';
    toast.appendChild(icon);
    
    // 添加消息
    const text = document.createElement('span');
    text.textContent = message;
    toast.appendChild(text);

    document.body.appendChild(toast);

    setTimeout(() => {
        toast.classList.add('show');
    }, 100);

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

// 确保表单提交时调用 handleLogin
document.getElementById('loginForm').addEventListener('submit', handleLogin); 