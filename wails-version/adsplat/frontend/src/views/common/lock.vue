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
}

.login-content {
    height: 100%;
}

.login-card {
    border-radius: 16px !important;
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
}
</style>
