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

                    <el-form-item label="版本" prop="version">
                        <el-input v-model="form.version" placeholder="请输入版本"></el-input>
                    </el-form-item>
                </el-col>

                <el-col :span="24">
                    <el-form-item>
                        <el-button type="primary" @click="onSubmit">{{ isEdit ? '更新客户' : '添加客户' }}</el-button>
                        <el-button @click="onReset">重置信息</el-button>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import api from '@/api'

const router = useRoute()
// Props for component mode and initial form data
const props = defineProps({
    mode: {
        type: String,
        default: 'add', // Default is 'add', can also be 'edit'
    },
    initialForm: {
        type: Object,
        default: () => ({
            project_id: '',
            project_name: '',
            contract_date: '',
            project_manage_url: '',
            project_api_url: '',
            api_port: '',
            front_port: '',
            version: '',
        }),
    },
    serverId: {
        type: String,
        required: false,
    },
})

const serverId = ref('')

// 编辑
const isEdit = props.mode === 'edit'
if (isEdit) {
    serverId.value = props.serverId as string
} else {
    serverId.value = router.params.id as string
}

const emit = defineEmits(['editSuccess'])

const formRef = ref<FormInstance>()
const form = reactive({ ...props.initialForm })

// Watch for changes to initialForm prop in edit mode
watch(
    () => props.initialForm,
    (newForm) => {
        if (isEdit) {
            Object.assign(form, newForm)
        }
    }
)

// Submit function
const onSubmit = () => {
    api('project_form', {
        serverId: serverId.value,
        projectInfo: JSON.stringify(form),
    })
        .then(() => {
            ElMessage.success(isEdit ? '更新成功' : '添加成功')
            if (isEdit) {
                emit('editSuccess') // Notify parent on edit success
            } else {
                // 刷新页面
                setTimeout(() => {
                    window.location.reload()
                }, 500)
            }
        })
        .catch(() => {
            ElMessage.error(isEdit ? '更新失败' : '添加失败')
        })
}

// Reset function
const onReset = () => {
    props.mode === 'edit'
        ? Object.assign(form, props.initialForm)
        : formRef.value?.resetFields()
}
</script>
