<template>
    <div class="login-bg select-none">
        <div class="login-container">
            <div class="login-header">
                <div class="login-title">脑筋急转弯</div>
            </div>
            <el-form size="large" @submit.prevent="submitForm">
                <el-form-item>
                    <el-input type="text" placeholder="动物园里面生气时谁最安静?" v-model="param.password">
                        <template #prepend>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item>
                    <el-input type="password" placeholder="授权秘钥" v-model="param.authorization">
                        <template #prepend>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-button class="login-btn" type="primary" size="large" @click="submitForm">尝试你的答案</el-button>
            </el-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import layoutRouter from '@/router/layout'

const router = useRouter()
const routes = router.getRoutes()

const param = reactive({
    password: '大猩猩',
    authorization: 'koyiq1c42xl9mdb80t5rwn7aehv6zjpf',
})

const layout = routes.find((item) => item.name === 'layout')
// 删除 layout路由
if (layout) {
    router.removeRoute('layout')
    //删除 authorization
    localStorage.removeItem('authorization')
}

const submitForm = () => {
    if (param.password === '大猩猩') {
        localStorage.setItem('authorization', param.authorization)
        ElMessage.success('🤙🤙🤙,你非常棒,居然猜对了')
        // 添加 layout路由
        if (layout) {
            router.addRoute(layout)
        } else {
            router.addRoute(layoutRouter)
        }
        router.replace({ path: '/' })
    } else {
        ElMessage.error('🖕🖕🖕换个姿势,再来一次')
    }
}
</script>

<style scoped>
.login-bg {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100vh;
    background: url(../../assets/img/login-bg.jpg) center/cover no-repeat;
}

.login-header {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 40px;
}

.logo {
    width: 35px;
}

.login-title {
    font-size: 22px;
    color: #333;
    font-weight: bold;
}

.login-container {
    width: 450px;
    border-radius: 5px;
    background: #fff;
    padding: 40px 50px 50px;
    box-sizing: border-box;
}

.login-btn {
    display: block;
    width: 100%;
}

.login-tips {
    font-size: 12px;
    color: #999;
}

.login-text {
    display: flex;
    align-items: center;
    margin-top: 20px;
    font-size: 14px;
    color: #787878;
}
</style>
