<template>
    <div>
        <el-form ref="formRef" :model="form" label-width="120px">
            <el-row :gutter="50">
                <el-col :span="12">
                    <el-form-item label="客户ID" prop="project_id">
                        <el-input v-model="form.project_id" placeholder="请输入客户ID"></el-input>
                    </el-form-item>

                    <el-form-item label="客户名称" prop="project_name">
                        <el-input v-model="form.project_name" placeholder="请输入客户名称"></el-input>
                    </el-form-item>

                    <el-form-item label="客户路径" prop="project_path">
                        <el-input v-model="form.project_path" placeholder="请输入客户路径"></el-input>
                    </el-form-item>

                    <el-form-item label="合同日期" prop="contract_date">
                        <el-date-picker v-model="form.contract_date" type="date" placeholder="选择日期"></el-date-picker>
                    </el-form-item>

                    <el-form-item label="客户管理地址" prop="project_manage_url">
                        <el-input v-model="form.project_manage_url" placeholder="请输入客户管理地址"></el-input>
                    </el-form-item>

                    <el-form-item label="客户API地址" prop="project_api_url">
                        <el-input v-model="form.project_api_url" placeholder="请输入客户API地址"></el-input>
                    </el-form-item>

                    <el-form-item label="API端口" prop="api_port">
                        <el-input v-model="form.api_port" placeholder="请输入API端口"></el-input>
                    </el-form-item>

                    <el-form-item label="前端端口" prop="front_port">
                        <el-input v-model="form.front_port" placeholder="请输入前端端口"></el-input>
                    </el-form-item>
                </el-col>

                <el-col :span="24">
                    <el-form-item>
                        <el-button type="primary" @click="onSubmit">添加客户</el-button>
                        <el-button @click="onReset">重置信息</el-button>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import api from '@/api'

const route = useRoute()
const serverId = ref(route.params.id)

const formRef = ref<FormInstance>()
const initialForm = {
    project_id: 's1p1',
    project_name: '测试客户',
    project_path: 'u1s1',
    contract_date: '2021-07-01',
    project_manage_url: 'https://www.baidu.com/',
    project_api_url: 'http://localhost:8848/v1',
    api_port: 3000,
    front_port: 3000,
}
const form = reactive({ ...initialForm })

// 提交
const onSubmit = () => {
    api('project_add', {
        serverId: serverId.value,
        projectInfo: JSON.stringify(form),
    })
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
        .then((res: any) => {
            ElMessage.success('添加成功')
            window.location.href = `/project/${form.project_id}`
        })
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
        .catch((err: any) => {
            ElMessage.error('添加失败')
        })
}

// 重置
const onReset = () => {
    Object.assign(form, initialForm) // Reset the form to its initial values
    formRef.value?.resetFields() // Optionally call resetFields for additional form handling
}
</script>
