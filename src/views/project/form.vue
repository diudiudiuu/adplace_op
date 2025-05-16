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

                    <el-form-item label="客户管理地址" prop="project_manage_url">
                        <el-input v-model="form.project_manage_url" placeholder="请输入客户管理地址"></el-input>
                    </el-form-item>

                    <el-form-item label="客户API地址" prop="project_api_url">
                        <el-input v-model="form.project_api_url" placeholder="请输入客户API地址"></el-input>
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
import { useSidebarStore } from '@/store/sidebar'
import api from '@/api'

const sidebar = useSidebarStore()
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
            project_manage_url: '',
            project_api_url: '',
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
        authorization: localStorage.getItem('authorization'),
    })
        .then(() => {
            ElMessage.success(isEdit ? '更新成功' : '添加成功')
            if (isEdit) {
                emit('editSuccess') // Notify parent on edit success
            } else {
                sidebar.setboolroute(true)
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
