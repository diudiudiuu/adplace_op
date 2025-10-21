<template>
    <div>
        <n-card title="服务器管理">
            <template #header-extra>
                <n-space>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="primary" @click="openServerForm(false)">
                                <template #icon>
                                    <n-icon><AddCircleOutline /></n-icon>
                                </template>
                                添加服务器
                            </n-button>
                        </template>
                        添加服务器
                    </n-tooltip>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="info" @click="refreshServers">
                                <template #icon>
                                    <n-icon><RefreshOutline /></n-icon>
                                </template>
                                刷新
                            </n-button>
                        </template>
                        刷新
                    </n-tooltip>
                </n-space>
            </template>

            <n-data-table
                :columns="serverColumns"
                :data="serverList"
                :pagination="pagination"
                striped
                class="special-table"
            />
        </n-card>

        <!-- 服务器表单弹窗 -->
        <n-modal v-model:show="isServerFormVisible" preset="dialog" :title="isEditMode ? '编辑服务器' : '添加服务器'">
            <n-form :model="serverFormData" label-placement="left" label-width="120px" ref="serverFormRef">
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
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="isServerFormVisible = false">
                        <template #icon>
                            <n-icon><CloseOutline /></n-icon>
                        </template>
                        取消
                    </n-button>
                    <n-button type="primary" @click="submitServerForm">
                        <template #icon>
                            <n-icon><CheckmarkOutline /></n-icon>
                        </template>
                        {{ isEditMode ? '更新' : '添加' }}
                    </n-button>
                </n-space>
            </template>
        </n-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed, h } from 'vue'
import { useMessage, useDialog, NButton, NIcon, NSpace, NTooltip, NBadge } from 'naive-ui'
import { 
    AddCircleOutline, 
    RefreshOutline, 
    CreateOutline, 
    TrashOutline, 
    CloseOutline, 
    CheckmarkOutline,
    ServerOutline
} from '@vicons/ionicons5'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import api from '@/api'

interface Project {
    project_id: string
    project_name: string
    project_api_url: string
    project_manage_url: string
}

interface Server {
    server_id: string
    server_name: string
    server_ip: string
    server_port: string
    server_user: string
    server_password: string
    project_list: Project[]
}

const message = useMessage()
const dialog = useDialog()
const sidebar = useSidebarStore()

const serverList = ref<Server[]>([])
const serverFormData = ref({
    server_id: '',
    server_name: '',
    server_ip: '',
    server_port: '',
    server_user: '',
    server_password: ''
})

const isServerFormVisible = ref(false)
const isEditMode = ref(false)
const serverFormRef = ref()

// 分页配置
const pagination = {
    pageSize: 10,
    showSizePicker: true,
    pageSizes: [10, 20, 50],
    showQuickJumper: true
}

// 服务器表格列配置
const serverColumns = computed(() => [
    {
        title: '操作',
        key: 'actions',
        width: 120,
        render: (row: Server) => {
            return h(NSpace, { size: 'small' }, {
                default: () => [
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'primary',
                            class: 'special-table-btn',
                            onClick: () => openServerForm(true, row)
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(CreateOutline) })
                        }),
                        default: () => '编辑'
                    }),
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'error',
                            class: 'special-table-btn',
                            onClick: () => confirmDeleteServer(row)
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(TrashOutline) })
                        }),
                        default: () => '删除'
                    })
                ]
            })
        }
    },
    {
        title: '服务器ID',
        key: 'server_id',
        width: 100
    },
    {
        title: '服务器名称',
        key: 'server_name',
        ellipsis: {
            tooltip: true
        }
    },
    {
        title: '服务器IP',
        key: 'server_ip',
        ellipsis: {
            tooltip: true
        }
    },
    {
        title: '端口',
        key: 'server_port',
        width: 80
    },
    {
        title: '用户名',
        key: 'server_user',
        width: 100
    },
    {
        title: '项目数量',
        key: 'project_count',
        width: 100,
        render: (row: Server) => {
            const count = row.project_list?.length || 0
            return h(NBadge, {
                value: count,
                type: count > 0 ? 'success' : 'default'
            }, {
                default: () => h(NIcon, { size: 16 }, { default: () => h(ServerOutline) })
            })
        }
    }
])

// 获取服务器列表
const fetchServers = async () => {
    try {
        const res = await api('server_list', {})
        if (res && Array.isArray(res)) {
            serverList.value = res
        } else {
            message.error('获取服务器列表失败')
        }
    } catch (error) {
        console.error('Failed to fetch servers:', error)
        message.error('获取服务器列表失败')
    }
}

// 刷新服务器列表
const refreshServers = () => {
    fetchServers()
}

// 打开服务器表单
const openServerForm = (editMode: boolean, server?: Server) => {
    isEditMode.value = editMode
    if (editMode && server) {
        serverFormData.value = {
            server_id: server.server_id,
            server_name: server.server_name,
            server_ip: server.server_ip,
            server_port: server.server_port,
            server_user: server.server_user,
            server_password: server.server_password
        }
    } else {
        serverFormData.value = {
            server_id: '',
            server_name: '',
            server_ip: '',
            server_port: '',
            server_user: '',
            server_password: ''
        }
    }
    isServerFormVisible.value = true
}

// 提交服务器表单
const submitServerForm = async () => {
    try {
        // 验证必填字段
        if (!serverFormData.value.server_id || !serverFormData.value.server_name) {
            message.error('服务器ID和服务器名称为必填项')
            return
        }
        
        const action = isEditMode.value ? 'server_update' : 'server_add'
        let requestData = { ...serverFormData.value }
        
        // 如果是编辑模式，需要传递原始服务器ID
        if (isEditMode.value) {
            requestData.old_server_id = serverFormData.value.server_id // 编辑时原ID就是当前ID
        }
        
        console.log('Submitting server form:', action, requestData)
        const res = await api(action, requestData)
        console.log('Server form response:', res)
        
        if (res && (res.code === 200 || res.success)) {
            message.success(isEditMode.value ? '服务器更新成功' : '服务器添加成功')
            isServerFormVisible.value = false
            await fetchServers()
            // 刷新左侧菜单
            await reloadMenus()
            sidebar.setboolroute(true)
        } else {
            const errorMsg = res?.msg || res?.message || (isEditMode.value ? '服务器更新失败' : '服务器添加失败')
            message.error(errorMsg)
        }
    } catch (error) {
        console.error('Server form submission error:', error)
        message.error(isEditMode.value ? '服务器更新失败' : '服务器添加失败')
    }
}

// 确认删除服务器
const confirmDeleteServer = (server: Server) => {
    const projectCount = server.project_list?.length || 0
    const content = projectCount > 0 
        ? `该服务器下有 ${projectCount} 个项目，删除服务器将同时删除所有项目。此操作不可恢复，是否继续？`
        : '确定要删除该服务器吗？此操作不可恢复。'
    
    dialog.warning({
        title: '确认删除',
        content,
        positiveText: '确定删除',
        negativeText: '取消',
        onPositiveClick: () => deleteServer(server)
    })
}

// 删除服务器
const deleteServer = async (server: Server) => {
    try {
        const res = await api('server_delete', {
            server_id: server.server_id
        })
        
        if (res.code === 200 || res.success) {
            message.success('服务器删除成功')
            await fetchServers()
            // 刷新左侧菜单
            await reloadMenus()
            sidebar.setboolroute(true)
        } else {
            message.error(res.message || '服务器删除失败')
        }
    } catch (error) {
        console.error('Server deletion error:', error)
        message.error('服务器删除失败')
    }
}

// 初始化数据
onMounted(() => {
    fetchServers()
})
</script>

<style scoped>
/* 特殊表格的24px按钮样式 */
:deep(.special-table .special-table-btn) {
    width: 24px !important;
    height: 24px !important;
    min-width: 24px !important;
    max-width: 24px !important;
    padding: 0 !important;
    border-radius: 4px !important;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
}

:deep(.special-table .special-table-btn .n-button__content) {
    padding: 0 !important;
    margin: 0 !important;
    width: 24px !important;
    height: 24px !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
}

:deep(.special-table .special-table-btn .n-icon) {
    font-size: 16px !important;
    width: 16px !important;
    height: 16px !important;
    margin: 0 !important;
}
</style>