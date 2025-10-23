<template>
    <div>
        <n-card title="服务器信息">
            <template #header-extra>
                <n-space>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="primary" @click="handleEdit">
                                <template #icon>
                                    <n-icon>
                                        <CreateOutline v-if="!editMode" />
                                        <CloseOutline v-else />
                                    </n-icon>
                                </template>
                                {{ !editMode ? '编辑' : '取消' }}
                            </n-button>
                        </template>
                        {{ !editMode ? '编辑' : '取消' }}
                    </n-tooltip>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="error" @click="handleDelete">
                                <template #icon>
                                    <n-icon><TrashOutline /></n-icon>
                                </template>
                                删除
                            </n-button>
                        </template>
                        删除
                    </n-tooltip>
                </n-space>
            </template>

            <!-- 编辑表单 -->
            <n-form v-if="editMode" :model="serverFormData" label-placement="left" label-width="120px" ref="serverFormRef">
                <n-form-item label="服务器ID" path="server_id" required>
                    <n-input v-model:value="serverFormData.server_id" placeholder="请输入服务器ID" />
                </n-form-item>
                <n-form-item label="服务器名称" path="server_name" required>
                    <n-input v-model:value="serverFormData.server_name" placeholder="请输入服务器名称" />
                </n-form-item>
                <n-form-item label="服务器IP" path="server_ip">
                    <n-input v-model:value="serverFormData.server_ip" placeholder="请输入服务器IP地址" />
                </n-form-item>
                <n-form-item label="端口" path="server_port">
                    <n-input v-model:value="serverFormData.server_port" placeholder="请输入端口号" />
                </n-form-item>
                <n-form-item label="用户名" path="server_user">
                    <n-input v-model:value="serverFormData.server_user" placeholder="请输入用户名" />
                </n-form-item>
                <n-form-item label="密码" path="server_password">
                    <n-input 
                        v-model:value="serverFormData.server_password" 
                        type="password" 
                        placeholder="请输入密码"
                        show-password-on="click"
                    />
                </n-form-item>
                <n-form-item>
                    <n-space>
                        <n-button type="primary" @click="submitEdit">
                            <template #icon>
                                <n-icon><CheckmarkOutline /></n-icon>
                            </template>
                            保存
                        </n-button>
                        <n-button @click="editMode = false">
                            <template #icon>
                                <n-icon><CloseOutline /></n-icon>
                            </template>
                            取消
                        </n-button>
                    </n-space>
                </n-form-item>
            </n-form>

            <!-- 显示信息 -->
            <n-descriptions v-if="!editMode" :column="1" bordered>
                <n-descriptions-item label="服务器ID">
                    {{ serverInfo.server_id }}
                </n-descriptions-item>
                <n-descriptions-item label="服务器名称">
                    {{ serverInfo.server_name }}
                </n-descriptions-item>
                <n-descriptions-item label="服务器IP">
                    <n-text type="info">{{ serverInfo.server_ip || '未设置' }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="端口">
                    <n-text type="info">{{ serverInfo.server_port || '未设置' }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="用户名">
                    <n-text type="info">{{ serverInfo.server_user || '未设置' }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="项目数量">
                    <n-badge :value="serverInfo.project_list?.length || 0" type="success">
                        <n-text>{{ serverInfo.project_list?.length || 0 }} 个项目</n-text>
                    </n-badge>
                </n-descriptions-item>
            </n-descriptions>

            <!-- 项目列表 -->
            <n-divider v-if="!editMode && serverInfo.project_list?.length > 0" />
            <n-card v-if="!editMode && serverInfo.project_list?.length > 0" title="项目列表" size="small">
                <n-list>
                    <n-list-item v-for="project in serverInfo.project_list" :key="project.project_id">
                        <n-thing>
                            <template #header>
                                {{ project.project_name }}
                            </template>
                            <template #description>
                                <n-space vertical size="small">
                                    <n-text depth="3">ID: {{ project.project_id }}</n-text>
                                    <n-text depth="3">API: {{ project.project_api_url }}</n-text>
                                    <n-text depth="3">管理: {{ project.project_manage_url }}</n-text>
                                </n-space>
                            </template>
                        </n-thing>
                    </n-list-item>
                </n-list>
            </n-card>
        </n-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, defineProps, onMounted, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import { CreateOutline, CloseOutline, TrashOutline, CheckmarkOutline } from '@vicons/ionicons5'
import api from '@/api'
import dataManager from '@/utils/dataManager'

const sidebar = useSidebarStore()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()

// 注入全局 loading
const globalLoading = inject('globalLoading') as any

// 定义接受 serverId 的 props
const props = defineProps({
    serverId: {
        type: String,
        required: true,
    },
})

interface Project {
    project_id: string
    project_name: string
    project_api_url: string
    project_manage_url: string
}

interface ServerInfo {
    server_id?: string
    server_name?: string
    server_ip?: string
    server_port?: string
    server_user?: string
    server_password?: string
    project_list?: Project[]
}

const serverInfo = ref<ServerInfo>({})
const serverFormData = ref<ServerInfo>({})
const editMode = ref(false)
const serverFormRef = ref()

// 编辑按钮点击事件
const handleEdit = () => {
    if (!editMode.value) {
        // 进入编辑模式，复制当前数据到表单
        serverFormData.value = { ...serverInfo.value }
    }
    editMode.value = !editMode.value
}

// 获取服务器信息
const getServerInfo = async () => {
    try {
        // 优先从缓存获取服务器信息
        const server = await dataManager.getServerById(props.serverId)
        if (server) {
            console.log('Server info from cache:', server)
            serverInfo.value = server
        } else {
            // 如果缓存中没有，则调用API
            const res = await api('server_info', {
                serverId: props.serverId,
            })
            console.log('Server info from API:', res)
            // ServerInfo API returns the server data directly
            if (res && typeof res === 'object' && !res.code) {
                serverInfo.value = res
            } else if (res && res.data) {
                serverInfo.value = res.data
            } else {
                serverInfo.value = res || {}
            }
        }
    } catch (error) {
        console.error('Failed to get server info:', error)
        message.error('获取服务器信息失败')
    }
}

// 提交编辑
const submitEdit = async () => {
    globalLoading.show('正在更新服务器信息...')
    try {
        const res = await api('server_update', {
            old_server_id: props.serverId,  // 原服务器ID
            server_id: serverFormData.value.server_id,  // 新服务器ID
            server_name: serverFormData.value.server_name,
            server_ip: serverFormData.value.server_ip,
            server_port: serverFormData.value.server_port,
            server_user: serverFormData.value.server_user,
            server_password: serverFormData.value.server_password,
        })
        
        if (res.code === 200 || res.success) {
            message.success('服务器信息更新成功')
            editMode.value = false
            getServerInfo() // 重新获取信息
            await reloadMenus() // 刷新菜单
            sidebar.setboolroute(true)
        } else {
            message.error(res.message || '服务器信息更新失败')
        }
    } catch (error) {
        console.error('Failed to update server info:', error)
        message.error('服务器信息更新失败')
    } finally {
        globalLoading.hide()
    }
}

// 删除按钮点击事件
const handleDelete = () => {
    const projectCount = serverInfo.value.project_list?.length || 0
    const content = projectCount > 0 
        ? `该服务器下有 ${projectCount} 个项目，删除服务器将同时删除所有项目。此操作不可恢复，是否继续？`
        : '确定要删除该服务器吗？此操作不可恢复。'
    
    dialog.warning({
        title: '确认删除',
        content,
        positiveText: '确定删除',
        negativeText: '取消',
        onPositiveClick: async () => {
            globalLoading.show('正在删除服务器...')
            try {
                const res = await api('server_delete', {
                    server_id: props.serverId,
                })
                
                if (res.code === 200 || res.success) {
                    message.success('服务器删除成功')
                    await reloadMenus()
                    sidebar.setboolroute(true)
                    router.push('/welcome')
                } else {
                    message.error(res.message || '服务器删除失败')
                }
            } catch (error) {
                console.error('Failed to delete server:', error)
                message.error('服务器删除失败')
            } finally {
                globalLoading.hide()
            }
        }
    })
}

// 初始化
onMounted(() => {
    getServerInfo()
})
</script>

<style scoped>
:deep(.n-descriptions-table-wrapper .n-descriptions-table th) {
    font-size: var(--table-font-size, 12px);
    padding: 4px 8px;
}

:deep(.n-descriptions-table-wrapper .n-descriptions-table td) {
    font-size: var(--table-font-size, 12px);
    padding: 4px 8px;
}

:deep(.n-list .n-list-item) {
    padding: 8px 0;
}

:deep(.n-thing .n-thing-header) {
    font-size: var(--menu-font-size, 14px);
    font-weight: 600;
}

:deep(.n-thing .n-thing-main) {
    font-size: var(--table-font-size, 12px);
}
</style>