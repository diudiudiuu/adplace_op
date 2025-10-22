<template>
    <div>
        <!-- 第一块：项目信息 -->
        <n-card title="📋 项目信息">
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
                                    <n-icon>
                                        <TrashOutline />
                                    </n-icon>
                                </template>
                                删除项目
                            </n-button>
                        </template>
                        删除项目
                    </n-tooltip>
                </n-space>
            </template>

            <Dform v-if="eidtmode" mode="edit" :serverId="serverId" :initialForm="projectInfo"
                @editSuccess="updateHandle" />

            <n-descriptions v-if="!eidtmode" :column="2" bordered>
                <n-descriptions-item label="项目ID">
                    <n-tag type="info">{{ projectInfo.project_id }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="项目名称">
                    <n-text strong>{{ projectInfo.project_name }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="管理地址">
                    <n-text type="info">{{ projectInfo.project_manage_url }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="API地址">
                    <n-text type="info">{{ projectInfo.project_api_url }}</n-text>
                </n-descriptions-item>
                <n-descriptions-item label="API端口">
                    <n-tag type="success">{{ projectInfo.api_port || '9000' }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="前端端口">
                    <n-tag type="success">{{ projectInfo.front_port || '3000' }}</n-tag>
                </n-descriptions-item>
            </n-descriptions>
        </n-card>

        <!-- 第二块：DNS 设置 -->
        <n-card v-if="!eidtmode" title="🌐 DNS 设置" style="margin-top: 16px;">
            <template #header-extra>
                <n-space>
                    <n-tooltip>
                        <template #trigger>
                            <n-button type="info" @click="showCloudflareConfig" size="small">
                                <template #icon>
                                    <n-icon>
                                        <SettingsOutline />
                                    </n-icon>
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
                                    <n-icon>
                                        <CloudOutline />
                                    </n-icon>
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
                        <n-icon>
                            <InformationCircleOutline />
                        </n-icon>
                        DNS 配置说明
                    </template>
                    将为以下域名配置 Cloudflare DNS 记录，所有记录都会开启代理（黄色小云朵）
                    <br>• 管理端域名：自动配置 CNAME 到 Pages 项目，并添加为 Pages 自定义域名
                    <br>• API 端域名：配置 A 记录直接指向服务器 IP
                </n-alert>

                <!-- DNS 记录表格 -->
                <div v-if="dnsStatus.length > 0">
                    <n-data-table :columns="dnsColumns" :data="dnsStatus" :pagination="false" striped size="small" />
                </div>
            </n-space>
        </n-card>

        <!-- Cloudflare 配置对话框 -->
        <n-modal v-model:show="showConfigModal" preset="dialog" title="Cloudflare 配置" style="width: 500px;">
            <n-form :model="cloudflareConfig" label-placement="left" label-width="120">
                <n-form-item label="API Token" required>
                    <n-input v-model:value="cloudflareConfig.apiToken" type="password"
                        placeholder="请输入 Cloudflare API Token" show-password-on="click" />
                </n-form-item>
                <n-form-item label="Zone ID" required>
                    <n-input v-model:value="cloudflareConfig.zoneId" placeholder="请输入域名的 Zone ID" />
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

        <!-- 第三块：项目部署 -->
        <n-card v-if="!eidtmode" title="🚀 项目部署" style="margin-top: 16px;">

            <n-space vertical size="large">
                <!-- 操作说明 -->
                <n-alert type="info" :show-icon="false">
                    <template #header>
                        <n-space align="center">
                            <n-icon>
                                <InformationCircleOutline />
                            </n-icon>
                            <span>部署操作说明</span>
                        </n-space>
                    </template>
                    以下三个操作相互独立，可根据需要单独执行，无需按顺序操作。
                    <br><strong>重要：</strong>所有部署操作都使用当前界面显示的项目数据，确保部署配置与界面一致。
                </n-alert>

                <!-- 独立操作功能卡片 -->
                <n-grid :cols="3" :x-gap="16" :y-gap="16">
                    <n-grid-item>
                        <n-card size="small" hoverable :bordered="false"
                            style="border: 2px solid #18a058; background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);">
                            <template #header>
                                <n-space align="center" justify="space-between">
                                    <n-space align="center">
                                        <n-icon size="24" color="#18a058">
                                            <DocumentOutline />
                                        </n-icon>
                                        <span style="font-weight: 600;">生成配置</span>
                                    </n-space>
                                    <n-button size="small" type="success" @click="generateCurrentProjectConfig"
                                        :loading="configLoading">
                                        执行
                                    </n-button>
                                </n-space>
                            </template>
                            <n-text depth="2" style="font-size: 13px; line-height: 1.5;">
                                • 检查并处理 release.zip 发布包<br>
                                • 使用当前项目数据生成配置文件<br>
                                • 只包含 api_port, web_port, api_domain 字段
                            </n-text>
                        </n-card>
                    </n-grid-item>
                    <n-grid-item>
                        <n-card size="small" hoverable :bordered="false"
                            style="border: 2px solid #2080f0; background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);">
                            <template #header>
                                <n-space align="center" justify="space-between">
                                    <n-space align="center">
                                        <n-icon size="24" color="#2080f0">
                                            <RocketOutline />
                                        </n-icon>
                                        <span style="font-weight: 600;">初始化项目</span>
                                    </n-space>
                                    <n-button size="small" type="primary" @click="executeInitCurrentProject"
                                        :loading="initLoading">
                                        执行
                                    </n-button>
                                </n-space>
                            </template>
                            <n-text depth="2" style="font-size: 13px; line-height: 1.5;">
                                • 初始化当前项目<br>
                                • 使用当前项目数据首次部署到服务器<br>
                                • 执行初始化脚本
                            </n-text>
                        </n-card>
                    </n-grid-item>
                    <n-grid-item>
                        <n-card size="small" hoverable :bordered="false"
                            style="border: 2px solid #f0a020; background: linear-gradient(135deg, #fffbf0 0%, #fef3c7 100%);">
                            <template #header>
                                <n-space align="center" justify="space-between">
                                    <n-space align="center">
                                        <n-icon size="24" color="#f0a020">
                                            <RefreshOutline />
                                        </n-icon>
                                        <span style="font-weight: 600;">更新项目</span>
                                    </n-space>
                                    <n-button size="small" type="warning" @click="executeUpdateCurrentProject"
                                        :loading="updateLoading">
                                        执行
                                    </n-button>
                                </n-space>
                            </template>
                            <n-text depth="2" style="font-size: 13px; line-height: 1.5;">
                                • 更新当前项目<br>
                                • 使用当前项目数据更新已部署的项目<br>
                                • 应用最新代码和配置
                            </n-text>
                        </n-card>
                    </n-grid-item>
                </n-grid>

                <!-- 操作状态显示 -->
                <div v-if="deploymentStatus">
                    <n-alert :type="deploymentStatus.type" :title="deploymentStatus.title" closable>
                        {{ deploymentStatus.message }}
                        <template #icon>
                            <n-icon>
                                <component :is="deploymentStatus.icon" />
                            </n-icon>
                        </template>
                    </n-alert>
                </div>

                <!-- 项目配置预览（简化显示） -->
                <div v-if="projectConfigPreview">
                    <n-card size="small" title="配置文件预览" style="margin-top: 16px;">
                        <template #header-extra>
                            <n-tag type="success" size="small">已生成</n-tag>
                        </template>
                        <n-scrollbar style="max-height: 200px;">
                            <pre style="font-size: 12px; line-height: 1.4; margin: 0;">{{ projectConfigPreview }}</pre>
                        </n-scrollbar>
                    </n-card>
                </div>
            </n-space>
        </n-card>

        <!-- 初始化项目选择对话框 -->
        <n-modal v-model:show="showInitProjectModal" preset="dialog" title="选择要初始化的项目" style="width: 500px;">
            <n-form label-placement="left" label-width="100">
                <n-form-item label="选择项目">
                    <n-select v-model:value="selectedInitProjectId" :options="serverProjects.map(p => ({
                        label: `${p.project_name} (${p.project_id})`,
                        value: p.project_id,
                        disabled: false
                    }))" placeholder="请选择要初始化的项目" clearable filterable />
                </n-form-item>
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="showInitProjectModal = false">取消</n-button>
                    <n-button type="primary" @click="executeInitProject" :loading="initLoading"
                        :disabled="!selectedInitProjectId">
                        开始初始化
                    </n-button>
                </n-space>
            </template>
        </n-modal>

        <!-- 更新项目选择对话框 -->
        <n-modal v-model:show="showUpdateProjectModal" preset="dialog" title="选择要更新的项目" style="width: 500px;">
            <n-form label-placement="left" label-width="100">
                <n-form-item label="选择项目">
                    <n-select v-model:value="selectedUpdateProjectId" :options="serverProjects.map(p => ({
                        label: `${p.project_name} (${p.project_id})`,
                        value: p.project_id,
                        disabled: false
                    }))" placeholder="请选择要更新的项目" clearable filterable />
                </n-form-item>
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="showUpdateProjectModal = false">取消</n-button>
                    <n-button type="warning" @click="executeUpdateProject" :loading="updateLoading"
                        :disabled="!selectedUpdateProjectId">
                        开始更新
                    </n-button>
                </n-space>
            </template>
        </n-modal>
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NButton, NIcon, NSpace, NTooltip, NTag, NText, NGrid, NGridItem, NScrollbar, NSelect, NModal, NForm, NFormItem } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import dataManager from '@/utils/dataManager'

import { CreateOutline, CloseOutline, TrashOutline, CloudOutline, InformationCircleOutline, PlayOutline, RefreshOutline, TrashBinOutline, SettingsOutline, DocumentOutline, RocketOutline, CheckmarkCircleOutline, AlertCircleOutline, TimeOutline } from '@vicons/ionicons5'
import Dform from './form.vue'
import api from '@/api'
import { getAuthorization } from '@/utils/auth'

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

// 项目部署相关状态
const configLoading = ref(false)
const initLoading = ref(false)
const updateLoading = ref(false)
const projectConfigPreview = ref('')
const deploymentStatus = ref<{
    type: 'success' | 'warning' | 'error' | 'info'
    title: string
    message: string
    icon: any
} | null>(null)

// 项目选择相关状态（保留以防其他地方使用）
const selectedInitProjectId = ref('')
const selectedUpdateProjectId = ref('')
const serverProjects = ref<any[]>([])
const showInitProjectModal = ref(false)
const showUpdateProjectModal = ref(false)

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
                active: { type: 'success' as const, text: '✓ 已配置' },
                pending: { type: 'warning' as const, text: '○ 待配置' },
                error: { type: 'error' as const, text: '✗ 错误' }
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
    // 加载服务器项目列表
    loadServerProjects()
})

// 获取服务器下的所有项目
const loadServerProjects = async () => {
    try {
        const server = await dataManager.getServerById(props.serverId)
        if (server && server.ProjectList) {
            serverProjects.value = server.ProjectList
        }
    } catch (error) {
        console.error('Failed to load server projects:', error)
    }
}

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
                    // 如果是 CNAME 记录，同时删除 Pages 自定义域名
                    if (record.type === 'CNAME') {
                        try {
                            console.log(`正在删除 Pages 自定义域名: ${record.name}`)
                            const pagesDeleteResult = await api('cloudflare_pages_delete_domain', {
                                api_token: apiToken,
                                zone_id: zoneId,
                                project_name: 'adswds',
                                domain: record.name
                            })

                            if (pagesDeleteResult.code === 200) {
                                console.log('Pages 自定义域名删除成功')
                            } else {
                                console.warn('Pages 自定义域名删除失败:', pagesDeleteResult.msg)
                            }
                        } catch (pagesError) {
                            console.warn('Pages 自定义域名删除出错:', pagesError)
                        }
                    }

                    // 更新记录状态为待配置
                    Object.assign(record, {
                        status: 'pending',
                        content: record.type === 'CNAME' ? 'adswds.pages.dev' : '待配置',
                        proxied: false,
                        recordId: null
                    })

                    const deleteMessage = record.type === 'CNAME'
                        ? `${record.name} 的 DNS 记录和 Pages 自定义域名删除成功！`
                        : `${record.name} 的 DNS 记录删除成功！`
                    message.success(deleteMessage)
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

// 生成当前项目配置文件 - 前端生成JSON，后端只负责上传
const generateCurrentProjectConfig = async () => {
    configLoading.value = true
    deploymentStatus.value = null

    try {
        message.loading('正在生成当前项目配置并上传到服务器...', { duration: 0 })

        // 获取当前服务器的完整数据
        const serverData = await dataManager.getServerById(props.serverId)
        if (!serverData) {
            throw new Error('无法获取服务器数据')
        }

        // 验证当前项目是否存在
        const currentProject = serverData.project_list?.find(p => p.project_id === props.projectId)
        if (!currentProject) {
            throw new Error('当前项目不存在于服务器数据中')
        }

        // 提取API域名
        const extractDomain = (url: string): string => {
            if (!url) return ''
            try {
                const urlObj = new URL(url.startsWith('http') ? url : `https://${url}`)
                return urlObj.hostname
            } catch {
                return url.replace(/^https?:\/\//, '').split('/')[0]
            }
        }

        // 前端生成JSON配置 - 只包含当前项目和指定字段
        const projectConfig = {
            [props.projectId]: {
                api_port: currentProject.api_port || '9000',
                web_port: currentProject.front_port || '3000',
                api_domain: extractDomain(currentProject.project_api_url)
            }
        }

        const projectConfigJson = JSON.stringify(projectConfig, null, 2)

        console.log('前端生成的项目配置:', {
            projectId: props.projectId,
            config: projectConfig,
            json: projectConfigJson
        })

        console.log('服务器数据结构:', {
            serverData: serverData,
            serverDataJson: JSON.stringify(serverData),
            keys: Object.keys(serverData)
        })

        // 调用后端API，只负责上传配置文件
        const result = await api('upload_project_config', {
            server_data_json: JSON.stringify(serverData),
            project_config_json: projectConfigJson,
            authorization: getAuthorization()
        })

        message.destroyAll()

        if (result.code === 200) {
            projectConfigPreview.value = projectConfigJson
            
            deploymentStatus.value = {
                type: 'success',
                title: '当前项目配置上传成功',
                message: `已生成并上传当前项目 ${currentProject.project_name} (${props.projectId}) 的配置文件。只包含 api_port, web_port, api_domain 三个字段。`,
                icon: CheckmarkCircleOutline
            }
            message.success(`当前项目配置文件上传成功`)
            
            console.log('配置上传成功:', result.data)
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '配置上传失败',
                message: result.msg || '上传配置文件失败',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '上传配置文件失败')
        }
    } catch (error) {
        console.error('Upload project config error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '配置上传异常',
            message: '上传配置文件时发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('上传配置文件失败：' + (error as Error).message)
    } finally {
        configLoading.value = false
    }
}

// 初始化项目
const initProject = async () => {
    if (!projectInfo.value.project_id) {
        message.error('项目ID不能为空')
        return
    }

    initLoading.value = true
    deploymentStatus.value = {
        type: 'info',
        title: '正在初始化项目',
        message: `正在为项目 ${projectInfo.value.project_id} 执行初始化操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在初始化项目 ${projectInfo.value.project_id}...`, { duration: 0 })

        const result = await api('project_init', {
            server_id: props.serverId,
            project_id: projectInfo.value.project_id
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目初始化成功',
                message: `项目 ${projectInfo.value.project_id} 已成功初始化，可以开始使用了`,
                icon: CheckmarkCircleOutline
            }
            message.success('项目初始化成功')
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目初始化失败',
                message: result.msg || '初始化过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目初始化失败')
        }
    } catch (error) {
        console.error('Project init error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '初始化异常',
            message: '初始化过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目初始化失败：' + (error as Error).message)
    } finally {
        initLoading.value = false
    }
}

// 执行初始化当前项目 - 使用前端当前数据
const executeInitCurrentProject = async () => {
    initLoading.value = true

    deploymentStatus.value = {
        type: 'info',
        title: '正在初始化项目',
        message: `正在为当前项目 ${projectInfo.value.project_name} (${props.projectId}) 执行初始化操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在初始化当前项目 ${projectInfo.value.project_name}...`, { duration: 0 })

        // 获取当前服务器的完整数据
        const serverData = await dataManager.getServerById(props.serverId)
        if (!serverData) {
            throw new Error('无法获取服务器数据')
        }

        console.log('使用前端数据初始化当前项目:', {
            serverId: serverData.server_id,
            projectId: props.projectId,
            projectName: projectInfo.value.project_name,
            serverIP: serverData.server_ip,
            defaultPath: serverData.default_path
        })

        // 使用新的API，传入序列化的服务器数据
        const result = await api('project_init_with_data', {
            server_id: props.serverId,
            project_id: props.projectId,
            server_data_json: JSON.stringify(serverData)
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目初始化成功',
                message: `当前项目 ${projectInfo.value.project_name} (${props.projectId}) 已成功初始化。使用当前界面数据，确保配置一致。`,
                icon: CheckmarkCircleOutline
            }
            message.success('当前项目初始化成功')
            
            console.log('初始化成功:', result.data)
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目初始化失败',
                message: result.msg || '初始化过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目初始化失败')
        }
    } catch (error) {
        console.error('Project init error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '初始化异常',
            message: '初始化过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目初始化失败：' + (error as Error).message)
    } finally {
        initLoading.value = false
    }
}

// 执行更新当前项目 - 使用前端当前数据
const executeUpdateCurrentProject = async () => {
    updateLoading.value = true

    deploymentStatus.value = {
        type: 'info',
        title: '正在更新项目',
        message: `正在为当前项目 ${projectInfo.value.project_name} (${props.projectId}) 执行更新操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在更新当前项目 ${projectInfo.value.project_name}...`, { duration: 0 })

        // 获取当前服务器的完整数据
        const serverData = await dataManager.getServerById(props.serverId)
        if (!serverData) {
            throw new Error('无法获取服务器数据')
        }

        console.log('使用前端数据更新当前项目:', {
            serverId: serverData.server_id,
            projectId: props.projectId,
            projectName: projectInfo.value.project_name,
            serverIP: serverData.server_ip,
            defaultPath: serverData.default_path
        })

        // 使用新的API，传入序列化的服务器数据
        const result = await api('project_update_with_data', {
            server_id: props.serverId,
            project_id: props.projectId,
            server_data_json: JSON.stringify(serverData)
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目更新成功',
                message: `当前项目 ${projectInfo.value.project_name} (${props.projectId}) 已成功更新到最新版本。使用当前界面数据，确保配置一致。`,
                icon: CheckmarkCircleOutline
            }
            message.success('当前项目更新成功')
            
            console.log('更新成功:', result.data)
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目更新失败',
                message: result.msg || '更新过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目更新失败')
        }
    } catch (error) {
        console.error('Project update error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '更新异常',
            message: '更新过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目更新失败：' + (error as Error).message)
    } finally {
        updateLoading.value = false
    }
}

// 执行初始化项目 - 使用前端当前数据（保留原方法以防其他地方调用）
const executeInitProject = async () => {
    if (!selectedInitProjectId.value) {
        message.error('请选择要初始化的项目')
        return
    }

    initLoading.value = true
    showInitProjectModal.value = false

    const selectedProject = serverProjects.value.find(p => p.project_id === selectedInitProjectId.value)
    const projectName = selectedProject ? selectedProject.project_name : selectedInitProjectId.value

    deploymentStatus.value = {
        type: 'info',
        title: '正在初始化项目',
        message: `正在为项目 ${projectName} (${selectedInitProjectId.value}) 执行初始化操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在初始化项目 ${projectName}...`, { duration: 0 })

        // 获取当前服务器的完整数据
        const serverData = await dataManager.getServerById(props.serverId)
        if (!serverData) {
            throw new Error('无法获取服务器数据')
        }

        console.log('使用前端数据初始化项目:', {
            serverId: serverData.server_id,
            projectId: selectedInitProjectId.value,
            projectName: projectName,
            serverIP: serverData.server_ip,
            defaultPath: serverData.default_path
        })

        // 使用新的API，传入序列化的服务器数据
        const result = await api('project_init_with_data', {
            server_id: props.serverId,
            project_id: selectedInitProjectId.value,
            server_data_json: JSON.stringify(serverData)
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目初始化成功',
                message: `项目 ${projectName} (${selectedInitProjectId.value}) 已成功初始化。使用当前界面数据，确保配置一致。`,
                icon: CheckmarkCircleOutline
            }
            message.success('项目初始化成功')
            
            console.log('初始化成功:', result.data)
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目初始化失败',
                message: result.msg || '初始化过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目初始化失败')
        }
    } catch (error) {
        console.error('Project init error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '初始化异常',
            message: '初始化过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目初始化失败：' + (error as Error).message)
    } finally {
        initLoading.value = false
        selectedInitProjectId.value = ''
    }
}

// 更新项目
const updateProject = async () => {
    if (!projectInfo.value.project_id) {
        message.error('项目ID不能为空')
        return
    }

    updateLoading.value = true
    deploymentStatus.value = {
        type: 'info',
        title: '正在更新项目',
        message: `正在为项目 ${projectInfo.value.project_id} 执行更新操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在更新项目 ${projectInfo.value.project_id}...`, { duration: 0 })

        const result = await api('project_update', {
            server_id: props.serverId,
            project_id: projectInfo.value.project_id
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目更新成功',
                message: `项目 ${projectInfo.value.project_id} 已成功更新到最新版本`,
                icon: CheckmarkCircleOutline
            }
            message.success('项目更新成功')
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目更新失败',
                message: result.msg || '更新过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目更新失败')
        }
    } catch (error) {
        console.error('Project update error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '更新异常',
            message: '更新过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目更新失败：' + (error as Error).message)
    } finally {
        updateLoading.value = false
    }
}

// 执行更新项目 - 使用前端当前数据
const executeUpdateProject = async () => {
    if (!selectedUpdateProjectId.value) {
        message.error('请选择要更新的项目')
        return
    }

    updateLoading.value = true
    showUpdateProjectModal.value = false

    const selectedProject = serverProjects.value.find(p => p.project_id === selectedUpdateProjectId.value)
    const projectName = selectedProject ? selectedProject.project_name : selectedUpdateProjectId.value

    deploymentStatus.value = {
        type: 'info',
        title: '正在更新项目',
        message: `正在为项目 ${projectName} (${selectedUpdateProjectId.value}) 执行更新操作...`,
        icon: TimeOutline
    }

    try {
        message.loading(`正在更新项目 ${projectName}...`, { duration: 0 })

        // 获取当前服务器的完整数据
        const serverData = await dataManager.getServerById(props.serverId)
        if (!serverData) {
            throw new Error('无法获取服务器数据')
        }

        console.log('使用前端数据更新项目:', {
            serverId: serverData.server_id,
            projectId: selectedUpdateProjectId.value,
            projectName: projectName,
            serverIP: serverData.server_ip,
            defaultPath: serverData.default_path
        })

        // 使用新的API，传入序列化的服务器数据
        const result = await api('project_update_with_data', {
            server_id: props.serverId,
            project_id: selectedUpdateProjectId.value,
            server_data_json: JSON.stringify(serverData)
        })

        message.destroyAll()

        if (result.code === 200) {
            deploymentStatus.value = {
                type: 'success',
                title: '项目更新成功',
                message: `项目 ${projectName} (${selectedUpdateProjectId.value}) 已成功更新到最新版本。使用当前界面数据，确保配置一致。`,
                icon: CheckmarkCircleOutline
            }
            message.success('项目更新成功')
            
            console.log('更新成功:', result.data)
        } else {
            deploymentStatus.value = {
                type: 'error',
                title: '项目更新失败',
                message: result.msg || '更新过程中发生错误，请检查服务器配置',
                icon: AlertCircleOutline
            }
            message.error(result.msg || '项目更新失败')
        }
    } catch (error) {
        console.error('Project update error:', error)
        message.destroyAll()
        deploymentStatus.value = {
            type: 'error',
            title: '更新异常',
            message: '更新过程中发生异常，请检查网络连接和服务器状态',
            icon: AlertCircleOutline
        }
        message.error('项目更新失败：' + (error as Error).message)
    } finally {
        updateLoading.value = false
        selectedUpdateProjectId.value = ''
    }
}
</script>

<style scoped>
/* 页面整体样式 */
.n-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
}

.n-card:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    transition: box-shadow 0.3s ease;
}

/* DNS 表格样式 */
:deep(.n-data-table) {
    border-radius: 8px;
    overflow: hidden;
}

:deep(.n-data-table .n-data-table-th) {
    background: rgba(0, 0, 0, 0.02);
    font-weight: 600;
}

/* 部署卡片样式 */
:deep(.n-grid-item .n-card .n-card__header) {
    padding-bottom: 8px;
}

/* 小卡片悬停效果 */
:deep(.n-grid-item .n-card) {
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

:deep(.n-grid-item .n-card:hover) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 配置预览样式 */
pre {
    background: #f8f9fa;
    border-radius: 6px;
    padding: 12px;
    border: 1px solid #e9ecef;
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