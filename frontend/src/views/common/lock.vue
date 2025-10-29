<template>
    <n-layout class="login-layout">
        <n-layout-content class="login-content">
            <n-space justify="center" align="center" style="height: 100vh">
                <n-card class="login-card" :bordered="true" size="large" style="width: 400px; max-width: 90vw">
                    <template #header>
                        <n-space justify="center">
                            <div class="hero-avatar">
                                <img src="@/assets/img/avatar.jpg" alt="用户头像" class="avatar-image" />
                                <div class="avatar-ring"></div>
                            </div>
                        </n-space>
                        <n-h2 style="text-align: center; margin: 16px 0 8px 0">
                            欢迎回来
                        </n-h2>
                        <p
                            style="text-align: center; margin: 0; color: #2d3748; font-size: 14px; font-weight: 500; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);">
                            请输入授权密钥和验证码以继续
                        </p>
                    </template>

                    <n-form size="large" @submit.prevent="submitForm" :show-label="false">
                        <n-form-item>
                            <n-input type="password" placeholder="请输入您的专属授权密钥" v-model:value="param.authorization"
                                size="large" show-password-on="mousedown" clearable>
                                <template #prefix>
                                    <n-icon>
                                        <LockClosedOutline />
                                    </n-icon>
                                </template>
                            </n-input>
                        </n-form-item>

                        <n-form-item>
                            <div class="captcha-row">
                                <n-input type="text" placeholder="请输入验证码" v-model:value="param.captcha" size="large"
                                    clearable maxlength="4" class="captcha-input">
                                    <template #prefix>
                                        <n-icon>
                                            <ShieldCheckmarkOutline />
                                        </n-icon>
                                    </template>
                                </n-input>
                                <Captcha @change="onCaptchaChange" ref="captchaRef" />
                            </div>
                        </n-form-item>

                        <n-button type="primary" size="large" @click="submitForm" block :loading="loading">
                            开始我的数字之旅
                        </n-button>
                    </n-form>
                </n-card>
            </n-space>
        </n-layout-content>
    </n-layout>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { LockClosedOutline, ShieldCheckmarkOutline } from '@vicons/ionicons5'
import { setAuthorization } from '@/utils/auth'
import Captcha from '@/components/Captcha.vue'

const router = useRouter()
const message = useMessage()
const loading = ref(false)
const captchaRef = ref()
const correctCaptcha = ref('')

const param = reactive({
    authorization: '',
    captcha: '',
})

const onCaptchaChange = (code: string) => {
    correctCaptcha.value = code
}

const submitForm = async () => {
    // 验证授权密钥
    if (!param.authorization || param.authorization.trim() === '') {
        message.error('请输入您的专属密钥')
        return
    }

    // 验证验证码
    if (!param.captcha || param.captcha.trim() === '') {
        message.error('请输入验证码')
        return
    }

    if (param.captcha.toUpperCase() !== correctCaptcha.value.toUpperCase()) {
        message.error('验证码错误，请重新输入')
        param.captcha = ''
        captchaRef.value?.refresh()
        return
    }

    loading.value = true
    const loadingMessage = message.loading('正在为您开启美好体验...', { duration: 0 })

    try {
        // 模拟验证过程
        await new Promise(resolve => setTimeout(resolve, 1000))

        loadingMessage.destroy()
        setAuthorization(param.authorization)
        message.success('欢迎回来！准备开始您的精彩旅程')

        await new Promise(resolve => setTimeout(resolve, 500))
        await router.push('/')
    } catch (error) {
        loadingMessage.destroy()
        message.error('验证遇到问题，请稍后再试')
    } finally {
        loading.value = false
    }
}
</script>

<style scoped>
.login-layout {
    height: 100vh;
    width: 100vw;
    background-image: url('@/assets/img/login-bg.jpg') !important;
    background-size: cover !important;
    background-position: center center !important;
    background-repeat: no-repeat !important;
    background-attachment: fixed !important;
    position: relative;
    overflow: hidden;
}

.login-content {
    height: 100%;
    width: 100%;
    background: transparent;
    position: relative;
    z-index: 1;
}

.login-card {
    border-radius: 16px !important;
    background: rgba(255, 255, 255, 0.4) !important;
    backdrop-filter: blur(8px) !important;
    -webkit-backdrop-filter: blur(8px) !important;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15) !important;
    border: 1px solid rgba(255, 255, 255, 0.3) !important;
}

:deep(.n-card-header) {
    padding: 32px 32px 16px;
}

:deep(.n-card__content) {
    padding: 16px 32px 32px;
}

:deep(.n-button) {
    height: 48px;
    font-weight: 600;
    transition: all 0.3s ease;
}

:deep(.n-button:hover) {
    transform: translateY(-1px);
}

:deep(.n-h2) {
    font-weight: 700;
    color: #1a202c;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* 登录卡片动画 */
.login-card {
    animation: fadeInUp 0.8s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(30px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* 输入框样式优化 */
:deep(.n-input) {
    background: rgba(255, 255, 255, 0.5) !important;
    border: 1px solid rgba(255, 255, 255, 0.4) !important;
    border-radius: 8px !important;
    backdrop-filter: blur(5px) !important;
    -webkit-backdrop-filter: blur(5px) !important;
    color: #1a202c !important;
    font-weight: 500 !important;
}

:deep(.n-input .n-input__input-el) {
    color: #1a202c !important;
    font-weight: 500 !important;
}

:deep(.n-input .n-input__placeholder) {
    color: #4a5568 !important;
    font-weight: 400 !important;
}

:deep(.n-input:hover) {
    background: rgba(255, 255, 255, 0.6) !important;
    border-color: rgba(99, 102, 241, 0.5) !important;
}

:deep(.n-input.n-input--focus) {
    background: rgba(255, 255, 255, 0.7) !important;
    border-color: rgba(99, 102, 241, 0.8) !important;
    box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2) !important;
}

/* 按钮样式优化 */
:deep(.n-button--primary-type) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
    border: none !important;
    border-radius: 8px !important;
}

:deep(.n-button--primary-type:hover) {
    background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%) !important;
    transform: translateY(-2px) !important;
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4) !important;
}

/* 头像样式 */
.hero-avatar {
    margin-bottom: 16px;
    position: relative;
    display: inline-block;
    animation: gentle-pulse 3s ease-in-out infinite;
}

.avatar-image {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    border: 3px solid rgba(255, 255, 255, 0.8);
    box-shadow: 0 8px 32px rgba(52, 152, 219, 0.3);
    transition: all 0.3s ease;
}

.avatar-ring {
    position: absolute;
    top: -6px;
    left: -6px;
    right: -6px;
    bottom: -6px;
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

@keyframes gentle-pulse {

    0%,
    100% {
        transform: scale(1);
        opacity: 1;
    }

    50% {
        transform: scale(1.05);
        opacity: 0.8;
    }
}

.hero-avatar:hover .avatar-image {
    transform: scale(1.05);
    box-shadow: 0 12px 40px rgba(52, 152, 219, 0.4);
}

/* 验证码行布局 */
.captcha-row {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
}

.captcha-input {
    flex: 1;
}

/* 移除额外的背景层，使用主背景 */

/* 响应式设计 */
@media (max-width: 768px) {
    .login-layout {
        background-attachment: scroll;
    }

    .login-layout::before {
        background-attachment: scroll;
    }

    .login-card {
        margin: 20px !important;
        width: calc(100vw - 40px) !important;
        max-width: none !important;
    }

    :deep(.n-card-header) {
        padding: 24px 24px 12px !important;
    }

    :deep(.n-card__content) {
        padding: 12px 24px 24px !important;
    }

    .avatar-image {
        width: 60px;
        height: 60px;
    }

    .captcha-row {
        gap: 8px;
    }
}
</style>
