<template>
    <div>
        <n-card :title="isEdit ? '编辑客户' : '添加客户'">
            <n-form ref="formRef" :model="form" label-placement="left" label-width="120">
                <n-grid :cols="24" :x-gap="24">
                    <n-grid-item :span="12">
                        <n-form-item label="服务器ID">
                            <n-input :value="serverId" placeholder="服务器ID" disabled />
                        </n-form-item>

                        <n-form-item label="客户ID" path="project_id" required>
                            <n-input v-model:value="form.project_id" placeholder="请输入客户ID（小写字母数字）" :disabled="isEdit"
                                @input="handleProjectIdInput">
                                <template #suffix v-if="!isEdit">
                                    <n-button text type="primary" @click="generateRandomId" size="small">
                                        随机
                                    </n-button>
                                </template>
                            </n-input>
                        </n-form-item>

                        <n-form-item label="客户名称" path="project_name" required>
                            <n-input v-model:value="form.project_name" placeholder="请输入客户名称" />
                        </n-form-item>

                        <n-form-item label="客户管理地址" path="project_manage_url" required>
                            <n-input v-model:value="form.project_manage_url" placeholder="请输入客户管理地址">
                                <template #suffix v-if="!isEdit">
                                    <n-button text type="info" @click="updateUrlsFromId" size="small">
                                        自动
                                    </n-button>
                                </template>
                            </n-input>
                        </n-form-item>

                        <n-form-item label="客户API地址" path="project_api_url" required>
                            <n-input v-model:value="form.project_api_url" placeholder="请输入客户API地址" />
                        </n-form-item>
                    </n-grid-item>

                    <n-grid-item :span="12">
                        <n-form-item label="API端口" path="api_port" required>
                            <n-input v-model:value="form.api_port" placeholder="9000" type="number">
                                <template #suffix v-if="!isEdit">
                                    <n-button text type="primary" @click="generateDefaultPorts" size="small">
                                        自动
                                    </n-button>
                                </template>
                            </n-input>
                        </n-form-item>

                        <n-form-item label="前端端口" path="front_port" required>
                            <n-input v-model:value="form.front_port" placeholder="3000" type="number" />
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
                                        <n-icon>
                                            <ReloadOutline />
                                        </n-icon>
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
import { reactive, ref, watch, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import dataManager from '@/utils/dataManager'
import { SaveOutline, RefreshOutline, ReloadOutline } from '@vicons/ionicons5'
import api from '@/api'

const sidebar = useSidebarStore()
const router = useRoute()
const route = useRouter()
const message = useMessage()

// 注入全局 loading
const globalLoading = inject('globalLoading') as any

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
            api_port: '',
            front_port: '',
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

// 获取服务器下所有项目的端口号
const getServerPorts = async (): Promise<{ apiPorts: number[], frontPorts: number[] }> => {
    try {
        // 优先从缓存获取服务器信息
        const serverInfo = await dataManager.getServerById(serverId.value)

        const apiPorts: number[] = []
        const frontPorts: number[] = []

        if (serverInfo && serverInfo.project_list) {
            serverInfo.project_list.forEach((project: any) => {
                if (project.api_port) {
                    const apiPort = parseInt(project.api_port)
                    if (!isNaN(apiPort)) {
                        apiPorts.push(apiPort)
                    }
                }
                if (project.front_port) {
                    const frontPort = parseInt(project.front_port)
                    if (!isNaN(frontPort)) {
                        frontPorts.push(frontPort)
                    }
                }
            })
        }

        return { apiPorts, frontPorts }
    } catch (error) {
        console.error('Failed to get server ports:', error)
        return { apiPorts: [], frontPorts: [] }
    }
}

// 自动生成默认端口
const generateDefaultPorts = async () => {
    try {
        const { apiPorts, frontPorts } = await getServerPorts()

        // 计算API端口：找到最大值+1，如果没有现有端口则使用9000
        let nextApiPort = 9000
        if (apiPorts.length > 0) {
            const maxApiPort = Math.max(...apiPorts)
            nextApiPort = maxApiPort + 1
        }

        // 计算前端端口：找到最大值+1，如果没有现有端口则使用3000
        let nextFrontPort = 3000
        if (frontPorts.length > 0) {
            const maxFrontPort = Math.max(...frontPorts)
            nextFrontPort = maxFrontPort + 1
        }

        form.api_port = nextApiPort.toString()
        form.front_port = nextFrontPort.toString()

        message.success(`已自动生成端口：API端口 ${nextApiPort}，前端端口 ${nextFrontPort}`)
    } catch (error) {
        console.error('Failed to generate default ports:', error)
        message.error('自动生成端口失败')
    }
}

// 监听客户ID变化，自动生成URL
watch(
    () => form.project_id,
    (newProjectId) => {
        if (!isEdit && newProjectId && validateProjectId(newProjectId)) {
            updateUrls(newProjectId)
        }
    }
)

// 在添加模式下，组件初始化时自动生成默认端口
if (!isEdit && serverId.value) {
    generateDefaultPorts()
}

// Submit function
const onSubmit = async () => {
    // 验证必填字段
    if (!form.project_id || !form.project_name || !form.project_manage_url || !form.project_api_url || !form.api_port || !form.front_port) {
        message.error('请填写所有必填字段')
        return
    }

    // 验证客户ID格式
    if (!validateProjectId(form.project_id)) {
        message.error('客户ID只能包含小写字母和数字')
        return
    }

    // 验证端口格式
    if (!validatePort(form.api_port)) {
        message.error('API端口必须是1-65535之间的数字')
        return
    }

    if (!validatePort(form.front_port)) {
        message.error('前端端口必须是1-65535之间的数字')
        return
    }

    if (!serverId.value) {
        message.error('服务器ID不能为空')
        return
    }

    try {
        globalLoading.show(props.mode === 'edit' ? '正在更新项目...' : '正在添加项目...')
        
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
            
            // 通知数据管理器数据已变更
            await dataManager.onDataChanged()
            
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
    } finally {
        globalLoading.hide()
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

// 验证端口格式
const validatePort = (value: string): boolean => {
    const port = parseInt(value)
    return !isNaN(port) && port >= 1 && port <= 65535
}

// Reset function
const onReset = async () => {
    if (props.mode === 'edit') {
        Object.assign(form, props.initialForm)
    } else {
        // 重置表单到初始状态
        form.project_id = ''
        form.project_name = ''
        form.project_manage_url = ''
        form.project_api_url = ''
        form.api_port = ''
        form.front_port = ''

        // 重新生成默认端口
        await generateDefaultPorts()
        message.success('表单已重置并重新生成默认端口')
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
