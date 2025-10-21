<template>
    <div class="welcome-container">
        <div class="welcome-content">
            <!-- 主标题区域 -->
            <div class="hero-section">
                <div class="hero-avatar">
                    <img src="@/assets/img/avatar.jpg" alt="用户头像" class="avatar-image" />
                    <div class="avatar-ring"></div>
                </div>
                <h1 class="hero-title">欢迎回来</h1>
                <p class="hero-subtitle">让我们一起创造美好的数字体验</p>
                <div class="hero-decoration"></div>
            </div>

            <!-- 温馨提示卡片 -->
            <div class="greeting-card">
                <div class="greeting-content">
                    <div class="greeting-icon">
                        <n-icon size="24">
                            <SunnyOutline />
                        </n-icon>
                    </div>
                    <div class="greeting-text">
                        <h3>今天是美好的一天</h3>
                        <p>从左侧菜单开始您的工作，每一个项目都值得用心对待</p>
                    </div>
                </div>
            </div>

            <!-- 功能卡片区域 -->
            <div class="features-section">
                <div class="feature-card clickable" @click="navigateToServerManagement">
                    <div class="feature-icon server-icon">
                        <n-icon size="28">
                            <ServerOutline />
                        </n-icon>
                    </div>
                    <h3>服务器管理</h3>
                    <p>添加、编辑和管理您的服务器</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon project-icon">
                        <n-icon size="28">
                            <FolderOutline />
                        </n-icon>
                    </div>
                    <h3>项目配置</h3>
                    <p>统一管理所有项目设置</p>
                </div>

                <div class="feature-card clickable logout-card" @click="handleLogout">
                    <div class="feature-icon logout-icon">
                        <n-icon size="28">
                            <LogOutOutline />
                        </n-icon>
                    </div>
                    <h3>退出登录</h3>
                    <p>安全退出当前账户</p>
                </div>
            </div>

            <!-- 温馨提示 -->
            <div class="tips-section">
                <div class="tip-item">
                    <n-icon size="20" class="tip-icon">
                        <BulbOutline />
                    </n-icon>
                    <span>小贴士：使用左侧菜单快速导航到不同的服务器和项目</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { handleUnauthorized } from '@/api'
import {
    ServerOutline,
    FolderOutline,
    TerminalOutline,
    HeartOutline,
    SunnyOutline,
    BulbOutline,
    LogOutOutline,
    RefreshOutline
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const navigateToServerManagement = () => {
    router.push('/server-management')
}


const handleLogout = () => {
    dialog.warning({
        title: '确认退出',
        content: '确定要退出登录吗？',
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: () => {
            message.success('已退出登录')
            handleUnauthorized()
        }
    })
}
</script>

<style scoped>
.welcome-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100%;
    padding: 40px 20px;
}

.welcome-content {
    max-width: 800px;
    width: 100%;
    text-align: center;
}

/* 主标题区域 */
.hero-section {
    margin-bottom: 48px;
    position: relative;
}

.hero-avatar {
    margin-bottom: 24px;
    position: relative;
    display: inline-block;
    animation: gentle-pulse 3s ease-in-out infinite;
}

.avatar-image {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    object-fit: cover;
    border: 4px solid rgba(255, 255, 255, 0.8);
    box-shadow: 0 8px 32px rgba(52, 152, 219, 0.3);
    transition: all 0.3s ease;
}

.avatar-ring {
    position: absolute;
    top: -8px;
    left: -8px;
    right: -8px;
    bottom: -8px;
    border: 2px solid transparent;
    border-radius: 50%;
    background: linear-gradient(45deg, #3498db, #9b59b6, #e74c3c, #f39c12) border-box;
    -webkit-mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: subtract;
    mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0);
    mask-composite: subtract;
    animation: rotate-ring 4s linear infinite;
}

@keyframes rotate-ring {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

.hero-avatar:hover .avatar-image {
    transform: scale(1.05);
    box-shadow: 0 12px 40px rgba(52, 152, 219, 0.4);
}

.hero-title {
    font-size: 48px;
    font-weight: 600;
    color: #2c3e50;
    margin: 0 0 16px 0;
    letter-spacing: -1px;
    line-height: 1.2;
}

.hero-subtitle {
    font-size: 18px;
    color: #7f8c8d;
    margin: 0 0 32px 0;
    font-weight: 400;
    line-height: 1.6;
}

.hero-decoration {
    width: 80px;
    height: 4px;
    background: #3498db;
    margin: 0 auto;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(52, 152, 219, 0.3);
}

/* 温馨提示卡片 */
.greeting-card {
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(254, 207, 239, 0.1));
    backdrop-filter: blur(20px);
    border-radius: 24px;
    padding: 32px;
    margin-bottom: 40px;
    box-shadow: 0 20px 60px rgba(255, 107, 157, 0.15);
    border: 1px solid rgba(255, 255, 255, 0.4);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.greeting-card:hover {
    transform: translateY(-8px) scale(1.02);
    box-shadow: 0 30px 80px rgba(255, 107, 157, 0.25);
    border-color: rgba(255, 107, 157, 0.3);
}

.greeting-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    text-align: left;
}

.greeting-icon {
    width: 56px;
    height: 56px;
    background: linear-gradient(135deg, #ffd93d, #ff9a9e);
    border-radius: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    flex-shrink: 0;
    box-shadow: 0 8px 25px rgba(255, 217, 61, 0.4);
}

.greeting-text h3 {
    margin: 0 0 8px 0;
    font-size: 20px;
    font-weight: 600;
    color: #ff6b9d;
}

.greeting-text p {
    margin: 0;
    font-size: 16px;
    color: #ff9a9e;
    line-height: 1.5;
    opacity: 0.9;
}

/* 功能卡片区域 */
.features-section {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 24px;
    margin-bottom: 40px;
}

.feature-card {
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(254, 207, 239, 0.1));
    backdrop-filter: blur(20px);
    border-radius: 20px;
    padding: 32px 28px;
    text-align: center;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.4);
    position: relative;
    overflow: hidden;
}

.feature-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
    transition: left 0.6s ease;
}

.feature-card:hover::before {
    left: 100%;
}

.feature-card:hover {
    transform: translateY(-12px) scale(1.03);
    box-shadow: 0 25px 60px rgba(255, 107, 157, 0.2);
    border-color: rgba(255, 107, 157, 0.3);
}

.feature-card.clickable {
    cursor: pointer;
}

.feature-card.clickable:hover {
    transform: translateY(-12px) scale(1.03);
    box-shadow: 0 25px 60px rgba(108, 92, 231, 0.3);
    border-color: rgba(108, 92, 231, 0.4);
}

.feature-card:nth-child(3):hover {
    transform: translateY(-12px) scale(1.03);
    box-shadow: 0 25px 60px rgba(0, 206, 201, 0.3);
    border-color: rgba(0, 206, 201, 0.4);
}



.feature-card.logout-card:hover {
    transform: translateY(-12px) scale(1.03);
    box-shadow: 0 25px 60px rgba(225, 112, 85, 0.3);
    border-color: rgba(225, 112, 85, 0.4);
}

.feature-icon {
    width: 64px;
    height: 64px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    transition: all 0.3s ease;
}

.server-icon {
    background: linear-gradient(135deg, #6c5ce7, #a29bfe);
    color: white;
    box-shadow: 0 8px 25px rgba(108, 92, 231, 0.4);
}

.project-icon {
    background: linear-gradient(135deg, #fd79a8, #fdcb6e);
    color: white;
    box-shadow: 0 8px 25px rgba(253, 121, 168, 0.4);
}

.terminal-icon {
    background: linear-gradient(135deg, #00b894, #55efc4);
    color: white;
    box-shadow: 0 8px 25px rgba(0, 184, 148, 0.4);
}

.test-icon {
    background: linear-gradient(135deg, #00cec9, #55efc4);
    color: white;
    box-shadow: 0 8px 25px rgba(0, 206, 201, 0.4);
}





.logout-icon {
    background: linear-gradient(135deg, #e17055, #fab1a0);
    color: white;
    box-shadow: 0 8px 25px rgba(225, 112, 85, 0.4);
}

.feature-card:hover .feature-icon {
    transform: scale(1.1) rotate(5deg);
}

.feature-card h3 {
    margin: 0 0 12px 0;
    font-size: 18px;
    font-weight: 600;
    color: #ff6b9d;
}

.feature-card p {
    margin: 0;
    font-size: 14px;
    color: #ff9a9e;
    line-height: 1.5;
    opacity: 0.9;
}

/* 温馨提示区域 */
.tips-section {
    margin-top: 32px;
}

.tip-item {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 18px 28px;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.7), rgba(254, 207, 239, 0.1));
    backdrop-filter: blur(20px);
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.4);
    font-size: 14px;
    color: #ff9a9e;
    box-shadow: 0 8px 25px rgba(255, 107, 157, 0.1);
}

.tip-icon {
    color: #ffd93d;
    flex-shrink: 0;
    filter: drop-shadow(0 2px 4px rgba(255, 217, 61, 0.3));
}

/* 动画 */
@keyframes gentle-pulse {
    0%, 100% {
        transform: scale(1);
        opacity: 1;
    }
    50% {
        transform: scale(1.05);
        opacity: 0.8;
    }
}

.feature-card:nth-child(1) {
    animation: fade-in-up 0.6s ease 0.1s both;
}

.feature-card:nth-child(2) {
    animation: fade-in-up 0.6s ease 0.2s both;
}

.feature-card:nth-child(3) {
    animation: fade-in-up 0.6s ease 0.3s both;
}

@keyframes fade-in-up {
    from {
        opacity: 0;
        transform: translateY(30px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* 响应式设计 */
@media (max-width: 768px) {
    .welcome-container {
        padding: 24px 16px;
    }

    .hero-title {
        font-size: 36px;
    }

    .hero-subtitle {
        font-size: 16px;
    }

    .greeting-content {
        flex-direction: column;
        text-align: center;
        gap: 16px;
    }

    .greeting-text {
        text-align: center;
    }

    .features-section {
        grid-template-columns: 1fr;
        gap: 16px;
    }

    .greeting-card {
        padding: 24px;
    }
}

@media (max-width: 480px) {
    .hero-title {
        font-size: 28px;
    }

    .hero-subtitle {
        font-size: 14px;
    }

    .greeting-card {
        padding: 20px;
    }

    .feature-card {
        padding: 20px;
    }
}
</style>