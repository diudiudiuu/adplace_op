<template>
    <div>
        <n-card :title="isEdit ? '编辑客户' : '添加客户'">
            <n-form ref="formRef" :model="form" label-placement="left" label-width="120">
                <n-grid :cols="24" :x-gap="24">
                    <n-grid-item :span="12">
                        <n-form-item label="客户ID" path="project_id" required>
                            <n-input 
                                v-model:value="form.project_id" 
                                placeholder="请输入客户ID"
                                :disabled="isEdit"
                            />
                        </n-form-item>

                        <n-form-item label="客户名称" path="project_name" required>
                            <n-input 
                                v-model:value="form.project_name" 
                                placeholder="请输入客户名称"
                            />
                        </n-form-item>

                        <n-form-item label="客户管理地址" path="project_manage_url" required>
                            <n-input 
                                v-model:value="form.project_manage_url" 
                                placeholder="请输入客户管理地址"
                            />
                        </n-form-item>

                        <n-form-item label="客户API地址" path="project_api_url" required>
                            <n-input 
                                v-model:value="form.project_api_url" 
                                placeholder="请输入客户API地址"
                            />
                        </n-form-item>
                    </n-grid-item>

                    <n-grid-item :span="24">
                        <n-form-item>
                            <n-space>
                                <n-button type="primary" @click="onSubmit">
                                    <template #icon>
                                        <n-icon>
                                            <SaveOutline v-if="!isEdit" />
                                            <RefreshOutline v-else />
                                        </n-icon>
                                    </template>
                                    {{ isEdit ? '更新客户' : '添加客户' }}
                                </n-button>
                                <n-button @click="onReset">
                                    <template #icon>
                                        <n-icon><ReloadOutline /></n-icon>
                                    </template>
                                    重置信息
                                </n-button>
                            </n-space>
                        </n-form-item>
                    </n-grid-item>
                </n-grid>
            </n-form>
        </n-card>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useRoute, useRouter  } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { SaveOutline, RefreshOutline, ReloadOutline } from '@vicons/ionicons5'
import api from '@/api'

const sidebar = useSidebarStore()
const router = useRoute()
const route = useRouter()
const message = useMessage()

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

const formRef = ref<FormInst>()
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
            message.success(isEdit ? '更新成功' : '添加成功')
            if (isEdit) {
                emit('editSuccess') // Notify parent on edit success
            } else {
                sidebar.setboolroute(true)
                // 跳转到welcome页面
                route.push(`/project/${serverId.value}/${form.project_id}`)
            }
        })
        .catch(() => {
            message.error(isEdit ? '更新失败' : '添加失败')
        })
}

// Reset function
const onReset = () => {
    props.mode === 'edit'
        ? Object.assign(form, props.initialForm)
        : formRef.value?.resetFields()
}
</script>
