<template>
    <n-layout class="login-layout">
        <n-layout-content class="login-content">
            <n-space justify="center" align="center" style="height: 100vh">
                <n-card 
                    class="login-card"
                    :bordered="true"
                    size="large"
                    style="width: 400px; max-width: 90vw"
                >
                    <template #header>
                        <n-space justify="center">
                            <n-icon size="32" color="var(--primary-color)">
                                <LockClosedOutline />
                            </n-icon>
                        </n-space>
                        <n-h2 style="text-align: center; margin: 16px 0 0 0">
                            脑筋急转弯
                        </n-h2>
                    </template>
                    
                    <n-form 
                        size="large" 
                        @submit.prevent="submitForm"
                        :show-label="false"
                    >
                        <n-form-item>
                            <n-input 
                                type="text" 
                                placeholder="动物园里面生气时谁最安静?" 
                                v-model:value="param.password" 
                                size="large"
                                clearable
                            >
                                <template #prefix>
                                    <n-icon>
                                        <LockClosedOutline />
                                    </n-icon>
                                </template>
                            </n-input>
                        </n-form-item>
                        
                        <n-form-item>
                            <n-input 
                                type="password" 
                                placeholder="授权秘钥信息" 
                                v-model:value="param.authorization" 
                                size="large"
                                show-password-on="mousedown"
                                clearable
                            >
                                <template #prefix>
                                    <n-icon>
                                        <LockClosedOutline />
                                    </n-icon>
                                </template>
                            </n-input>
                        </n-form-item>
                        
                        <n-button 
                            type="primary" 
                            size="large" 
                            @click="submitForm" 
                            block
                            :loading="loading"
                        >
                            尝试你的答案
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
import { LockClosedOutline } from '@vicons/ionicons5'
import { setAuthorization } from '@/utils/auth'

const router = useRouter()
const message = useMessage()
const loading = ref(false)

const param = reactive({
    password: '',
    authorization: '',
})

const submitForm = async () => {
    if (param.password === '大猩猩') {
        if (!param.authorization || param.authorization.trim() === '') {
            message.error('请输入授权秘钥')
            return
        }

        loading.value = true
        const loadingMessage = message.loading('验证中...', { duration: 0 })
        
        try {
            // 模拟验证过程
            await new Promise(resolve => setTimeout(resolve, 1000))
            
            loadingMessage.destroy()
            setAuthorization(param.authorization)
            message.success('登录成功')
            
            await new Promise(resolve => setTimeout(resolve, 500))
            await router.push('/')
        } catch (error) {
            loadingMessage.destroy()
            message.error('登录失败，请重试')
        } finally {
            loading.value = false
        }
    } else {
        message.error('答案错误，请重试')
    }
}
</script>

<style scoped>
.login-layout {
    height: 100vh;
    background-image: url('@/assets/img/login-bg.jpg');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    background-attachment: fixed;
}

.login-content {
    height: 100%;
    background: rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(2px);
    -webkit-backdrop-filter: blur(2px);
}

.login-card {
    border-radius: 16px !important;
    background: rgba(255, 255, 255, 0.95) !important;
    backdrop-filter: blur(20px) !important;
    -webkit-backdrop-filter: blur(20px) !important;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3) !important;
    border: 1px solid rgba(255, 255, 255, 0.8) !important;
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
    font-weight: 600;
    color: #2c3e50;
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
    background: rgba(255, 255, 255, 0.9) !important;
    border: 1px solid rgba(255, 255, 255, 0.6) !important;
    border-radius: 8px !important;
}

:deep(.n-input:hover) {
    background: rgba(255, 255, 255, 1) !important;
    border-color: rgba(99, 102, 241, 0.5) !important;
}

:deep(.n-input.n-input--focus) {
    background: rgba(255, 255, 255, 1) !important;
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

/* 响应式设计 */
@media (max-width: 768px) {
    .login-layout {
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
}
</style>
