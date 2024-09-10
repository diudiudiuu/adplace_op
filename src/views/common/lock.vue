<template>
    <div class="login-bg">
        <div class="login-container">
            <div class="login-header">
                <div class="login-title">服务器管理</div>
            </div>
            <el-form :model="param" :rules="rules" ref="login" size="large">
                <el-form-item prop="password">
                    <el-input type="password" placeholder="请输入锁屏幕密码" v-model="param.password">
                        <template #prepend>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>

                <el-button
                    class="login-btn"
                    type="primary"
                    size="large"
                    @click="submitForm(login)"
                >解锁</el-button>
            </el-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

interface LoginInfo {
    password: string
}

const lgStr = localStorage.getItem('login-param')
const defParam = lgStr ? JSON.parse(lgStr) : null

const router = useRouter()
const param = reactive<LoginInfo>({
    password: defParam ? defParam.password : '',
})

const rules: FormRules = {
    password: [{ required: true, message: '请输入锁屏密码', trigger: 'blur' }],
}

const login = ref<FormInstance>()
const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid: boolean): void => {
        if (valid) {
            ElMessage.success('解锁成功')
            router.push('/')
        }
    })
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
