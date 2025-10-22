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
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { useSidebarStore } from '@/store/sidebar'
import { reloadMenus } from '@/components/menu'
import dataManager from '@/utils/dataManager'
import { CreateOutline, CloseOutline, TrashOutline } from '@vicons/ionicons5'
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

const eidtmode = ref(false)
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

getProjectInfo()
const updateHandle = () => {
    eidtmode.value = false
    getProjectInfo()
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