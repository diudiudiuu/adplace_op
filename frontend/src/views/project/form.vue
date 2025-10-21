<template>
    <div>
        <n-card :title="isEdit ? '编辑客户' : '添加客户'">
            <n-form ref="formRef" :model="form" label-placement="left" label-width="120">
                <n-grid :cols="24" :x-gap="24">
                    <n-grid-item :span="12">
                        <n-form-item label="服务器ID">
                            <n-input 
                                :value="serverId" 
                                placeholder="服务器ID"
                                disabled
                            />
                        </n-form-item>

                        <n-form-item label="客户ID" path="project_id" required>
                            <n-input 
                                v-model:value="form.project_id" 
                                placeholder="请输入客户ID（小写字母数字）"
                                :disabled="isEdit"
                                @input="handleProjectIdInput"
                            >
                                <template #suffix v-if="!isEdit">
                                    <n-button 
                                        text 
                                        type="primary" 
                                        @click="generateRandomId"
                                        size="small"
                                    >
                                        随机
                                    </n-button>
                                </template>
                            </n-input>
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
                            >
                                <template #suffix v-if="!isEdit">
                                    <n-button 
                                        text 
                                        type="info" 
                                        @click="updateUrlsFromId"
                                        size="small"
                                    >
                                        自动
                                    </n-button>
                                </template>
                            </n-input>
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
import { reloadMenus } from '@/components/menu'
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
    serverId.value = router.params.pid as string
}

console.log('Project form initialized:', {
    isEdit,
    serverId: serverId.value,
    routerParams: router.params,
    propsServerId: props.serverId
})

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

// 监听客户ID变化，自动生成URL
watch(
    () => form.project_id,
    (newProjectId) => {
        if (!isEdit && newProjectId && validateProjectId(newProjectId)) {
            updateUrls(newProjectId)
        }
    }
)

// Submit function
const onSubmit = async () => {
    // 验证必填字段
    if (!form.project_id || !form.project_name || !form.project_manage_url || !form.project_api_url) {
        message.error('请填写所有必填字段')
        return
    }
    
    // 验证客户ID格式
    if (!validateProjectId(form.project_id)) {
        message.error('客户ID只能包含小写字母和数字')
        return
    }
    
    if (!serverId.value) {
        message.error('服务器ID不能为空')
        return
    }
    
    try {
        console.log('Submitting project form:', {
            serverId: serverId.value,
            form: form
        })
        
        const res = await api('project_form', {
            serverId: serverId.value,
            projectInfo: JSON.stringify(form),
        })
        
        console.log('Project form response:', res)
        
        if (res && (res.code === 200 || res.success)) {
            message.success(isEdit ? '更新成功' : '添加成功')
            if (isEdit) {
                emit('editSuccess') // Notify parent on edit success
            } else {
                // 刷新菜单
                await reloadMenus()
                sidebar.setboolroute(true)
                // 跳转到项目页面
                route.push(`/project/${serverId.value}/${form.project_id}`)
            }
        } else {
            const errorMsg = res?.msg || res?.message || (isEdit ? '更新失败' : '添加失败')
            message.error(errorMsg)
        }
    } catch (error) {
        console.error('Project form submission error:', error)
        message.error(isEdit ? '更新失败' : '添加失败')
    }
}

// 根据客户ID自动生成URL
const updateUrls = (projectId: string) => {
    if (projectId && projectId.length > 0) {
        form.project_manage_url = `https://manage-${projectId}.adswds.com`
        form.project_api_url = `https://api-${projectId}.adswds.com/v1`
    }
}

// 手动触发URL更新
const updateUrlsFromId = () => {
    if (!form.project_id) {
        message.warning('请先输入客户ID')
        return
    }
    
    if (!validateProjectId(form.project_id)) {
        message.error('客户ID格式不正确')
        return
    }
    
    updateUrls(form.project_id)
    message.success('已根据客户ID自动生成URL')
}

// 客户ID输入验证 - 只允许小写字母和数字
const handleProjectIdInput = (value: string) => {
    // 只保留小写字母和数字，并转换大写为小写
    let filteredValue = value.toLowerCase().replace(/[^a-z0-9]/g, '')
    
    if (filteredValue !== value) {
        form.project_id = filteredValue
        message.warning('客户ID只能包含小写字母和数字')
    }
}

// 生成随机客户ID
const generateRandomId = () => {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
    
    // 随机生成3-6位长度
    const length = Math.floor(Math.random() * 4) + 3 // 3-6位
    let result = ''
    
    // 生成随机字符
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    
    form.project_id = result
    // 自动生成URL
    updateUrls(result)
    message.success(`已生成${length}位随机客户ID并更新URL`)
}

// 验证客户ID格式
const validateProjectId = (value: string): boolean => {
    // 检查是否只包含小写字母和数字，不限制长度
    const regex = /^[a-z0-9]+$/
    return regex.test(value) && value.length > 0
}

// Reset function
const onReset = () => {
    if (props.mode === 'edit') {
        Object.assign(form, props.initialForm)
    } else {
        formRef.value?.resetFields()
        // 重置时清空客户ID
        form.project_id = ''
    }
}
</script>

<style scoped>
/* 表单样式优化 */
:deep(.n-form-item-label) {
    font-size: var(--table-font-size, 12px);
    font-weight: 500;
}

:deep(.n-input) {
    font-size: var(--table-font-size, 12px);
}

:deep(.n-button) {
    font-size: var(--table-font-size, 12px);
}

/* 输入框内按钮样式优化 */
:deep(.n-input .n-input-wrapper .n-input__suffix .n-button) {
    margin-right: 4px;
    padding: 0 8px;
    height: 24px;
    font-size: 11px;
}
</style>
