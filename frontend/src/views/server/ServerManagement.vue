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
                <n-form-item label="默认路径" path="default_path">
                    <n-input 
                        v-model:value="serverFormData.default_path" 
                        placeholder="请输入默认路径"
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
    ServerOutline,
    LinkOutline,
    CloudUploadOutline
} from '@vicons/ionicons5'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import api from '@/api'
import { getAuthorization } from '@/utils/auth'
import dataManager from '@/utils/dataManager'

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
    default_path?: string
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
    server_password: '',
    default_path: '/adplace'
})

const isServerFormVisible = ref(false)
const isEditMode = ref(false)
const serverFormRef = ref()

const testingServers = ref(new Set<string>())
const updatingServers = ref(new Set<string>())

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
        width: 200,
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
                            type: 'info',
                            class: 'special-table-btn',
                            onClick: () => testStoredServerSSH(row.server_id),
                            loading: testingServers.value.has(row.server_id)
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(LinkOutline) })
                        }),
                        default: () => '测试'
                    }),
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'warning',
                            class: 'special-table-btn',
                            onClick: () => updateAllProjects(row),
                            loading: updatingServers.value.has(row.server_id),
                            disabled: !row.project_list || row.project_list.length === 0
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(CloudUploadOutline) })
                        }),
                        default: () => '全部更新'
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
        title: '默认路径',
        key: 'default_path',
        width: 120,
        ellipsis: {
            tooltip: true
        }
    },
    {
        title: '连接状态',
        key: 'connection_status',
        width: 100,
        render: (row: Server) => {
            const status = (row as any).connection_status
            if (status === 'connected') {
                return h('span', { style: 'color: #18a058;' }, '✓ 已连接')
            } else if (status === 'disconnected') {
                return h('span', { style: 'color: #d03050;' }, '✗ 连接失败')
            } else {
                return h('span', { style: 'color: #909399;' }, '- 未测试')
            }
        }
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
const fetchServers = async (forceRefresh = false) => {
    try {
        // 使用DataManager的缓存机制
        const res = await dataManager.getServerData(forceRefresh)
        if (res && Array.isArray(res)) {
            serverList.value = res
            console.log('ServerManagement: Loaded servers from cache/API:', res.length)
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
    fetchServers(true) // 强制刷新
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
        console.log('Default path value:', requestData.default_path)
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
            server_password: server.server_password,
            default_path: server.default_path || '/adplace'
        }
    } else {
        serverFormData.value = {
            server_id: '',
            server_name: '',
            server_ip: '',
            server_port: '',
            server_user: '',
            server_password: '',
            default_path: '/adplace'
        }
    }
    isServerFormVisible.value = true
}

// 测试已存储服务器的SSH连接
const testStoredServerSSH = async (serverId: string) => {
    testingServers.value.add(serverId)
    
    try {
        const res = await api('test_stored_ssh', {
            server_id: serverId
        })
        
        console.log('Stored SSH test result:', res)
        
        if (res && res.code === 200) {
            message.success(`服务器 ${serverId} SSH连接测试成功`)
        } else {
            message.error(`服务器 ${serverId} SSH连接测试失败: ${res?.msg || '未知错误'}`)
        }
        
        // 刷新服务器列表以显示最新状态
        await fetchServers()
    } catch (error) {
        console.error('Stored SSH test error:', error)
        message.error(`服务器 ${serverId} SSH连接测试异常`)
    } finally {
        testingServers.value.delete(serverId)
    }
}

// 全部更新项目
const updateAllProjects = async (server: Server) => {
    if (!server.project_list || server.project_list.length === 0) {
        message.warning('该服务器下没有项目')
        return
    }

    const projectCount = server.project_list.length
    
    // 确认对话框
    dialog.warning({
        title: '确认全部更新',
        content: `即将更新服务器 "${server.server_name}" 下的 ${projectCount} 个项目。此操作将：

1. 生成所有项目的配置文件
2. 逐个执行项目更新操作
3. 可能需要较长时间完成

是否继续？`,
        positiveText: '确认更新',
        negativeText: '取消',
        onPositiveClick: async () => {
            await executeUpdateAllProjects(server)
        }
    })
}

// 执行全部更新项目
const executeUpdateAllProjects = async (server: Server) => {
    updatingServers.value.add(server.server_id)
    
    try {
        message.loading(`正在更新服务器 "${server.server_name}" 下的所有项目...`, { duration: 0 })
        
        // 1. 生成所有项目的配置JSON
        const extractDomain = (url: string): string => {
            if (!url) return ''
            try {
                const urlObj = new URL(url.startsWith('http') ? url : `https://${url}`)
                return urlObj.hostname
            } catch {
                return url.replace(/^https?:\/\//, '').split('/')[0]
            }
        }

        // 生成所有项目的配置
        const allProjectsConfig: Record<string, any> = {}
        server.project_list.forEach(project => {
            allProjectsConfig[project.project_id] = {
                api_port: project.api_port || '9000',
                web_port: project.front_port || '3000',
                api_domain: extractDomain(project.project_api_url)
            }
        })

        const configJson = JSON.stringify(allProjectsConfig, null, 2)
        
        console.log('生成的所有项目配置:', {
            serverId: server.server_id,
            projectCount: server.project_list.length,
            config: allProjectsConfig
        })

        // 2. 上传配置文件
        const configResult = await api('upload_project_config', {
            server_data_json: JSON.stringify(server),
            project_config_json: configJson,
            authorization: getAuthorization()
        })

        if (configResult.code !== 200) {
            throw new Error(`配置文件上传失败: ${configResult.msg}`)
        }

        message.success('配置文件上传成功，开始更新项目...')

        // 3. 循环更新每个项目
        let successCount = 0
        let failCount = 0
        const results: Array<{project: any, success: boolean, error?: string}> = []

        for (const project of server.project_list) {
            try {
                console.log(`正在更新项目: ${project.project_name} (${project.project_id})`)
                
                const updateResult = await api('project_update_with_data', {
                    server_id: server.server_id,
                    project_id: project.project_id,
                    server_data_json: JSON.stringify(server)
                })

                if (updateResult.code === 200) {
                    successCount++
                    results.push({ project, success: true })
                    console.log(`项目 ${project.project_id} 更新成功`)
                } else {
                    failCount++
                    results.push({ project, success: false, error: updateResult.msg })
                    console.error(`项目 ${project.project_id} 更新失败:`, updateResult.msg)
                }
            } catch (error) {
                failCount++
                results.push({ project, success: false, error: (error as Error).message })
                console.error(`项目 ${project.project_id} 更新异常:`, error)
            }
        }

        message.destroyAll()

        // 4. 显示结果
        if (failCount === 0) {
            message.success(`全部更新完成！成功更新 ${successCount} 个项目`)
        } else if (successCount === 0) {
            message.error(`全部更新失败！${failCount} 个项目更新失败`)
        } else {
            message.warning(`部分更新完成：成功 ${successCount} 个，失败 ${failCount} 个`)
        }

        // 显示详细结果
        console.log('更新结果详情:', results)

    } catch (error) {
        console.error('Update all projects error:', error)
        message.destroyAll()
        message.error(`全部更新失败：${(error as Error).message}`)
    } finally {
        updatingServers.value.delete(server.server_id)
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