<template>
    <div>
        <n-card title="项目信息">
            <template #header-extra>
                <n-space>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="primary" @click="handleEdit">
                                <template #icon>
                                    <n-icon>
                                        <CreateOutline v-if="!eidtmode" />
                                        <CloseOutline v-else />
                                    </n-icon>
                                </template>
                                {{ !eidtmode ? '编辑' : '取消' }}
                            </n-button>
                        </template>
                        {{ !eidtmode ? '编辑' : '取消' }}
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

            <Dform v-if="eidtmode" mode="edit" :serverId="serverId" :initialForm="projectInfo"
                @editSuccess="updateHandle" />

            <n-descriptions v-if="!eidtmode" :column="1" bordered>
                <n-descriptions-item label="客户ID">
                    {{ projectInfo.project_id }}
                </n-descriptions-item>
                <n-descriptions-item label="客户名称">
                    {{ projectInfo.project_name }}
                </n-descriptions-item>
                <n-descriptions-item label="客户管理地址">
                    <n-text type="info">{{ projectInfo.project_manage_url }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="客户API地址">
                    <n-text type="info">{{ projectInfo.project_api_url }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="API端口">
                    <n-text type="success">{{ projectInfo.api_port || '8080' }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="前端端口">
                    <n-text type="success">{{ projectInfo.front_port || '3000' }}</n-text>
                </n-descriptions-item>
            </n-descriptions>
        </n-card>

        <!-- Cloudflare DNS 配置卡片 - 只在非编辑状态下显示 -->
        <n-card v-if="!eidtmode" title="Cloudflare DNS 配置" style="margin-top: 16px;">
            <template #header-extra>
                <n-space>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="info" @click="showCloudflareConfig" size="small">
                                <template #icon>
                                    <n-icon><SettingsOutline /></n-icon>
                                </template>
                                配置
                            </n-button>
                        </template>
                        配置 Cloudflare API
                    </n-tooltip>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="primary" @click="batchConfigureDNS" :loading="dnsLoading">
                                <template #icon>
                                    <n-icon><CloudOutline /></n-icon>
                                </template>
                                批量配置
                            </n-button>
                        </template>
                        批量配置所有 DNS 记录
                    </n-tooltip>
                </n-space>
            </template>

            <n-space vertical>
                <n-alert type="info" :show-icon="false">
                    <template #header>
                        <n-icon><InformationCircleOutline /></n-icon>
                        DNS 配置说明
                    </template>
                    将为以下域名配置 Cloudflare DNS 记录，所有记录都会开启代理（黄色小云朵）
                    <br>• 管理端域名：自动配置 CNAME 到 Pages 项目，并添加为 Pages 自定义域名
                    <br>• API 端域名：配置 A 记录直接指向服务器 IP
                </n-alert>

                <!-- DNS 记录表格 -->
                <div v-if="dnsStatus.length > 0">
                    <n-data-table
                        :columns="dnsColumns"
                        :data="dnsStatus"
                        :pagination="false"
                        striped
                        size="small"
                    />
                </div>
            </n-space>
        </n-card>

        <!-- Cloudflare 配置对话框 -->
        <n-modal v-model:show="showConfigModal" preset="dialog" title="Cloudflare 配置" style="width: 500px;">
            <n-form :model="cloudflareConfig" label-placement="left" label-width="120">
                <n-form-item label="API Token" required>
                    <n-input 
                        v-model:value="cloudflareConfig.apiToken" 
                        type="password" 
                        placeholder="请输入 Cloudflare API Token"
                        show-password-on="click"
                    />
                </n-form-item>
                <n-form-item label="Zone ID" required>
                    <n-input 
                        v-model:value="cloudflareConfig.zoneId" 
                        placeholder="请输入域名的 Zone ID"
                    />
                </n-form-item>
                <n-alert type="info" style="margin-top: 16px;">
                    <template #header>配置说明</template>
                    <ul style="margin: 8px 0; padding-left: 20px;">
                        <li>API Token 需要有 Zone:Edit 权限</li>
                        <li>Zone ID 可在 Cloudflare 域名概览页面找到</li>
                        <li>配置信息将保存在本地浏览器中</li>
                        <li>DNS 操作通过后端安全调用 Cloudflare API</li>
                    </ul>
                </n-alert>
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="showConfigModal = false">取消</n-button>
                    <n-button type="primary" @click="saveCloudflareConfig">保存配置</n-button>
                </n-space>
            </template>
        </n-modal>
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NButton, NIcon, NSpace, NTooltip, NTag, NText } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import dataManager from '@/utils/dataManager'

import { CreateOutline, CloseOutline, TrashOutline, CloudOutline, InformationCircleOutline, PlayOutline, RefreshOutline, TrashBinOutline, SettingsOutline } from '@vicons/ionicons5'
import Dform from './form.vue'
import api from '@/api'

const sidebar = useSidebarStore()
const route = useRouter()
const message = useMessage()
const dialog = useDialog()

// 定义接受 projectId 的 props
const props = defineProps({
    serverId: {
        type: String,
        required: true,
    },
    projectId: {
        type: String,
        required: true,
    },
})

interface ProjectInfo {
    project_id?: string;
    project_name?: string;
    project_manage_url?: string;
    project_api_url?: string;
    api_port?: string;
    front_port?: string;
}

const projectInfo = ref<ProjectInfo>({})
const dnsLoading = ref(false)
const dnsStatus = ref<any[]>([])
const showConfigModal = ref(false)
const cloudflareConfig = ref({
    apiToken: localStorage.getItem('cloudflare_api_token') || '',
    zoneId: localStorage.getItem('cloudflare_zone_id') || ''
})

const eidtmode = ref(false)

// 从 URL 中提取域名
const extractDomain = (url: string): string => {
    try {
        const urlObj = new URL(url.startsWith('http') ? url : `https://${url}`)
        return urlObj.hostname
    } catch (error) {
        // 如果不是完整 URL，假设它就是域名
        return url.replace(/^https?:\/\//, '').split('/')[0]
    }
}

// DNS 表格列配置
const dnsColumns = computed(() => [
    {
        title: '域名',
        key: 'name',
        width: 200,
        render: (row: any) => {
            return h(NText, { strong: true }, { default: () => row.name })
        }
    },
    {
        title: '记录类型',
        key: 'type',
        width: 100,
        render: (row: any) => {
            return h(NTag, {
                type: row.type === 'CNAME' ? 'success' : 'warning',
                size: 'small'
            }, { default: () => row.type })
        }
    },
    {
        title: '记录值',
        key: 'content',
        width: 180,
        render: (row: any) => {
            return h(NText, { depth: 3 }, { default: () => row.content })
        }
    },
    {
        title: '代理状态',
        key: 'proxied',
        width: 100,
        render: (row: any) => {
            return h(NTag, {
                type: row.proxied ? 'info' : 'default',
                size: 'small'
            }, { default: () => row.proxied ? '🟡 已代理' : '⚪ 未代理' })
        }
    },
    {
        title: '状态',
        key: 'status',
        width: 100,
        render: (row: any) => {
            const statusMap = {
                active: { type: 'success', text: '✓ 已配置' },
                pending: { type: 'warning', text: '○ 待配置' },
                error: { type: 'error', text: '✗ 错误' }
            }
            const status = statusMap[row.status as keyof typeof statusMap] || statusMap.pending
            return h(NTag, {
                type: status.type,
                size: 'small'
            }, { default: () => status.text })
        }
    },
    {
        title: '操作',
        key: 'actions',
        width: 200,
        render: (row: any) => {
            return h(NSpace, { size: 'small' }, {
                default: () => [
                    // 配置/更新按钮
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: row.status === 'active' ? 'info' : 'primary',
                            loading: row.loading,
                            onClick: () => configureSingleDNS(row)
                        }, {
                            icon: () => h(NIcon, {}, { default: () => h(row.status === 'active' ? RefreshOutline : PlayOutline) }),
                            default: () => row.status === 'active' ? '更新' : '配置'
                        }),
                        default: () => row.status === 'active' ? '更新 DNS 记录' : '配置 DNS 记录'
                    }),
                    // 删除按钮（只有已配置的记录才显示）
                    row.status === 'active' ? h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'error',
                            loading: row.loading,
                            onClick: () => deleteSingleDNS(row)
                        }, {
                            icon: () => h(NIcon, {}, { default: () => h(TrashBinOutline) })
                        }),
                        default: () => '删除 DNS 记录'
                    }) : null
                ].filter(Boolean)
            })
        }
    }
])
// 编辑按钮点击事件
const handleEdit = () => {
    // 切换到编辑模式
    eidtmode.value = !eidtmode.value
}
// 获取项目信息
const getProjectInfo = async () => {
    try {
        // 优先从数据管理器获取项目信息
        const project = await dataManager.getProjectById(props.projectId)
        if (project) {
            console.log('Project info from cache:', project)
            projectInfo.value = project
        } else {
            // 如果缓存中没有，则调用 API
            const res = await api('project_info', {
                projectId: props.projectId,
            })
            console.log('Project info from API:', res)
            projectInfo.value = res
        }
    } catch (error) {
        console.error('Failed to get project info:', error)
        message.error('获取项目信息失败')
    }
}

getProjectInfo().then(() => {
    // 项目信息加载完成后检查 DNS 状态
    checkDNSStatus()
})

const updateHandle = () => {
    eidtmode.value = false
    getProjectInfo().then(() => {
        checkDNSStatus()
    })
}

// 批量配置 Cloudflare DNS
const batchConfigureDNS = async () => {
    if (!projectInfo.value.project_manage_url || !projectInfo.value.project_api_url) {
        message.error('项目管理地址或API地址不能为空')
        return
    }

    // 检查 Cloudflare 配置
    const apiToken = localStorage.getItem('cloudflare_api_token')
    const zoneId = localStorage.getItem('cloudflare_zone_id')
    
    if (!apiToken || !zoneId) {
        dialog.warning({
            title: 'Cloudflare 配置',
            content: '请先配置 Cloudflare API Token 和 Zone ID',
            positiveText: '去配置',
            negativeText: '取消',
            onPositiveClick: () => {
                showCloudflareConfig()
            }
        })
        return
    }

    dnsLoading.value = true
    
    try {
        // 获取服务器信息以获取IP地址
        const serverInfo = await dataManager.getServerById(props.serverId)
        if (!serverInfo) {
            message.error('无法获取服务器信息')
            return
        }

        const serverIP = serverInfo.server_ip
        if (!serverIP) {
            message.error('服务器IP地址不能为空')
            return
        }

        // 解析域名
        const manageDomain = extractDomain(projectInfo.value.project_manage_url)
        const apiDomain = extractDomain(projectInfo.value.project_api_url)

        message.loading('正在批量配置 Cloudflare DNS 记录和 Pages 自定义域名...', { duration: 0 })

        // 1. 首先配置 Pages 自定义域名
        let pagesConfigSuccess = false
        try {
            console.log(`正在为 Pages 项目 'adswds' 添加自定义域名: ${manageDomain}`)
            const pagesResult = await api('cloudflare_pages_add_domain', {
                api_token: apiToken,
                zone_id: zoneId,
                project_name: 'adswds', // Pages 项目名称
                domain: manageDomain
            })
            
            console.log('Pages API 响应:', pagesResult)
            
            if (pagesResult.code === 200) {
                console.log('Pages 自定义域名配置成功:', pagesResult.data)
                pagesConfigSuccess = true
            } else {
                console.error('Pages 自定义域名配置失败:', pagesResult.msg)
                if (pagesResult.msg && pagesResult.msg.includes('Authentication error')) {
                    message.warning('Pages 自定义域名配置失败：API Token 缺少 Cloudflare Pages:Edit 权限，请手动在 Pages 控制台添加自定义域名')
                } else {
                    message.warning(`Pages 自定义域名配置失败: ${pagesResult.msg}`)
                }
            }
        } catch (pagesError) {
            console.error('Pages 自定义域名配置出错:', pagesError)
            message.warning(`Pages 自定义域名配置出错，将继续配置 DNS 记录`)
        }

        // 2. 配置 DNS 记录 - 管理端CNAME到Pages，API端A记录到服务器
        const records = [
            {
                name: manageDomain,
                type: 'CNAME',
                content: 'adswds.pages.dev',
                proxied: true
            },
            {
                name: apiDomain,
                type: 'A',
                content: serverIP,
                proxied: true
            }
        ]

        // 调用后端批量配置 DNS API
        const result = await api('cloudflare_batch_configure', {
            api_token: apiToken,
            zone_id: zoneId,
            records_json: JSON.stringify(records)
        })

        message.destroyAll()

        if (result.code === 200) {
            const results = result.data || []
            
            // 更新 DNS 状态
            dnsStatus.value = dnsStatus.value.map((status: any) => {
                const matchResult = results.find((r: any) => 
                    r.record?.name === status.name || r.name === status.name
                )
                
                if (matchResult) {
                    if (matchResult.error) {
                        return {
                            ...status,
                            status: 'error',
                            error: matchResult.error,
                            loading: false
                        }
                    } else {
                        return {
                            ...status,
                            status: 'active',
                            content: matchResult.record.content,
                            proxied: matchResult.record.proxied,
                            recordId: matchResult.record.id,
                            action: matchResult.action,
                            loading: false
                        }
                    }
                }
                return status
            })

            const successCount = results.filter((r: any) => !r.error).length
            const errorCount = results.filter((r: any) => r.error).length
            
            if (errorCount === 0) {
                message.success(`Cloudflare DNS 批量配置完成！成功配置 ${successCount} 条记录`)
            } else {
                message.warning(`部分配置完成：成功 ${successCount} 条，失败 ${errorCount} 条`)
            }
        } else {
            message.error(result.msg || 'DNS 配置失败')
        }
        
    } catch (error) {
        console.error('DNS batch configuration error:', error)
        message.destroyAll()
        message.error('DNS 配置失败：' + (error as Error).message)
    } finally {
        dnsLoading.value = false
    }
}

// 显示 Cloudflare 配置对话框
const showCloudflareConfig = () => {
    showConfigModal.value = true
}

// 保存 Cloudflare 配置
const saveCloudflareConfig = () => {
    if (!cloudflareConfig.value.apiToken || !cloudflareConfig.value.zoneId) {
        message.error('请填写完整的配置信息')
        return
    }

    // 保存到本地存储
    localStorage.setItem('cloudflare_api_token', cloudflareConfig.value.apiToken)
    localStorage.setItem('cloudflare_zone_id', cloudflareConfig.value.zoneId)
    
    showConfigModal.value = false
    message.success('Cloudflare 配置已保存')
    
    // 重新检查 DNS 状态
    checkDNSStatus()
}

// 检查 DNS 记录状态
const checkDNSStatus = async () => {
    if (!projectInfo.value.project_manage_url || !projectInfo.value.project_api_url) {
        return
    }

    const apiToken = localStorage.getItem('cloudflare_api_token')
    const zoneId = localStorage.getItem('cloudflare_zone_id')
    
    const manageDomain = extractDomain(projectInfo.value.project_manage_url)
    const apiDomain = extractDomain(projectInfo.value.project_api_url)
    
    if (!apiToken || !zoneId) {
        // 如果没有配置，显示待配置状态
        dnsStatus.value = [
            {
                name: manageDomain,
                type: 'CNAME',
                content: 'adswds.pages.dev',
                status: 'pending',
                proxied: false,
                loading: false
            },
            {
                name: apiDomain,
                type: 'A',
                content: '待配置',
                status: 'pending',
                proxied: false,
                loading: false
            }
        ]
        return
    }

    try {
        // 查询现有 DNS 记录 - 管理端CNAME，API端A记录
        const [manageResult, apiResult] = await Promise.all([
            api('cloudflare_get_dns', {
                api_token: apiToken,
                zone_id: zoneId,
                name: manageDomain,
                type: 'CNAME'
            }),
            api('cloudflare_get_dns', {
                api_token: apiToken,
                zone_id: zoneId,
                name: apiDomain,
                type: 'A'
            })
        ])
        
        const status = []
        
        // 检查管理端域名 CNAME 记录
        const manageRecords = manageResult.code === 200 ? manageResult.data : []
        const manageRecord = manageRecords.find((r: any) => r.name === manageDomain && r.type === 'CNAME')
        status.push({
            name: manageDomain,
            type: 'CNAME',
            content: manageRecord?.content || 'adswds.pages.dev',
            status: manageRecord ? 'active' : 'pending',
            proxied: manageRecord?.proxied || false,
            recordId: manageRecord?.id || null,
            loading: false
        })
        
        // 检查 API 域名 A 记录
        const apiRecords = apiResult.code === 200 ? apiResult.data : []
        const apiRecord = apiRecords.find((r: any) => r.name === apiDomain && r.type === 'A')
        status.push({
            name: apiDomain,
            type: 'A',
            content: apiRecord?.content || '待配置',
            status: apiRecord ? 'active' : 'pending',
            proxied: apiRecord?.proxied || false,
            recordId: apiRecord?.id || null,
            loading: false
        })
        
        dnsStatus.value = status
        
    } catch (error) {
        console.error('Failed to check DNS status:', error)
        // 出错时显示错误状态
        dnsStatus.value = [
            {
                name: manageDomain,
                type: 'CNAME',
                content: 'adswds.pages.dev',
                status: 'error',
                proxied: false,
                loading: false,
                error: (error as Error).message
            },
            {
                name: apiDomain,
                type: 'A',
                content: '查询失败',
                status: 'error',
                proxied: false,
                loading: false,
                error: (error as Error).message
            }
        ]
    }
}



// 配置单个 DNS 记录
const configureSingleDNS = async (record: any) => {
    const apiToken = localStorage.getItem('cloudflare_api_token')
    const zoneId = localStorage.getItem('cloudflare_zone_id')
    
    if (!apiToken || !zoneId) {
        dialog.warning({
            title: 'Cloudflare 配置',
            content: '请先配置 Cloudflare API Token 和 Zone ID',
            positiveText: '去配置',
            negativeText: '取消',
            onPositiveClick: () => {
                showCloudflareConfig()
            }
        })
        return
    }

    // 设置单个记录的加载状态
    record.loading = true
    
    try {
        let content = record.content
        
        // 如果是 A 记录且需要配置，获取服务器 IP
        if (record.type === 'A' && (record.content === '待配置' || record.content === '查询失败')) {
            const serverInfo = await dataManager.getServerById(props.serverId)
            if (!serverInfo || !serverInfo.server_ip) {
                message.error('无法获取服务器IP地址')
                return
            }
            content = serverInfo.server_ip
        }

        message.loading(`正在配置 ${record.name} 的 ${record.type} 记录...`, { duration: 0 })

        // 如果是 CNAME 记录，先配置 Pages 自定义域名
        if (record.type === 'CNAME') {
            try {
                console.log(`正在为 Pages 项目 'adswds' 添加自定义域名: ${record.name}`)
                const pagesResult = await api('cloudflare_pages_add_domain', {
                    api_token: apiToken,
                    zone_id: zoneId,
                    project_name: 'adswds', // Pages 项目名称
                    domain: record.name
                })
                
                console.log('Pages API 响应:', pagesResult)
                
                if (pagesResult.code === 200) {
                    console.log('Pages 自定义域名配置成功:', pagesResult.data)
                    message.success(`Pages 自定义域名 ${record.name} 配置成功`)
                } else {
                    console.error('Pages 自定义域名配置失败:', pagesResult.msg)
                    if (pagesResult.msg && pagesResult.msg.includes('Authentication error')) {
                        message.warning('Pages 自定义域名配置失败：API Token 缺少权限，请手动添加')
                    } else {
                        message.warning(`Pages 自定义域名配置失败: ${pagesResult.msg}`)
                    }
                }
            } catch (pagesError) {
                console.error('Pages 自定义域名配置出错:', pagesError)
                message.warning('Pages 自定义域名配置出错，将继续配置 DNS 记录')
            }
        }

        // 调用后端 API 配置单个记录
        const result = await api('cloudflare_configure_dns', {
            api_token: apiToken,
            zone_id: zoneId,
            name: record.name,
            type: record.type,
            content: content,
            proxied: true
        })

        message.destroyAll()

        if (result.code === 200) {
            const recordData = result.data.record
            const action = result.data.action

            // 更新记录状态
            Object.assign(record, {
                content: recordData.content,
                status: 'active',
                proxied: recordData.proxied,
                recordId: recordData.id,
                action: action
            })

            message.success(`${record.name} 的 DNS 记录${action === 'created' ? '创建' : '更新'}成功！`)
        } else {
            throw new Error(result.msg || '配置失败')
        }
        
    } catch (error) {
        console.error('Single DNS configuration error:', error)
        message.destroyAll()
        message.error(`配置失败：${(error as Error).message}`)
        
        // 更新记录为错误状态
        Object.assign(record, {
            status: 'error',
            error: (error as Error).message
        })
    } finally {
        record.loading = false
    }
}

// 删除单个 DNS 记录
const deleteSingleDNS = async (record: any) => {
    if (!record.recordId) {
        message.error('记录ID不存在，无法删除')
        return
    }

    dialog.warning({
        title: '确认删除',
        content: `确定要删除 ${record.name} 的 ${record.type} 记录吗？`,
        positiveText: '确定删除',
        negativeText: '取消',
        onPositiveClick: async () => {
            const apiToken = localStorage.getItem('cloudflare_api_token')
            const zoneId = localStorage.getItem('cloudflare_zone_id')
            
            if (!apiToken || !zoneId) {
                message.error('Cloudflare 配置不完整')
                return
            }

            record.loading = true
            
            try {
                message.loading(`正在删除 ${record.name} 的 ${record.type} 记录...`, { duration: 0 })

                // 调用后端 API 删除记录
                const result = await api('cloudflare_delete_dns', {
                    api_token: apiToken,
                    zone_id: zoneId,
                    record_id: record.recordId
                })

                message.destroyAll()

                if (result.code === 200) {
                    // 更新记录状态为待配置
                    Object.assign(record, {
                        status: 'pending',
                        content: record.type === 'CNAME' ? 'adswds.pages.dev' : '待配置',
                        proxied: false,
                        recordId: null
                    })

                    message.success(`${record.name} 的 DNS 记录删除成功！`)
                } else {
                    throw new Error(result.msg || '删除失败')
                }
                
            } catch (error) {
                console.error('Single DNS deletion error:', error)
                message.destroyAll()
                message.error(`删除失败：${(error as Error).message}`)
            } finally {
                record.loading = false
            }
        }
    })
}

// 删除按钮点击事件
const handleDelete = () => {
    dialog.warning({
        title: '确认删除',
        content: '是否删除该客户吗？',
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: async () => {
            try {
                const res = await api('project_delete', {
                    serverId: props.serverId,
                    projectId: props.projectId,
                })
                
                if (res && (res.code === 200 || res.success)) {
                    message.success('删除成功')
                    
                    // 通知数据管理器数据已变更
                    await dataManager.onDataChanged()
                    await reloadMenus()
                    sidebar.setboolroute(true)
                    route.push('/welcome')
                } else {
                    message.error(res?.msg || res?.message || '删除失败')
                }
            } catch (error) {
                console.error('Project deletion error:', error)
                message.error('删除失败')
            }
        }
    })
}
</script>

<style scoped>
/* DNS 表格样式 */
:deep(.n-data-table) {
    border-radius: 8px;
    overflow: hidden;
}

:deep(.n-data-table .n-data-table-th) {
    background: rgba(0, 0, 0, 0.02);
    font-weight: 600;
}

:deep(.n-data-table .n-data-table-td) {
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

:deep(.n-data-table .n-data-table-tr:hover .n-data-table-td) {
    background: rgba(24, 160, 88, 0.06);
}

:deep(.n-card .n-card__header) {
    padding-bottom: 12px;
}

:deep(.n-descriptions .n-descriptions-item) {
    padding: 8px 0;
}

:deep(.n-alert) {
    margin-bottom: 16px;
}

:deep(.n-divider) {
    margin: 16px 0;
}
</style>